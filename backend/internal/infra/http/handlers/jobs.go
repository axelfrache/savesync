package handlers

import (
	"net/http"
	"strconv"

	"github.com/axelfrache/savesync/internal/app/jobservice"
	"github.com/axelfrache/savesync/internal/domain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// JobHandler handles job-related requests
type JobHandler struct {
	service *jobservice.Service
	logger  *zap.Logger
}

// NewJobHandler creates a new job handler
func NewJobHandler(service *jobservice.Service, logger *zap.Logger) *JobHandler {
	return &JobHandler{
		service: service,
		logger:  logger,
	}
}

// List handles GET /api/jobs
func (h *JobHandler) List(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("failed to list jobs", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to list jobs")
		return
	}

	WriteJSON(w, http.StatusOK, jobs)
}

// Get handles GET /api/jobs/:id
func (h *JobHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	job, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Job not found")
			return
		}
		h.logger.Error("failed to get job", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to get job")
		return
	}

	WriteJSON(w, http.StatusOK, job)
}
