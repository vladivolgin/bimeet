package eventrepo

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

func (r *Repository) Create(ctx context.Context, req model.CreateEventRequest, organizerID uuid.UUID) (model.Event, error) {
	var e model.Event
	err := r.pool.QueryRow(ctx,
		`INSERT INTO events (title, description, date_start, date_end, location, category, dress_code, is_public, max_guests, organizer_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, title, description, date_start, date_end, location, category, dress_code, is_public, max_guests, status, organizer_id, invite_token, created_at, updated_at`,
		req.Title, req.Description, req.DateStart, req.DateEnd, req.Location, req.Category, req.DressCode, req.IsPublic, req.MaxGuests, organizerID,
	).Scan(&e.ID, &e.Title, &e.Description, &e.DateStart, &e.DateEnd, &e.Location,
		&e.Category, &e.DressCode, &e.IsPublic, &e.MaxGuests, &e.Status, &e.OrganizerID, &e.InviteToken, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return model.Event{}, fmt.Errorf("create event: %w", err)
	}
	return e, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (model.Event, error) {
	var e model.Event
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, description, date_start, date_end, location, category, dress_code, is_public, max_guests, status, organizer_id, invite_token, created_at, updated_at
		 FROM events WHERE id=$1`,
		id,
	).Scan(&e.ID, &e.Title, &e.Description, &e.DateStart, &e.DateEnd, &e.Location,
		&e.Category, &e.DressCode, &e.IsPublic, &e.MaxGuests, &e.Status, &e.OrganizerID, &e.InviteToken, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return model.Event{}, fmt.Errorf("get event by id: %w", err)
	}
	return e, nil
}

func (r *Repository) GetByInviteToken(ctx context.Context, token uuid.UUID) (model.InviteEventInfo, error) {
	var e model.InviteEventInfo
	err := r.pool.QueryRow(ctx,
		`SELECT e.id, e.title, e.description, e.date_start, e.date_end, e.location,
		        e.category, e.dress_code, e.is_public,
		        u.id, u.name, u.email,
		        COUNT(ep.id) FILTER (WHERE ep.status = 'confirmed') AS confirmed_count
		 FROM events e
		 JOIN users u ON u.id = e.organizer_id
		 LEFT JOIN event_participants ep ON ep.event_id = e.id
		 WHERE e.invite_token = $1
		 GROUP BY e.id, u.id`,
		token,
	).Scan(&e.ID, &e.Title, &e.Description, &e.DateStart, &e.DateEnd, &e.Location,
		&e.Category, &e.DressCode, &e.IsPublic,
		&e.Organizer.ID, &e.Organizer.Name, &e.Organizer.Email,
		&e.ConfirmedCount)
	if err != nil {
		return model.InviteEventInfo{}, fmt.Errorf("get event by invite token: %w", err)
	}
	return e, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req model.UpdateEventRequest) (model.Event, error) {
	var e model.Event
	err := r.pool.QueryRow(ctx,
		`UPDATE events SET
		   title       = COALESCE($2, title),
		   description = COALESCE($3, description),
		   date_start  = COALESCE($4, date_start),
		   date_end    = COALESCE($5, date_end),
		   location    = COALESCE($6, location),
		   category    = COALESCE($7, category),
		   dress_code  = COALESCE($8, dress_code),
		   is_public   = COALESCE($9, is_public),
		   max_guests  = COALESCE($10, max_guests),
		   updated_at  = NOW()
		 WHERE id=$1
		 RETURNING id, title, description, date_start, date_end, location, category, dress_code, is_public, max_guests, status, organizer_id, invite_token, created_at, updated_at`,
		id, req.Title, req.Description, req.DateStart, req.DateEnd, req.Location, req.Category, req.DressCode, req.IsPublic, req.MaxGuests,
	).Scan(&e.ID, &e.Title, &e.Description, &e.DateStart, &e.DateEnd, &e.Location,
		&e.Category, &e.DressCode, &e.IsPublic, &e.MaxGuests, &e.Status, &e.OrganizerID, &e.InviteToken, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return model.Event{}, fmt.Errorf("update event: %w", err)
	}
	return e, nil
}

func (r *Repository) Complete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE events SET status='completed', updated_at=NOW() WHERE id=$1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("complete event: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM events WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	return nil
}

// ─── Participants ─────────────────────────────────────────────────────────

func (r *Repository) AddParticipant(ctx context.Context, eventID, userID uuid.UUID, status string) (model.EventParticipant, error) {
	var p model.EventParticipant
	err := r.pool.QueryRow(ctx,
		`INSERT INTO event_participants (event_id, user_id, status)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (event_id, user_id) DO UPDATE SET status = EXCLUDED.status
		 RETURNING id, event_id, user_id, status`,
		eventID, userID, status,
	).Scan(&p.ID, &p.EventID, &p.UserID, &p.Status)
	if err != nil {
		return model.EventParticipant{}, fmt.Errorf("add participant: %w", err)
	}
	return p, nil
}

func (r *Repository) UpdateParticipantStatus(ctx context.Context, eventID, userID uuid.UUID, status string) (model.EventParticipant, error) {
	var p model.EventParticipant
	err := r.pool.QueryRow(ctx,
		`UPDATE event_participants SET status=$3
		 WHERE event_id=$1 AND user_id=$2
		 RETURNING id, event_id, user_id, status`,
		eventID, userID, status,
	).Scan(&p.ID, &p.EventID, &p.UserID, &p.Status)
	if err != nil {
		return model.EventParticipant{}, fmt.Errorf("update participant status: %w", err)
	}
	return p, nil
}

func (r *Repository) GetParticipant(ctx context.Context, eventID, userID uuid.UUID) (model.EventParticipant, error) {
	var p model.EventParticipant
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_id, user_id, status FROM event_participants WHERE event_id=$1 AND user_id=$2`,
		eventID, userID,
	).Scan(&p.ID, &p.EventID, &p.UserID, &p.Status)
	if err != nil {
		return model.EventParticipant{}, fmt.Errorf("get participant: %w", err)
	}
	return p, nil
}

