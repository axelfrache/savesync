package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/axelfrache/savesync/internal/app/backupservice"
	"github.com/axelfrache/savesync/internal/app/jobservice"
	"github.com/axelfrache/savesync/internal/app/sourceservice"
	"github.com/axelfrache/savesync/internal/app/targetservice"
	"github.com/axelfrache/savesync/internal/config"
	"github.com/axelfrache/savesync/internal/infra/backends"
	"github.com/axelfrache/savesync/internal/infra/db"
	"github.com/axelfrache/savesync/internal/infra/db/repositories"
	httpinfra "github.com/axelfrache/savesync/internal/infra/http"
	"github.com/axelfrache/savesync/internal/infra/observability"
	"go.uber.org/zap"
)

// @title SaveSync API
// @version 1.0
// @description API de gestion de backups automatisés avec déduplication et versioning
// @termsOfService http://swagger.io/terms/

// @contact.name Support SaveSync
// @contact.email support@savesync.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
// @schemes http https

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid config: %v\n", err)
		os.Exit(1)
	}

	logger, err := observability.NewLogger(cfg.Log.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting savesync",
		zap.String("version", "0.1.0"),
		zap.String("db_path", cfg.Database.Path),
		zap.Int("port", cfg.Server.Port),
	)

	database, err := db.New(cfg.Database.Path)
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	defer database.Close()

	logger.Info("database initialized")

	sourceRepo := repositories.NewSourceRepo(database.DB)
	targetRepo := repositories.NewTargetRepo(database.DB)
	snapshotRepo := repositories.NewSnapshotRepo(database.DB)
	jobRepo := repositories.NewJobRepo(database.DB)

	backendRegistry := backends.NewRegistry()

	sourceService := sourceservice.New(sourceRepo, logger)
	targetService := targetservice.New(targetRepo, backendRegistry, logger)
	jobService := jobservice.New(jobRepo, logger)
	backupService := backupservice.New(sourceRepo, targetRepo, snapshotRepo, jobRepo, logger)

	logger.Info("services initialized")

	router := httpinfra.NewRouter(
		sourceService,
		targetService,
		backupService,
		jobService,
		logger,
	)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("starting http server", zap.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server stopped")
}
