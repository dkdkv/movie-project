// internal/service/movie_service_test.go
package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"movie-project/internal/model"
	"movie-project/pkg/logger"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, movie *model.Movie) error {
	args := m.Called(ctx, movie)
	return args.Error(0)
}

// ... implement other mock methods ...

func TestCreateMovie(t *testing.T) {
	mockRepo := new(MockRepository)
	mockLogger := logger.NewLogger()
	service := NewMovieService(mockRepo, mockLogger)

	movie := &model.Movie{
		Title:       "Test Movie",
		Director:    "Test Director",
		ReleaseDate: time.Now(),
		Genre:       "Test Genre",
		Rating:      8.5,
	}

	mockRepo.On("Create", mock.Anything, movie).Return(nil)

	err := service.CreateMovie(context.Background(), movie)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// ... add more test cases ...
