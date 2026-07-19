package application

import (
	"context"
)

type UpdateMovieUsecase struct {
	repository MovieRepository
}

func NewUpdateMovieUsecase(repository MovieRepository) *UpdateMovieUsecase {
	return &UpdateMovieUsecase{
		repository: repository,
	}
}

func (u *UpdateMovieUsecase) Execute(ctx context.Context, command UpdateMovieCommand) (UpdateMovieResult, error) {
	// Get the existing movie
	movie, err := u.repository.GetByID(ctx, command.ID)
	if err != nil {
		return UpdateMovieResult{}, err
	}

	if movie.Title() == "" {
		return UpdateMovieResult{}, ErrMovieNotFound
	}

	// Perform partial update with only provided fields
	err = movie.PartialUpdate(command.Title, command.Description, command.Genres, command.DurationMinutes, command.Rating, command.ReleaseDate, command.DueDate)
	if err != nil {
		return UpdateMovieResult{}, err
	}

	// Save the updated movie
	err = u.repository.Update(ctx, movie)
	if err != nil {
		return UpdateMovieResult{}, err
	}

	// Convert genres to string slice
	var genresStr []string
	for _, genre := range movie.Genres() {
		genresStr = append(genresStr, genre.String())
	}

	return UpdateMovieResult{
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
