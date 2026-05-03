package profilehandler

import (
	"log"
	"net/http"
	"strings"

	"bimeet/internal/handler/response"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
)

type Handler struct {
	svc ProfileService
}

func New(svc ProfileService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.svc.GetMe(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get profile")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req model.UpdateProfileRequest
	if err := response.ParseJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required")
		return
	}

	user, err := h.svc.UpdateMe(r.Context(), userID, req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid birth_date") {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	stats, err := h.svc.GetStats(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get stats")
		return
	}

	response.JSON(w, http.StatusOK, stats)
}

func (h *Handler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.svc.DeleteAvatar(r.Context(), userID)
	if err != nil {
		log.Printf("delete avatar: %v", err)
		response.Error(w, http.StatusInternalServerError, "failed to delete avatar")
		return
	}

	response.JSON(w, http.StatusOK, user)
}

const maxAvatarSize = 5 << 20 // 5 MB

func (h *Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := r.ParseMultipartForm(maxAvatarSize); err != nil {
		response.Error(w, http.StatusBadRequest, "file too large or invalid form")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "avatar file is required")
		return
	}
	defer file.Close()

	user, err := h.svc.UploadAvatar(r.Context(), userID, file, header)
	if err != nil {
		log.Printf("upload avatar: %v", err)
		if strings.Contains(err.Error(), "unsupported image type") {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if strings.Contains(err.Error(), "not configured") {
			response.Error(w, http.StatusServiceUnavailable, "avatar storage is not configured")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to upload avatar")
		return
	}

	response.JSON(w, http.StatusOK, user)
}
