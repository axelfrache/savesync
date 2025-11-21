package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/axelfrache/savesync/internal/app/sourceservice"
	"github.com/axelfrache/savesync/internal/domain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// SourceHandler handles source-related requests
type SourceHandler struct {
	service *sourceservice.Service
	logger  *zap.Logger
}

// NewSourceHandler creates a new source handler
func NewSourceHandler(service *sourceservice.Service, logger *zap.Logger) *SourceHandler {
	return &SourceHandler{
		service: service,
		logger:  logger,
	}
}

// List godoc
// @Summary Lister toutes les sources
// @Description Récupère la liste de toutes les sources de backup configurées
// @Tags sources
// @Produce json
// @Success 200 {array} handlers.SourceResponse
// @Failure 500 {object} handlers.ErrorInfo
// @Router /sources [get]
func (h *SourceHandler) List(w http.ResponseWriter, r *http.Request) {
	sources, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Error("failed to list sources", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to list sources")
		return
	}

	WriteJSON(w, http.StatusOK, sources)
}

// Get godoc
// @Summary Récupérer une source
// @Description Récupère les détails d'une source par son ID
// @Tags sources
// @Produce json
// @Param id path int true "Source ID"
// @Success 200 {object} handlers.SourceResponse
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 404 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /sources/{id} [get]
func (h *SourceHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid source ID")
		return
	}

	source, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Source not found")
			return
		}
		h.logger.Error("failed to get source", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to get source")
		return
	}

	WriteJSON(w, http.StatusOK, source)
}

// Create godoc
// @Summary Créer une source
// @Description Crée une nouvelle source de backup
// @Tags sources
// @Accept json
// @Produce json
// @Param source body handlers.CreateSourceRequest true "Source à créer"
// @Success 201 {object} handlers.SourceResponse
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /sources [post]
func (h *SourceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	source := &domain.Source{
		Name:       req.Name,
		Path:       req.Path,
		Exclusions: req.Exclusions,
		TargetID:   req.TargetID,
		ScheduleID: req.ScheduleID,
	}

	if err := h.service.Create(r.Context(), source); err != nil {
		if err == domain.ErrInvalidPath {
			WriteError(w, http.StatusBadRequest, "Invalid path: directory does not exist")
			return
		}
		if err == domain.ErrInvalidInput {
			WriteError(w, http.StatusBadRequest, "Invalid input")
			return
		}
		h.logger.Error("failed to create source", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to create source")
		return
	}

	WriteJSON(w, http.StatusCreated, source)
}

// Update handles PUT /api/sources/:id
// Update godoc
// @Summary Mettre à jour une source
// @Description Met à jour une source existante
// @Tags sources
// @Accept json
// @Produce json
// @Param id path int true "Source ID"
// @Param source body handlers.UpdateSourceRequest true "Données de la source"
// @Success 200 {object} handlers.SourceResponse
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 404 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /sources/{id} [put]
func (h *SourceHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid source ID")
		return
	}

	var req UpdateSourceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	source := &domain.Source{
		ID:         id,
		Name:       req.Name,
		Path:       req.Path,
		Exclusions: req.Exclusions,
		TargetID:   req.TargetID,
		ScheduleID: req.ScheduleID,
	}

	if err := h.service.Update(r.Context(), source); err != nil {
		if err == domain.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Source not found")
			return
		}
		if err == domain.ErrInvalidPath {
			WriteError(w, http.StatusBadRequest, "Invalid path: directory does not exist")
			return
		}
		h.logger.Error("failed to update source", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to update source")
		return
	}

	WriteJSON(w, http.StatusOK, source)
}

// Delete godoc
// @Summary Supprimer une source
// @Description Supprime une source de backup
// @Tags sources
// @Param id path int true "Source ID"
// @Success 204
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 404 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /sources/{id} [delete]
func (h *SourceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid source ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err == domain.ErrNotFound {
			WriteError(w, http.StatusNotFound, "Source not found")
			return
		}
		h.logger.Error("failed to delete source", zap.Error(err), zap.Int64("id", id))
		WriteError(w, http.StatusInternalServerError, "Failed to delete source")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
