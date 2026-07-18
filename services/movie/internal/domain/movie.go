package domain

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	id          uuid.UUID
	title       string
	description string

	genres   []Genre
	duration Duration
	rating   Rating

	releaseDate time.Time
	dueDate     time.Time

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewMovie(title, description string, genres []string, duration int, rating string, releaseDate time.Time, dueDate time.Time) (Movie, error) {
	
	var validatedGenres []Genre

	for _, genre := range genres {
		validatedGenre, err := ParseGenre(genre)

		if err != nil {
			return Movie{}, err
		}

		validatedGenres = append(validatedGenres, validatedGenre)
	}

	validatedDuration, err := NewDuration(duration)

	if err != nil {
		return Movie{}, err
	}

	validatedRating, err := ParseRating(rating)

	if err != nil {
		return Movie{}, err
	}
	return Movie{
		id:          uuid.New(),
		title:       title,
		description: description,
		genres:      validatedGenres,
		duration:    validatedDuration,
		rating:      validatedRating,
		releaseDate: releaseDate,
		dueDate:     dueDate,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	},nil
}