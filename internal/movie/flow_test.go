package movie

import (
	"context"
	"fmt"
	"roketin-case-study-challenge2/internal/entity"
	"testing"
)

type MockMovieRepository struct {
	movies []entity.Movie
	err    error
}

func (m *MockMovieRepository) CreateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	if m.err != nil {
		return nil, m.err
	}

	movie.ID = len(m.movies) + 1
	m.movies = append(m.movies, *movie)

	return movie, nil
}

func (m *MockMovieRepository) ListMovies(ctx context.Context, filter *entity.MovieFilter) ([]entity.Movie, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	return m.movies, int64(len(m.movies)), nil
}

func (m *MockMovieRepository) UpdateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	if m.err != nil {
		return nil, m.err
	}

	if len(m.movies) == 0 {
		m.movies = []entity.Movie{
			{
				ID:          1,
				Title:       "Original Movie",
				Description: "Original Description",
				Duration:    120,
				Artists:     "Original Artist",
				Genres:      "Action",
			},
		}
	}

	for i, mov := range m.movies {
		if mov.ID == movie.ID {
			m.movies[i] = *movie
			return movie, nil
		}
	}

	return nil, fmt.Errorf("movie with ID %d not found", movie.ID)
}

func (m *MockMovieRepository) DeleteMovie(ctx context.Context, id int) error {
	if m.err != nil {
		return m.err
	}

	if len(m.movies) == 0 {
		m.movies = []entity.Movie{
			{
				ID:          1,
				Title:       "Movie to Delete",
				Description: "Description",
				Duration:    120,
				Artists:     "Artist",
				Genres:      "Action",
			},
		}
	}

	for i, mov := range m.movies {
		if mov.ID == id {
			m.movies = append(m.movies[:i], m.movies[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("movie with ID %d not found", id)
}

func TestCreateMovie(t *testing.T) {
	tests := []struct {
		name      string
		movie     entity.Movie
		mockError error
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success create movie",
			movie: entity.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Duration:    25,
				Artists:     "Test Artist",
				Genres:      "Horror",
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name: "failed create movie - empty title",
			movie: entity.Movie{
				Title:       "",
				Description: "Test Description",
				Duration:    25,
				Artists:     "Test Artist",
				Genres:      "Horror",
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "title is required",
		},
		{
			name: "failed create movie - error creating movie",
			movie: entity.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Duration:    25,
				Artists:     "Test Artist",
				Genres:      "Horror",
			},
			mockError: fmt.Errorf("failed to create movie"),
			wantErr:   true,
			errMsg:    "failed to create movie",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := &MockMovieRepository{
				err: test.mockError,
			}

			flow := NewMovieFlow(mockRepo)

			movie, err := flow.CreateMovie(context.Background(), &test.movie)

			if (err != nil) != test.wantErr {
				t.Errorf("CreateMovie() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				if err == nil {
					t.Error("CreateMovie() expected error but got nil")
					return
				}

				if err.Error() != test.errMsg {
					t.Errorf("CreateMovie() error message = %v, want %v", err.Error(), test.errMsg)
				}

				return
			}

			if movie == nil {
				t.Error("CreateMovie() got nil, want non-nil")
			}
		})
	}
}

func TestListMovies(t *testing.T) {
	tests := []struct {
		name      string
		filter    *entity.MovieFilter
		mockData  []entity.Movie
		mockError error
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success list movies",
			filter: &entity.MovieFilter{
				Page:  1,
				Limit: 10,
			},
			mockData: []entity.Movie{
				{
					ID:          1,
					Title:       "Movie 1",
					Description: "Description 1",
					Duration:    120,
					Artists:     "Artist 1",
					Genres:      "Action",
				},
				{
					ID:          2,
					Title:       "Movie 2",
					Description: "Description 2",
					Duration:    130,
					Artists:     "Artist 2",
					Genres:      "Drama",
				},
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name: "success list movies with filter",
			filter: &entity.MovieFilter{
				Title: "Movie 1",
				Page:  1,
				Limit: 10,
			},
			mockData: []entity.Movie{
				{
					ID:          1,
					Title:       "Movie 1",
					Description: "Description 1",
					Duration:    120,
					Artists:     "Artist 1",
					Genres:      "Action",
				},
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name: "failed list movies - error getting total",
			filter: &entity.MovieFilter{
				Page:  1,
				Limit: 10,
			},
			mockData:  nil,
			mockError: fmt.Errorf("failed to get total movies"),
			wantErr:   true,
			errMsg:    "failed to get total movies",
		},
		{
			name: "failed list movies - error getting movies",
			filter: &entity.MovieFilter{
				Page:  1,
				Limit: 10,
			},
			mockData:  nil,
			mockError: fmt.Errorf("failed to get movies"),
			wantErr:   true,
			errMsg:    "failed to get movies",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := &MockMovieRepository{
				movies: test.mockData,
				err:    test.mockError,
			}
			flow := NewMovieFlow(mockRepo)

			movies, total, err := flow.ListMovies(context.Background(), test.filter)

			if (err != nil) != test.wantErr {
				t.Errorf("ListMovies() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				if err == nil {
					t.Error("ListMovies() expected error but got nil")
					return
				}
				if err.Error() != test.errMsg {
					t.Errorf("ListMovies() error message = %v, want %v", err.Error(), test.errMsg)
				}
				return
			}

			if len(movies) != len(test.mockData) {
				t.Errorf("ListMovies() got %v items, want %v", len(movies), len(test.mockData))
			}

			if total != int64(len(test.mockData)) {
				t.Errorf("ListMovies() total = %v, want %v", total, len(test.mockData))
			}
		})
	}
}

func TestUpdateMovie(t *testing.T) {
	tests := []struct {
		name      string
		movie     *entity.Movie
		mockError error
		wantErr   bool
		errMsg    string
	}{
		{
			name: "success update movie",
			movie: &entity.Movie{
				ID:          1,
				Title:       "Updated Movie",
				Description: "Updated Description",
				Duration:    150,
				Artists:     "Updated Artist",
				Genres:      "Action, Drama",
			},
			mockError: nil,
			wantErr:   false,
		},
		{
			name: "fail update movie - empty ID",
			movie: &entity.Movie{
				ID:          0,
				Title:       "Updated Movie",
				Description: "Updated Description",
				Duration:    150,
				Artists:     "Updated Artist",
				Genres:      "Action, Drama",
			},
			mockError: nil,
			wantErr:   true,
			errMsg:    "movie ID is required",
		},
		{
			name: "fail update movie - not found",
			movie: &entity.Movie{
				ID:          999,
				Title:       "Updated Movie",
				Description: "Updated Description",
				Duration:    150,
				Artists:     "Updated Artist",
				Genres:      "Action, Drama",
			},
			mockError: fmt.Errorf("movie with ID 999 not found"),
			wantErr:   true,
			errMsg:    "movie with ID 999 not found",
		},
		{
			name: "fail update movie - database error",
			movie: &entity.Movie{
				ID:          1,
				Title:       "Updated Movie",
				Description: "Updated Description",
				Duration:    150,
				Artists:     "Updated Artist",
				Genres:      "Action, Drama",
			},
			mockError: fmt.Errorf("failed to update movie"),
			wantErr:   true,
			errMsg:    "failed to update movie",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := &MockMovieRepository{
				err: test.mockError,
			}
			flow := NewMovieFlow(mockRepo)

			movie, err := flow.UpdateMovie(context.Background(), test.movie)

			if (err != nil) != test.wantErr {
				t.Errorf("UpdateMovie() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				if err == nil {
					t.Error("UpdateMovie() expected error but got nil")
					return
				}
				if err.Error() != test.errMsg {
					t.Errorf("UpdateMovie() error message = %v, want %v", err.Error(), test.errMsg)
				}
				return
			}

			if movie == nil {
				t.Error("UpdateMovie() got nil, want non-nil")
				return
			}
		})
	}
}

func TestDeleteMovie(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		mockError error
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "success delete movie",
			id:        1,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "fail delete movie - not found",
			id:        999,
			mockError: fmt.Errorf("movie with ID 999 not found"),
			wantErr:   true,
			errMsg:    "movie with ID 999 not found",
		},
		{
			name:      "fail delete movie - database error",
			id:        1,
			mockError: fmt.Errorf("failed to delete movie"),
			wantErr:   true,
			errMsg:    "failed to delete movie",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := &MockMovieRepository{
				err: test.mockError,
			}
			flow := NewMovieFlow(mockRepo)

			err := flow.DeleteMovie(context.Background(), test.id)

			if (err != nil) != test.wantErr {
				t.Errorf("DeleteMovie() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				if err == nil {
					t.Error("DeleteMovie() expected error but got nil")
					return
				}
				if err.Error() != test.errMsg {
					t.Errorf("DeleteMovie() error message = %v, want %v", err.Error(), test.errMsg)
				}
				return
			}
		})
	}
}
