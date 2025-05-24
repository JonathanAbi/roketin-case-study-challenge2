package movie

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"strconv"
	"testing"
)

func TestParseCreateMovie(t *testing.T) {
	tests := []struct {
		name       string
		formData   map[string]string
		fileData   string
		fileName   string
		wantErr    bool
		errMessage string
	}{
		{
			name: "success parse create movie",
			formData: map[string]string{
				"title":            "Test Movie",
				"description":      "Test Description",
				"duration_minutes": "120",
				"artists":          "Test Artist",
				"genres":           "Action",
			},
			fileData: "test content",
			fileName: "test.mp4",
			wantErr:  false,
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
			fileData:   "test content",
			fileName:   "test.mp4",
			wantErr:    true,
			errMessage: "title is required",
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
			fileData:   "test content",
			fileName:   "test.mp4",
			wantErr:    true,
			errMessage: "duration must be a number",
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
			wantErr:    true,
			errMessage: "movie file is required",
		},
		{
			name: "fail - invalid file extension",
			formData: map[string]string{
				"title":            "Test Movie",
				"description":      "Test Description",
				"duration_minutes": "120",
				"artists":          "Test Artist",
				"genres":           "Action",
			},
			fileData:   "test content",
			fileName:   "test.txt",
			wantErr:    true,
			errMessage: "file extension .txt is not allowed",
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

			req, err := http.NewRequest("POST", "/", body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			parser := NewMovieParser()
			movie, file, err := parser.ParseCreateMovie(req)

			if (err != nil) != test.wantErr {
				t.Errorf("ParseCreateMovie() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				if err == nil {
					t.Error("ParseCreateMovie() expected error but got nil")
					return
				}
				if err.Error() != test.errMessage {
					t.Errorf("ParseCreateMovie() error message = %v, want %v", err.Error(), test.errMessage)
				}
				return
			}

			if movie.Title != test.formData["title"] {
				t.Errorf("ParseCreateMovie() title = %v, want %v", movie.Title, test.formData["title"])
			}

			if movie.Description != test.formData["description"] {
				t.Errorf("ParseCreateMovie() description = %v, want %v", movie.Description, test.formData["description"])
			}

			if file == nil {
				t.Error("ParseCreateMovie() file is nil")
			}
		})
	}
}

func TestParseMovieFilter(t *testing.T) {
	tests := []struct {
		name        string
		queryParams map[string][]string
		wantErr     bool
		errMessage  string
		wantPage    int
		wantLimit   int
		wantTitle   string
		wantGenres  []string
	}{
		{
			name: "success - with all parameters",
			queryParams: map[string][]string{
				"page":  {"2"},
				"limit": {"20"},
				"title": {"Test Movie"},
				"genre": {"Action", "Drama"},
			},
			wantErr:    false,
			wantPage:   2,
			wantLimit:  20,
			wantTitle:  "Test Movie",
			wantGenres: []string{"Action", "Drama"},
		},
		{
			name: "success - without page and limit",
			queryParams: map[string][]string{
				"title": {"Test Movie"},
			},
			wantErr:    false,
			wantPage:   1,
			wantLimit:  10,
			wantTitle:  "Test Movie",
			wantGenres: []string{},
		},
		{
			name: "fail - invalid page",
			queryParams: map[string][]string{
				"page":  {"-1"},
				"title": {"Test Movie"},
			},
			wantErr:    true,
			errMessage: "page number must be greater than 0: -1",
		},
		{
			name: "fail - invalid limit",
			queryParams: map[string][]string{
				"limit": {"-1"},
				"title": {"Test Movie"},
			},
			wantErr:    true,
			errMessage: "limit number must be greater than 0: -1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			q := req.URL.Query()
			for key, values := range test.queryParams {
				for _, value := range values {
					q.Add(key, value)
				}
			}
			req.URL.RawQuery = q.Encode()

			parser := NewMovieParser()
			filter, err := parser.ParseMovieFilter(req)

			if (err != nil) != test.wantErr {
				t.Errorf("ParseMovieFilter() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				if err == nil {
					t.Error("ParseMovieFilter() expected error but got nil")
					return
				}
				if err.Error() != test.errMessage {
					t.Errorf("ParseMovieFilter() error message = %v, want %v", err.Error(), test.errMessage)
				}
				return
			}

			if filter.GetPage() != test.wantPage {
				t.Errorf("ParseMovieFilter() page = %v, want %v", filter.GetPage(), test.wantPage)
			}

			if filter.GetLimit() != test.wantLimit {
				t.Errorf("ParseMovieFilter() limit = %v, want %v", filter.GetLimit(), test.wantLimit)
			}

			if filter.Title != test.wantTitle {
				t.Errorf("ParseMovieFilter() title = %v, want %v", filter.Title, test.wantTitle)
			}

			if len(filter.Genres) != len(test.wantGenres) {
				t.Errorf("ParseMovieFilter() genres length = %v, want %v", len(filter.Genres), len(test.wantGenres))
			}

			for i, genre := range filter.Genres {
				if genre != test.wantGenres[i] {
					t.Errorf("ParseMovieFilter() genre[%d] = %v, want %v", i, genre, test.wantGenres[i])
				}
			}
		})
	}
}

func TestParseUpdateMovie(t *testing.T) {
	tests := []struct {
		name       string
		formData   map[string]string
		wantErr    bool
		errMessage string
	}{
		{
			name: "success parse update movie - all fields",
			formData: map[string]string{
				"title":            "Updated Movie",
				"description":      "Updated Description",
				"duration_minutes": "150",
				"artists":          "Updated Artist",
				"genres":           "Action, Drama",
			},
			wantErr: false,
		},
		{
			name: "success parse update movie - partial update",
			formData: map[string]string{
				"title":       "Updated Movie",
				"description": "Updated Description",
			},
			wantErr: false,
		},
		{
			name: "fail - invalid duration",
			formData: map[string]string{
				"title":            "Updated Movie",
				"duration_minutes": "abc",
			},
			wantErr:    true,
			errMessage: "duration must be a number",
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

			req, err := http.NewRequest("PUT", "/", body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			parser := NewMovieParser()
			movie, err := parser.ParseUpdateMovie(req)

			if (err != nil) != test.wantErr {
				t.Errorf("ParseUpdateMovie() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if test.wantErr {
				if err == nil {
					t.Error("ParseUpdateMovie() expected error but got nil")
					return
				}
				if err.Error() != test.errMessage {
					t.Errorf("ParseUpdateMovie() error message = %v, want %v", err.Error(), test.errMessage)
				}
				return
			}

			if movie == nil {
				t.Error("ParseUpdateMovie() movie is nil")
				return
			}

			for key, value := range test.formData {
				switch key {
				case "title":
					if movie.Title != value {
						t.Errorf("ParseUpdateMovie() title = %v, want %v", movie.Title, value)
					}
				case "description":
					if movie.Description != value {
						t.Errorf("ParseUpdateMovie() description = %v, want %v", movie.Description, value)
					}
				case "duration_minutes":
					duration, _ := strconv.Atoi(value)
					if movie.Duration != duration {
						t.Errorf("ParseUpdateMovie() duration = %v, want %v", movie.Duration, duration)
					}
				case "artists":
					if movie.Artists != value {
						t.Errorf("ParseUpdateMovie() artists = %v, want %v", movie.Artists, value)
					}
				case "genres":
					if movie.Genres != value {
						t.Errorf("ParseUpdateMovie() genres = %v, want %v", movie.Genres, value)
					}
				}
			}
		})
	}
}
