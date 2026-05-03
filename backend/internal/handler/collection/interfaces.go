package collectionhandler

//go:generate mockgen -destination=mock/mock.go -package=mock . CollectionService

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type CollectionService interface {
	List(ctx context.Context, eventID, userID uuid.UUID) ([]model.Collection, error)
	Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreateCollectionRequest) (model.Collection, error)
	Delete(ctx context.Context, eventID, collectionID, userID uuid.UUID) error
	SubmitContribution(ctx context.Context, eventID, collectionID, userID uuid.UUID, file io.Reader, header *multipart.FileHeader) (model.CollectionContribution, error)
	ConfirmContribution(ctx context.Context, eventID, collectionID, contributionID, organizerID uuid.UUID) (model.CollectionContribution, error)
	RejectContribution(ctx context.Context, eventID, collectionID, contributionID, organizerID uuid.UUID) (model.CollectionContribution, error)
	MarkPaid(ctx context.Context, eventID, collectionID, targetUserID, organizerID uuid.UUID) (model.CollectionContribution, error)
	Summary(ctx context.Context, eventID, userID uuid.UUID) (model.EventCollectionsSummaryResponse, error)
}
