package middleware

import (
	"context"
	"net/http"

	"github.com/axelfrache/savesync/internal/app/userservice"
	"go.uber.org/zap"
)

type adminContextKey string

const AdminKey adminContextKey = "isAdmin"

// AdminMiddleware validates that the user is an administrator
func AdminMiddleware(userService *userservice.Service, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			user, err := userService.GetByID(r.Context(), userID)
			if err != nil {
				logger.Error("failed to get user", zap.Error(err), zap.Int64("userID", userID))
				writeError(w, http.StatusInternalServerError, "Failed to verify admin status")
				return
			}

			if !user.IsAdmin {
				logger.Warn("non-admin user attempted admin action", zap.Int64("userID", userID))
				writeError(w, http.StatusForbidden, "Admin access required")
				return
			}

			// Add admin flag to context for handlers
			ctx := context.WithValue(r.Context(), AdminKey, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// IsAdmin checks if the current user is an admin
func IsAdmin(ctx context.Context) bool {
	isAdmin, ok := ctx.Value(AdminKey).(bool)
	return ok && isAdmin
}
