package movie

import (
	"context"
	"roketin-case-study-challenge2/internal/entity"
)

type MovieRepository interface {
	ListMovies(ctx context.Context, filter *entity.MovieFilter) ([]entity.Movie, int64, error)
	CreateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
	UpdateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
	DeleteMovie(ctx context.Context, id int) error
}