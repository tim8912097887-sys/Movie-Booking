package domain

import "errors"

var (
	ErrInvalidRating = errors.New("invalid rating")
	ErrInvalidGenre = errors.New("invalid genres")
	ErrInvalidDuration = errors.New("invalid duration")
)