func (r *Repository) ListParticipants(ctx context.Context, eventID uuid.UUID) ([]model.EventParticipant, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, user_id, status FROM event_participants WHERE event_id=$1`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("list participants: %w", err)
	}
	defer rows.Close()

	var participants []model.EventParticipant
	for rows.Next() {
		var p model.EventParticipant
		if err := rows.Scan(&p.ID, &p.EventID, &p.UserID, &p.Status); err != nil {
			return nil, fmt.Errorf("scan participant: %w", err)
		}
		participants = append(participants, p)
	}
	return participants, rows.Err()
}

func (r *Repository) CountConfirmedParticipants(ctx context.Context, eventID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM (
			SELECT organizer_id AS user_id FROM events WHERE id=$1
			UNION
			SELECT user_id FROM event_participants WHERE event_id=$1 AND status='confirmed'
		) t`,
		eventID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count confirmed participants: %w", err)
	}
	return count, nil
}

// ─── Change log ───────────────────────────────────────────────────────────

func (r *Repository) AddChangeLog(ctx context.Context, eventID, changedBy uuid.UUID, field, oldVal, newVal string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO event_change_logs (event_id, changed_by, field_name, old_value, new_value)
		 VALUES ($1, $2, $3, $4, $5)`,
		eventID, changedBy, field, oldVal, newVal,
	)
	if err != nil {
		return fmt.Errorf("add change log: %w", err)
	}
	return nil
}

func (r *Repository) ListChangeLogs(ctx context.Context, eventID uuid.UUID) ([]model.EventChangeLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, event_id, changed_by, field_name, old_value, new_value, changed_at
		 FROM event_change_logs WHERE event_id=$1 ORDER BY changed_at`,
		eventID,
	)
	if err != nil {
		return nil, fmt.Errorf("list change logs: %w", err)
	}
	defer rows.Close()

	var logs []model.EventChangeLog
	for rows.Next() {
		var l model.EventChangeLog
		if err := rows.Scan(&l.ID, &l.EventID, &l.ChangedBy, &l.FieldName, &l.OldValue, &l.NewValue, &l.ChangedAt); err != nil {
			return nil, fmt.Errorf("scan change log: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// GetDetail loads an event with all related data.
func (r *Repository) GetDetail(ctx context.Context, eventID, userID uuid.UUID) (model.EventDetail, error) {
	var detail model.EventDetail
	var orgID uuid.UUID
	err := r.pool.QueryRow(ctx,
		`SELECT e.id, e.title, e.description, e.date_start, e.date_end, e.location,
		        e.category, e.dress_code, e.status, e.is_public, e.max_guests,
		        e.organizer_id, u.name, u.email,
		        e.invite_token, e.created_at, e.updated_at
		 FROM events e
		 JOIN users u ON u.id = e.organizer_id
		 WHERE e.id = $1`,
		eventID,
	).Scan(
		&detail.ID, &detail.Title, &detail.Description, &detail.DateStart, &detail.DateEnd, &detail.Location,
		&detail.Category, &detail.DressCode, &detail.Status, &detail.IsPublic, &detail.MaxGuests,
		&orgID, &detail.Organizer.Name, &detail.Organizer.Email,
		&detail.InviteToken, &detail.CreatedAt, &detail.UpdatedAt,
	)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("get event detail: %w", err)
	}
	detail.Organizer.ID = orgID

	detail.Participants = []model.ParticipantInfo{}
	detail.Collections = []model.CollectionInfo{}
	detail.Polls = []model.PollInfo{}
	detail.Items = []model.ItemInfo{}
	detail.Carpools = []model.CarpoolInfo{}
	detail.Links = []model.EventLinkInfo{}
	detail.ChangeLog = []model.ChangeLogInfo{}

	// Participants
	pRows, err := r.pool.Query(ctx,
		`SELECT ep.id, ep.user_id, u.name, u.email, ep.status
		 FROM event_participants ep
		 JOIN users u ON u.id = ep.user_id
		 WHERE ep.event_id = $1`,
		eventID,
	)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("list participants detail: %w", err)
	}
	for pRows.Next() {
		var p model.ParticipantInfo
		if err := pRows.Scan(&p.ID, &p.User.ID, &p.User.Name, &p.User.Email, &p.Status); err != nil {
			pRows.Close()
			return model.EventDetail{}, fmt.Errorf("scan participant detail: %w", err)
		}
		detail.Participants = append(detail.Participants, p)
	}
	pRows.Close()
	if err := pRows.Err(); err != nil {
		return model.EventDetail{}, err
	}

	// ConfirmedCount: organizer + confirmed participants
	detail.ConfirmedCount = 1
	for _, p := range detail.Participants {
		if p.Status == "confirmed" {
			detail.ConfirmedCount++
		}
	}

	// Collections
	colRows, err := r.pool.Query(ctx,
		`SELECT id, title, per_person_amount, created_by, created_at FROM collections WHERE event_id = $1 ORDER BY created_at`,
		eventID,
	)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("list collections detail: %w", err)
	}
	for colRows.Next() {
		var c model.CollectionInfo
		if err := colRows.Scan(&c.ID, &c.Title, &c.PerPersonAmount, &c.CreatedBy, &c.CreatedAt); err != nil {
			colRows.Close()
			return model.EventDetail{}, fmt.Errorf("scan collection detail: %w", err)
		}
		detail.Collections = append(detail.Collections, c)
	}
	colRows.Close()
	if err := colRows.Err(); err != nil {
		return model.EventDetail{}, err
	}
	for i, col := range detail.Collections {
		ctRows, err := r.pool.Query(ctx,
			`SELECT id, user_id, paid, paid_at, status, receipt_url FROM collection_contributions WHERE collection_id = $1`,
			col.ID,
		)
		if err != nil {
			return model.EventDetail{}, fmt.Errorf("list contributions detail: %w", err)
		}
		detail.Collections[i].Contributions = []model.ContributionInfo{}
		for ctRows.Next() {
			var ct model.ContributionInfo
			if err := ctRows.Scan(&ct.ID, &ct.UserID, &ct.Paid, &ct.PaidAt, &ct.Status, &ct.ReceiptURL); err != nil {
				ctRows.Close()
				return model.EventDetail{}, fmt.Errorf("scan contribution detail: %w", err)
			}
			detail.Collections[i].Contributions = append(detail.Collections[i].Contributions, ct)
		}
		ctRows.Close()
		if err := ctRows.Err(); err != nil {
			return model.EventDetail{}, err
		}
	}

	// Polls
	pollRows, err := r.pool.Query(ctx,
		`SELECT id, question, created_by, created_at FROM polls WHERE event_id = $1 ORDER BY created_at`,
		eventID,
	)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("list polls detail: %w", err)
	}
	for pollRows.Next() {
		var p model.PollInfo
		if err := pollRows.Scan(&p.ID, &p.Question, &p.CreatedBy, &p.CreatedAt); err != nil {
			pollRows.Close()
			return model.EventDetail{}, fmt.Errorf("scan poll detail: %w", err)
		}
		detail.Polls = append(detail.Polls, p)
	}
	pollRows.Close()
	if err := pollRows.Err(); err != nil {
		return model.EventDetail{}, err
	}
	for i, poll := range detail.Polls {
		optRows, err := r.pool.Query(ctx,
			`SELECT po.id, po.label FROM poll_options po WHERE po.poll_id = $1 ORDER BY po.id`,
			poll.ID,
		)
		if err != nil {
			return model.EventDetail{}, fmt.Errorf("list poll options detail: %w", err)
		}
		detail.Polls[i].Options = []model.PollOptionInfo{}
		for optRows.Next() {
			var opt model.PollOptionInfo
			if err := optRows.Scan(&opt.ID, &opt.Label); err != nil {
				optRows.Close()
				return model.EventDetail{}, fmt.Errorf("scan poll option detail: %w", err)
			}
			detail.Polls[i].Options = append(detail.Polls[i].Options, opt)
		}
		optRows.Close()
		if err := optRows.Err(); err != nil {
			return model.EventDetail{}, err
		}
		for j, opt := range detail.Polls[i].Options {
			voteRows, err := r.pool.Query(ctx,
				`SELECT user_id FROM poll_votes WHERE poll_option_id = $1`,
				opt.ID,
			)
			if err != nil {
				return model.EventDetail{}, fmt.Errorf("list poll votes detail: %w", err)
			}
			detail.Polls[i].Options[j].Votes = []uuid.UUID{}
			for voteRows.Next() {
				var voterID uuid.UUID
				if err := voteRows.Scan(&voterID); err != nil {
					voteRows.Close()
					return model.EventDetail{}, fmt.Errorf("scan poll vote detail: %w", err)
				}
				detail.Polls[i].Options[j].Votes = append(detail.Polls[i].Options[j].Votes, voterID)
			}
			voteRows.Close()
			if err := voteRows.Err(); err != nil {
				return model.EventDetail{}, err
			}
		}
	}

	// Items
	itemRows, err := r.pool.Query(ctx,
		`SELECT i.id, i.name, i.assigned_to, u.name, u.email
		 FROM items i
		 LEFT JOIN users u ON u.id = i.assigned_to
		 WHERE i.event_id = $1 ORDER BY i.id`,
		eventID,
	)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("list items detail: %w", err)
	}
	for itemRows.Next() {
		var item model.ItemInfo
		var assignedID *uuid.UUID
		var assignedName, assignedEmail *string
		if err := itemRows.Scan(&item.ID, &item.Name, &assignedID, &assignedName, &assignedEmail); err != nil {
			itemRows.Close()
			return model.EventDetail{}, fmt.Errorf("scan item detail: %w", err)
		}
		if assignedID != nil {
			item.AssignedTo = &model.UserInfo{
				ID:    *assignedID,
				Name:  *assignedName,
				Email: *assignedEmail,
			}
		}
		detail.Items = append(detail.Items, item)
	}
	itemRows.Close()
	if err := itemRows.Err(); err != nil {
		return model.EventDetail{}, err
	}

	// Carpools
	cpRows, err := r.pool.Query(ctx,
		`SELECT c.id, c.driver_id, u.name, u.email, c.seats_available, c.departure_point
		 FROM carpools c
		 JOIN users u ON u.id = c.driver_id
		 WHERE c.event_id = $1 ORDER BY c.id`,
		eventID,
	)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("list carpools detail: %w", err)
	}
	for cpRows.Next() {
		var cp model.CarpoolInfo
		if err := cpRows.Scan(&cp.ID, &cp.Driver.ID, &cp.Driver.Name, &cp.Driver.Email, &cp.SeatsAvailable, &cp.DeparturePoint); err != nil {
			cpRows.Close()
			return model.EventDetail{}, fmt.Errorf("scan carpool detail: %w", err)
		}
		detail.Carpools = append(detail.Carpools, cp)
	}
	cpRows.Close()
	if err := cpRows.Err(); err != nil {
		return model.EventDetail{}, err
	}
	for i, cp := range detail.Carpools {
		passRows, err := r.pool.Query(ctx,
			`SELECT u.id, u.name, u.email FROM carpool_passengers cp JOIN users u ON u.id = cp.user_id WHERE cp.carpool_id = $1`,
			cp.ID,
		)
		if err != nil {
			return model.EventDetail{}, fmt.Errorf("list carpool passengers detail: %w", err)
		}
		detail.Carpools[i].Passengers = []model.UserInfo{}
		for passRows.Next() {
			var passenger model.UserInfo
			if err := passRows.Scan(&passenger.ID, &passenger.Name, &passenger.Email); err != nil {
				passRows.Close()
				return model.EventDetail{}, fmt.Errorf("scan carpool passenger detail: %w", err)
			}
			detail.Carpools[i].Passengers = append(detail.Carpools[i].Passengers, passenger)
		}
		passRows.Close()
		if err := passRows.Err(); err != nil {
			return model.EventDetail{}, err
		}
	}

	// Links
	linkRows, err := r.pool.Query(ctx,
		`SELECT id, title, url, created_by, created_at FROM event_links WHERE event_id = $1 ORDER BY created_at`,
		eventID,
	)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("list links detail: %w", err)
	}
	for linkRows.Next() {
		var l model.EventLinkInfo
		if err := linkRows.Scan(&l.ID, &l.Title, &l.URL, &l.CreatedBy, &l.CreatedAt); err != nil {
			linkRows.Close()
			return model.EventDetail{}, fmt.Errorf("scan link detail: %w", err)
		}
		detail.Links = append(detail.Links, l)
	}
	linkRows.Close()
	if err := linkRows.Err(); err != nil {
		return model.EventDetail{}, err
	}

	// Change log
	clRows, err := r.pool.Query(ctx,
		`SELECT cl.id, cl.changed_by, u.name, u.email, cl.field_name, cl.old_value, cl.new_value, cl.changed_at
		 FROM event_change_logs cl
		 JOIN users u ON u.id = cl.changed_by
		 WHERE cl.event_id = $1 ORDER BY cl.changed_at`,
		eventID,
	)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("list change logs detail: %w", err)
	}
	for clRows.Next() {
		var cl model.ChangeLogInfo
		if err := clRows.Scan(&cl.ID, &cl.ChangedBy.ID, &cl.ChangedBy.Name, &cl.ChangedBy.Email, &cl.FieldName, &cl.OldValue, &cl.NewValue, &cl.ChangedAt); err != nil {
			clRows.Close()
			return model.EventDetail{}, fmt.Errorf("scan change log detail: %w", err)
		}
		detail.ChangeLog = append(detail.ChangeLog, cl)
	}
	clRows.Close()
	if err := clRows.Err(); err != nil {
		return model.EventDetail{}, err
	}

	return detail, nil
}

