package movie

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"roketin-case-study-challenge2/internal/entity"
	"roketin-case-study-challenge2/internal/response"
	"testing"

	"github.com/go-chi/chi"
)

type MockMovieFlow struct {
	movies     []entity.Movie
	err        error
	totalItems int64
}

func (m *MockMovieFlow) CreateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	if m.err != nil {
		return nil, m.err
	}
	movie.ID = len(m.movies) + 1
	m.movies = append(m.movies, *movie)
	return movie, nil
}

func (m *MockMovieFlow) ListMovies(ctx context.Context, filter *entity.MovieFilter) ([]entity.Movie, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.movies, m.totalItems, nil
}

func (m *MockMovieFlow) UpdateMovie(ctx context.Context, movie *entity.Movie) (*entity.Movie, error) {
	if m.err != nil {
		return nil, m.err
	}
	for i, mov := range m.movies {
		if mov.ID == movie.ID {
			m.movies[i] = *movie
			return movie, nil
		}
	}
	return movie, nil
}

func (m *MockMovieFlow) DeleteMovie(ctx context.Context, id int) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestCreateMovieHandler(t *testing.T) {
	tests := []struct {
		name         string
		formData     map[string]string
		fileData     string
		fileName     string
		mockError    error
		wantStatus   int
		wantResponse bool
		wantErrorMsg string
	}{
		{
			name: "success create movie",
			formData: map[string]string{
				"title":            "Test Movie",
				"description":      "Test Description",
				"duration_minutes": "120",
				"artists":          "Test Artist",
				"genres":           "Action",
			},
			fileData:     "test content",
			fileName:     "test.mp4",
			mockError:    nil,
			wantStatus:   http.StatusOK,
			wantResponse: true,
		},
		{
			name: "fail - empty title",
			formData: map[string]string{
				"title":            "",
				"description":      "Test Description",
				"duration_minutes": "120",
				"artists":          "Test Artist",
				"genres":           "Action",
			},
			fileData:     "test content",
			fileName:     "test.mp4",
			mockError:    nil,
			wantStatus:   http.StatusBadRequest,
			wantResponse: false,
			wantErrorMsg: "title is required",
		},
		{
			name: "fail - invalid duration",
			formData: map[string]string{
				"title":            "Test Movie",
				"description":      "Test Description",
				"duration_minutes": "abc",
				"artists":          "Test Artist",
				"genres":           "Action",
			},
			fileData:     "test content",
			fileName:     "test.mp4",
			mockError:    nil,
			wantStatus:   http.StatusBadRequest,
			wantResponse: false,
			wantErrorMsg: "duration must be a number",
		},
		{
			name: "fail - no file",
			formData: map[string]string{
				"title":            "Test Movie",
				"description":      "Test Description",
				"duration_minutes": "120",
				"artists":          "Test Artist",
				"genres":           "Action",
			},
			mockError:    nil,
			wantStatus:   http.StatusBadRequest,
			wantResponse: false,
			wantErrorMsg: "movie file is required",
		},
		{
			name: "fail - flow error",
			formData: map[string]string{
				"title":            "Test Movie",
				"description":      "Test Description",
				"duration_minutes": "120",
				"artists":          "Test Artist",
				"genres":           "Action",
			},
			fileData:     "test content",
			fileName:     "test.mp4",
			mockError:    fmt.Errorf("failed to create movie"),
			wantStatus:   http.StatusInternalServerError,
			wantResponse: false,
			wantErrorMsg: "failed to create movie",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range test.formData {
				err := writer.WriteField(key, value)
				if err != nil {
					t.Fatal(err)
				}
			}

			if test.fileName != "" {
				part, err := writer.CreateFormFile("movie_file", test.fileName)
				if err != nil {
					t.Fatal(err)
				}
				part.Write([]byte(test.fileData))
			}

			writer.Close()

			req := httptest.NewRequest("POST", "/api/movies", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rr := httptest.NewRecorder()

			mockFlow := &MockMovieFlow{
				err: test.mockError,
			}

			handler := NewMovieHandler(NewMovieParser(), mockFlow)

			handler.CreateMovie(rr, req)

			if rr.Code != test.wantStatus {
				t.Errorf("CreateMovie() status = %v, want %v", rr.Code, test.wantStatus)
			}

			var resp response.Response
			err := json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				t.Fatal(err)
			}

			if test.wantResponse {
				if resp.Status != "success" {
					t.Errorf("CreateMovie() response status = %v, want success", resp.Status)
				}
				if resp.Data == nil {
					t.Error("CreateMovie() response data is nil")
				}
			} else {
				if resp.Status != "error" {
					t.Errorf("CreateMovie() response status = %v, want error", resp.Status)
				}
				if resp.Message != test.wantErrorMsg {
					t.Errorf("CreateMovie() error message = %v, want %v", resp.Message, test.wantErrorMsg)
				}
			}
		})
	}
}

