package application

import (
	"context"
)

type DeleteMovieUsecase struct {
	repository MovieRepository
}

func NewDeleteMovieUsecase(repository MovieRepository) *DeleteMovieUsecase {
	return &DeleteMovieUsecase{
		repository: repository,
	}
}

func (u *DeleteMovieUsecase) Execute(ctx context.Context, command DeleteMovieCommand) (DeleteMovieResult, error) {
	// Check if movie exists
	existMovie, err := u.repository.GetByID(ctx, command.ID)
	if err != nil {
		return DeleteMovieResult{}, err
	}

	if existMovie.Title() == "" {
		return DeleteMovieResult{}, ErrMovieNotFound
	}

	// Delete the movie
	err = u.repository.Delete(ctx, command.ID)
	if err != nil {
		return DeleteMovieResult{}, err
	}

	return DeleteMovieResult{
		ID: command.ID,
	}, nil
}
