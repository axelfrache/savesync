package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/axelfrache/savesync/internal/app/targetservice"
	"github.com/axelfrache/savesync/internal/domain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// TargetHandler handles target-related requests
type TargetHandler struct {
	service *targetservice.Service
	logger  *zap.Logger
}

// NewTargetHandler creates a new target handler
func NewTargetHandler(service *targetservice.Service, logger *zap.Logger) *TargetHandler {
	return &TargetHandler{
		service: service,
		logger:  logger,
	}
}

// List handles GET /api/targets
func (h *TargetHandler) List(w http.ResponseWriter, r *http.Request) {
	targets, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("failed to list targets", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to list targets")
		return
	}

	WriteJSON(w, http.StatusOK, targets)
}

// Get handles GET /api/targets/:id
func (h *TargetHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid target ID")
		return
	}

	target, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Target not found")
			return
		}
		h.logger.Error("failed to get target", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to get target")
		return
	}

	WriteJSON(w, http.StatusOK, target)
}

// Create handles POST /api/targets
func (h *TargetHandler) Create(w http.ResponseWriter, r *http.Request) {
	var target domain.Target
	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(r.Context(), &target); err != nil {
		if err == domain.ErrBackendInit {
			WriteError(w, http.StatusBadRequest, "Failed to initialize backend with provided configuration")
			return
		}
		h.logger.Error("failed to create target", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to create target")
		return
	}

	WriteJSON(w, http.StatusCreated, target)
}

// Update handles PUT /api/targets/:id
func (h *TargetHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid target ID")
		return
	}

	var target domain.Target
	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	target.ID = id

	if err := h.service.Update(r.Context(), &target); err != nil {
		if err == domain.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Target not found")
			return
		}
		if err == domain.ErrBackendInit {
			WriteError(w, http.StatusBadRequest, "Failed to initialize backend with provided configuration")
			return
		}
		h.logger.Error("failed to update target", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to update target")
		return
	}

	WriteJSON(w, http.StatusOK, target)
}

// Delete handles DELETE /api/targets/:id
func (h *TargetHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid target ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Target not found")
			return
		}
		h.logger.Error("failed to delete target", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to delete target")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
