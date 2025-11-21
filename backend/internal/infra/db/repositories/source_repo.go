package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/axelfrache/savesync/internal/domain"
)

// SourceRepo implements domain.SourceRepository
type SourceRepo struct {
	db *sql.DB
}

// NewSourceRepo creates a new source repository
func NewSourceRepo(db *sql.DB) *SourceRepo {
	return &SourceRepo{db: db}
}

// Create creates a new source
func (r *SourceRepo) Create(ctx context.Context, source *domain.Source) error {
	exclusionsJSON, err := json.Marshal(source.Exclusions)
	if err != nil {
		return fmt.Errorf("failed to marshal exclusions: %w", err)
	}

	query := `
		INSERT INTO sources (name, path, exclusions, target_id, schedule_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		source.Name,
		source.Path,
		string(exclusionsJSON),
		source.TargetID,
		source.ScheduleID,
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create source: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	source.ID = id
	source.CreatedAt = now
	source.UpdatedAt = now

	return nil
}

// GetByID retrieves a source by ID
func (r *SourceRepo) GetByID(ctx context.Context, id int64) (*domain.Source, error) {
	query := `
		SELECT id, name, path, exclusions, target_id, schedule_id, created_at, updated_at
		FROM sources
		WHERE id = ?
	`

	var source domain.Source
	var exclusionsJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&source.ID,
		&source.Name,
		&source.Path,
		&exclusionsJSON,
		&source.TargetID,
		&source.ScheduleID,
		&source.CreatedAt,
		&source.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	if err := json.Unmarshal([]byte(exclusionsJSON), &source.Exclusions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exclusions: %w", err)
	}

	return &source, nil
}

// GetAll retrieves all sources
func (r *SourceRepo) GetAll(ctx context.Context) ([]*domain.Source, error) {
	query := `
		SELECT id, name, path, exclusions, target_id, schedule_id, created_at, updated_at
		FROM sources
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sources: %w", err)
	}
	defer rows.Close()

	var sources []*domain.Source
	for rows.Next() {
		var source domain.Source
		var exclusionsJSON string

		err := rows.Scan(
			&source.ID,
			&source.Name,
			&source.Path,
			&exclusionsJSON,
			&source.TargetID,
			&source.ScheduleID,
			&source.CreatedAt,
			&source.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan source: %w", err)
		}

		if err := json.Unmarshal([]byte(exclusionsJSON), &source.Exclusions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal exclusions: %w", err)
		}

		sources = append(sources, &source)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return sources, nil
}

// Update updates a source
func (r *SourceRepo) Update(ctx context.Context, source *domain.Source) error {
	exclusionsJSON, err := json.Marshal(source.Exclusions)
	if err != nil {
		return fmt.Errorf("failed to marshal exclusions: %w", err)
	}

	query := `
		UPDATE sources
		SET name = ?, path = ?, exclusions = ?, target_id = ?, schedule_id = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		source.Name,
		source.Path,
		string(exclusionsJSON),
		source.TargetID,
		source.ScheduleID,
		now,
		source.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update source: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	source.UpdatedAt = now
	return nil
}

// Delete deletes a source
func (r *SourceRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sources WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete source: %w", err)
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
