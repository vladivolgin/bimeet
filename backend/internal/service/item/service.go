package itemsvc

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type Service struct {
	itemRepo  ItemRepo
	eventRepo EventRepo
}

func New(itemRepo ItemRepo, eventRepo EventRepo) *Service {
	return &Service{itemRepo: itemRepo, eventRepo: eventRepo}
}

func (s *Service) List(ctx context.Context, eventID, userID uuid.UUID) ([]model.Item, error) {
	ok, err := s.eventRepo.IsParticipant(ctx, eventID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("forbidden")
	}
	return s.itemRepo.ListByEvent(ctx, eventID)
}

func (s *Service) Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreateItemRequest) (model.Item, error) {
	ok, err := s.eventRepo.IsParticipant(ctx, eventID, userID)
	if err != nil {
		return model.Item{}, err
	}
	if !ok {
		return model.Item{}, fmt.Errorf("forbidden")
	}
	if req.Name == "" {
		return model.Item{}, fmt.Errorf("name is required")
	}
	return s.itemRepo.Create(ctx, eventID, req)
}

func (s *Service) UpdateAssignment(ctx context.Context, eventID, itemID, userID uuid.UUID, req model.UpdateItemRequest) (model.Item, error) {
	ok, err := s.eventRepo.IsConfirmedParticipant(ctx, eventID, userID)
	if err != nil {
		return model.Item{}, err
	}
	if !ok {
		return model.Item{}, fmt.Errorf("forbidden: must be a confirmed participant")
	}

	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return model.Item{}, fmt.Errorf("item not found")
	}
	if item.EventID != eventID {
		return model.Item{}, fmt.Errorf("item does not belong to this event")
	}

	// Participants may only assign/unassign themselves
	if req.AssignedTo != nil && *req.AssignedTo != userID {
		event, err := s.eventRepo.GetByID(ctx, eventID)
		if err != nil {
			return model.Item{}, err
		}
		if event.OrganizerID != userID {
			return model.Item{}, fmt.Errorf("forbidden: can only assign items to yourself")
		}
	}

	return s.itemRepo.UpdateAssignment(ctx, itemID, req.AssignedTo)
}
