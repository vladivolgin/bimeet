package pollsvc_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"bimeet/internal/model"
	pollsvc "bimeet/internal/service/poll"
	"bimeet/internal/service/poll/mock"
)

func newService(t *testing.T) (*pollsvc.Service, *mock.MockPollRepo, *mock.MockEventRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	pr := mock.NewMockPollRepo(ctrl)
	er := mock.NewMockEventRepo(ctrl)
	return pollsvc.New(pr, er), pr, er
}

// ─── List ──────────────────────────────────────────────────────────────────

func TestList_Forbidden(t *testing.T) {
	svc, _, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(false, nil)

	_, err := svc.List(context.Background(), eventID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestList_HappyPath(t *testing.T) {
	svc, pr, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	pr.EXPECT().ListByEvent(gomock.Any(), eventID).Return([]model.Poll{}, nil)

	_, err := svc.List(context.Background(), eventID, userID)
	require.NoError(t, err)
}

// ─── Create ────────────────────────────────────────────────────────────────

func TestCreate_Forbidden(t *testing.T) {
	svc, _, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: uuid.New()}, nil)

	_, err := svc.Create(context.Background(), eventID, userID, model.CreatePollRequest{Question: "Q?", Options: []string{"A", "B"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestCreate_MissingQuestion(t *testing.T) {
	svc, _, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID}, nil)

	_, err := svc.Create(context.Background(), eventID, organizerID, model.CreatePollRequest{Options: []string{"A", "B"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "question is required")
}

func TestCreate_TooFewOptions(t *testing.T) {
	svc, _, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID}, nil)

	_, err := svc.Create(context.Background(), eventID, organizerID, model.CreatePollRequest{Question: "Q?", Options: []string{"A"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "at least 2 options")
}

func TestCreate_HappyPath(t *testing.T) {
	svc, pr, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	req := model.CreatePollRequest{Question: "Where?", Options: []string{"Paris", "Berlin"}}
	expected := model.Poll{ID: uuid.New(), EventID: eventID}

	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID}, nil)
	pr.EXPECT().Create(gomock.Any(), eventID, req, organizerID).Return(expected, nil)

	result, err := svc.Create(context.Background(), eventID, organizerID, req)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

// ─── Vote ──────────────────────────────────────────────────────────────────

func TestVote_NotConfirmed(t *testing.T) {
	svc, _, er := newService(t)
	eventID, pollID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(false, nil)

	err := svc.Vote(context.Background(), eventID, pollID, userID, model.VoteRequest{OptionID: uuid.New()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestVote_OptionWrongPoll(t *testing.T) {
	svc, pr, er := newService(t)
	eventID, pollID, userID := uuid.New(), uuid.New(), uuid.New()
	optionID := uuid.New()

	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	pr.EXPECT().GetByID(gomock.Any(), pollID).Return(model.Poll{ID: pollID, EventID: eventID}, nil)
	pr.EXPECT().GetOption(gomock.Any(), optionID).Return(model.PollOption{PollID: uuid.New()}, nil) // wrong poll

	err := svc.Vote(context.Background(), eventID, pollID, userID, model.VoteRequest{OptionID: optionID})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong")
}

func TestVote_HappyPath(t *testing.T) {
	svc, pr, er := newService(t)
	eventID, pollID, userID := uuid.New(), uuid.New(), uuid.New()
	optionID := uuid.New()

	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	pr.EXPECT().GetByID(gomock.Any(), pollID).Return(model.Poll{ID: pollID, EventID: eventID}, nil)
	pr.EXPECT().GetOption(gomock.Any(), optionID).Return(model.PollOption{ID: optionID, PollID: pollID}, nil)
	pr.EXPECT().Vote(gomock.Any(), optionID, userID).Return(nil)

	require.NoError(t, svc.Vote(context.Background(), eventID, pollID, userID, model.VoteRequest{OptionID: optionID}))
}

func TestVote_PollWrongEvent(t *testing.T) {
	svc, pr, er := newService(t)
	eventID, pollID, userID := uuid.New(), uuid.New(), uuid.New()

	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	pr.EXPECT().GetByID(gomock.Any(), pollID).Return(model.Poll{ID: pollID, EventID: uuid.New()}, nil)

	err := svc.Vote(context.Background(), eventID, pollID, userID, model.VoteRequest{OptionID: uuid.New()})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong")
}

