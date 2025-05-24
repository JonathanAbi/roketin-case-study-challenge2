package movie

import (
	"net/http"
	"roketin-case-study-challenge2/internal"
	"roketin-case-study-challenge2/internal/response"
	"roketin-case-study-challenge2/internal/constant"

	"strconv"

	"github.com/go-chi/chi"
)

type MovieHandler struct {
	movieParser MovieParserInterface
	movieFlow   MovieFlowInterface
}

func NewMovieHandler(movieParser MovieParserInterface, movieFlow MovieFlowInterface) *MovieHandler {
	return &MovieHandler{
		movieParser: movieParser,
		movieFlow:   movieFlow,
	}
}

func (h *MovieHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateMovie)
	r.Get("/", h.ListMovies)
	r.Get("/search", h.SearchMovies)
	r.Put("/{id}", h.UpdateMovie)
	r.Delete("/{id}", h.DeleteMovie)

	return r
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	movieData, file, err := h.movieParser.ParseCreateMovie(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	filePath, err := internal.SaveUploadedFile(file, "uploads")
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	movieData.FilePath = filePath

	createdMovie, err := h.movieFlow.CreateMovie(ctx, movieData)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
	}

	response.Success(w, createdMovie)
}

func (h *MovieHandler) ListMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := h.movieParser.ParseMovieFilter(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	movies, total, err := h.movieFlow.ListMovies(ctx, filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := response.Pagination{
		CurrentPage: filter.GetPage(),
		PerPage: filter.GetLimit(),
		TotalItems: total,
		TotalPages: int((total + int64(filter.GetLimit()) - 1) / int64(filter.GetLimit())),
	}

	response.SuccessWithPagination(w, movies, pagination)
}

func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	filter, err := h.movieParser.ParseMovieFilter(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	movies, total, err := h.movieFlow.ListMovies(ctx, filter)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	pagination := response.Pagination{
		CurrentPage: filter.GetPage(),
		PerPage: filter.GetLimit(),
		TotalItems: total,
		TotalPages: int((total + int64(filter.GetLimit()) - 1) / int64(filter.GetLimit())),
	}

	response.SuccessWithPagination(w, movies, pagination)
}

func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, constant.ERROR_INVALID_MOVIE_ID)
		return
	}

	request, err := h.movieParser.ParseUpdateMovie(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	request.ID = id

	updatedMovie, err := h.movieFlow.UpdateMovie(ctx, request)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, updatedMovie)
}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, constant.ERROR_INVALID_MOVIE_ID)
		return
	}

	err = h.movieFlow.DeleteMovie(ctx, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(w, constant.MOVIE_DELETED_SUCCESSFULLY)
}
