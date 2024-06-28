// internal/model/movie.go
package model

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "movie-project/proto/movie"
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null" validate:"required,min=1,max=255"`
	Director    string    `json:"director" gorm:"not null" validate:"required,min=1,max=255"`
	ReleaseDate time.Time `json:"release_date" validate:"required"`
	Genre       string    `json:"genre" validate:"required,min=1,max=100"`
	Rating      float32   `json:"rating" gorm:"type:decimal(3,1)" validate:"required,min=0,max=10"`
}

func modelToProto(movie *Movie) *pb.Movie {
	return &pb.Movie{
		Id:          int64(uint32(movie.ID)),
		Title:       movie.Title,
		Director:    movie.Director,
		ReleaseDate: timestamppb.New(movie.ReleaseDate),
		Genre:       movie.Genre,
		Rating:      movie.Rating,
	}
}

func protoToModel(movie *pb.Movie) *Movie {
	return &Movie{
		Model: gorm.Model{
			ID:        uint(movie.Id),
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Title:       movie.Title,
		Director:    movie.Director,
		ReleaseDate: movie.ReleaseDate.AsTime(),
		Genre:       movie.Genre,
		Rating:      movie.Rating,
	}
}
