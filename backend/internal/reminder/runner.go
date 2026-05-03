package reminder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Runner struct {
	eventRepo        EventRepo
	notificationRepo NotificationRepo
	interval         time.Duration
}

func New(eventRepo EventRepo, notificationRepo NotificationRepo, interval time.Duration) *Runner {
	return &Runner{
		eventRepo:        eventRepo,
		notificationRepo: notificationRepo,
		interval:         interval,
	}
}

func (r *Runner) Start(ctx context.Context) {
	r.run(ctx)
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.run(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (r *Runner) run(ctx context.Context) {
	events, err := r.eventRepo.ListForReminder(ctx)
	if err != nil {
		log.Printf("reminder: list events: %v", err)
		return
	}

	for _, e := range events {
		recipients := uniqueRecipients(e.ParticipantIDs, e.OrganizerID)

		if e.Needs3d {
			msg := fmt.Sprintf("Встреча «%s» заканчивается через 3 дня", e.Title)
			r.sendToAll(ctx, recipients, e.ID, "event_reminder_3d", msg)
			if err := r.eventRepo.MarkReminder3dSent(ctx, e.ID); err != nil {
				log.Printf("reminder: mark 3d sent for %s: %v", e.ID, err)
			}
		}

		if e.Needs1d {
			msg := fmt.Sprintf("Встреча «%s» заканчивается завтра", e.Title)
			r.sendToAll(ctx, recipients, e.ID, "event_reminder_1d", msg)
			if err := r.eventRepo.MarkReminder1dSent(ctx, e.ID); err != nil {
				log.Printf("reminder: mark 1d sent for %s: %v", e.ID, err)
			}
		}
	}
}

func (r *Runner) sendToAll(ctx context.Context, recipients []uuid.UUID, eventID uuid.UUID, notifType, msg string) {
	for _, uid := range recipients {
		if _, err := r.notificationRepo.Create(ctx, uid, &eventID, notifType, msg); err != nil {
			log.Printf("reminder: create notification for user %s: %v", uid, err)
		}
	}
}

func uniqueRecipients(participants []uuid.UUID, organizer uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(participants)+1)
	result := make([]uuid.UUID, 0, len(participants)+1)
	for _, id := range append(participants, organizer) {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			result = append(result, id)
		}
	}
	return result
}
