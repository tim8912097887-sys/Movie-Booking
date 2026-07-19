package http

import (
	"context"
	"time"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/application"
)

type CreateMovieUsecase interface {
	Execute(ctx context.Context, movie application.CreateMovieCommand) (application.CreateMovieResult, error)
}

type GetMovieUsecase interface {
	Execute(ctx context.Context, id string) (application.GetMovieResult, error)
}

type GetMoviesUsecase interface {
	Execute(ctx context.Context) ([]application.GetMoviesResult, error)
}

type UpdateMovieUsecase interface {
	Execute(ctx context.Context, command application.UpdateMovieCommand) (application.UpdateMovieResult, error)
}

type DeleteMovieUsecase interface {
	Execute(ctx context.Context, command application.DeleteMovieCommand) (application.DeleteMovieResult, error)
}

type CreateMovieRequest struct {
	Title           string    `json:"title" validate:"required,min=1,max=255"`
	Description     string    `json:"description" validate:"required"`
	Genres          []string  `json:"genres" validate:"required,min=1,dive,required"`
	DurationMinutes int       `json:"duration_minutes" validate:"required,min=1,max=240"`
	Rating          string    `json:"rating" validate:"required"`
	ReleaseDate     time.Time `json:"release_date" validate:"required"`
	DueDate         time.Time `json:"due_date" validate:"required"`
}

type CreateMovieResponse struct {
    ID              string    `json:"id"`
    Title           string    `json:"title"`
    Description     string    `json:"description"`
    Genres          []string  `json:"genres"`
    DurationMinutes int       `json:"duration_minutes"`
    Rating          string    `json:"rating"`
    ReleaseDate     time.Time `json:"release_date"`
    DueDate         time.Time `json:"due_date"`
}

type GetMovieResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Genres          []string  `json:"genres"`
	DurationMinutes int       `json:"duration_minutes"`
	Rating          string    `json:"rating"`
	ReleaseDate     time.Time `json:"release_date"`
	DueDate         time.Time `json:"due_date"`
}

type GetMoviesResponse struct {
	Movies []GetMovieResponse `json:"movies"`
}

type UUIDRequest struct {
	ID string `json:"id" validate:"required,uuid4"`
}

type UpdateMovieRequest struct {
	Title           *string    `json:"title" validate:"omitempty,min=1,max=255"`
	Description     *string    `json:"description" validate:"omitempty"`
	Genres          []string   `json:"genres" validate:"omitempty,min=1,dive,required"`
	DurationMinutes *int       `json:"duration_minutes" validate:"omitempty,min=1,max=240"`
	Rating          *string    `json:"rating" validate:"omitempty"`
	ReleaseDate     *time.Time `json:"release_date" validate:"omitempty"`
	DueDate         *time.Time `json:"due_date" validate:"omitempty"`
}

type UpdateMovieResponse struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Genres          []string  `json:"genres"`
	DurationMinutes int       `json:"duration_minutes"`
	Rating          string    `json:"rating"`
	ReleaseDate     time.Time `json:"release_date"`
	DueDate         time.Time `json:"due_date"`
}