// ListForUserEnriched lists events with organizer info and participants for a given user.
func (r *Repository) ListForUserEnriched(ctx context.Context, userID uuid.UUID) ([]model.EventListItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT e.id, e.title, e.description, e.date_start, e.date_end, e.location,
		        e.category, e.dress_code, e.status, e.is_public, e.max_guests,
		        e.organizer_id, ou.name, ou.email,
		        e.invite_token, e.created_at,
		        COALESCE(ep2.status, '') AS my_status
		 FROM events e
		 JOIN users ou ON ou.id = e.organizer_id
		 LEFT JOIN event_participants ep ON ep.event_id = e.id AND ep.user_id = $1
		 LEFT JOIN event_participants ep2 ON ep2.event_id = e.id AND ep2.user_id = $1
		 WHERE e.organizer_id = $1 OR ep.user_id = $1
		 ORDER BY e.date_start`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list events enriched: %w", err)
	}
	defer rows.Close()

	var events []model.EventListItem
	for rows.Next() {
		var e model.EventListItem
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.DateStart, &e.DateEnd, &e.Location,
			&e.Category, &e.DressCode, &e.Status, &e.IsPublic, &e.MaxGuests,
			&e.Organizer.ID, &e.Organizer.Name, &e.Organizer.Email,
			&e.InviteToken, &e.CreatedAt, &e.MyStatus,
		); err != nil {
			return nil, fmt.Errorf("scan event list item: %w", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, e := range events {
		pRows, err := r.pool.Query(ctx,
			`SELECT ep.id, ep.user_id, u.name, u.email, ep.status
			 FROM event_participants ep
			 JOIN users u ON u.id = ep.user_id
			 WHERE ep.event_id = $1`,
			e.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("list participants for event list: %w", err)
		}
		events[i].Participants = []model.ParticipantInfo{}
		for pRows.Next() {
			var p model.ParticipantInfo
			if err := pRows.Scan(&p.ID, &p.User.ID, &p.User.Name, &p.User.Email, &p.Status); err != nil {
				pRows.Close()
				return nil, fmt.Errorf("scan participant list item: %w", err)
			}
			events[i].Participants = append(events[i].Participants, p)
		}
		pRows.Close()
		if err := pRows.Err(); err != nil {
			return nil, err
		}

		// ConfirmedCount: organizer + confirmed participants
		events[i].ConfirmedCount = 1
		for _, p := range events[i].Participants {
			if p.Status == "confirmed" {
				events[i].ConfirmedCount++
			}
		}
	}

	return events, nil
}

// IsParticipant checks whether the user is organizer OR has a participant record (any status).
func (r *Repository) IsParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM events WHERE id=$1 AND organizer_id=$2
			UNION ALL
			SELECT 1 FROM event_participants WHERE event_id=$1 AND user_id=$2
		)`,
		eventID, userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("is participant: %w", err)
	}
	return exists, nil
}

