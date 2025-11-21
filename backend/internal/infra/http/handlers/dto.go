package handlers

import "time"

type CreateSourceRequest struct {
	Name       string   `json:"name" example:"mes-documents"`
	Path       string   `json:"path" example:"/home/user/documents"`
	Exclusions []string `json:"exclusions" example:"*.tmp,*.log"`
	TargetID   *int64   `json:"target_id" example:"1"`
	ScheduleID *int64   `json:"schedule_id,omitempty"`
}

type UpdateSourceRequest struct {
	Name       string   `json:"name" example:"mes-documents"`
	Path       string   `json:"path" example:"/home/user/documents"`
	Exclusions []string `json:"exclusions" example:"*.tmp,*.log"`
	TargetID   *int64   `json:"target_id" example:"1"`
	ScheduleID *int64   `json:"schedule_id,omitempty"`
}

type SourceResponse struct {
	ID         int64     `json:"id" example:"1"`
	Name       string    `json:"name" example:"mes-documents"`
	Path       string    `json:"path" example:"/home/user/documents"`
	Exclusions []string  `json:"exclusions" example:"*.tmp,*.log"`
	TargetID   *int64    `json:"target_id" example:"1"`
	ScheduleID *int64    `json:"schedule_id,omitempty"`
	CreatedAt  time.Time `json:"created_at" example:"2025-01-21T10:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2025-01-21T10:00:00Z"`
}

type CreateTargetRequest struct {
	Name   string            `json:"name" example:"backup-local"`
	Type   string            `json:"type" example:"local"`
	Config map[string]string `json:"config" example:"path:/tmp/backups"`
}

type UpdateTargetRequest struct {
	Name   string            `json:"name" example:"backup-local"`
	Type   string            `json:"type" example:"local"`
	Config map[string]string `json:"config" example:"path:/tmp/backups"`
}

type TargetResponse struct {
	ID        int64             `json:"id" example:"1"`
	Name      string            `json:"name" example:"backup-local"`
	Type      string            `json:"type" example:"local"`
	Config    map[string]string `json:"config"`
	CreatedAt time.Time         `json:"created_at" example:"2025-01-21T10:00:00Z"`
	UpdatedAt time.Time         `json:"updated_at" example:"2025-01-21T10:00:00Z"`
}

type JobResponse struct {
	ID         int64      `json:"id" example:"1"`
	Type       string     `json:"type" example:"backup"`
	SourceID   *int64     `json:"source_id,omitempty" example:"1"`
	SnapshotID *int64     `json:"snapshot_id,omitempty"`
	Status     string     `json:"status" example:"success"`
	Error      *string    `json:"error,omitempty"`
	StartedAt  time.Time  `json:"started_at" example:"2025-01-21T10:00:00Z"`
	EndedAt    *time.Time `json:"ended_at,omitempty" example:"2025-01-21T10:05:00Z"`
}

type BackupResponse struct {
	JobID  int64  `json:"job_id" example:"1"`
	Status string `json:"status" example:"pending"`
}
