package handlers

import (
	"net/http"
	"strconv"

	"github.com/axelfrache/savesync/internal/app/backupservice"
	"github.com/axelfrache/savesync/internal/app/jobservice"
	"github.com/axelfrache/savesync/internal/app/targetservice"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// BackupHandler handles backup-related requests
type BackupHandler struct {
	backupService *backupservice.Service
	targetService *targetservice.Service
	jobService    *jobservice.Service
	logger        *zap.Logger
}

// NewBackupHandler creates a new backup handler
func NewBackupHandler(
	backupService *backupservice.Service,
	targetService *targetservice.Service,
	jobService *jobservice.Service,
	logger *zap.Logger,
) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
		targetService: targetService,
		jobService:    jobService,
		logger:        logger,
	}
}

// Run handles POST /api/sources/:id/run
func (h *BackupHandler) Run(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	sourceID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid source ID")
		return
	}

	// Create job
	job, err := h.jobService.CreateBackupJob(r.Context(), sourceID)
	if err != nil {
		h.logger.Error("failed to create backup job", zap.Error(err), zap.Int64("source_id", sourceID))
		WriteError(w, http.StatusInternalServerError, "Failed to create backup job")
		return
	}

	// Run backup asynchronously
	go func() {
		ctx := r.Context()

		// Update job status to running
		h.jobService.UpdateStatus(ctx, job.ID, "running", nil)

		// Get source to find target
		// Note: In a real implementation, we'd pass the source through or fetch it
		// For now, we'll use a simplified approach

		// This is a placeholder - in the actual main.go, we'll wire this properly
		h.logger.Info("backup job started", zap.Int64("job_id", job.ID), zap.Int64("source_id", sourceID))

		// The actual backup execution will be wired in main.go
		// For now, just update the job as pending
		h.jobService.UpdateStatus(ctx, job.ID, "pending", nil)
	}()

	WriteJSON(w, http.StatusAccepted, map[string]interface{}{
		"job_id": job.ID,
		"status": job.Status,
	})
}
