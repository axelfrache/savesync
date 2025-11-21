package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/axelfrache/savesync/internal/domain"
)

// JobRepo implements domain.JobRepository
type JobRepo struct {
	db *sql.DB
}

// NewJobRepo creates a new job repository
func NewJobRepo(db *sql.DB) *JobRepo {
	return &JobRepo{db: db}
}

// Create creates a new job
func (r *JobRepo) Create(ctx context.Context, job *domain.Job) error {
	query := `
		INSERT INTO jobs (type, source_id, snapshot_id, status, error, started_at, ended_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		job.Type,
		job.SourceID,
		job.SnapshotID,
		job.Status,
		job.Error,
		job.StartedAt,
		job.EndedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	job.ID = id
	return nil
}

// GetByID retrieves a job by ID
func (r *JobRepo) GetByID(ctx context.Context, id int64) (*domain.Job, error) {
	query := `
		SELECT id, type, source_id, snapshot_id, status, error, started_at, ended_at
		FROM jobs
		WHERE id = ?
	`

	var job domain.Job
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID,
		&job.Type,
		&job.SourceID,
		&job.SnapshotID,
		&job.Status,
		&job.Error,
		&job.StartedAt,
		&job.EndedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return &job, nil
}

// GetAll retrieves all jobs
func (r *JobRepo) GetAll(ctx context.Context) ([]*domain.Job, error) {
	query := `
		SELECT id, type, source_id, snapshot_id, status, error, started_at, ended_at
		FROM jobs
		ORDER BY started_at DESC
		LIMIT 100
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*domain.Job
	for rows.Next() {
		var job domain.Job
		err := rows.Scan(
			&job.ID,
			&job.Type,
			&job.SourceID,
			&job.SnapshotID,
			&job.Status,
			&job.Error,
			&job.StartedAt,
			&job.EndedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan job: %w", err)
		}
		jobs = append(jobs, &job)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return jobs, nil
}

// Update updates a job
func (r *JobRepo) Update(ctx context.Context, job *domain.Job) error {
	query := `
		UPDATE jobs
		SET status = ?, error = ?, ended_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		job.Status,
		job.Error,
		job.EndedAt,
		job.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *JobRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM jobs WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
