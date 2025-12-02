package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/axelfrache/savesync/internal/domain"
)

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
	Name   string                 `json:"name" example:"minio-homelab"`
	Type   string                 `json:"type" example:"s3_generic"`
	Config map[string]interface{} `json:"config"`
}

type UpdateTargetRequest struct {
	Name   string                 `json:"name" example:"minio-homelab"`
	Type   string                 `json:"type" example:"s3_generic"`
	Config map[string]interface{} `json:"config"`
}

type TargetResponse struct {
	ID        int64                  `json:"id" example:"1"`
	Name      string                 `json:"name" example:"minio-homelab"`
	Type      string                 `json:"type" example:"s3_generic"`
	Config    map[string]interface{} `json:"config"`
	CreatedAt time.Time              `json:"created_at" example:"2025-01-21T10:00:00Z"`
	UpdatedAt time.Time              `json:"updated_at" example:"2025-01-21T10:00:00Z"`
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

// ToTargetDomain converts CreateTargetRequest to domain.Target
func (r *CreateTargetRequest) ToTargetDomain() (*domain.Target, error) {
	configJSON, err := json.Marshal(r.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	return &domain.Target{
		Name:       r.Name,
		Type:       domain.TargetType(r.Type),
		ConfigJSON: string(configJSON),
	}, nil
}

// ToTargetResponse converts domain.Target to TargetResponse
func ToTargetResponse(target *domain.Target) (*TargetResponse, error) {
	var config map[string]interface{}
	if target.ConfigJSON != "" {
		if err := json.Unmarshal([]byte(target.ConfigJSON), &config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	}

	return &TargetResponse{
		ID:        target.ID,
		Name:      target.Name,
		Type:      string(target.Type),
		Config:    config,
		CreatedAt: target.CreatedAt,
		UpdatedAt: target.UpdatedAt,
	}, nil
}

type BackupResponse struct {
	JobID  int64  `json:"job_id" example:"1"`
	Status string `json:"status" example:"pending"`
}
