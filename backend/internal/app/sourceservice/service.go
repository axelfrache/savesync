package sourceservice

import (
	"context"
	"fmt"
	"os"

	"github.com/axelfrache/savesync/internal/domain"
	"go.uber.org/zap"
)

// Service handles source management operations
type Service struct {
	repo   domain.SourceRepository
	logger *zap.Logger
}

// New creates a new source service
func New(repo domain.SourceRepository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new source
func (s *Service) Create(ctx context.Context, source *domain.Source) error {
	// Validate path exists
	if _, err := os.Stat(source.Path); err != nil {
		if os.IsNotExist(err) {
			return domain.ErrInvalidPath
		}
		return fmt.Errorf("failed to check path: %w", err)
	}

	// Validate name is not empty
	if source.Name == "" {
		return domain.ErrInvalidInput
	}

	// Initialize exclusions if nil
	if source.Exclusions == nil {
		source.Exclusions = []string{}
	}

	if err := s.repo.Create(ctx, source); err != nil {
		s.logger.Error("failed to create source", zap.Error(err), zap.String("name", source.Name))
		return err
	}

	s.logger.Info("source created", zap.Int64("id", source.ID), zap.String("name", source.Name))
	return nil
}

// GetByID retrieves a source by ID
func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Source, error) {
	source, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return source, nil
}

// GetAll retrieves all sources
func (s *Service) GetAll(ctx context.Context) ([]*domain.Source, error) {
	sources, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return sources, nil
}

// Update updates a source
func (s *Service) Update(ctx context.Context, source *domain.Source) error {
	// Validate path exists
	if _, err := os.Stat(source.Path); err != nil {
		if os.IsNotExist(err) {
			return domain.ErrInvalidPath
		}
		return fmt.Errorf("failed to check path: %w", err)
	}

	// Validate name is not empty
	if source.Name == "" {
		return domain.ErrInvalidInput
	}

	if err := s.repo.Update(ctx, source); err != nil {
		s.logger.Error("failed to update source", zap.Error(err), zap.Int64("id", source.ID))
		return err
	}

	s.logger.Info("source updated", zap.Int64("id", source.ID), zap.String("name", source.Name))
	return nil
}

// Delete deletes a source
func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete source", zap.Error(err), zap.Int64("id", id))
		return err
	}

	s.logger.Info("source deleted", zap.Int64("id", id))
	return nil
}
