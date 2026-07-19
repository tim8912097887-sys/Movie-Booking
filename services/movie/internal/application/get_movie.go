package application

import (
	"context"
)

type GetMovieUsecase struct {
	repository MovieRepository
}

func NewGetMovieUsecase(repository MovieRepository) *GetMovieUsecase {
	return &GetMovieUsecase{
		repository: repository,
	}
}

func (u *GetMovieUsecase) Execute(ctx context.Context, id string) (GetMovieResult, error) {
	movie, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return GetMovieResult{}, err
	}

	if movie.Title() == "" {
		return GetMovieResult{}, ErrMovieNotFound
	}

	// Convert genres to string slice
	var genresStr []string
	for _, genre := range movie.Genres() {
		genresStr = append(genresStr, genre.String())
	}

	return GetMovieResult{
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
