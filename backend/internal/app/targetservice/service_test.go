package targetservice

import (
	"context"
	"testing"

	"github.com/axelfrache/savesync/internal/domain"
	"github.com/axelfrache/savesync/internal/infra/backends"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockTargetRepository is a mock implementation of domain.TargetRepository
type MockTargetRepository struct {
	mock.Mock
}

func (m *MockTargetRepository) Create(ctx context.Context, target *domain.Target) error {
	args := m.Called(ctx, target)
	if args.Error(0) == nil {
		target.ID = 1 // Simulate ID assignment
	}
	return args.Error(0)
}

func (m *MockTargetRepository) GetByID(ctx context.Context, id int64) (*domain.Target, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Target), args.Error(1)
}

func (m *MockTargetRepository) GetAll(ctx context.Context) ([]*domain.Target, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Target), args.Error(1)
}

func (m *MockTargetRepository) Update(ctx context.Context, target *domain.Target) error {
	args := m.Called(ctx, target)
	return args.Error(0)
}

func (m *MockTargetRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTargetService_Create(t *testing.T) {
	// Setup
	mockRepo := new(MockTargetRepository)
	logger, _ := zap.NewDevelopment()
	registry := backends.NewRegistry()
	// Register a dummy backend for testing
	registry.Register("local", func() domain.Backend {
		return &MockBackend{}
	})

	service := New(mockRepo, registry, logger)

	target := &domain.Target{
		Name:       "test-target",
		Type:       "local",
		ConfigJSON: `{"path":"/tmp/backup"}`,
	}

	mockRepo.On("Create", mock.Anything, target).Return(nil)

	// Execute
	err := service.Create(context.Background(), target)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), target.ID)
	mockRepo.AssertExpectations(t)
}

func TestTargetService_Create_InvalidType(t *testing.T) {
	// Setup
	mockRepo := new(MockTargetRepository)
	logger, _ := zap.NewDevelopment()
	registry := backends.NewRegistry()
	service := New(mockRepo, registry, logger)

	target := &domain.Target{
		Name: "test-target",
		Type: "invalid-type",
	}

	// Execute
	err := service.Create(context.Background(), target)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported backend type")
}

func TestTargetService_GetByID(t *testing.T) {
	// Setup
	mockRepo := new(MockTargetRepository)
	logger, _ := zap.NewDevelopment()
	registry := backends.NewRegistry()
	service := New(mockRepo, registry, logger)

	expectedTarget := &domain.Target{
		ID:   1,
		Name: "test-target",
		Type: "local",
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedTarget, nil)

	// Execute
	target, err := service.GetByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTarget, target)
	mockRepo.AssertExpectations(t)
}

func TestTargetService_Delete(t *testing.T) {
	// Setup
	mockRepo := new(MockTargetRepository)
	logger, _ := zap.NewDevelopment()
	registry := backends.NewRegistry()
	service := New(mockRepo, registry, logger)

	mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

	// Execute
	err := service.Delete(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// MockBackend is a mock implementation of domain.Backend
type MockBackend struct{}

func (m *MockBackend) Init(config map[string]string) error                            { return nil }
func (m *MockBackend) Close() error                                                   { return nil }
func (m *MockBackend) StoreChunk(ctx context.Context, hash string, data []byte) error { return nil }
func (m *MockBackend) GetChunk(ctx context.Context, hash string) ([]byte, error)      { return nil, nil }
func (m *MockBackend) LoadChunk(ctx context.Context, hash string) ([]byte, error)     { return nil, nil }
func (m *MockBackend) ChunkExists(ctx context.Context, hash string) (bool, error)     { return false, nil }
func (m *MockBackend) DeleteChunk(ctx context.Context, hash string) error             { return nil }
func (m *MockBackend) StoreManifest(ctx context.Context, snapshotID string, manifest []byte) error {
	return nil
}
func (m *MockBackend) LoadManifest(ctx context.Context, snapshotID string) ([]byte, error) {
	return nil, nil
}
func (m *MockBackend) DeleteManifest(ctx context.Context, snapshotID string) error { return nil }
