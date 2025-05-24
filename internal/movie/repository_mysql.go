package movie

import (
	"context"
	"fmt"
	"roketin-case-study-challenge2/internal/entity"
	"strings"

	"gorm.io/gorm"
)

type mySQLMovieRepository struct {
	db *gorm.DB
}

func NewMySQLMovieRepository(db *gorm.DB) MovieRepository {
	return &mySQLMovieRepository{
		db: db,
	}
}

func (r *mySQLMovieRepository) CreateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	result := r.db.WithContext(ctx).Create(movie)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create movie: %w", result.Error)
	}

	return movie, nil
}

func (r *mySQLMovieRepository) ListMovies(ctx context.Context, filter *entity.MovieFilter) ([]entity.Movie, int64, error) {
	var movies []entity.Movie
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Movie{})

	if filter.Title != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(filter.Title)+"%")
	}

	if filter.Description != "" {
		query = query.Where("LOWER(description) LIKE ?", "%"+strings.ToLower(filter.Description)+"%")
	}

	if len(filter.Genres) > 0 {
		var conditions []string
		var values []interface{}

		for _, genre := range filter.Genres {
			conditions = append(conditions, "genres LIKE ?")
			values = append(values, "%"+genre+"%")
		}

		query = query.Where(strings.Join(conditions, " OR "), values...)
	}

	if len(filter.Artists) > 0 {
		var conditions []string
		var values []interface{}

		for _, artist := range filter.Artists {
			conditions = append(conditions, "artists LIKE ?")
			values = append(values, "%"+artist+"%")
		}

		query = query.Where(strings.Join(conditions, " OR "), values...)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get total movies: %w", err)
	}

	page := filter.GetPage()
	limit := filter.GetLimit()
	offset := (page - 1) * limit

	result := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&movies)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to get movies: %w", result.Error)
	}

	return movies, total, nil
}

func (r *mySQLMovieRepository) UpdateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	if movie.ID == 0 {
		return nil, fmt.Errorf("movie ID is required")
	}

	result := r.db.WithContext(ctx).Model(&entity.Movie{}).Where("id = ?", movie.ID).Updates(movie)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update movie: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("movie with ID %d not found", movie.ID)
	}

	var updatedMovie entity.Movie
	if errDb := r.db.WithContext(ctx).First(&updatedMovie, movie.ID).Error; errDb != nil {
		return nil, fmt.Errorf("failed to get updated movie: %w", errDb)
	}

	return &updatedMovie, nil
}

func (r *mySQLMovieRepository) DeleteMovie(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&entity.Movie{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete movie: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("movie with ID %d not found", id)
	}

	return nil
}