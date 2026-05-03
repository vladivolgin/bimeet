package collectionrepo

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

func (r *Repository) Create(ctx context.Context, eventID uuid.UUID, req model.CreateCollectionRequest, createdBy uuid.UUID) (model.Collection, error) {
	var c model.Collection
	err := r.pool.QueryRow(ctx,
		`INSERT INTO collections (event_id, title, per_person_amount, created_by)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, event_id, title, per_person_amount, created_by, created_at`,
		eventID, req.Title, req.PerPersonAmount, createdBy,
	).Scan(&c.ID, &c.EventID, &c.Title, &c.PerPersonAmount, &c.CreatedBy, &c.CreatedAt)
	if err != nil {
		return model.Collection{}, fmt.Errorf("create collection: %w", err)
	}
	return c, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (model.Collection, error) {
	var c model.Collection
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_id, title, per_person_amount, created_by, created_at FROM collections WHERE id=$1`,
		id,
	).Scan(&c.ID, &c.EventID, &c.Title, &c.PerPersonAmount, &c.CreatedBy, &c.CreatedAt)
	if err != nil {
		return model.Collection{}, fmt.Errorf("get collection: %w", err)
	}
	return c, nil
}

func (r *Repository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Collection, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, title, per_person_amount, created_by, created_at
		 FROM collections WHERE event_id=$1 ORDER BY created_at`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}
	defer rows.Close()

	var cols []model.Collection
	for rows.Next() {
		var c model.Collection
		if err := rows.Scan(&c.ID, &c.EventID, &c.Title, &c.PerPersonAmount, &c.CreatedBy, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan collection: %w", err)
		}
		cols = append(cols, c)
	}
	return cols, rows.Err()
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM collections WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	return nil
}

func scanContribution(row interface {
	Scan(...any) error
}) (model.CollectionContribution, error) {
	var cc model.CollectionContribution
	err := row.Scan(&cc.ID, &cc.CollectionID, &cc.UserID, &cc.Paid, &cc.PaidAt, &cc.Status, &cc.ReceiptURL)
	return cc, err
}

func (r *Repository) GetContribution(ctx context.Context, collectionID, userID uuid.UUID) (model.CollectionContribution, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, collection_id, user_id, paid, paid_at, status, receipt_url
		 FROM collection_contributions
		 WHERE collection_id=$1 AND user_id=$2`,
		collectionID, userID,
	)
	cc, err := scanContribution(row)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("get contribution: %w", err)
	}
	return cc, nil
}

func (r *Repository) GetContributionByID(ctx context.Context, id uuid.UUID) (model.CollectionContribution, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, collection_id, user_id, paid, paid_at, status, receipt_url
		 FROM collection_contributions WHERE id=$1`,
		id,
	)
	cc, err := scanContribution(row)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("get contribution by id: %w", err)
	}
	return cc, nil
}

// SubmitContribution sets status=pending and saves receipt URL.
// Only allowed when current status is not_paid.
func (r *Repository) SubmitContribution(ctx context.Context, collectionID, userID uuid.UUID, receiptURL string) (model.CollectionContribution, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO collection_contributions (collection_id, user_id, paid, status, receipt_url)
		 VALUES ($1, $2, FALSE, 'pending', $3)
		 ON CONFLICT (collection_id, user_id) DO UPDATE
		   SET status = 'pending',
		       receipt_url = EXCLUDED.receipt_url,
		       paid = FALSE,
		       paid_at = NULL
		 RETURNING id, collection_id, user_id, paid, paid_at, status, receipt_url`,
		collectionID, userID, receiptURL,
	)
	cc, err := scanContribution(row)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("submit contribution: %w", err)
	}
	return cc, nil
}

// ConfirmContribution marks a pending contribution as paid.
func (r *Repository) ConfirmContribution(ctx context.Context, id uuid.UUID) (model.CollectionContribution, error) {
	row := r.pool.QueryRow(ctx,
		`UPDATE collection_contributions
		 SET status = 'paid', paid = TRUE, paid_at = NOW()
		 WHERE id = $1
		 RETURNING id, collection_id, user_id, paid, paid_at, status, receipt_url`,
		id,
	)
	cc, err := scanContribution(row)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("confirm contribution: %w", err)
	}
	return cc, nil
}

// RejectContribution resets a pending contribution back to not_paid.
func (r *Repository) RejectContribution(ctx context.Context, id uuid.UUID) (model.CollectionContribution, error) {
	row := r.pool.QueryRow(ctx,
		`UPDATE collection_contributions
		 SET status = 'not_paid', paid = FALSE, paid_at = NULL, receipt_url = NULL
		 WHERE id = $1
		 RETURNING id, collection_id, user_id, paid, paid_at, status, receipt_url`,
		id,
	)
	cc, err := scanContribution(row)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("reject contribution: %w", err)
	}
	return cc, nil
}

// MarkPaid forcibly marks any contribution as paid (organizer action, no receipt needed).
func (r *Repository) MarkPaid(ctx context.Context, collectionID, userID uuid.UUID) (model.CollectionContribution, error) {
	row := r.pool.QueryRow(ctx,
		`INSERT INTO collection_contributions (collection_id, user_id, paid, paid_at, status)
		 VALUES ($1, $2, TRUE, NOW(), 'paid')
		 ON CONFLICT (collection_id, user_id) DO UPDATE
		   SET status = 'paid', paid = TRUE, paid_at = NOW(), receipt_url = collection_contributions.receipt_url
		 RETURNING id, collection_id, user_id, paid, paid_at, status, receipt_url`,
		collectionID, userID,
	)
	cc, err := scanContribution(row)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("mark paid: %w", err)
	}
	return cc, nil
}

func (r *Repository) ListContributions(ctx context.Context, collectionID uuid.UUID) ([]model.CollectionContribution, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, collection_id, user_id, paid, paid_at, status, receipt_url
		 FROM collection_contributions WHERE collection_id=$1`,
		collectionID,
	)
	if err != nil {
		return nil, fmt.Errorf("list contributions: %w", err)
	}
	defer rows.Close()

	var contribs []model.CollectionContribution
	for rows.Next() {
		cc, err := scanContribution(rows)
		if err != nil {
			return nil, fmt.Errorf("scan contribution: %w", err)
		}
		contribs = append(contribs, cc)
	}
	return contribs, rows.Err()
}

func (r *Repository) CountPaidContributions(ctx context.Context, collectionID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM collection_contributions WHERE collection_id=$1 AND status='paid'`,
		collectionID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count paid contributions: %w", err)
	}
	return count, nil
}