// IsConfirmedParticipant checks if user has confirmed status or is organizer.
func (r *Repository) IsConfirmedParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM events WHERE id=$1 AND organizer_id=$2
			UNION ALL
			SELECT 1 FROM event_participants WHERE event_id=$1 AND user_id=$2 AND status='confirmed'
		)`,
		eventID, userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("is confirmed participant: %w", err)
	}
	return exists, nil
}

func (r *Repository) ListForReminder(ctx context.Context) ([]model.ReminderEvent, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			e.id,
			e.title,
			e.date_end,
			e.organizer_id,
			COALESCE(array_agg(ep.user_id) FILTER (WHERE ep.user_id IS NOT NULL), '{}') AS participant_ids,
			(NOT e.reminder_3d_sent AND e.date_end <= NOW() + INTERVAL '3 days') AS needs_3d,
			(NOT e.reminder_1d_sent AND e.date_end <= NOW() + INTERVAL '1 day')  AS needs_1d
		FROM events e
		LEFT JOIN event_participants ep ON ep.event_id = e.id AND ep.status = 'confirmed'
		WHERE e.status = 'active'
		  AND e.date_end > NOW()
		  AND (
		    (NOT e.reminder_3d_sent AND e.date_end <= NOW() + INTERVAL '3 days')
		    OR
		    (NOT e.reminder_1d_sent AND e.date_end <= NOW() + INTERVAL '1 day')
		  )
		GROUP BY e.id, e.title, e.date_end, e.organizer_id, e.reminder_3d_sent, e.reminder_1d_sent
	`)
	if err != nil {
		return nil, fmt.Errorf("list for reminder: %w", err)
	}
	defer rows.Close()

	var result []model.ReminderEvent
	for rows.Next() {
		var e model.ReminderEvent
		if err := rows.Scan(&e.ID, &e.Title, &e.DateEnd, &e.OrganizerID, &e.ParticipantIDs, &e.Needs3d, &e.Needs1d); err != nil {
			return nil, fmt.Errorf("scan reminder event: %w", err)
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *Repository) MarkReminder3dSent(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE events SET reminder_3d_sent=TRUE WHERE id=$1`, id)
	return err
}

func (r *Repository) MarkReminder1dSent(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE events SET reminder_1d_sent=TRUE WHERE id=$1`, id)
	return err
}

