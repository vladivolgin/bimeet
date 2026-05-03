package eventlinkrepo

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

func (r *Repository) Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreateEventLinkRequest) (model.EventLink, error) {
	var l model.EventLink
	err := r.pool.QueryRow(ctx,
		`INSERT INTO event_links (event_id, title, url, created_by)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, event_id, title, url, created_by, created_at`,
		eventID, req.Title, req.URL, userID,
	).Scan(&l.ID, &l.EventID, &l.Title, &l.URL, &l.CreatedBy, &l.CreatedAt)
	if err != nil {
		return model.EventLink{}, fmt.Errorf("create event link: %w", err)
	}
	return l, nil
}

func (r *Repository) List(ctx context.Context, eventID uuid.UUID) ([]model.EventLink, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, title, url, created_by, created_at
		 FROM event_links WHERE event_id=$1 ORDER BY created_at`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("list event links: %w", err)
	}
	defer rows.Close()

	var links []model.EventLink
	for rows.Next() {
		var l model.EventLink
		if err := rows.Scan(&l.ID, &l.EventID, &l.Title, &l.URL, &l.CreatedBy, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event link: %w", err)
		}
		links = append(links, l)
	}
	return links, rows.Err()
}

func (r *Repository) Delete(ctx context.Context, linkID, eventID uuid.UUID) error {
	result, err := r.pool.Exec(ctx,
		`DELETE FROM event_links WHERE id=$1 AND event_id=$2`,
		linkID, eventID,
	)
	if err != nil {
		return fmt.Errorf("delete event link: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("event link not found")
	}
	return nil
}
