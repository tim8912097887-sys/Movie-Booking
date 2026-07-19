package application

import (
	"context"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/domain"
)

type CreateMovieUsecase struct {
	repository MovieRepository
}

func NewCreateMovieUsecase(repository MovieRepository) *CreateMovieUsecase {
	return &CreateMovieUsecase{
		repository: repository,
	}
}

func (u *CreateMovieUsecase) Execute(ctx context.Context, command CreateMovieCommand) (CreateMovieResult, error) {
	movie, err := domain.NewMovie(command.Title, command.Description, command.Genres, command.DurationMinutes, command.Rating, command.ReleaseDate, command.DueDate)
	if err != nil {
		return CreateMovieResult{}, err
	}

	createdMovie, err := u.repository.Create(ctx, movie)
	if err != nil {
		return CreateMovieResult{}, err
	}

	// Convert genres to string slice
	var genresStr []string
	for _, genre := range createdMovie.Genres() {
		genresStr = append(genresStr, genre.String())
	}

	return CreateMovieResult{
		ID:              createdMovie.ID(),
		Title:           createdMovie.Title(),
		Description:     createdMovie.Description(),
		Genres:          genresStr,
		DurationMinutes: createdMovie.Duration().Minutes(),
		Rating:          createdMovie.Rating().String(),
		ReleaseDate:     createdMovie.ReleaseDate(),
		DueDate:         createdMovie.DueDate(),
	}, nil
}