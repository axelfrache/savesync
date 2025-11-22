package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/axelfrache/savesync/internal/app/authservice"
	"github.com/axelfrache/savesync/internal/app/settingsservice"
	"github.com/axelfrache/savesync/internal/app/userservice"
	"github.com/axelfrache/savesync/internal/infra/http/middleware"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userService     *userservice.Service
	authService     *authservice.Service
	settingsService *settingsservice.Service
	logger          *zap.Logger
}

func NewAuthHandler(userService *userservice.Service, authService *authservice.Service, settingsService *settingsservice.Service, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userService:     userService,
		authService:     authService,
		settingsService: settingsService,
		logger:          logger,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if registration is enabled
	enabled, err := h.settingsService.IsRegistrationEnabled(r.Context())
	if err != nil {
		h.logger.Error("failed to check registration status", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to check registration status")
		return
	}

	if !enabled {
		WriteError(w, http.StatusForbidden, "Registration is currently disabled. Contact an administrator.")
		return
	}

	// Create user
	user, err := h.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Error("failed to register user", zap.Error(err))
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("failed to generate token", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	WriteJSON(w, http.StatusCreated, response)
}

// Login godoc
// @Summary Login
// @Description Authenticate user and receive JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} handlers.ErrorInfo
// @Failure 401 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Authenticate user
	user, err := h.userService.Authenticate(r.Context(), req.Email, req.Password)
	if err != nil {
		h.logger.Warn("failed login attempt", zap.String("email", req.Email), zap.Error(err))
		WriteError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	token, err := h.authService.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("failed to generate token", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to generate authentication token")
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	WriteJSON(w, http.StatusOK, response)
}

// Me godoc
// @Summary Get current user
// @Description Get authenticated user's information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.User
// @Failure 401 {object} handlers.ErrorInfo
// @Failure 404 {object} handlers.ErrorInfo
// @Failure 500 {object} handlers.ErrorInfo
// @Router /auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get user
	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err), zap.Int64("userID", userID))
		WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	WriteJSON(w, http.StatusOK, user)
}

// GetPublicSettings godoc
// @Summary Get public application settings
// @Description Get public settings like registration status
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} handlers.ErrorInfo
// @Router /auth/settings [get]
func (h *AuthHandler) GetPublicSettings(w http.ResponseWriter, r *http.Request) {
	enabled, err := h.settingsService.IsRegistrationEnabled(r.Context())
	if err != nil {
		h.logger.Error("failed to check registration status", zap.Error(err))
		WriteError(w, http.StatusInternalServerError, "Failed to get settings")
		return
	}

	response := map[string]interface{}{
		"registration_enabled": enabled,
	}

	WriteJSON(w, http.StatusOK, response)
}
