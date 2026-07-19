package db

import "errors"

var (
	ErrMovieNotFound = errors.New("movie not found")
	ErrMovieGenreNotFound = errors.New("movie genre not found")
	ErrInvalidGenre = errors.New("invalid genre")
)