package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/axelfrache/savesync/internal/domain"
)

// TargetRepo implements domain.TargetRepository
type TargetRepo struct {
	db *sql.DB
}

// NewTargetRepo creates a new target repository
func NewTargetRepo(db *sql.DB) *TargetRepo {
	return &TargetRepo{db: db}
}

// Create creates a new target
func (r *TargetRepo) Create(ctx context.Context, target *domain.Target) error {
	configJSON, err := json.Marshal(target.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	query := `
		INSERT INTO targets (name, type, config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		target.Name,
		target.Type,
		string(configJSON),
		now,
		now,
	)
	if err != nil {
		return fmt.Errorf("failed to create target: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	target.ID = id
	target.CreatedAt = now
	target.UpdatedAt = now

	return nil
}

// GetByID retrieves a target by ID
func (r *TargetRepo) GetByID(ctx context.Context, id int64) (*domain.Target, error) {
	query := `
		SELECT id, name, type, config, created_at, updated_at
		FROM targets
		WHERE id = ?
	`

	var target domain.Target
	var configJSON string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&target.ID,
		&target.Name,
		&target.Type,
		&configJSON,
		&target.CreatedAt,
		&target.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get target: %w", err)
	}

	if err := json.Unmarshal([]byte(configJSON), &target.Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &target, nil
}

// GetAll retrieves all targets
func (r *TargetRepo) GetAll(ctx context.Context) ([]*domain.Target, error) {
	query := `
		SELECT id, name, type, config, created_at, updated_at
		FROM targets
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query targets: %w", err)
	}
	defer rows.Close()

	var targets []*domain.Target
	for rows.Next() {
		var target domain.Target
		var configJSON string

		err := rows.Scan(
			&target.ID,
			&target.Name,
			&target.Type,
			&configJSON,
			&target.CreatedAt,
			&target.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan target: %w", err)
		}

		if err := json.Unmarshal([]byte(configJSON), &target.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		targets = append(targets, &target)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return targets, nil
}

// Update updates a target
func (r *TargetRepo) Update(ctx context.Context, target *domain.Target) error {
	configJSON, err := json.Marshal(target.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	query := `
		UPDATE targets
		SET name = ?, type = ?, config = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query,
		target.Name,
		target.Type,
		string(configJSON),
		now,
		target.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update target: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}

	target.UpdatedAt = now
	return nil
}

// Delete deletes a target
func (r *TargetRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM targets WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
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
