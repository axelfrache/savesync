package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/axelfrache/savesync/internal/app/userservice"
	"github.com/axelfrache/savesync/internal/domain"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AdminHandler struct {
	userService *userservice.Service
	logger      *zap.Logger
}

func NewAdminHandler(userService *userservice.Service, logger *zap.Logger) *AdminHandler {
	return &AdminHandler{
		userService: userService,
		logger:      logger,
	}
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of all registered users (admin only)
// @Tags admin
// @Produce json
// @Success 200 {array} domain.User
// @Failure 500 {object} handlers.ErrorInfo
// @Router /api/admin/users [get]
// @Security BearerAuth
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAll(r.Context())
	if err != nil {
		h.logger.Error("failed to list users", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	WriteJSON(w, http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User details"
// @Success 201 {object} domain.User
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /api/admin/users [post]
// @Security BearerAuth
func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("failed to create user", zap.Error(err))
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	WriteJSON(w, http.StatusCreated, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user account (admin only)
// @Tags admin
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /api/admin/users/{id} [delete]
// @Security BearerAuth
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Prevent deleting self (optional but recommended, though middleware handles auth check)
	// For now, we rely on frontend or service logic if needed.

	if err := h.userService.Delete(r.Context(), id); err != nil {
		h.logger.Error("failed to delete user", zap.Error(err), zap.Int64("userID", id))
		WriteError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToggleAdmin godoc
// @Summary Toggle admin status
// @Description Promote or demote a user (admin only)
// @Tags admin
// @Param id path int true "User ID"
// @Param request body object{is_admin=boolean} true "Admin status"
// @Success 200 {object} handlers.SuccessInfo
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /api/admin/users/{id}/admin [put]
// @Security BearerAuth
func (h *AdminHandler) ToggleAdmin(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		IsAdmin bool `json:"is_admin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.userService.SetAdminStatus(r.Context(), id, req.IsAdmin); err != nil {
		if err == domain.ErrNotFound {
			WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		h.logger.Error("failed to update user admin status", zap.Error(err), zap.Int64("userID", id))
		WriteError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}
