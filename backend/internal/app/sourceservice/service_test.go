package sourceservice

import (
	"context"
	"testing"

	"github.com/axelfrache/savesync/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockSourceRepository is a mock implementation of domain.SourceRepository
type MockSourceRepository struct {
	mock.Mock
}

func (m *MockSourceRepository) Create(ctx context.Context, source *domain.Source) error {
	args := m.Called(ctx, source)
	if args.Get(0) != nil {
		source.ID = 1 // Simulate ID assignment
	}
	return args.Error(0)
}

func (m *MockSourceRepository) GetByID(ctx context.Context, id int64) (*domain.Source, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Source), args.Error(1)
}

func (m *MockSourceRepository) GetAll(ctx context.Context) ([]*domain.Source, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Source), args.Error(1)
}

func (m *MockSourceRepository) Update(ctx context.Context, source *domain.Source) error {
	args := m.Called(ctx, source)
	return args.Error(0)
}

func (m *MockSourceRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSourceService_GetAll(t *testing.T) {
	// Setup
	mockRepo := new(MockSourceRepository)
	logger, _ := zap.NewDevelopment()
	service := New(mockRepo, logger)

	expectedSources := []*domain.Source{
		{ID: 1, Name: "source1", Path: "/tmp/test1"},
		{ID: 2, Name: "source2", Path: "/tmp/test2"},
	}

	mockRepo.On("GetAll", mock.Anything).Return(expectedSources, nil)

	// Execute
	sources, err := service.GetAll(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(sources))
	assert.Equal(t, "source1", sources[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestSourceService_GetByID(t *testing.T) {
	// Setup
	mockRepo := new(MockSourceRepository)
	logger, _ := zap.NewDevelopment()
	service := New(mockRepo, logger)

	expectedSource := &domain.Source{
		ID:   1,
		Name: "test-source",
		Path: "/tmp/test",
	}

	mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedSource, nil)

	// Execute
	source, err := service.GetByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, source)
	assert.Equal(t, "test-source", source.Name)
	mockRepo.AssertExpectations(t)
}

func TestSourceService_GetByID_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockSourceRepository)
	logger, _ := zap.NewDevelopment()
	service := New(mockRepo, logger)

	mockRepo.On("GetByID", mock.Anything, int64(999)).Return(nil, domain.ErrNotFound)

	// Execute
	source, err := service.GetByID(context.Background(), 999)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
	assert.Nil(t, source)
	mockRepo.AssertExpectations(t)
}
