package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/axelfrache/savesync/internal/domain"
)

// SnapshotRepo implements domain.SnapshotRepository
type SnapshotRepo struct {
	db *sql.DB
}

// NewSnapshotRepo creates a new snapshot repository
func NewSnapshotRepo(db *sql.DB) *SnapshotRepo {
	return &SnapshotRepo{db: db}
}

// Create creates a new snapshot
func (r *SnapshotRepo) Create(ctx context.Context, snapshot *domain.Snapshot) error {
	query := `
		INSERT INTO snapshots (source_id, target_id, status, file_count, total_bytes, delta_bytes, error, created_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		snapshot.SourceID,
		snapshot.TargetID,
		snapshot.Status,
		snapshot.FileCount,
		snapshot.TotalBytes,
		snapshot.DeltaBytes,
		snapshot.Error,
		snapshot.CreatedAt,
		snapshot.CompletedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	snapshot.ID = id
	return nil
}

// GetByID retrieves a snapshot by ID
func (r *SnapshotRepo) GetByID(ctx context.Context, id int64) (*domain.Snapshot, error) {
	query := `
		SELECT id, source_id, target_id, status, file_count, total_bytes, delta_bytes, error, created_at, completed_at
		FROM snapshots
		WHERE id = ?
	`

	var snapshot domain.Snapshot
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&snapshot.ID,
		&snapshot.SourceID,
		&snapshot.TargetID,
		&snapshot.Status,
		&snapshot.FileCount,
		&snapshot.TotalBytes,
		&snapshot.DeltaBytes,
		&snapshot.Error,
		&snapshot.CreatedAt,
		&snapshot.CompletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	return &snapshot, nil
}

// GetBySourceID retrieves all snapshots for a source
func (r *SnapshotRepo) GetBySourceID(ctx context.Context, sourceID int64) ([]*domain.Snapshot, error) {
	query := `
		SELECT id, source_id, target_id, status, file_count, total_bytes, delta_bytes, error, created_at, completed_at
		FROM snapshots
		WHERE source_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query snapshots: %w", err)
	}
	defer rows.Close()

	return r.scanSnapshots(rows)
}

// GetAll retrieves all snapshots
func (r *SnapshotRepo) GetAll(ctx context.Context) ([]*domain.Snapshot, error) {
	query := `
		SELECT id, source_id, target_id, status, file_count, total_bytes, delta_bytes, error, created_at, completed_at
		FROM snapshots
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query snapshots: %w", err)
	}
	defer rows.Close()

	return r.scanSnapshots(rows)
}

// Update updates a snapshot
func (r *SnapshotRepo) Update(ctx context.Context, snapshot *domain.Snapshot) error {
	query := `
		UPDATE snapshots
		SET status = ?, file_count = ?, total_bytes = ?, delta_bytes = ?, error = ?, completed_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		snapshot.Status,
		snapshot.FileCount,
		snapshot.TotalBytes,
		snapshot.DeltaBytes,
		snapshot.Error,
		snapshot.CompletedAt,
		snapshot.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update snapshot: %w", err)
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

// Delete deletes a snapshot
func (r *SnapshotRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM snapshots WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot: %w", err)
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

// scanSnapshots is a helper to scan multiple snapshots
func (r *SnapshotRepo) scanSnapshots(rows *sql.Rows) ([]*domain.Snapshot, error) {
	var snapshots []*domain.Snapshot
	for rows.Next() {
		var snapshot domain.Snapshot
		err := rows.Scan(
			&snapshot.ID,
			&snapshot.SourceID,
			&snapshot.TargetID,
			&snapshot.Status,
			&snapshot.FileCount,
			&snapshot.TotalBytes,
			&snapshot.DeltaBytes,
			&snapshot.Error,
			&snapshot.CreatedAt,
			&snapshot.CompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan snapshot: %w", err)
		}
		snapshots = append(snapshots, &snapshot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return snapshots, nil
}
