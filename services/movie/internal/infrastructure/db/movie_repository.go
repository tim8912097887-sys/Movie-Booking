package db

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/domain"
)

type MovieRepository struct {
	movies map[string]domain.Movie
	mu     sync.RWMutex
}

func NewMovieRepository() *MovieRepository {
	repo := &MovieRepository{
		movies: make(map[string]domain.Movie),
	}
	repo.seedFakeData()
	return repo
}

// seedFakeData populates the in-memory repository with fake movie data
func (r *MovieRepository) seedFakeData() {
	fakeMovies := []struct {
		id          string
		title       string
		description string
		genres      []string
		duration    int
		rating      string
		releaseDate time.Time
		dueDate     time.Time
	}{
		{
			id:          "550e8400-e29b-41d4-a716-446655440001",
			title:       "The Shawshank Redemption",
			description: "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.",
			genres:      []string{"DRAMA"},
			duration:    142,
			rating:      "R",
			releaseDate: time.Date(1994, 10, 14, 0, 0, 0, 0, time.UTC),
			dueDate:     time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			id:          "550e8400-e29b-41d4-a716-446655440002",
			title:       "The Godfather",
			description: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant youngest son.",
			genres:      []string{"DRAMA", "THRILLER"},
			duration:    175,
			rating:      "R",
			releaseDate: time.Date(1972, 3, 24, 0, 0, 0, 0, time.UTC),
			dueDate:     time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			id:          "550e8400-e29b-41d4-a716-446655440003",
			title:       "The Dark Knight",
			description: "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological tests.",
			genres:      []string{"ACTION", "THRILLER"},
			duration:    152,
			rating:      "PG13",
			releaseDate: time.Date(2008, 7, 18, 0, 0, 0, 0, time.UTC),
			dueDate:     time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			id:          "550e8400-e29b-41d4-a716-446655440004",
			title:       "Forrest Gump",
			description: "The presidencies of Kennedy and Johnson unfold from the perspective of an Alabama man with an IQ of 75.",
			genres:      []string{"DRAMA", "ROMANCE"},
			duration:    142,
			rating:      "PG",
			releaseDate: time.Date(1994, 7, 6, 0, 0, 0, 0, time.UTC),
			dueDate:     time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			id:          "550e8400-e29b-41d4-a716-446655440005",
			title:       "Inception",
			description: "A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea.",
			genres:      []string{"ACTION", "SCIENCE_FICTION", "THRILLER"},
			duration:    148,
			rating:      "PG13",
			releaseDate: time.Date(2010, 7, 16, 0, 0, 0, 0, time.UTC),
			dueDate:     time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, fm := range fakeMovies {
		movie, err := domain.RehydrateMovie(fm.id, fm.title, fm.description, fm.genres, fm.duration, fm.rating, fm.releaseDate, fm.dueDate, time.Now(), time.Now(), nil)
		if err == nil {
			r.movies[fm.id] = movie
		}
	}
}

func (r *MovieRepository) Create(ctx context.Context, movie domain.Movie) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.movies[movie.ID()] = movie
	return nil
}

func (r *MovieRepository) GetByID(ctx context.Context, id string) (domain.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	movie, exists := r.movies[id]
	if !exists {
		return domain.Movie{}, nil
	}

	return movie, nil
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]domain.Movie, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	movies := make([]domain.Movie, 0, len(r.movies))
	for _, movie := range r.movies {
		movies = append(movies, movie)
	}

	return movies, nil
}

func (r *MovieRepository) Update(ctx context.Context, movie domain.Movie) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.movies[movie.ID()]; !exists {
		return errors.New("not found")
	}

	r.movies[movie.ID()] = movie
	return nil
}

func (r *MovieRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.movies[id]; !exists {
		return errors.New("not found")
	}

	delete(r.movies, id)
	return nil
}
