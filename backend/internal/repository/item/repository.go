package itemrepo

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

func (r *Repository) Create(ctx context.Context, eventID uuid.UUID, req model.CreateItemRequest) (model.Item, error) {
	var item model.Item
	err := r.pool.QueryRow(ctx,
		`INSERT INTO items (event_id, name) VALUES ($1, $2)
		 RETURNING id, event_id, name, assigned_to`,
		eventID, req.Name,
	).Scan(&item.ID, &item.EventID, &item.Name, &item.AssignedTo)
	if err != nil {
		return model.Item{}, fmt.Errorf("create item: %w", err)
	}
	return item, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (model.Item, error) {
	var item model.Item
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_id, name, assigned_to FROM items WHERE id=$1`,
		id,
	).Scan(&item.ID, &item.EventID, &item.Name, &item.AssignedTo)
	if err != nil {
		return model.Item{}, fmt.Errorf("get item: %w", err)
	}
	return item, nil
}

func (r *Repository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Item, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, name, assigned_to FROM items WHERE event_id=$1 ORDER BY id`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.EventID, &item.Name, &item.AssignedTo); err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) UpdateAssignment(ctx context.Context, id uuid.UUID, assignedTo *uuid.UUID) (model.Item, error) {
	var item model.Item
	err := r.pool.QueryRow(ctx,
		`UPDATE items SET assigned_to=$2 WHERE id=$1
		 RETURNING id, event_id, name, assigned_to`,
		id, assignedTo,
	).Scan(&item.ID, &item.EventID, &item.Name, &item.AssignedTo)
	if err != nil {
		return model.Item{}, fmt.Errorf("update item assignment: %w", err)
	}
	return item, nil
}
