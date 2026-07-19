package application

import (
	"context"
)

type GetMoviesUsecase struct {
	repository MovieRepository
}

func NewGetMoviesUsecase(repository MovieRepository) *GetMoviesUsecase {
	return &GetMoviesUsecase{
		repository: repository,
	}
}

func (u *GetMoviesUsecase) Execute(ctx context.Context) ([]GetMoviesResult, error) {
	movies, err := u.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var results []GetMoviesResult
	for _, movie := range movies {
		// Convert genres to string slice
		var genresStr []string
		for _, genre := range movie.Genres() {
			genresStr = append(genresStr, genre.String())
		}

		result := GetMoviesResult{
			ID:              movie.ID(),
			Title:           movie.Title(),
			Description:     movie.Description(),
			Genres:          genresStr,
			DurationMinutes: movie.Duration().Minutes(),
			Rating:          movie.Rating().String(),
			ReleaseDate:     movie.ReleaseDate(),
			DueDate:         movie.DueDate(),
		}
		results = append(results, result)
	}

	return results, nil
}
