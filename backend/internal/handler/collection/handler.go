package collectionhandler

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"bimeet/internal/handler/response"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
)

type Handler struct {
	svc CollectionService
}

func New(svc CollectionService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
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

	collections, err := h.svc.List(r.Context(), eventID, userID)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to list collections")
		return
	}
	if collections == nil {
		collections = []model.Collection{}
	}
	response.JSON(w, http.StatusOK, collections)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	var req model.CreateCollectionRequest
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	col, err := h.svc.Create(r.Context(), eventID, userID, req)
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

	response.JSON(w, http.StatusCreated, col)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
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

	collectionID, err := uuid.Parse(chi.URLParam(r, "collectionId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid collection id")
		return
	}

	if err := h.svc.Delete(r.Context(), eventID, collectionID, userID); err != nil {
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		if strings.Contains(msg, "not found") {
			response.Error(w, http.StatusNotFound, msg)
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete collection")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

const maxReceiptSize = 10 << 20 // 10 MB

func (h *Handler) SubmitContribution(w http.ResponseWriter, r *http.Request) {
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

	collectionID, err := uuid.Parse(chi.URLParam(r, "collectionId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid collection id")
		return
	}

	if err := r.ParseMultipartForm(maxReceiptSize); err != nil {
		response.Error(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("receipt")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "receipt file is required")
		return
	}
	defer file.Close()

	contrib, err := h.svc.SubmitContribution(r.Context(), eventID, collectionID, userID, file, header)
	if err != nil {
		log.Printf("submit contribution: %v", err)
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, msg)
			return
		}
		if strings.Contains(msg, "not found") {
			response.Error(w, http.StatusNotFound, msg)
			return
		}
		if strings.Contains(msg, "already paid") {
			response.Error(w, http.StatusConflict, msg)
			return
		}
		if strings.Contains(msg, "unsupported file type") {
			response.Error(w, http.StatusBadRequest, msg)
			return
		}
		if strings.Contains(msg, "not configured") {
			response.Error(w, http.StatusServiceUnavailable, "receipt storage is not configured")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to submit contribution")
		return
	}

	response.JSON(w, http.StatusOK, contrib)
}

func (h *Handler) ConfirmContribution(w http.ResponseWriter, r *http.Request) {
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

	collectionID, err := uuid.Parse(chi.URLParam(r, "collectionId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid collection id")
		return
	}

	contribID, err := uuid.Parse(chi.URLParam(r, "contribId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid contribution id")
		return
	}

	contrib, err := h.svc.ConfirmContribution(r.Context(), eventID, collectionID, contribID, userID)
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

	response.JSON(w, http.StatusOK, contrib)
}

func (h *Handler) RejectContribution(w http.ResponseWriter, r *http.Request) {
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

	collectionID, err := uuid.Parse(chi.URLParam(r, "collectionId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid collection id")
		return
	}

	contribID, err := uuid.Parse(chi.URLParam(r, "contribId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid contribution id")
		return
	}

	contrib, err := h.svc.RejectContribution(r.Context(), eventID, collectionID, contribID, userID)
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

	response.JSON(w, http.StatusOK, contrib)
}

type markPaidRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

func (h *Handler) MarkPaid(w http.ResponseWriter, r *http.Request) {
	organizerID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	eventID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid event id")
		return
	}

	collectionID, err := uuid.Parse(chi.URLParam(r, "collectionId"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid collection id")
		return
	}

	var req markPaidRequest
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UserID == uuid.Nil {
		response.Error(w, http.StatusBadRequest, "user_id is required")
		return
	}

	contrib, err := h.svc.MarkPaid(r.Context(), eventID, collectionID, req.UserID, organizerID)
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

	response.JSON(w, http.StatusOK, contrib)
}

func (h *Handler) Summary(w http.ResponseWriter, r *http.Request) {
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

	summary, err := h.svc.Summary(r.Context(), eventID, userID)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "forbidden") {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get summary")
		return
	}

	response.JSON(w, http.StatusOK, summary)
}
