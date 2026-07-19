package application

import "time"

type CreateMovieCommand struct {
	Title           string
	Description     string
	Genres          []string
	DurationMinutes int
	Rating          string
	ReleaseDate     time.Time
	DueDate         time.Time
}

type CreateMovieResult struct {
    ID              string
    Title           string
    Description     string
    Genres          []string
    DurationMinutes int
    Rating          string
    ReleaseDate     time.Time
    DueDate         time.Time
}

type GetMovieResult struct {
	ID              string
	Title           string
	Description     string
	Genres          []string
	DurationMinutes int
	Rating          string
	ReleaseDate     time.Time
	DueDate         time.Time
}

type GetMoviesResult struct {
	ID              string
	Title           string
	Description     string
	Genres          []string
	DurationMinutes int
	Rating          string
	ReleaseDate     time.Time
	DueDate         time.Time
}

type UpdateMovieCommand struct {
	ID              string
	Title           *string
	Description     *string
	Genres          []string
	DurationMinutes *int
	Rating          *string
	ReleaseDate     *time.Time
	DueDate         *time.Time
}

type UpdateMovieResult struct {
	ID              string
	Title           string
	Description     string
	Genres          []string
	DurationMinutes int
	Rating          string
	ReleaseDate     time.Time
	DueDate         time.Time
}

type DeleteMovieCommand struct {
	ID string
}

type DeleteMovieResult struct {
	ID string
}