// ListPublic returns active public events with confirmed_count and is_participant flag.
func (r *Repository) ListPublic(ctx context.Context, userID uuid.UUID) ([]model.PublicEventListItem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT e.id, e.title, e.description, e.date_start, e.date_end, e.location,
		        e.category, e.dress_code, e.status, e.is_public, e.max_guests,
		        e.organizer_id, ou.name, ou.email,
		        e.invite_token, e.created_at,
		        EXISTS(
		            SELECT 1 FROM event_participants ep
		            WHERE ep.event_id = e.id AND ep.user_id = $1
		        ) AS is_participant
		 FROM events e
		 JOIN users ou ON ou.id = e.organizer_id
		 WHERE e.is_public = TRUE AND e.status = 'active'
		 ORDER BY e.date_start`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list public events: %w", err)
	}
	defer rows.Close()

	var events []model.PublicEventListItem
	for rows.Next() {
		var e model.PublicEventListItem
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.DateStart, &e.DateEnd, &e.Location,
			&e.Category, &e.DressCode, &e.Status, &e.IsPublic, &e.MaxGuests,
			&e.Organizer.ID, &e.Organizer.Name, &e.Organizer.Email,
			&e.InviteToken, &e.CreatedAt,
			&e.IsParticipant,
		); err != nil {
			return nil, fmt.Errorf("scan public event list item: %w", err)
		}
		e.MyStatus = ""
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i, e := range events {
		pRows, err := r.pool.Query(ctx,
			`SELECT ep.id, ep.user_id, u.name, u.email, ep.status
			 FROM event_participants ep
			 JOIN users u ON u.id = ep.user_id
			 WHERE ep.event_id = $1`,
			e.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("list participants for public event: %w", err)
		}
		events[i].Participants = []model.ParticipantInfo{}
		for pRows.Next() {
			var p model.ParticipantInfo
			if err := pRows.Scan(&p.ID, &p.User.ID, &p.User.Name, &p.User.Email, &p.Status); err != nil {
				pRows.Close()
				return nil, fmt.Errorf("scan participant public event: %w", err)
			}
			events[i].Participants = append(events[i].Participants, p)
		}
		pRows.Close()
		if err := pRows.Err(); err != nil {
			return nil, err
		}

		events[i].ConfirmedCount = 1
		for _, p := range events[i].Participants {
			if p.Status == "confirmed" {
				events[i].ConfirmedCount++
			}
		}
	}

	return events, nil
}

// JoinPublic adds a user as a confirmed participant of a public event.
func (r *Repository) JoinPublic(ctx context.Context, eventID, userID uuid.UUID) (model.EventParticipant, error) {
	var p model.EventParticipant
	err := r.pool.QueryRow(ctx,
		`INSERT INTO event_participants (event_id, user_id, status)
		 VALUES ($1, $2, 'confirmed')
		 ON CONFLICT (event_id, user_id) DO UPDATE SET status = 'confirmed'
		 RETURNING id, event_id, user_id, status`,
		eventID, userID,
	).Scan(&p.ID, &p.EventID, &p.UserID, &p.Status)
	if err != nil {
		return model.EventParticipant{}, fmt.Errorf("join public event: %w", err)
	}
	return p, nil
}
