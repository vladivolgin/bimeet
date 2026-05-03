package eventhandler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"bimeet/internal/handler/response"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
	eventsvc "bimeet/internal/service/event"
)

type Handler struct {
	svc EventService
}

func New(svc EventService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	events, err := h.svc.List(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list events")
		return
	}
	if events == nil {
		events = []model.EventListItem{}
	}
	response.JSON(w, http.StatusOK, events)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.CreateEventRequest
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.svc.Create(r.Context(), req, userID)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "required") || strings.Contains(msg, "must be") {
			response.Error(w, http.StatusBadRequest, msg)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create event")
		return
	}

	response.JSON(w, http.StatusCreated, event)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	event, err := h.svc.GetByID(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, eventsvc.ErrForbidden) || strings.Contains(err.Error(), "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		response.Error(w, http.StatusNotFound, "event not found")
		return
	}

	response.JSON(w, http.StatusOK, event)
}

func (h *Handler) GetByInviteToken(w http.ResponseWriter, r *http.Request) {
	tokenStr := chi.URLParam(r, "token")
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid invite token")
		return
	}

	event, err := h.svc.GetByInviteToken(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusNotFound, "event not found")
		return
	}

	response.JSON(w, http.StatusOK, event)
}

func (h *Handler) JoinByInviteToken(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tokenStr := chi.URLParam(r, "token")
	token, err := uuid.Parse(tokenStr)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid invite token")
		return
	}

	var req model.JoinByInviteTokenRequest
	// Body is optional — ignore parse errors (empty body = join)
	_ = response.ParseJSON(r, &req)

	participant, err := h.svc.JoinByInviteToken(r.Context(), token, userID, req)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "organizer") {
			response.Error(w, http.StatusConflict, msg)
			return
		}
		if strings.Contains(msg, "action") {
			response.Error(w, http.StatusBadRequest, msg)
			return
		}
		response.Error(w, http.StatusNotFound, "event not found")
		return
	}

	response.JSON(w, http.StatusOK, participant)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	var req model.UpdateEventRequest
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event, err := h.svc.Update(r.Context(), id, userID, req)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		if strings.Contains(msg, "not found") {
			response.Error(w, http.StatusNotFound, "event not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update event")
		return
	}

	response.JSON(w, http.StatusOK, event)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	if err := h.svc.Delete(r.Context(), id, userID); err != nil {
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		if strings.Contains(msg, "not found") {
			response.Error(w, http.StatusNotFound, "event not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete event")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Complete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	if err := h.svc.Complete(r.Context(), id, userID); err != nil {
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		if strings.Contains(msg, "not found") {
			response.Error(w, http.StatusNotFound, "event not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to complete event")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) InviteParticipant(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	eventID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	var req model.InviteParticipantRequest
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.svc.InviteParticipant(r.Context(), eventID, userID, req)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		if strings.Contains(msg, "not found") {
			response.Error(w, http.StatusNotFound, msg)
			return
		}
		response.Error(w, http.StatusBadRequest, msg)
		return
	}

	response.JSON(w, http.StatusCreated, p)
}

func (h *Handler) UpdateParticipantStatus(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	eventID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	targetUserID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req model.UpdateParticipantStatusRequest
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	participant, err := h.svc.UpdateParticipantStatus(r.Context(), eventID, targetUserID, userID, req)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		if strings.Contains(msg, "not found") {
			response.Error(w, http.StatusNotFound, msg)
			return
		}
		response.Error(w, http.StatusBadRequest, msg)
		return
	}

	response.JSON(w, http.StatusOK, participant)
}

func (h *Handler) ListPublic(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	events, err := h.svc.ListPublic(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list public events")
		return
	}
	if events == nil {
		events = []model.PublicEventListItem{}
	}
	response.JSON(w, http.StatusOK, events)
}

func (h *Handler) JoinPublic(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	eventID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	if err := h.svc.JoinPublic(r.Context(), eventID, userID); err != nil {
		if errors.Is(err, eventsvc.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		if errors.Is(err, eventsvc.ErrFull) {
			response.Error(w, http.StatusConflict, "event is full")
			return
		}
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "joined"})
}
