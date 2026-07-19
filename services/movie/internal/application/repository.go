package application

import (
	"context"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/domain"
)

type MovieRepository interface {
	Create(ctx context.Context, movie domain.Movie) error
	GetByID(ctx context.Context, id string) (domain.Movie, error)
	GetAll(ctx context.Context) ([]domain.Movie, error)
	Update(ctx context.Context, movie domain.Movie) error
	Delete(ctx context.Context, id string) error
}
