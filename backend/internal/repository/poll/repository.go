package pollrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"bimeet/internal/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, eventID uuid.UUID, req model.CreatePollRequest, createdBy uuid.UUID) (model.Poll, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return model.Poll{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var p model.Poll
	err = tx.QueryRow(ctx,
		`INSERT INTO polls (event_id, question, created_by)
		 VALUES ($1, $2, $3)
		 RETURNING id, event_id, question, created_by, created_at`,
		eventID, req.Question, createdBy,
	).Scan(&p.ID, &p.EventID, &p.Question, &p.CreatedBy, &p.CreatedAt)
	if err != nil {
		return model.Poll{}, fmt.Errorf("insert poll: %w", err)
	}

	for _, label := range req.Options {
		var opt model.PollOption
		err = tx.QueryRow(ctx,
			`INSERT INTO poll_options (poll_id, label) VALUES ($1, $2)
			 RETURNING id, poll_id, label`,
			p.ID, label,
		).Scan(&opt.ID, &opt.PollID, &opt.Label)
		if err != nil {
			return model.Poll{}, fmt.Errorf("insert poll option: %w", err)
		}
		p.Options = append(p.Options, opt)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Poll{}, fmt.Errorf("commit poll tx: %w", err)
	}
	return p, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (model.Poll, error) {
	var p model.Poll
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_id, question, created_by, created_at FROM polls WHERE id=$1`,
		id,
	).Scan(&p.ID, &p.EventID, &p.Question, &p.CreatedBy, &p.CreatedAt)
	if err != nil {
		return model.Poll{}, fmt.Errorf("get poll: %w", err)
	}

	opts, err := r.listOptionsWithVotes(ctx, p.ID)
	if err != nil {
		return model.Poll{}, err
	}
	p.Options = opts
	return p, nil
}

func (r *Repository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Poll, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, question, created_by, created_at FROM polls WHERE event_id=$1 ORDER BY created_at`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("list polls: %w", err)
	}
	defer rows.Close()

	var polls []model.Poll
	for rows.Next() {
		var p model.Poll
		if err := rows.Scan(&p.ID, &p.EventID, &p.Question, &p.CreatedBy, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan poll: %w", err)
		}
		opts, err := r.listOptionsWithVotes(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		p.Options = opts
		polls = append(polls, p)
	}
	return polls, rows.Err()
}

func (r *Repository) listOptionsWithVotes(ctx context.Context, pollID uuid.UUID) ([]model.PollOption, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT po.id, po.poll_id, po.label, COUNT(pv.user_id) AS votes
		 FROM poll_options po
		 LEFT JOIN poll_votes pv ON pv.poll_option_id = po.id
		 WHERE po.poll_id=$1
		 GROUP BY po.id, po.poll_id, po.label
		 ORDER BY po.id`,
		pollID,
	)
	if err != nil {
		return nil, fmt.Errorf("list options with votes: %w", err)
	}
	defer rows.Close()

	var opts []model.PollOption
	for rows.Next() {
		var o model.PollOption
		if err := rows.Scan(&o.ID, &o.PollID, &o.Label, &o.Votes); err != nil {
			return nil, fmt.Errorf("scan poll option: %w", err)
		}
		opts = append(opts, o)
	}
	return opts, rows.Err()
}

func (r *Repository) GetOption(ctx context.Context, optionID uuid.UUID) (model.PollOption, error) {
	var o model.PollOption
	err := r.pool.QueryRow(ctx,
		`SELECT id, poll_id, label FROM poll_options WHERE id=$1`,
		optionID,
	).Scan(&o.ID, &o.PollID, &o.Label)
	if err != nil {
		return model.PollOption{}, fmt.Errorf("get poll option: %w", err)
	}
	return o, nil
}

func (r *Repository) Vote(ctx context.Context, optionID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM poll_votes
		 WHERE user_id=$1
		   AND poll_option_id IN (
		       SELECT id FROM poll_options WHERE poll_id=(
		           SELECT poll_id FROM poll_options WHERE id=$2
		       )
		   )`,
		userID, optionID,
	)
	if err != nil {
		return fmt.Errorf("delete previous vote: %w", err)
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO poll_votes (poll_option_id, user_id) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`,
		optionID, userID,
	)
	if err != nil {
		return fmt.Errorf("insert vote: %w", err)
	}
	return nil
}
