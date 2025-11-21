package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/axelfrache/savesync/internal/app/backupservice"
	"github.com/axelfrache/savesync/internal/app/sourceservice"
	"github.com/axelfrache/savesync/internal/app/targetservice"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type SnapshotHandler struct {
	service       *backupservice.Service
	sourceService *sourceservice.Service
	targetService *targetservice.Service
	logger        *zap.Logger
}

func NewSnapshotHandler(
	service *backupservice.Service,
	sourceService *sourceservice.Service,
	targetService *targetservice.Service,
	logger *zap.Logger,
) *SnapshotHandler {
	return &SnapshotHandler{
		service:       service,
		sourceService: sourceService,
		targetService: targetService,
		logger:        logger,
	}
}

// List godoc
// @Summary Lister tous les snapshots
// @Description Récupère la liste de tous les snapshots
// @Tags snapshots
// @Produce json
// @Success 200 {array} domain.Snapshot
// @Failure 500 {object} handlers.ErrorInfo
// @Router /snapshots [get]
func (h *SnapshotHandler) List(w http.ResponseWriter, r *http.Request) {
	snapshots, err := h.service.ListSnapshots(r.Context())
	if err != nil {
		h.logger.Error("failed to list snapshots", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to list snapshots")
		return
	}

	WriteJSON(w, http.StatusOK, snapshots)
}

// Get godoc
// @Summary Récupérer un snapshot
// @Description Récupère les détails d'un snapshot par son ID
// @Tags snapshots
// @Produce json
// @Param id path int true "Snapshot ID"
// @Success 200 {object} domain.Snapshot
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 404 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /snapshots/{id} [get]
func (h *SnapshotHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid snapshot ID")
		return
	}

	snapshot, err := h.service.GetSnapshot(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to get snapshot", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to get snapshot")
		return
	}

	WriteJSON(w, http.StatusOK, snapshot)
}

// Restore godoc
// @Summary Restaurer un snapshot
// @Description Restaure un snapshot spécifique
// @Tags snapshots
// @Param id path int true "Snapshot ID"
// @Success 202
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /snapshots/{id}/restore [post]
func (h *SnapshotHandler) Restore(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid snapshot ID")
		return
	}

	if err := h.service.RestoreSnapshot(r.Context(), id); err != nil {
		h.logger.Error("failed to restore snapshot", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to restore snapshot")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// GetManifest godoc
// @Summary Télécharger le manifeste
// @Description Télécharge le fichier manifeste d'un snapshot
// @Tags snapshots
// @Produce application/json
// @Param id path int true "Snapshot ID"
// @Success 200 {file} []byte
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /snapshots/{id}/manifest [get]
func (h *SnapshotHandler) GetManifest(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid snapshot ID")
		return
	}

	ctx := r.Context()

	// 1. Get Snapshot
	snapshot, err := h.service.GetSnapshot(ctx, id)
	if err != nil {
		h.logger.Error("failed to get snapshot", zap.Error(err))
		WriteError(w, http.StatusNotFound, "Snapshot not found")
		return
	}

	// 2. Get Source
	source, err := h.sourceService.GetByID(ctx, snapshot.SourceID)
	if err != nil {
		h.logger.Error("failed to get source", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to get source")
		return
	}

	if source.TargetID == nil {
		WriteError(w, http.StatusBadRequest, "Source has no target")
		return
	}

	// 3. Get Backend
	backend, err := h.targetService.GetBackend(ctx, *source.TargetID)
	if err != nil {
		h.logger.Error("failed to get backend", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to initialize backend")
		return
	}
	defer backend.Close()

	// 4. Get Manifest
	manifestJSON, err := h.service.GetManifestWithBackend(ctx, id, backend)
	if err != nil {
		h.logger.Error("failed to get manifest", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve manifest")
		return
	}

	// 5. Serve File
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"manifest-%d.json\"", id))
	w.Write(manifestJSON)
}

// GetFiles godoc
// @Summary Get file tree
// @Description Récupère l'arborescence des fichiers d'un snapshot
// @Tags snapshots
// @Produce json
// @Param id path int true "Snapshot ID"
// @Success 200 {object} backupservice.FileNode
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /snapshots/{id}/files [get]
func (h *SnapshotHandler) GetFiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid snapshot ID")
		return
	}

	ctx := r.Context()

	// 1. Get Snapshot
	snapshot, err := h.service.GetSnapshot(ctx, id)
	if err != nil {
		h.logger.Error("failed to get snapshot", zap.Error(err))
		WriteError(w, http.StatusNotFound, "Snapshot not found")
		return
	}

	// 2. Get Source
	source, err := h.sourceService.GetByID(ctx, snapshot.SourceID)
	if err != nil {
		h.logger.Error("failed to get source", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to get source")
		return
	}

	if source.TargetID == nil {
		WriteError(w, http.StatusBadRequest, "Source has no target")
		return
	}

	// 3. Get Backend
	backend, err := h.targetService.GetBackend(ctx, *source.TargetID)
	if err != nil {
		h.logger.Error("failed to get backend", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to initialize backend")
		return
	}
	defer backend.Close()

	// 4. Get and Parse Manifest
	tree, err := h.service.GetSnapshotFileTree(ctx, id, backend)
	if err != nil {
		h.logger.Error("failed to get file tree", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to retrieve file tree")
		return
	}

	WriteJSON(w, http.StatusOK, tree)
}
