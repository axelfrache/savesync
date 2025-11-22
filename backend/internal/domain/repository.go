package domain

import "context"

// UserRepository handles user data persistence
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

// SettingsRepository handles application settings persistence
type SettingsRepository interface {
	Get(ctx context.Context, key string) (*Setting, error)
	Set(ctx context.Context, key, value string) error
	GetAll(ctx context.Context) (map[string]string, error)
}

type SourceRepository interface {
	Create(ctx context.Context, source *Source) error
	GetByID(ctx context.Context, id int64) (*Source, error)
	GetAll(ctx context.Context) ([]*Source, error)
	Update(ctx context.Context, source *Source) error
	Delete(ctx context.Context, id int64) error
}

type TargetRepository interface {
	Create(ctx context.Context, target *Target) error
	GetByID(ctx context.Context, id int64) (*Target, error)
	GetAll(ctx context.Context) ([]*Target, error)
	Update(ctx context.Context, target *Target) error
	Delete(ctx context.Context, id int64) error
}

type SnapshotRepository interface {
	Create(ctx context.Context, snapshot *Snapshot) error
	GetByID(ctx context.Context, id int64) (*Snapshot, error)
	GetAll(ctx context.Context) ([]*Snapshot, error)
	GetBySourceID(ctx context.Context, sourceID int64) ([]*Snapshot, error)
	Update(ctx context.Context, snapshot *Snapshot) error
	Delete(ctx context.Context, id int64) error
}

type SnapshotFileRepository interface {
	Create(ctx context.Context, file *SnapshotFile) error
	GetBySnapshotID(ctx context.Context, snapshotID int64) ([]*SnapshotFile, error)
	Delete(ctx context.Context, id int64) error
}

type JobRepository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id int64) (*Job, error)
	GetAll(ctx context.Context) ([]*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id int64) error
}

type ScheduleRepository interface {
	Create(ctx context.Context, schedule *Schedule) error
	GetByID(ctx context.Context, id int64) (*Schedule, error)
	GetBySourceID(ctx context.Context, sourceID int64) (*Schedule, error)
	GetAll(ctx context.Context) ([]*Schedule, error)
	GetEnabled(ctx context.Context) ([]*Schedule, error)
	Update(ctx context.Context, schedule *Schedule) error
	Delete(ctx context.Context, id int64) error
}

type ChunkRepository interface {
	Create(ctx context.Context, chunk *Chunk) error
	GetByHash(ctx context.Context, hash string) (*Chunk, error)
	Exists(ctx context.Context, hash string) (bool, error)
	IncrementRefCount(ctx context.Context, hash string) error
	DecrementRefCount(ctx context.Context, hash string) error
	GetUnreferenced(ctx context.Context) ([]*Chunk, error)
	Delete(ctx context.Context, hash string) error
}

type ManifestRepository interface {
	Create(ctx context.Context, manifest *Manifest) error
	GetBySnapshotID(ctx context.Context, snapshotID int64) (*Manifest, error)
	Delete(ctx context.Context, snapshotID int64) error
}

type BackupRepository interface {
	Create(ctx context.Context, backup *Snapshot) error
	GetByID(ctx context.Context, id int64) (*Snapshot, error)
	GetBySourceID(ctx context.Context, sourceID int64) ([]*Snapshot, error)
	GetLatest(ctx context.Context, sourceID int64) (*Snapshot, error)
	Update(ctx context.Context, backup *Snapshot) error
	Delete(ctx context.Context, id int64) error
}

type RestoreRepository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id int64) (*Job, error)
	GetAll(ctx context.Context) ([]*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id int64) error
}
