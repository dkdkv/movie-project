// internal/handler/movie_handler.go
package handler

import (
	"context"
	"gorm.io/gorm"
	"movie-project/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"movie-project/internal/model"
	"movie-project/pkg/logger"
	pb "movie-project/proto/movie"
)

type MovieHandler struct {
	pb.UnimplementedMovieServiceServer
	service service.MovieService
	logger  logger.Logger
}

func NewMovieHandler(service service.MovieService, logger logger.Logger) MovieHandler {
	return MovieHandler{service: service, logger: logger}
}

func (h *MovieHandler) CreateMovie(ctx context.Context, req *pb.CreateMovieRequest) (*pb.Movie, error) {
	movie := &model.Movie{
		Title:       req.Title,
		Director:    req.Director,
		ReleaseDate: req.ReleaseDate.AsTime(),
		Genre:       req.Genre,
		Rating:      req.Rating,
	}

	err := h.service.CreateMovie(ctx, movie)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create movie: %v", err)
	}

	return modelToProto(movie), nil
}

func (h *MovieHandler) GetMovie(ctx context.Context, req *pb.GetMovieRequest) (*pb.Movie, error) {
	movie, err := h.service.GetMovie(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Movie not found: %v", err)
	}

	return modelToProto(movie), nil
}

func (h *MovieHandler) ListMovies(ctx context.Context, req *pb.ListMoviesRequest) (*pb.ListMoviesResponse, error) {
	if req.PageNumber < 1 {
		req.PageNumber = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10 // или любое другое значение по умолчанию
	}

	h.logger.InfoContext(ctx, "Listing movies request", "pageNumber", req.PageNumber, "pageSize", req.PageSize)

	movies, total, err := h.service.ListMovies(ctx, int(req.PageNumber), int(req.PageSize))
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to list movies", "error", err)
		return nil, status.Errorf(codes.Internal, "Failed to list movies: %v", err)
	}

	h.logger.InfoContext(ctx, "Movies retrieved", "count", len(movies), "total", total)

	pbMovies := make([]*pb.Movie, len(movies))
	for i, movie := range movies {
		pbMovies[i] = modelToProto(movie)
		h.logger.InfoContext(ctx, "Movie converted to proto", "id", movie.ID, "title", movie.Title)
	}

	response := &pb.ListMoviesResponse{
		Movies:     pbMovies,
		TotalCount: int32(total),
	}

	h.logger.InfoContext(ctx, "Response prepared", "moviesCount", len(response.Movies), "totalCount", response.TotalCount)

	return response, nil
}

func (h *MovieHandler) UpdateMovie(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.Movie, error) {
	movie := &model.Movie{
		Model:       gorm.Model{ID: uint(req.Id)},
		Title:       req.Title,
		Director:    req.Director,
		ReleaseDate: req.ReleaseDate.AsTime(),
		Genre:       req.Genre,
		Rating:      req.Rating,
	}

	err := h.service.UpdateMovie(ctx, movie)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update movie: %v", err)
	}

	return modelToProto(movie), nil
}

func (h *MovieHandler) DeleteMovie(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	err := h.service.DeleteMovie(ctx, uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete movie: %v", err)
	}

	return &pb.DeleteMovieResponse{Success: true}, nil
}

func modelToProto(movie *model.Movie) *pb.Movie {
	return &pb.Movie{
		Id:          int64(movie.ID),
		Title:       movie.Title,
		Director:    movie.Director,
		ReleaseDate: timestamppb.New(movie.ReleaseDate),
		Genre:       movie.Genre,
		Rating:      movie.Rating,
	}
}
