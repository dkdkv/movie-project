// internal/service/movie_service.go
package service

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"

	"gorm.io/gorm"

	"movie-project/internal/model"
	"movie-project/internal/repository"
	"movie-project/pkg/logger"
	"movie-project/pkg/metrics"
)

type MovieService struct {
	repo     *repository.MovieRepository
	logger   *logger.Logger
	validate *validator.Validate
}

func NewMovieService(repo *repository.MovieRepository, logger *logger.Logger) *MovieService {
	return &MovieService{
		repo:     repo,
		logger:   logger,
		validate: validator.New(),
	}
}

func (s *MovieService) CreateMovie(ctx context.Context, movie *model.Movie) error {
	if err := s.validate.Struct(movie); err != nil {
		s.logger.ErrorContext(ctx, "Invalid movie data", "error", err)
		return err
	}

	err := s.repo.Create(ctx, movie)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to create movie", "error", err)
		return err
	}

	metrics.MovieCreations.Inc()
	s.logger.InfoContext(ctx, "Created new movie", "id", movie.ID, "title", movie.Title)
	return nil
}

func (s *MovieService) GetMovie(ctx context.Context, id uint) (*model.Movie, error) {
	movie, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WarnContext(ctx, "Movie not found", "id", id)
			return nil, err
		}
		s.logger.ErrorContext(ctx, "Failed to get movie", "error", err, "id", id)
		return nil, err
	}

	metrics.MovieRetrievals.Inc()
	s.logger.InfoContext(ctx, "Retrieved movie", "id", id)
	return movie, nil
}

func (s *MovieService) ListMovies(ctx context.Context, page, pageSize int) ([]*model.Movie, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10 // или любое другое значение по умолчанию
	}

	offset := (page - 1) * pageSize

	s.logger.InfoContext(ctx, "Listing movies", "page", page, "pageSize", pageSize, "offset", offset)

	movies, total, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to list movies", "error", err, "page", page, "pageSize", pageSize)
		return nil, 0, err
	}

	s.logger.InfoContext(ctx, "Listed movies", "page", page, "pageSize", pageSize, "total", total, "retrieved", len(movies))
	return movies, total, nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *model.Movie) error {
	if err := s.validate.Struct(movie); err != nil {
		s.logger.ErrorContext(ctx, "Invalid movie data for update", "error", err)
		return err
	}

	err := s.repo.Update(ctx, movie)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to update movie", "error", err, "id", movie.ID)
		return err
	}

	metrics.MovieUpdates.Inc()
	s.logger.InfoContext(ctx, "Updated movie", "id", movie.ID)
	return nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "Failed to delete movie", "error", err, "id", id)
		return err
	}

	metrics.MovieDeletions.Inc()
	s.logger.InfoContext(ctx, "Deleted movie", "id", id)
	return nil
}