func TestListMoviesHandler(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockMovies     []entity.Movie
		mockTotalItems int64
		mockError      error
		wantStatus     int
		wantResponse   bool
		wantErrorMsg   string
	}{
		{
			name: "success list movies",
			mockMovies: []entity.Movie{
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
			mockTotalItems: 2,
			wantStatus:     http.StatusOK,
			wantResponse:   true,
		},
		{
			name: "success list movies with filter",
			queryParams: map[string]string{
				"title": "Movie 1",
				"page":  "1",
				"limit": "10",
				"genre": "Action",
			},
			mockMovies: []entity.Movie{
				{
					ID:          1,
					Title:       "Movie 1",
					Description: "Description 1",
					Duration:    120,
					Artists:     "Artist 1",
					Genres:      "Action",
				},
			},
			mockTotalItems: 1,
			wantStatus:     http.StatusOK,
			wantResponse:   true,
		},
		{
			name: "fail - invalid page",
			queryParams: map[string]string{
				"page": "-1",
			},
			mockError:    fmt.Errorf("page number must be greater than 0: -1"),
			wantStatus:   http.StatusBadRequest,
			wantResponse: false,
			wantErrorMsg: "page number must be greater than 0: -1",
		},
		{
			name:         "fail - flow error",
			mockError:    fmt.Errorf("failed to get list movies"),
			wantStatus:   http.StatusInternalServerError,
			wantResponse: false,
			wantErrorMsg: "failed to get list movies",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			baseURL := "/api/movies"
			if len(test.queryParams) > 0 {
				values := url.Values{}
				for key, value := range test.queryParams {
					values.Add(key, value)
				}
				baseURL = baseURL + "?" + values.Encode()
			}

			req := httptest.NewRequest(http.MethodGet, baseURL, nil)

			rr := httptest.NewRecorder()

			mockFlow := &MockMovieFlow{
				movies:     test.mockMovies,
				err:        test.mockError,
				totalItems: test.mockTotalItems,
			}

			handler := NewMovieHandler(NewMovieParser(), mockFlow)

			handler.ListMovies(rr, req)

			if rr.Code != test.wantStatus {
				t.Errorf("ListMovies() status = %v, want %v", rr.Code, test.wantStatus)
			}

			var resp response.Response
			err := json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				t.Fatal(err)
			}

			if test.wantResponse {
				if resp.Status != "success" {
					t.Errorf("ListMovies() response status = %v, want success", resp.Status)
				}
				if resp.Data == nil {
					t.Error("ListMovies() response data is nil")
				}

				responseData := resp.Data.(map[string]interface{})
				pagination, ok := responseData["pagination"].(map[string]interface{})
				if !ok {
					t.Error("ListMovies() response does not contain pagination")
				}

				if pagination["total_items"].(float64) != float64(test.mockTotalItems) {
					t.Errorf("ListMovies() total items = %v, want %v", pagination["total_items"], test.mockTotalItems)
				}
			} else {
				if resp.Status != "error" {
					t.Errorf("ListMovies() response status = %v, want error", resp.Status)
				}
				if resp.Message != test.wantErrorMsg {
					t.Errorf("ListMovies() error message = %v, want %v", resp.Message, test.wantErrorMsg)
				}
			}
		})
	}
}

