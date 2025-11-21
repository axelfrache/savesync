package http

import (
	"net/http"

	"github.com/axelfrache/savesync/internal/app/backupservice"
	"github.com/axelfrache/savesync/internal/app/jobservice"
	"github.com/axelfrache/savesync/internal/app/sourceservice"
	"github.com/axelfrache/savesync/internal/app/targetservice"
	"github.com/axelfrache/savesync/internal/infra/http/handlers"
	"github.com/axelfrache/savesync/internal/infra/http/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "github.com/axelfrache/savesync/docs"
)

// Router creates and configures the HTTP router
func NewRouter(
	sourceService *sourceservice.Service,
	targetService *targetservice.Service,
	backupService *backupservice.Service,
	jobService *jobservice.Service,
	logger *zap.Logger,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))
	r.Use(chimiddleware.Compress(5))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check
	healthHandler := handlers.NewHealthHandler()
	r.Get("/health", healthHandler.ServeHTTP)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Metrics
	r.Handle("/metrics", promhttp.Handler())

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Sources
		sourceHandler := handlers.NewSourceHandler(sourceService, logger)
		r.Route("/sources", func(r chi.Router) {
			r.Get("/", sourceHandler.List)
			r.Post("/", sourceHandler.Create)
			r.Get("/{id}", sourceHandler.Get)
			r.Put("/{id}", sourceHandler.Update)
			r.Delete("/{id}", sourceHandler.Delete)
		})

		// Targets
		targetHandler := handlers.NewTargetHandler(targetService, logger)
		r.Route("/targets", func(r chi.Router) {
			r.Get("/", targetHandler.List)
			r.Post("/", targetHandler.Create)
			r.Get("/{id}", targetHandler.Get)
			r.Put("/{id}", targetHandler.Update)
			r.Delete("/{id}", targetHandler.Delete)
		})

		// Jobs
		jobHandler := handlers.NewJobHandler(jobService, logger)
		r.Route("/jobs", func(r chi.Router) {
			r.Get("/", jobHandler.List)
			r.Get("/{id}", jobHandler.Get)
		})

		// Backup trigger
		backupHandler := handlers.NewBackupHandler(backupService, targetService, jobService, sourceService, logger)
		r.Post("/sources/{id}/run", backupHandler.Run)

		// Snapshots
		snapshotHandler := handlers.NewSnapshotHandler(backupService, sourceService, targetService, logger)
		r.Route("/snapshots", func(r chi.Router) {
			r.Get("/", snapshotHandler.List)
			r.Get("/{id}", snapshotHandler.Get)
			r.Get("/{id}/manifest", snapshotHandler.GetManifest)
			r.Post("/{id}/restore", snapshotHandler.Restore)
		})

		// System (File Explorer)
		systemHandler := handlers.NewSystemHandler(logger)
		r.Get("/system/files", systemHandler.ListFiles)
	})

	return r
}
