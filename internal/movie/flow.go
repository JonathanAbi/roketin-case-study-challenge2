package movie

import (
	"context"
	"fmt"
	"roketin-case-study-challenge2/internal/entity"
	"time"
)

type MovieFlowInterface interface {
	CreateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
	ListMovies(ctx context.Context, filter *entity.MovieFilter) ([]entity.Movie, int64, error)
	UpdateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error)
	DeleteMovie(ctx context.Context, id int) error
}

type movieFlow struct {
	movieRepo MovieRepository
}

func NewMovieFlow(movieRepo MovieRepository) MovieFlowInterface {
	return &movieFlow{
		movieRepo: movieRepo,
	}
}

func (f *movieFlow) CreateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	if movie.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	currentTime := time.Now()
	movie.CreatedAt = currentTime
	movie.UpdatedAt = currentTime

	createdMovie, err := f.movieRepo.CreateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	return createdMovie, nil
}

func (f *movieFlow) ListMovies(ctx context.Context, filter *entity.MovieFilter) ([]entity.Movie, int64, error) {
	movies, total, err := f.movieRepo.ListMovies(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return movies, total, nil
}

func (f *movieFlow) UpdateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	if movie.ID == 0 {
		return nil, fmt.Errorf("movie ID is required")
	}

	movie.UpdatedAt = time.Now()

	updatedMovie, err := f.movieRepo.UpdateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	return updatedMovie, nil
}

func (f *movieFlow) DeleteMovie(ctx context.Context, id int) error {
	err := f.movieRepo.DeleteMovie(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
