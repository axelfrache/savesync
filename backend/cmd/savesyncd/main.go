package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/axelfrache/savesync/internal/app/authservice"
	"github.com/axelfrache/savesync/internal/app/backupservice"
	"github.com/axelfrache/savesync/internal/app/jobservice"
	"github.com/axelfrache/savesync/internal/app/settingsservice"
	"github.com/axelfrache/savesync/internal/app/sourceservice"
	"github.com/axelfrache/savesync/internal/app/targetservice"
	"github.com/axelfrache/savesync/internal/app/userservice"
	"github.com/axelfrache/savesync/internal/config"
	"github.com/axelfrache/savesync/internal/infra/backends"
	"github.com/axelfrache/savesync/internal/infra/db"
	"github.com/axelfrache/savesync/internal/infra/db/repositories"
	httpinfra "github.com/axelfrache/savesync/internal/infra/http"
	"github.com/axelfrache/savesync/internal/infra/observability"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

	// Create default admin user if no users exist
	defaultPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Fatal("failed to hash default password", zap.Error(err))
	}

	var userCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		logger.Fatal("failed to check user count", zap.Error(err))
	}

	if userCount == 0 {
		_, err = database.DB.Exec(
			"INSERT INTO users (email, password_hash, is_admin, created_at, updated_at) VALUES (?, ?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			"admin@savesync.local",
			string(hashedPassword),
		)
		if err != nil {
			logger.Fatal("failed to create default admin user", zap.Error(err))
		}
		logger.Info("created default admin user",
			zap.String("email", "admin@savesync.local"),
			zap.String("password", defaultPassword),
		)
	}

	// Initialize default settings
	_, err = database.DB.Exec(
		"INSERT OR IGNORE INTO settings (key, value) VALUES ('registration_enabled', 'false')",
	)
	if err != nil {
		logger.Warn("failed to initialize default settings", zap.Error(err))
	}

	// Initialize repositories
	userRepo := db.NewUserRepository(database.DB)
	settingsRepo := db.NewSettingsRepository(database.DB)
	sourceRepo := repositories.NewSourceRepo(database.DB)
	targetRepo := repositories.NewTargetRepo(database.DB)
	snapshotRepo := repositories.NewSnapshotRepo(database.DB)
	jobRepo := repositories.NewJobRepo(database.DB)

	// Initialize backend registry
	backendRegistry := backends.NewRegistry()

	// Initialize services
	userService := userservice.New(userRepo, logger)
	settingsService := settingsservice.New(settingsRepo, logger)
	authService := authservice.New("your-secret-key-change-in-production", 24*time.Hour)
	sourceService := sourceservice.New(sourceRepo, logger)
	targetService := targetservice.New(targetRepo, backendRegistry, logger)
	jobService := jobservice.New(jobRepo, logger)
	backupService := backupservice.New(sourceRepo, targetRepo, snapshotRepo, jobRepo, logger)

	logger.Info("services initialized")

	// Create HTTP router with all services
	router := httpinfra.NewRouter(
		userService,
		authService,
		settingsService,
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
