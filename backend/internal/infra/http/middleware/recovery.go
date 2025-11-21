package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

// Recovery is a middleware that recovers from panics
func Recovery(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered",
						zap.Any("error", err),
						zap.String("path", r.URL.Path),
						zap.String("method", r.Method),
					)

					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":{"message":"Internal server error"}}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
