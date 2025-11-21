package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/axelfrache/savesync/internal/app/backupservice"
	"github.com/axelfrache/savesync/internal/app/jobservice"
	"github.com/axelfrache/savesync/internal/app/sourceservice"
	"github.com/axelfrache/savesync/internal/app/targetservice"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// BackupHandler handles backup-related requests
type BackupHandler struct {
	backupService *backupservice.Service
	targetService *targetservice.Service
	jobService    *jobservice.Service
	sourceService *sourceservice.Service
	logger        *zap.Logger
}

// NewBackupHandler creates a new backup handler
func NewBackupHandler(
	backupService *backupservice.Service,
	targetService *targetservice.Service,
	jobService *jobservice.Service,
	sourceService *sourceservice.Service,
	logger *zap.Logger,
) *BackupHandler {
	return &BackupHandler{
		backupService: backupService,
		targetService: targetService,
		jobService:    jobService,
		sourceService: sourceService,
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
		// Create a background context since the request context will be cancelled
		ctx := context.Background()

		// Update job status to running
		h.jobService.UpdateStatus(ctx, job.ID, "running", nil)

		// Get source to find target
		source, err := h.sourceService.GetByID(ctx, sourceID)
		if err != nil {
			h.logger.Error("failed to get source for backup", zap.Error(err), zap.Int64("source_id", sourceID))
			h.jobService.UpdateStatus(ctx, job.ID, "failed", fmt.Errorf("failed to get source: %w", err))
			return
		}

		if source.TargetID == nil {
			h.logger.Error("source has no target", zap.Int64("source_id", sourceID))
			h.jobService.UpdateStatus(ctx, job.ID, "failed", fmt.Errorf("source has no target configured"))
			return
		}

		// Initialize backend
		backend, err := h.targetService.GetBackend(ctx, *source.TargetID)
		if err != nil {
			h.logger.Error("failed to initialize backend", zap.Error(err), zap.Int64("target_id", *source.TargetID))
			h.jobService.UpdateStatus(ctx, job.ID, "failed", fmt.Errorf("failed to initialize backend: %w", err))
			return
		}
		defer backend.Close()

		h.logger.Info("backup job started", zap.Int64("job_id", job.ID), zap.Int64("source_id", sourceID))

		// Run backup
		if err := h.backupService.RunBackup(ctx, sourceID, backend); err != nil {
			h.logger.Error("backup failed", zap.Error(err), zap.Int64("job_id", job.ID))
			h.jobService.UpdateStatus(ctx, job.ID, "failed", fmt.Errorf("backup failed: %w", err))
			return
		}

		// Update job as success
		h.jobService.UpdateStatus(ctx, job.ID, "success", nil)
		h.logger.Info("backup job completed successfully", zap.Int64("job_id", job.ID))
	}()

	WriteJSON(w, http.StatusAccepted, map[string]interface{}{
		"job_id": job.ID,
		"status": job.Status,
	})
}
