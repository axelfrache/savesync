package targetservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/axelfrache/savesync/internal/domain"
	"github.com/axelfrache/savesync/internal/infra/backends"
	"go.uber.org/zap"
)

type Service struct {
	repo     domain.TargetRepository
	registry *backends.Registry
	logger   *zap.Logger
}

func New(repo domain.TargetRepository, registry *backends.Registry, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		registry: registry,
		logger:   logger,
	}
}

func (s *Service) Create(ctx context.Context, target *domain.Target) error {
	if target.Name == "" {
		return domain.ErrInvalidInput
	}

	supportedTypes := s.registry.SupportedTypes()
	typeSupported := false
	for _, t := range supportedTypes {
		if t == string(target.Type) {
			typeSupported = true
			break
		}
	}
	if !typeSupported {
		return fmt.Errorf("unsupported backend type: %s", target.Type)
	}

	var config map[string]interface{}
	if target.ConfigJSON != "" {
		if err := json.Unmarshal([]byte(target.ConfigJSON), &config); err != nil {
			s.logger.Error("failed to unmarshal config", zap.Error(err))
			return domain.ErrInvalidInput
		}
	}

	stringConfig := make(map[string]string)
	for k, v := range config {
		stringConfig[k] = fmt.Sprintf("%v", v)
	}

	backend, err := s.registry.Create(string(target.Type), stringConfig)
	if err != nil {
		s.logger.Error("failed to initialize backend", zap.Error(err), zap.String("type", string(target.Type)))
		return domain.ErrBackendInit
	}
	defer backend.Close()

	if err := s.repo.Create(ctx, target); err != nil {
		s.logger.Error("failed to create target", zap.Error(err), zap.String("name", target.Name))
		return err
	}

	s.logger.Info("target created", zap.Int64("id", target.ID), zap.String("name", target.Name), zap.String("type", string(target.Type)))
	return nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Target, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*domain.Target, error) {
	targets, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return targets, nil
}

func (s *Service) Update(ctx context.Context, target *domain.Target) error {
	if target.Name == "" {
		return domain.ErrInvalidInput
	}

	var config map[string]interface{}
	if target.ConfigJSON != "" {
		if err := json.Unmarshal([]byte(target.ConfigJSON), &config); err != nil {
			s.logger.Error("failed to unmarshal config", zap.Error(err))
			return domain.ErrInvalidInput
		}
	}

	stringConfig := make(map[string]string)
	for k, v := range config {
		stringConfig[k] = fmt.Sprintf("%v", v)
	}

	backend, err := s.registry.Create(string(target.Type), stringConfig)
	if err != nil {
		s.logger.Error("failed to initialize backend", zap.Error(err), zap.String("type", string(target.Type)))
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

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete target", zap.Error(err), zap.Int64("id", id))
		return err
	}

	s.logger.Info("target deleted", zap.Int64("id", id))
	return nil
}

func (s *Service) GetBackend(ctx context.Context, id int64) (domain.Backend, error) {
	target, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	if target.ConfigJSON != "" {
		if err := json.Unmarshal([]byte(target.ConfigJSON), &config); err != nil {
			s.logger.Error("failed to unmarshal config", zap.Error(err))
			return nil, domain.ErrInvalidInput
		}
	}

	stringConfig := make(map[string]string)
	for k, v := range config {
		stringConfig[k] = fmt.Sprintf("%v", v)
	}

	backend, err := s.registry.Create(string(target.Type), stringConfig)
	if err != nil {
		s.logger.Error("failed to create backend", zap.Error(err), zap.Int64("id", id))
		return nil, err
	}

	return backend, nil
}
