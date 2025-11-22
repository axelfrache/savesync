package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/axelfrache/savesync/internal/app/authservice"
	"go.uber.org/zap"
)

type contextKey string

const UserIDKey contextKey = "userID"

// writeError writes a JSON error response
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": message,
		},
	})
}

// AuthMiddleware validates JWT tokens and adds user ID to context
func AuthMiddleware(authService *authservice.Service, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("missing authorization header", zap.String("path", r.URL.Path))
				writeError(w, http.StatusUnauthorized, "Missing authorization token")
				return
			}

			// Check for Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("invalid authorization header format", zap.String("header", authHeader))
				writeError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token := parts[1]

			// Validate token
			userID, err := authService.ValidateToken(token)
			if err != nil {
				logger.Warn("invalid token", zap.Error(err), zap.String("path", r.URL.Path))
				writeError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from request context
func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
