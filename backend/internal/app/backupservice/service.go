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

// ListSnapshots returns all snapshots
func (s *Service) ListSnapshots(ctx context.Context) ([]*domain.Snapshot, error) {
	return s.snapshotRepo.GetAll(ctx)
}

// GetSnapshot returns a snapshot by ID
func (s *Service) GetSnapshot(ctx context.Context, id int64) (*domain.Snapshot, error) {
	return s.snapshotRepo.GetByID(ctx, id)
}

// RestoreSnapshot restores a snapshot
func (s *Service) RestoreSnapshot(ctx context.Context, id int64) error {
	// TODO: Implement restore logic
	return fmt.Errorf("restore not implemented")
}

// GetManifest returns the manifest for a snapshot
func (s *Service) GetManifest(ctx context.Context, id int64) ([]byte, error) {
	snapshot, err := s.snapshotRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot: %w", err)
	}

	source, err := s.sourceRepo.GetByID(ctx, snapshot.SourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source: %w", err)
	}

	if source.TargetID == nil {
		return nil, fmt.Errorf("source has no target")
	}

	// We need the backend to retrieve the manifest
	// Ideally, we should inject the target service or have a way to get the backend
	// For now, let's assume we can't easily get the backend here without circular deps if we inject TargetService
	// Refactoring might be needed, but for now, let's change the signature or dependency.

	// WAIT: Service struct doesn't have TargetService. It has TargetRepo.
	// But TargetRepo doesn't give us the backend client.
	// The BackupHandler has both services.
	// So the Handler should get the backend and pass it to the Service?
	// Or the Service should have a way to create the backend.

	// In RunBackup, we pass the backend.
	// Let's do the same here.
	return nil, fmt.Errorf("not implemented: requires backend injection")
}

// GetManifestWithBackend returns the manifest using the provided backend
func (s *Service) GetManifestWithBackend(ctx context.Context, id int64, backend domain.Backend) ([]byte, error) {
	// The manifest is stored with the snapshot ID as key (or similar)
	// In RunBackup: backend.StoreManifest(ctx, strconv.FormatInt(snapshot.ID, 10), manifestJSON)

	manifestJSON, err := backend.LoadManifest(ctx, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve manifest: %w", err)
	}

	return manifestJSON, nil
}

// FileNode represents a node in the file tree
type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"is_dir"`
	Size     int64       `json:"size,omitempty"`
	ModTime  string      `json:"mod_time,omitempty"`
	Children []*FileNode `json:"children,omitempty"`
}

// GetSnapshotFileTree builds a hierarchical file tree from the manifest
func (s *Service) GetSnapshotFileTree(ctx context.Context, id int64, backend domain.Backend) (*FileNode, error) {
	// Get manifest
	manifestJSON, err := s.GetManifestWithBackend(ctx, id, backend)
	if err != nil {
		return nil, err
	}

	var manifest domain.Manifest
	if err := json.Unmarshal(manifestJSON, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Create root node
	root := &FileNode{
		Name:     filepath.Base(manifest.SourcePath),
		Path:     manifest.SourcePath,
		IsDir:    true,
		Children: make([]*FileNode, 0),
	}

	// Build tree from flat file list
	for _, file := range manifest.Files {
		s.insertFileIntoTree(root, file, manifest.SourcePath)
	}

	return root, nil
}

// insertFileIntoTree inserts a file into the tree structure
func (s *Service) insertFileIntoTree(root *FileNode, file domain.ManifestFile, sourcePath string) {
	// Get relative path
	relPath, err := filepath.Rel(sourcePath, file.Path)
	if err != nil {
		// If we can't get relative path, use the full path
		relPath = file.Path
	}

	// Split path into components
	parts := filepath.SplitList(relPath)
	if len(parts) == 0 {
		parts = []string{relPath}
	}

	// For paths with separators, split manually
	if len(parts) == 1 && filepath.Separator != ' ' {
		parts = splitPath(relPath)
	}

	current := root

	// Navigate/create directories
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if part == "." || part == "" {
			continue
		}

		// Find or create directory
		found := false
		for _, child := range current.Children {
			if child.Name == part && child.IsDir {
				current = child
				found = true
				break
			}
		}

		if !found {
			// Create new directory node
			newDir := &FileNode{
				Name:     part,
				Path:     filepath.Join(current.Path, part),
				IsDir:    true,
				Children: make([]*FileNode, 0),
			}
			current.Children = append(current.Children, newDir)
			current = newDir
		}
	}

	// Add the file
	fileName := parts[len(parts)-1]
	if fileName != "." && fileName != "" {
		fileNode := &FileNode{
			Name:    fileName,
			Path:    file.Path,
			IsDir:   false,
			Size:    file.Size,
			ModTime: file.ModTime.Format(time.RFC3339),
		}
		current.Children = append(current.Children, fileNode)
	}
}

// splitPath splits a file path into its components
func splitPath(path string) []string {
	var parts []string
	for {
		dir, file := filepath.Split(path)
		if file != "" {
			parts = append([]string{file}, parts...)
		}
		if dir == "" || dir == string(filepath.Separator) {
			break
		}
		path = filepath.Clean(dir)
	}
	return parts
}
