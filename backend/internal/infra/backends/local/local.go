package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/axelfrache/savesync/internal/domain"
)

// Backend implements domain.Backend for local filesystem storage
type Backend struct {
	basePath string
}

// New creates a new local backend
func New() *Backend {
	return &Backend{}
}

// Init initializes the backend with configuration
func (b *Backend) Init(config map[string]string) error {
	path, ok := config["path"]
	if !ok {
		return fmt.Errorf("path is required in config")
	}

	b.basePath = path

	// Create base directories
	if err := os.MkdirAll(filepath.Join(b.basePath, "chunks"), 0755); err != nil {
		return fmt.Errorf("failed to create chunks directory: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(b.basePath, "manifests"), 0755); err != nil {
		return fmt.Errorf("failed to create manifests directory: %w", err)
	}

	return nil
}

// StoreChunk stores a chunk with content-addressable path
func (b *Backend) StoreChunk(ctx context.Context, hash string, data []byte) error {
	if len(hash) < 4 {
		return fmt.Errorf("invalid hash: too short")
	}

	// Create nested directory structure: chunks/ab/cd/abcd1234...
	dir := filepath.Join(b.basePath, "chunks", hash[:2], hash[2:4])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create chunk directory: %w", err)
	}

	chunkPath := filepath.Join(dir, hash)

	// Check if chunk already exists (deduplication)
	if _, err := os.Stat(chunkPath); err == nil {
		return nil // Chunk already exists
	}

	if err := os.WriteFile(chunkPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write chunk: %w", err)
	}

	return nil
}

// LoadChunk loads a chunk by hash
func (b *Backend) LoadChunk(ctx context.Context, hash string) ([]byte, error) {
	if len(hash) < 4 {
		return nil, fmt.Errorf("invalid hash: too short")
	}

	chunkPath := filepath.Join(b.basePath, "chunks", hash[:2], hash[2:4], hash)

	data, err := os.ReadFile(chunkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to read chunk: %w", err)
	}

	return data, nil
}

// DeleteChunk deletes a chunk by hash
func (b *Backend) DeleteChunk(ctx context.Context, hash string) error {
	if len(hash) < 4 {
		return fmt.Errorf("invalid hash: too short")
	}

	chunkPath := filepath.Join(b.basePath, "chunks", hash[:2], hash[2:4], hash)

	if err := os.Remove(chunkPath); err != nil {
		if os.IsNotExist(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete chunk: %w", err)
	}

	return nil
}

// ChunkExists checks if a chunk exists
func (b *Backend) ChunkExists(ctx context.Context, hash string) (bool, error) {
	if len(hash) < 4 {
		return false, fmt.Errorf("invalid hash: too short")
	}

	chunkPath := filepath.Join(b.basePath, "chunks", hash[:2], hash[2:4], hash)

	_, err := os.Stat(chunkPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check chunk: %w", err)
	}

	return true, nil
}

// StoreManifest stores a snapshot manifest
func (b *Backend) StoreManifest(ctx context.Context, snapshotID string, manifest []byte) error {
	manifestPath := filepath.Join(b.basePath, "manifests", snapshotID+".json")

	if err := os.WriteFile(manifestPath, manifest, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// LoadManifest loads a snapshot manifest
func (b *Backend) LoadManifest(ctx context.Context, snapshotID string) ([]byte, error) {
	manifestPath := filepath.Join(b.basePath, "manifests", snapshotID+".json")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	return data, nil
}

// DeleteManifest deletes a snapshot manifest
func (b *Backend) DeleteManifest(ctx context.Context, snapshotID string) error {
	manifestPath := filepath.Join(b.basePath, "manifests", snapshotID+".json")

	if err := os.Remove(manifestPath); err != nil {
		if os.IsNotExist(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to delete manifest: %w", err)
	}

	return nil
}

// Close closes the backend (no-op for local filesystem)
func (b *Backend) Close() error {
	return nil
}
