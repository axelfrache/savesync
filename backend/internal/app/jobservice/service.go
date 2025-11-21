package jobservice

import (
	"context"
	"time"

	"github.com/axelfrache/savesync/internal/domain"
	"go.uber.org/zap"
)

// Service handles job tracking operations
type Service struct {
	repo   domain.JobRepository
	logger *zap.Logger
}

// New creates a new job service
func New(repo domain.JobRepository, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreateBackupJob creates a new backup job
func (s *Service) CreateBackupJob(ctx context.Context, sourceID int64) (*domain.Job, error) {
	job := &domain.Job{
		Type:      "backup",
		SourceID:  &sourceID,
		Status:    "pending",
		StartedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, job); err != nil {
		s.logger.Error("failed to create backup job", zap.Error(err), zap.Int64("source_id", sourceID))
		return nil, err
	}

	s.logger.Info("backup job created", zap.Int64("job_id", job.ID), zap.Int64("source_id", sourceID))
	return job, nil
}

// UpdateStatus updates a job's status
func (s *Service) UpdateStatus(ctx context.Context, jobID int64, status string, err error) error {
	job, getErr := s.repo.GetByID(ctx, jobID)
	if getErr != nil {
		return getErr
	}

	job.Status = status
	if err != nil {
		errMsg := err.Error()
		job.Error = &errMsg
	}

	if status == "success" || status == "failed" {
		now := time.Now()
		job.EndedAt = &now
	}

	if updateErr := s.repo.Update(ctx, job); updateErr != nil {
		s.logger.Error("failed to update job", zap.Error(updateErr), zap.Int64("job_id", jobID))
		return updateErr
	}

	s.logger.Info("job status updated", zap.Int64("job_id", jobID), zap.String("status", status))
	return nil
}

// GetByID retrieves a job by ID
func (s *Service) GetByID(ctx context.Context, id int64) (*domain.Job, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAll retrieves all jobs
func (s *Service) GetAll(ctx context.Context) ([]*domain.Job, error) {
	return s.repo.GetAll(ctx)
}
