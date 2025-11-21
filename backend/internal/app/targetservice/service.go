package targetservice

import (
	"context"
	"fmt"

	"github.com/axelfrache/savesync/internal/domain"
	"github.com/axelfrache/savesync/internal/infra/backends"
	"go.uber.org/zap"
)

// Service handles target management operations
type Service struct {
	repo     domain.TargetRepository
	registry *backends.Registry
	logger   *zap.Logger
}

// New creates a new target service
func New(repo domain.TargetRepository, registry *backends.Registry, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		registry: registry,
		logger:   logger,
	}
}

// Create creates a new target
func (s *Service) Create(ctx context.Context, target *domain.Target) error {
	// Validate name is not empty
	if target.Name == "" {
		return domain.ErrInvalidInput
	}

	// Validate type is supported
	supportedTypes := s.registry.SupportedTypes()
	typeSupported := false
	for _, t := range supportedTypes {
		if t == target.Type {
			typeSupported = true
			break
		}
	}
	if !typeSupported {
		return fmt.Errorf("unsupported backend type: %s", target.Type)
	}

	// Validate backend can be initialized
	backend, err := s.registry.Create(target.Type, target.Config)
	if err != nil {
		s.logger.Error("failed to initialize backend", zap.Error(err), zap.String("type", target.Type))
		return domain.ErrBackendInit
	}
	defer backend.Close()

	// Create target in database
	if err := s.repo.Create(ctx, target); err != nil {
		s.logger.Error("failed to create target", zap.Error(err), zap.String("name", target.Name))
		return err
	}

	s.logger.Info("target created", zap.Int64("id", target.ID), zap.String("name", target.Name), zap.String("type", target.Type))
	return nil
}

// GetByID retrieves a target by ID
func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Target, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return target, nil
}

// GetAll retrieves all targets
func (s *Service) GetAll(ctx context.Context) ([]*domain.Target, error) {
	targets, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return targets, nil
}

// Update updates a target
func (s *Service) Update(ctx context.Context, target *domain.Target) error {
	// Validate name is not empty
	if target.Name == "" {
		return domain.ErrInvalidInput
	}

	// Validate backend can be initialized
	backend, err := s.registry.Create(target.Type, target.Config)
	if err != nil {
		s.logger.Error("failed to initialize backend", zap.Error(err), zap.String("type", target.Type))
		return domain.ErrBackendInit
	}
	defer backend.Close()

	if err := s.repo.Update(ctx, target); err != nil {
		s.logger.Error("failed to update target", zap.Error(err), zap.Int64("id", target.ID))
		return err
	}

	s.logger.Info("target updated", zap.Int64("id", target.ID), zap.String("name", target.Name))
	return nil
}

// Delete deletes a target
func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete target", zap.Error(err), zap.Int64("id", id))
		return err
	}

	s.logger.Info("target deleted", zap.Int64("id", id))
	return nil
}

// GetBackend creates and returns a backend instance for a target
func (s *Service) GetBackend(ctx context.Context, targetID int64) (domain.Backend, error) {
	target, err := s.repo.GetByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	backend, err := s.registry.Create(target.Type, target.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to create backend: %w", err)
	}

	return backend, nil
}
