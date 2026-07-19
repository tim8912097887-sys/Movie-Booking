package middleware

import (
	"net/http"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/application"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/domain"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/json"
)

func ErrorHandler(w http.ResponseWriter, err error) error {
	var responseErr error
	switch err {
	case application.ErrMovieNotFound:
		responseErr = json.ErrorJson(w, http.StatusNotFound, "movie_not_found", err.Error())
	case domain.ErrInvalidDuration:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "invalid_duration", err.Error())
	case domain.ErrInvalidRating:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "invalid_rating", err.Error())
	case domain.ErrInvalidGenre:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "invalid_genres", err.Error())
	case domain.ErrReleaseDateAfterDueDate:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "release_date_cannot_be_after_due_date", err.Error())
	default:
		responseErr = json.ErrorJson(w, http.StatusInternalServerError, "internal_error", err.Error())
	}
	return responseErr
}