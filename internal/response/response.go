package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseWithPagination struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

func Success(w http.ResponseWriter, data interface{}) {
	response := Response{
		Status: "success",
		Data:   data,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func SuccessWithPagination(w http.ResponseWriter, data interface{}, pagination Pagination) {
	response := ResponseWithPagination{
		Data:       data,
		Pagination: pagination,
	}

	Success(w, response)
}

func Error(w http.ResponseWriter, status int, message string) {
	response := Response{
		Status:  "error",
		Message: message,
	}

	respondWithJSON(w, status, response)
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
