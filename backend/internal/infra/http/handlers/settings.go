package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/axelfrache/savesync/internal/app/settingsservice"
	"go.uber.org/zap"
)

type SettingsHandler struct {
	settingsService *settingsservice.Service
	logger          *zap.Logger
}

func NewSettingsHandler(settingsService *settingsservice.Service, logger *zap.Logger) *SettingsHandler {
	return &SettingsHandler{
		settingsService: settingsService,
		logger:          logger,
	}
}

// GetSettings godoc
// @Summary Get all settings
// @Description Get all application settings (admin only)
// @Tags settings
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} handlers.ErrorInfo
// @Router /api/settings [get]
// @Security BearerAuth
func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settingsService.GetAll(r.Context())
	if err != nil {
		h.logger.Error("failed to get settings", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to get settings")
		return
	}

	WriteJSON(w, http.StatusOK, settings)
}

// UpdateSetting godoc
// @Summary Update a setting
// @Description Update an application setting (admin only)
// @Tags settings
// @Accept json
// @Produce json
// @Param setting body object{key=string,value=string} true "Setting to update"
// @Success 200 {object} handlers.SuccessInfo
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /api/settings [put]
// @Security BearerAuth
func (h *SettingsHandler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Key == "" {
		WriteError(w, http.StatusBadRequest, "Key is required")
		return
	}

	if err := h.settingsService.Set(r.Context(), req.Key, req.Value); err != nil {
		h.logger.Error("failed to update setting", zap.Error(err), zap.String("key", req.Key))
		WriteError(w, http.StatusInternalServerError, "Failed to update setting")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Setting updated successfully",
	})
}
