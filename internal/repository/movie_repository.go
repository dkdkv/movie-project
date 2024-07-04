// internal/repository/movie_repository.go
package repository

import (
	"context"

	"gorm.io/gorm"
	"movie-project/internal/model"
	"movie-project/pkg/logger"
)

type IMovieRepository interface {
	Create(ctx context.Context, movie *model.Movie) error
	GetByID(ctx context.Context, id uint) (*model.Movie, error)
	List(ctx context.Context, offset, limit int) ([]*model.Movie, int64, error)
	Update(ctx context.Context, movie *model.Movie) error
	Delete(ctx context.Context, id uint) error
}

type MovieRepository struct {
	db     gorm.DB
	logger logger.Logger
}

func NewMovieRepository(db gorm.DB, logger logger.Logger) MovieRepository {
	return MovieRepository{db: db, logger: logger}
}

func (r *MovieRepository) Create(ctx context.Context, movie *model.Movie) error {
	result := r.db.WithContext(ctx).Create(movie)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to create movie", "error", result.Error)
		return result.Error
	}
	return nil
}

func (r *MovieRepository) GetByID(ctx context.Context, id uint) (*model.Movie, error) {
	var movie model.Movie
	result := r.db.WithContext(ctx).First(&movie, id)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to get movie", "error", result.Error, "id", id)
		return nil, result.Error
	}
	return &movie, nil
}

func (r *MovieRepository) List(ctx context.Context, offset, limit int) ([]*model.Movie, int64, error) {
	var movies []*model.Movie
	var total int64

	if offset < 0 {
		offset = 0
	}
	if limit < 1 {
		limit = 10 // или любое другое значение по умолчанию
	}

	result := r.db.WithContext(ctx).Model(&model.Movie{}).Count(&total)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to count movies", "error", result.Error)
		return nil, 0, result.Error
	}

	r.logger.InfoContext(ctx, "Querying movies", "offset", offset, "limit", limit, "total", total)

	result = r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&movies)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to list movies", "error", result.Error)
		return nil, 0, result.Error
	}

	r.logger.InfoContext(ctx, "Retrieved movies", "count", len(movies), "total", total)

	return movies, total, nil
}

func (r *MovieRepository) Update(ctx context.Context, movie *model.Movie) error {
	result := r.db.WithContext(ctx).Save(movie)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to update movie", "error", result.Error, "id", movie.ID)
		return result.Error
	}
	return nil
}

func (r *MovieRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Movie{}, id)
	if result.Error != nil {
		r.logger.ErrorContext(ctx, "Failed to delete movie", "error", result.Error, "id", id)
		return result.Error
	}
	return nil
}
