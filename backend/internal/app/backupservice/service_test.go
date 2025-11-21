package backupservice

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/axelfrache/savesync/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mocks
type MockSourceRepository struct{ mock.Mock }

func (m *MockSourceRepository) GetByID(ctx context.Context, id int64) (*domain.Source, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Source), args.Error(1)
}
func (m *MockSourceRepository) Create(ctx context.Context, source *domain.Source) error {
	return m.Called(ctx, source).Error(0)
}
func (m *MockSourceRepository) GetAll(ctx context.Context) ([]*domain.Source, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Source), args.Error(1)
}
func (m *MockSourceRepository) Update(ctx context.Context, source *domain.Source) error {
	return m.Called(ctx, source).Error(0)
}
func (m *MockSourceRepository) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

type MockTargetRepository struct{ mock.Mock }

func (m *MockTargetRepository) GetByID(ctx context.Context, id int64) (*domain.Target, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Target), args.Error(1)
}
func (m *MockTargetRepository) Create(ctx context.Context, target *domain.Target) error {
	return m.Called(ctx, target).Error(0)
}
func (m *MockTargetRepository) GetAll(ctx context.Context) ([]*domain.Target, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Target), args.Error(1)
}
func (m *MockTargetRepository) Update(ctx context.Context, target *domain.Target) error {
	return m.Called(ctx, target).Error(0)
}
func (m *MockTargetRepository) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

type MockSnapshotRepository struct{ mock.Mock }

func (m *MockSnapshotRepository) Create(ctx context.Context, snapshot *domain.Snapshot) error {
	args := m.Called(ctx, snapshot)
	if args.Error(0) == nil {
		snapshot.ID = 1
	}
	return args.Error(0)
}
func (m *MockSnapshotRepository) Update(ctx context.Context, snapshot *domain.Snapshot) error {
	return m.Called(ctx, snapshot).Error(0)
}
func (m *MockSnapshotRepository) GetByID(ctx context.Context, id int64) (*domain.Snapshot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Snapshot), args.Error(1)
}
func (m *MockSnapshotRepository) GetAll(ctx context.Context) ([]*domain.Snapshot, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Snapshot), args.Error(1)
}
func (m *MockSnapshotRepository) GetBySourceID(ctx context.Context, sourceID int64) ([]*domain.Snapshot, error) {
	args := m.Called(ctx, sourceID)
	return args.Get(0).([]*domain.Snapshot), args.Error(1)
}
func (m *MockSnapshotRepository) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

type MockJobRepository struct{ mock.Mock }

func (m *MockJobRepository) Create(ctx context.Context, job *domain.Job) error {
	return m.Called(ctx, job).Error(0)
}
func (m *MockJobRepository) Update(ctx context.Context, job *domain.Job) error {
	return m.Called(ctx, job).Error(0)
}
func (m *MockJobRepository) GetByID(ctx context.Context, id int64) (*domain.Job, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Job), args.Error(1)
}
func (m *MockJobRepository) GetAll(ctx context.Context) ([]*domain.Job, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Job), args.Error(1)
}
func (m *MockJobRepository) Delete(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

type MockBackend struct{ mock.Mock }

func (m *MockBackend) Init(config map[string]string) error { return nil }
func (m *MockBackend) Close() error                        { return nil }
func (m *MockBackend) StoreChunk(ctx context.Context, hash string, data []byte) error {
	return m.Called(ctx, hash, data).Error(0)
}
func (m *MockBackend) GetChunk(ctx context.Context, hash string) ([]byte, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *MockBackend) LoadChunk(ctx context.Context, hash string) ([]byte, error) {
	args := m.Called(ctx, hash)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *MockBackend) ChunkExists(ctx context.Context, hash string) (bool, error) {
	args := m.Called(ctx, hash)
	return args.Bool(0), args.Error(1)
}
func (m *MockBackend) DeleteChunk(ctx context.Context, hash string) error {
	return m.Called(ctx, hash).Error(0)
}
func (m *MockBackend) StoreManifest(ctx context.Context, snapshotID string, manifest []byte) error {
	return m.Called(ctx, snapshotID, manifest).Error(0)
}
func (m *MockBackend) LoadManifest(ctx context.Context, snapshotID string) ([]byte, error) {
	args := m.Called(ctx, snapshotID)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *MockBackend) DeleteManifest(ctx context.Context, snapshotID string) error {
	return m.Called(ctx, snapshotID).Error(0)
}

func TestBackupService_RunBackup(t *testing.T) {
	// Setup
	mockSourceRepo := new(MockSourceRepository)
	mockTargetRepo := new(MockTargetRepository)
	mockSnapshotRepo := new(MockSnapshotRepository)
	mockJobRepo := new(MockJobRepository)
	mockBackend := new(MockBackend)
	logger, _ := zap.NewDevelopment()

	service := New(mockSourceRepo, mockTargetRepo, mockSnapshotRepo, mockJobRepo, logger)

	// Create temp dir for source
	tmpDir, err := os.MkdirTemp("", "savesync-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	assert.NoError(t, err)

	targetID := int64(2)
	source := &domain.Source{
		ID:       1,
		Name:     "test-source",
		Path:     tmpDir,
		TargetID: &targetID,
	}

	mockSourceRepo.On("GetByID", mock.Anything, int64(1)).Return(source, nil)
	mockSnapshotRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Snapshot")).Return(nil)
	mockSnapshotRepo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Snapshot")).Return(nil)

	// Mock backend calls
	mockBackend.On("ChunkExists", mock.Anything, mock.Anything).Return(false, nil)
	mockBackend.On("StoreChunk", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockBackend.On("StoreManifest", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Execute
	err = service.RunBackup(context.Background(), 1, mockBackend)

	// Assert
	assert.NoError(t, err)
	mockSourceRepo.AssertExpectations(t)
	mockSnapshotRepo.AssertExpectations(t)
	mockBackend.AssertExpectations(t)
}

func TestBackupService_ListSnapshots(t *testing.T) {
	// Setup
	mockSourceRepo := new(MockSourceRepository)
	mockTargetRepo := new(MockTargetRepository)
	mockSnapshotRepo := new(MockSnapshotRepository)
	mockJobRepo := new(MockJobRepository)
	logger, _ := zap.NewDevelopment()

	service := New(mockSourceRepo, mockTargetRepo, mockSnapshotRepo, mockJobRepo, logger)

	expectedSnapshots := []*domain.Snapshot{
		{ID: 1, Status: "success"},
		{ID: 2, Status: "failed"},
	}

	mockSnapshotRepo.On("GetAll", mock.Anything).Return(expectedSnapshots, nil)

	// Execute
	snapshots, err := service.ListSnapshots(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(snapshots))
	mockSnapshotRepo.AssertExpectations(t)
}
