package carpoolrepo

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

func (r *Repository) Create(ctx context.Context, eventID, driverID uuid.UUID, req model.CreateCarpoolRequest) (model.Carpool, error) {
	var c model.Carpool
	err := r.pool.QueryRow(ctx,
		`INSERT INTO carpools (event_id, driver_id, seats_available, departure_point)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, event_id, driver_id, seats_available, departure_point`,
		eventID, driverID, req.SeatsAvailable, req.DeparturePoint,
	).Scan(&c.ID, &c.EventID, &c.DriverID, &c.SeatsAvailable, &c.DeparturePoint)
	if err != nil {
		return model.Carpool{}, fmt.Errorf("create carpool: %w", err)
	}
	return c, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (model.Carpool, error) {
	var c model.Carpool
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_id, driver_id, seats_available, departure_point FROM carpools WHERE id=$1`,
		id,
	).Scan(&c.ID, &c.EventID, &c.DriverID, &c.SeatsAvailable, &c.DeparturePoint)
	if err != nil {
		return model.Carpool{}, fmt.Errorf("get carpool: %w", err)
	}
	return c, nil
}

func (r *Repository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Carpool, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, driver_id, seats_available, departure_point FROM carpools WHERE event_id=$1 ORDER BY id`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("list carpools: %w", err)
	}
	defer rows.Close()

	var carpools []model.Carpool
	for rows.Next() {
		var c model.Carpool
		if err := rows.Scan(&c.ID, &c.EventID, &c.DriverID, &c.SeatsAvailable, &c.DeparturePoint); err != nil {
			return nil, fmt.Errorf("scan carpool: %w", err)
		}
		carpools = append(carpools, c)
	}
	return carpools, rows.Err()
}

func (r *Repository) AddPassenger(ctx context.Context, carpoolID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO carpool_passengers (carpool_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		carpoolID, userID,
	)
	if err != nil {
		return fmt.Errorf("add carpool passenger: %w", err)
	}
	return nil
}

func (r *Repository) CountPassengers(ctx context.Context, carpoolID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM carpool_passengers WHERE carpool_id=$1`,
		carpoolID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count carpool passengers: %w", err)
	}
	return count, nil
}

func (r *Repository) IsPassenger(ctx context.Context, carpoolID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM carpool_passengers WHERE carpool_id=$1 AND user_id=$2)`,
		carpoolID, userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("is carpool passenger: %w", err)
	}
	return exists, nil
}