func TestUpdateMovieHandler(t *testing.T) {
	tests := []struct {
		name         string
		movieID      string
		formData     map[string]string
		mockError    error
		wantStatus   int
		wantResponse bool
		wantErrorMsg string
	}{
		{
			name:    "success update movie",
			movieID: "1",
			formData: map[string]string{
				"title":            "Updated Movie",
				"description":      "Updated Description",
				"duration_minutes": "150",
				"artists":          "Updated Artist",
				"genres":           "Action, Drama",
			},
			wantStatus:   http.StatusOK,
			wantResponse: true,
		},
		{
			name:    "success partial update",
			movieID: "1",
			formData: map[string]string{
				"title":       "Updated Movie",
				"description": "Updated Description",
			},
			wantStatus:   http.StatusOK,
			wantResponse: true,
		},
		{
			name:         "fail - invalid movie ID",
			movieID:      "invalid",
			wantStatus:   http.StatusBadRequest,
			wantResponse: false,
			wantErrorMsg: "invalid movie ID",
		},
		{
			name:    "fail - invalid duration",
			movieID: "1",
			formData: map[string]string{
				"title":            "Updated Movie",
				"duration_minutes": "abc",
			},
			wantStatus:   http.StatusBadRequest,
			wantResponse: false,
			wantErrorMsg: "duration must be a number",
		},
		{
			name:    "fail - flow error",
			movieID: "1",
			formData: map[string]string{
				"title": "Updated Movie",
			},
			mockError:    fmt.Errorf("failed to update movie"),
			wantStatus:   http.StatusInternalServerError,
			wantResponse: false,
			wantErrorMsg: "failed to update movie",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			for key, value := range test.formData {
				err := writer.WriteField(key, value)
				if err != nil {
					t.Fatal(err)
				}
			}
			writer.Close()

			req := httptest.NewRequest("PUT", "/api/movies/"+test.movieID, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.movieID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			mockFlow := &MockMovieFlow{
				err: test.mockError,
			}

			handler := NewMovieHandler(NewMovieParser(), mockFlow)

			handler.UpdateMovie(rr, req)

			if rr.Code != test.wantStatus {
				t.Errorf("UpdateMovie() status = %v, want %v", rr.Code, test.wantStatus)
			}

			var resp response.Response
			err := json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				t.Fatal(err)
			}

			if test.wantResponse {
				if resp.Status != "success" {
					t.Errorf("UpdateMovie() response status = %v, want success", resp.Status)
				}
				if resp.Data == nil {
					t.Error("UpdateMovie() response data is nil")
				}
			} else {
				if resp.Status != "error" {
					t.Errorf("UpdateMovie() response status = %v, want error", resp.Status)
				}
				if resp.Message != test.wantErrorMsg {
					t.Errorf("UpdateMovie() error message = %v, want %v", resp.Message, test.wantErrorMsg)
				}
			}
		})
	}
}

func TestDeleteMovieHandler(t *testing.T) {
	tests := []struct {
		name         string
		movieID      string
		mockError    error
		wantStatus   int
		wantResponse bool
		wantErrorMsg string
	}{
		{
			name:         "success delete movie",
			movieID:      "1",
			wantStatus:   http.StatusOK,
			wantResponse: true,
		},
		{
			name:         "fail - invalid movie ID",
			movieID:      "invalid",
			wantStatus:   http.StatusBadRequest,
			wantResponse: false,
			wantErrorMsg: "invalid movie ID",
		},
		{
			name:         "fail - flow error",
			movieID:      "1",
			mockError:    fmt.Errorf("failed to delete movie"),
			wantStatus:   http.StatusInternalServerError,
			wantResponse: false,
			wantErrorMsg: "failed to delete movie",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/api/movies/"+test.movieID, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", test.movieID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			mockFlow := &MockMovieFlow{
				err: test.mockError,
			}

			handler := NewMovieHandler(NewMovieParser(), mockFlow)

			handler.DeleteMovie(rr, req)

			if rr.Code != test.wantStatus {
				t.Errorf("DeleteMovie() status = %v, want %v", rr.Code, test.wantStatus)
			}

			var resp response.Response
			err := json.NewDecoder(rr.Body).Decode(&resp)
			if err != nil {
				t.Fatal(err)
			}

			if test.wantResponse {
				if resp.Status != "success" {
					t.Errorf("DeleteMovie() response status = %v, want success", resp.Status)
				}
			} else {
				if resp.Status != "error" {
					t.Errorf("DeleteMovie() response status = %v, want error", resp.Status)
				}
				if resp.Message != test.wantErrorMsg {
					t.Errorf("DeleteMovie() error message = %v, want %v", resp.Message, test.wantErrorMsg)
				}
			}
		})
	}
}
