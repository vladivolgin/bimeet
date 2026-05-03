package eventlinksvc

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type Service struct {
	linkRepo  EventLinkRepo
	eventRepo EventRepo
}

func New(linkRepo EventLinkRepo, eventRepo EventRepo) *Service {
	return &Service{linkRepo: linkRepo, eventRepo: eventRepo}
}

func (s *Service) List(ctx context.Context, eventID, userID uuid.UUID) ([]model.EventLink, error) {
	ok, err := s.eventRepo.IsParticipant(ctx, eventID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("forbidden")
	}
	links, err := s.linkRepo.List(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if links == nil {
		links = []model.EventLink{}
	}
	return links, nil
}

func (s *Service) Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreateEventLinkRequest) (model.EventLink, error) {
	if req.Title == "" {
		return model.EventLink{}, fmt.Errorf("title is required")
	}
	if req.URL == "" {
		return model.EventLink{}, fmt.Errorf("url is required")
	}

	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return model.EventLink{}, fmt.Errorf("event not found")
	}
	if event.OrganizerID != userID {
		return model.EventLink{}, fmt.Errorf("forbidden")
	}
	if event.Category != "business" {
		return model.EventLink{}, fmt.Errorf("links are only available for business events")
	}

	return s.linkRepo.Create(ctx, eventID, userID, req)
}

func (s *Service) Delete(ctx context.Context, eventID, linkID, userID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if event.OrganizerID != userID {
		return fmt.Errorf("forbidden")
	}

	if err := s.linkRepo.Delete(ctx, linkID, eventID); err != nil {
		if err.Error() == "event link not found" {
			return fmt.Errorf("link not found")
		}
		return err
	}
	return nil
}
