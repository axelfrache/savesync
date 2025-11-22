package settingsservice

import (
	"context"
	"strconv"

	"github.com/axelfrache/savesync/internal/domain"
	"go.uber.org/zap"
)

type Service struct {
	settingsRepo domain.SettingsRepository
	logger       *zap.Logger
}

func New(settingsRepo domain.SettingsRepository, logger *zap.Logger) *Service {
	return &Service{
		settingsRepo: settingsRepo,
		logger:       logger,
	}
}

func (s *Service) IsRegistrationEnabled(ctx context.Context) (bool, error) {
	setting, err := s.settingsRepo.Get(ctx, "registration_enabled")
	if err == domain.ErrNotFound {
		return false, nil // Default to disabled
	}
	if err != nil {
		return false, err
	}

	enabled, _ := strconv.ParseBool(setting.Value)
	return enabled, nil
}

func (s *Service) SetRegistrationEnabled(ctx context.Context, enabled bool) error {
	value := "false"
	if enabled {
		value = "true"
	}
	return s.settingsRepo.Set(ctx, "registration_enabled", value)
}

func (s *Service) GetAll(ctx context.Context) (map[string]string, error) {
	return s.settingsRepo.GetAll(ctx)
}

func (s *Service) Set(ctx context.Context, key, value string) error {
	return s.settingsRepo.Set(ctx, key, value)
}
