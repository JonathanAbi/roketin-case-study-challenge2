package movie

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"roketin-case-study-challenge2/internal/entity"
	"strconv"
	"strings"
)

type MovieParserInterface interface {
	ParseCreateMovie(r *http.Request) (*entity.Movie, *multipart.FileHeader, error)
	ParseMovieFilter(r *http.Request) (*entity.MovieFilter, error)
	ParseUpdateMovie(r *http.Request) (*entity.Movie, error)
}

type MovieParser struct {
}

func NewMovieParser() MovieParserInterface {
	return &MovieParser{}
}

func (p *MovieParser) ParseCreateMovie(r *http.Request) (*entity.Movie, *multipart.FileHeader, error) {
	title := r.PostFormValue("title")
	if title == "" {
		return nil, nil, fmt.Errorf("title is required")
	}

	description := r.PostFormValue("description")
	durationStr := r.PostFormValue("duration_minutes")
	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		return nil, nil, fmt.Errorf("duration must be a number")
	}

	artists := r.PostFormValue("artists")
	genres := r.PostFormValue("genres")

	_, file, err := r.FormFile("movie_file")
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, nil, fmt.Errorf("movie file is required")
		}

		return nil, nil, fmt.Errorf("failed to get movie file: %w", err)
	}

	allowedExtensions := map[string]bool{".mp4": true, ".mov": true, ".mkv": true, ".avi": true}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return nil, nil, fmt.Errorf("file extension %s is not allowed", ext)
	}

	movieData := &entity.Movie{
		Title: title,
		Description: description,
		Duration: duration,
		Artists: artists,
		Genres: genres,
	}

	return movieData, file, nil
}

func (p *MovieParser)ParseMovieFilter(r *http.Request) (*entity.MovieFilter, error) {
	query := r.URL.Query()

	title := query.Get("title")
	description := query.Get("description")
	genres := query["genre"]
	artists := query["artist"]
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	var page int 
	if pageStr != "" { 
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, fmt.Errorf("page number is not valid: '%s'", pageStr)
		}
		if p <= 0 {
			return nil, fmt.Errorf("page number must be greater than 0: %d", p)
		}
		page = p
	}

	var limit int
	if limitStr != "" { 
		l, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("limit number is not valid: '%s'", limitStr)
		}
		if l < 0 {
			return nil, fmt.Errorf("limit number must be greater than 0: %d", l)
		}
		limit = l
	}

	return &entity.MovieFilter {
		Title: title,
		Description: description,
		Genres: genres,
		Artists: artists,
		Page: page,
		Limit: limit,
	}, nil
}

func (p *MovieParser) ParseUpdateMovie(r *http.Request) (*entity.Movie, error) {
	title := r.PostFormValue("title")
	description := r.PostFormValue("description")
	durationStr := r.PostFormValue("duration_minutes")
	artists := r.PostFormValue("artists")
	genres := r.PostFormValue("genres")

	movieData := &entity.Movie{
		Title: title,
		Description: description,
		Artists: artists,
		Genres: genres,
	}

	if durationStr != "" { 
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return nil, fmt.Errorf("duration must be a number: '%s'", durationStr)
		}

		movieData.Duration = duration
	}

	return movieData, nil
}