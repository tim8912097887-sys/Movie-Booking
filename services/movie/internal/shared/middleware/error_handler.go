package middleware

import (
	"net/http"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/application"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/domain"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/infrastructure/db"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/shared/json"
)

func ErrorHandler(w http.ResponseWriter, err error) error {
	var responseErr error
	switch err {
	// Application error
	case application.ErrMovieNotFound:
		responseErr = json.ErrorJson(w, http.StatusNotFound, "movie_not_found", err.Error())
	// Domain error
	case domain.ErrInvalidDuration:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "invalid_duration", err.Error())
	case domain.ErrInvalidRating:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "invalid_rating", err.Error())
	case domain.ErrInvalidGenre:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "invalid_genres", err.Error())
	case domain.ErrReleaseDateAfterDueDate:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "release_date_cannot_be_after_due_date", err.Error())
	// Repository error
	case db.ErrMovieNotFound:
		responseErr = json.ErrorJson(w, http.StatusNotFound, "movie_not_found", err.Error())
	case db.ErrMovieGenreNotFound:
		responseErr = json.ErrorJson(w, http.StatusNotFound, "movie_genre_not_found", err.Error())
	case db.ErrInvalidGenre:
		responseErr = json.ErrorJson(w, http.StatusBadRequest, "invalid_genre", err.Error())
	default:
		responseErr = json.ErrorJson(w, http.StatusInternalServerError, "internal_error", err.Error())
	}
	return responseErr
}