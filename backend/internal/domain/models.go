package domain

import "time"

// Source represents a directory to be backed up
type Source struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Exclusions []string  `json:"exclusions"` // Glob patterns to exclude
	TargetID   *int64    `json:"target_id"`
	ScheduleID *int64    `json:"schedule_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Target represents a storage backend configuration
type Target struct {
	ID        int64             `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"` // local, s3, sftp
	Config    map[string]string `json:"config"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Snapshot represents a backup snapshot at a point in time
type Snapshot struct {
	ID          int64      `json:"id"`
	SourceID    int64      `json:"source_id"`
	TargetID    int64      `json:"target_id"`
	Status      string     `json:"status"` // pending, running, success, failed
	FileCount   int        `json:"file_count"`
	TotalBytes  int64      `json:"total_bytes"`
	DeltaBytes  int64      `json:"delta_bytes"` // New bytes uploaded
	Error       *string    `json:"error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// SnapshotFile represents a file within a snapshot
type SnapshotFile struct {
	ID         int64     `json:"id"`
	SnapshotID int64     `json:"snapshot_id"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	Hash       string    `json:"hash"`   // SHA256
	Chunks     []string  `json:"chunks"` // List of chunk hashes
	ModTime    time.Time `json:"mod_time"`
	CreatedAt  time.Time `json:"created_at"`
}

// Job represents a backup job execution
type Job struct {
	ID         int64      `json:"id"`
	Type       string     `json:"type"` // backup, restore
	SourceID   *int64     `json:"source_id,omitempty"`
	SnapshotID *int64     `json:"snapshot_id,omitempty"`
	Status     string     `json:"status"` // pending, running, success, failed
	Error      *string    `json:"error,omitempty"`
	StartedAt  time.Time  `json:"started_at"`
	EndedAt    *time.Time `json:"ended_at,omitempty"`
}

// Schedule represents a backup schedule
type Schedule struct {
	ID        int64     `json:"id"`
	SourceID  int64     `json:"source_id"`
	Frequency string    `json:"frequency"` // manual, hourly, daily, weekly, cron
	CronExpr  *string   `json:"cron_expr,omitempty"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Chunk represents a content-addressable chunk
type Chunk struct {
	Hash string `json:"hash"` // SHA256
	Size int64  `json:"size"`
	Data []byte `json:"-"` // Not serialized
}

// Manifest represents a snapshot manifest stored in the backend
type Manifest struct {
	SnapshotID int64          `json:"snapshot_id"`
	SourcePath string         `json:"source_path"`
	Files      []ManifestFile `json:"files"`
	CreatedAt  time.Time      `json:"created_at"`
}

// ManifestFile represents a file entry in a manifest
type ManifestFile struct {
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	Hash    string    `json:"hash"`
	Chunks  []string  `json:"chunks"`
	ModTime time.Time `json:"mod_time"`
}
