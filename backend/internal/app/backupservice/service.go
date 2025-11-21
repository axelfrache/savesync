package backupservice

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/axelfrache/savesync/internal/domain"
	"github.com/axelfrache/savesync/internal/infra/observability"
	"go.uber.org/zap"
)

// Service handles backup operations
type Service struct {
	sourceRepo   domain.SourceRepository
	targetRepo   domain.TargetRepository
	snapshotRepo domain.SnapshotRepository
	jobRepo      domain.JobRepository
	logger       *zap.Logger
	chunker      *Chunker
}

// New creates a new backup service
func New(
	sourceRepo domain.SourceRepository,
	targetRepo domain.TargetRepository,
	snapshotRepo domain.SnapshotRepository,
	jobRepo domain.JobRepository,
	logger *zap.Logger,
) *Service {
	return &Service{
		sourceRepo:   sourceRepo,
		targetRepo:   targetRepo,
		snapshotRepo: snapshotRepo,
		jobRepo:      jobRepo,
		logger:       logger,
		chunker:      NewChunker(DefaultChunkSize),
	}
}

// RunBackup executes a backup for a source
func (s *Service) RunBackup(ctx context.Context, sourceID int64, backend domain.Backend) error {
	startTime := time.Now()

	// Get source
	source, err := s.sourceRepo.GetByID(ctx, sourceID)
	if err != nil {
		return fmt.Errorf("failed to get source: %w", err)
	}

	// Validate target
	if source.TargetID == nil {
		return fmt.Errorf("source has no target configured")
	}

	// Create snapshot
	snapshot := &domain.Snapshot{
		SourceID:  sourceID,
		TargetID:  *source.TargetID,
		Status:    "running",
		CreatedAt: time.Now(),
	}

	if err := s.snapshotRepo.Create(ctx, snapshot); err != nil {
		return fmt.Errorf("failed to create snapshot: %w", err)
	}

	s.logger.Info("starting backup",
		zap.Int64("snapshot_id", snapshot.ID),
		zap.Int64("source_id", sourceID),
		zap.String("path", source.Path),
	)

	// Scan and backup files
	var manifestFiles []domain.ManifestFile
	var totalBytes int64
	var deltaBytes int64
	fileCount := 0

	err = filepath.Walk(source.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check exclusions
		relPath, _ := filepath.Rel(source.Path, path)
		if s.shouldExclude(relPath, source.Exclusions) {
			s.logger.Debug("excluding file", zap.String("path", relPath))
			return nil
		}

		// Hash file
		fileHash, err := HashFile(path)
		if err != nil {
			s.logger.Warn("failed to hash file", zap.Error(err), zap.String("path", path))
			return nil // Skip file but continue
		}

		// Chunk file
		chunks, err := s.chunker.ChunkFile(path)
		if err != nil {
			s.logger.Warn("failed to chunk file", zap.Error(err), zap.String("path", path))
			return nil // Skip file but continue
		}

		// Upload chunks with deduplication
		var chunkHashes []string
		for _, chunk := range chunks {
			// Check if chunk already exists
			exists, err := backend.ChunkExists(ctx, chunk.Hash)
			if err != nil {
				return fmt.Errorf("failed to check chunk existence: %w", err)
			}

			if !exists {
				// Upload new chunk
				if err := backend.StoreChunk(ctx, chunk.Hash, chunk.Data); err != nil {
					return fmt.Errorf("failed to store chunk: %w", err)
				}
				deltaBytes += chunk.Size
			}

			chunkHashes = append(chunkHashes, chunk.Hash)
			totalBytes += chunk.Size
		}

		// Add to manifest
		manifestFiles = append(manifestFiles, domain.ManifestFile{
			Path:    relPath,
			Size:    info.Size(),
			Hash:    fileHash,
			Chunks:  chunkHashes,
			ModTime: info.ModTime(),
		})

		fileCount++

		if fileCount%100 == 0 {
			s.logger.Info("backup progress",
				zap.Int("files", fileCount),
				zap.Int64("bytes", totalBytes),
			)
		}

		return nil
	})

	if err != nil {
		// Update snapshot as failed
		snapshot.Status = "failed"
		errMsg := err.Error()
		snapshot.Error = &errMsg
		now := time.Now()
		snapshot.CompletedAt = &now
		s.snapshotRepo.Update(ctx, snapshot)

		observability.ErrorCountTotal.WithLabelValues("backup").Inc()
		return fmt.Errorf("backup failed: %w", err)
	}

	// Create and store manifest
	manifest := domain.Manifest{
		SnapshotID: snapshot.ID,
		SourcePath: source.Path,
		Files:      manifestFiles,
		CreatedAt:  time.Now(),
	}

	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := backend.StoreManifest(ctx, strconv.FormatInt(snapshot.ID, 10), manifestJSON); err != nil {
		return fmt.Errorf("failed to store manifest: %w", err)
	}

	// Update snapshot as success
	snapshot.Status = "success"
	snapshot.FileCount = fileCount
	snapshot.TotalBytes = totalBytes
	snapshot.DeltaBytes = deltaBytes
	now := time.Now()
	snapshot.CompletedAt = &now

	if err := s.snapshotRepo.Update(ctx, snapshot); err != nil {
		return fmt.Errorf("failed to update snapshot: %w", err)
	}

	// Update metrics
	duration := time.Since(startTime).Seconds()
	observability.BackupDuration.WithLabelValues(
		strconv.FormatInt(sourceID, 10),
		source.Name,
	).Observe(duration)

	observability.BackupLastRunTimestamp.WithLabelValues(
		strconv.FormatInt(sourceID, 10),
		source.Name,
	).SetToCurrentTime()

	observability.BackupStatus.WithLabelValues(
		strconv.FormatInt(sourceID, 10),
		source.Name,
	).Set(1)

	observability.BytesTransferredTotal.WithLabelValues(
		strconv.FormatInt(sourceID, 10),
		source.Name,
	).Add(float64(deltaBytes))

	s.logger.Info("backup completed",
		zap.Int64("snapshot_id", snapshot.ID),
		zap.Int("files", fileCount),
		zap.Int64("total_bytes", totalBytes),
		zap.Int64("delta_bytes", deltaBytes),
		zap.Float64("duration_seconds", duration),
	)

	return nil
}

// shouldExclude checks if a file should be excluded based on patterns
func (s *Service) shouldExclude(path string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			s.logger.Warn("invalid exclusion pattern", zap.String("pattern", pattern), zap.Error(err))
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
