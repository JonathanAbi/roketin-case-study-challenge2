package movie

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"roketin-case-study-challenge2/internal/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}

	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open GORM DB: %v", err)
	}

	return db, mock, err
}

func TestCreateMovieRepository(t *testing.T) {
	db, mock, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	repo := NewMySQLMovieRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		movie   *entity.Movie
		mockSQL func()
		wantErr bool
	}{
		{
			name: "success create movie",
			movie: &entity.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Duration:    120,
				Artists:     "Test Artist",
				Genres:      "Action",
				FilePath:    "uploads/test.mp4",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `movies`")).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "fail create movie - database error",
			movie: &entity.Movie{
				Title:       "Test Movie",
				Description: "Test Description",
				Duration:    120,
				Artists:     "Test Artist",
				Genres:      "Action",
				FilePath:    "uploads/test.mp4",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			mockSQL: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `movies`")).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSQL()

			movie, err := repo.CreateMovie(ctx, test.movie)

			if (err != nil) != test.wantErr {
				t.Errorf("CreateMovie() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

			if !test.wantErr {
				if movie == nil {
					t.Error("CreateMovie() returned nil movie")
				}
			}
		})
	}
}

func TestListMoviesRepository(t *testing.T) {
	db, mock, err := setupTestDB(t)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	repo := NewMySQLMovieRepository(db)
	ctx := context.Background()

	tests := []struct {
		name      string
		filter    *entity.MovieFilter
		mockSQL   func()
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "success list movies without filter",
			filter: &entity.MovieFilter{
				Page:  1,
				Limit: 10,
			},
			mockSQL: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `movies`")).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

				rows := sqlmock.NewRows([]string{"id", "title", "description", "duration", "artists", "genres", "file_path", "created_at", "updated_at"}).
					AddRow(1, "Movie 1", "Desc 1", 120, "Artist 1", "Action", "path1.mp4", time.Now(), time.Now()).
					AddRow(2, "Movie 2", "Desc 2", 130, "Artist 2", "Drama", "path2.mp4", time.Now(), time.Now())

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `movies`")).
					WillReturnRows(rows)
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "success list movies with title filter",
			filter: &entity.MovieFilter{
				Title: "Movie 1",
				Page:  1,
				Limit: 10,
			},
			mockSQL: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `movies`")).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				rows := sqlmock.NewRows([]string{"id", "title", "description", "duration", "artists", "genres", "file_path", "created_at", "updated_at"}).
					AddRow(1, "Movie 1", "Desc 1", 120, "Artist 1", "Action", "path1.mp4", time.Now(), time.Now())

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `movies`")).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantTotal: 1,
			wantErr:   false,
		},
		{
			name: "fail list movies - database error",
			filter: &entity.MovieFilter{
				Page:  1,
				Limit: 10,
			},
			mockSQL: func() {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `movies`")).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockSQL()

			movies, total, err := repo.ListMovies(ctx, test.filter)

			if (err != nil) != test.wantErr {
				t.Errorf("ListMovies() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %s", err)
			}

			if !test.wantErr {
				if len(movies) != test.wantCount {
					t.Errorf("ListMovies() got %v movies, want %v", len(movies), test.wantCount)
				}
				if total != test.wantTotal {
					t.Errorf("ListMovies() got total = %v, want %v", total, test.wantTotal)
				}
			}
		})
	}
}
