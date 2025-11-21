package domain

import "context"

type Backend interface {
	Init(config map[string]string) error
	StoreChunk(ctx context.Context, hash string, data []byte) error
	LoadChunk(ctx context.Context, hash string) ([]byte, error)
	DeleteChunk(ctx context.Context, hash string) error
	ChunkExists(ctx context.Context, hash string) (bool, error)
	StoreManifest(ctx context.Context, snapshotID string, manifest []byte) error
	LoadManifest(ctx context.Context, snapshotID string) ([]byte, error)
	DeleteManifest(ctx context.Context, snapshotID string) error
	Close() error
}
