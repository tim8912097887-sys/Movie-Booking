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

	err = u.repository.Create(ctx, movie)
	if err != nil {
		return CreateMovieResult{}, err
	}

	// Convert genres to string slice
	var genresStr []string
	for _, genre := range movie.Genres() {
		genresStr = append(genresStr, genre.String())
	}

	return CreateMovieResult{
		ID:              movie.ID(),
		Title:           movie.Title(),
		Description:     movie.Description(),
		Genres:          genresStr,
		DurationMinutes: movie.Duration().Minutes(),
		Rating:          movie.Rating().String(),
		ReleaseDate:     movie.ReleaseDate(),
		DueDate:         movie.DueDate(),
	}, nil
}