package db

import (
	"context"
	"fmt"
)

// migrate runs database migrations
func (db *DB) migrate(ctx context.Context) error {
	migrations := []string{
		// Sources table
		`CREATE TABLE IF NOT EXISTS sources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			path TEXT NOT NULL,
			exclusions TEXT, -- JSON array
			target_id INTEGER,
			schedule_id INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE SET NULL,
			FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE SET NULL
		)`,

		// Targets table
		`CREATE TABLE IF NOT EXISTS targets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL, -- local, s3, sftp
			config TEXT NOT NULL, -- JSON object
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Snapshots table
		`CREATE TABLE IF NOT EXISTS snapshots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source_id INTEGER NOT NULL,
			target_id INTEGER NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending', -- pending, running, success, failed
			file_count INTEGER DEFAULT 0,
			total_bytes INTEGER DEFAULT 0,
			delta_bytes INTEGER DEFAULT 0,
			error TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			completed_at TIMESTAMP,
			FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE,
			FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE CASCADE
		)`,

		// Snapshot files table
		`CREATE TABLE IF NOT EXISTS snapshot_files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			snapshot_id INTEGER NOT NULL,
			path TEXT NOT NULL,
			size INTEGER NOT NULL,
			hash TEXT NOT NULL,
			chunks TEXT NOT NULL, -- JSON array of chunk hashes
			mod_time TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (snapshot_id) REFERENCES snapshots(id) ON DELETE CASCADE
		)`,

		// Jobs table
		`CREATE TABLE IF NOT EXISTS jobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL, -- backup, restore
			source_id INTEGER,
			snapshot_id INTEGER,
			status TEXT NOT NULL DEFAULT 'pending', -- pending, running, success, failed
			error TEXT,
			started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			ended_at TIMESTAMP,
			FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE SET NULL,
			FOREIGN KEY (snapshot_id) REFERENCES snapshots(id) ON DELETE SET NULL
		)`,

		// Schedules table
		`CREATE TABLE IF NOT EXISTS schedules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source_id INTEGER NOT NULL,
			frequency TEXT NOT NULL, -- manual, hourly, daily, weekly, cron
			cron_expr TEXT,
			enabled BOOLEAN DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (source_id) REFERENCES sources(id) ON DELETE CASCADE
		)`,

		// Indexes for performance
		`CREATE INDEX IF NOT EXISTS idx_snapshots_source_id ON snapshots(source_id)`,
		`CREATE INDEX IF NOT EXISTS idx_snapshots_status ON snapshots(status)`,
		`CREATE INDEX IF NOT EXISTS idx_snapshot_files_snapshot_id ON snapshot_files(snapshot_id)`,
		`CREATE INDEX IF NOT EXISTS idx_snapshot_files_hash ON snapshot_files(hash)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_source_id ON jobs(source_id)`,
		`CREATE INDEX IF NOT EXISTS idx_schedules_source_id ON schedules(source_id)`,
		`CREATE INDEX IF NOT EXISTS idx_schedules_enabled ON schedules(enabled)`,
	}

	for i, migration := range migrations {
		if _, err := db.ExecContext(ctx, migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i, err)
		}
	}

	return nil
}
