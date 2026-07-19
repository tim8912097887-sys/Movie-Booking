package domain

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	id          uuid.UUID
	title       string
	description string

	genres   []Genre
	duration Duration
	rating   Rating

	releaseDate time.Time
	dueDate     time.Time

	createdAt time.Time
	updatedAt time.Time
	deletedAt *time.Time
}

func NewMovie(title, description string, genres []string, duration int, rating string, releaseDate time.Time, dueDate time.Time) (Movie, error) {

	var validatedGenres []Genre

	for _, genre := range genres {
		validatedGenre, err := ParseGenre(genre)

		if err != nil {
			return Movie{}, err
		}

		validatedGenres = append(validatedGenres, validatedGenre)
	}

	validatedDuration, err := NewDuration(duration)

	if err != nil {
		return Movie{}, err
	}

	validatedRating, err := ParseRating(rating)

	if err != nil {
		return Movie{}, err
	}

	if releaseDate.After(dueDate) {
		return Movie{}, ErrReleaseDateAfterDueDate
	}
	return Movie{
		id:          uuid.New(),
		title:       title,
		description: description,
		genres:      validatedGenres,
		duration:    validatedDuration,
		rating:      validatedRating,
		releaseDate: releaseDate,
		dueDate:     dueDate,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	},nil
}

// NewMovieWithID is used to reconstruct a movie from database records
func RehydrateMovie(id string, title, description string, genres []string, duration int, rating string, releaseDate time.Time, dueDate time.Time, createdAt time.Time, updatedAt time.Time, deletedAt *time.Time) (Movie, error) {
	
	var validatedGenres []Genre

	for _, genre := range genres {
		validatedGenre, err := ParseGenre(genre)

		if err != nil {
			return Movie{}, err
		}

		validatedGenres = append(validatedGenres, validatedGenre)
	}

	validatedDuration, err := NewDuration(duration)

	if err != nil {
		return Movie{}, err
	}

	validatedRating, err := ParseRating(rating)

	if err != nil {
		return Movie{}, err
	}

	if releaseDate.After(dueDate) {
		return Movie{}, ErrReleaseDateAfterDueDate
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Movie{}, err
	}

	return Movie{
		id:          parsedID,
		title:       title,
		description: description,
		genres:      validatedGenres,
		duration:    validatedDuration,
		rating:      validatedRating,
		releaseDate: releaseDate,
		dueDate:     dueDate,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		deletedAt:   deletedAt,
	}, nil
}
// Getter methods
func (m *Movie) ID() string {
	return m.id.String()
}

func (m *Movie) Title() string {
	return m.title
}

func (m *Movie) Description() string {
	return m.description
}

func (m *Movie) Genres() []Genre {
	return m.genres
}

func (m *Movie) Duration() Duration {
	return m.duration
}

func (m *Movie) Rating() Rating {
	return m.rating
}

func (m *Movie) ReleaseDate() time.Time {
	return m.releaseDate
}

func (m *Movie) DueDate() time.Time {
	return m.dueDate
}

func (m *Movie) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Movie) UpdatedAt() time.Time {
	return m.updatedAt
}

func (m *Movie) DeletedAt() *time.Time {
	return m.deletedAt
}

// Update method for updating movie details
func (m *Movie) Update(title, description string, genres []string, duration int, rating string, releaseDate time.Time, dueDate time.Time) error {
	var validatedGenres []Genre

	for _, genre := range genres {
		validatedGenre, err := ParseGenre(genre)
		if err != nil {
			return err
		}
		validatedGenres = append(validatedGenres, validatedGenre)
	}

	validatedDuration, err := NewDuration(duration)
	if err != nil {
		return err
	}

	validatedRating, err := ParseRating(rating)
	if err != nil {
		return err
	}

	if releaseDate.After(dueDate) {
		return ErrReleaseDateAfterDueDate
	}

	m.title = title
	m.description = description
	m.genres = validatedGenres
	m.duration = validatedDuration
	m.rating = validatedRating
	m.releaseDate = releaseDate
	m.dueDate = dueDate
	m.updatedAt = time.Now()

	return nil
}

// PartialUpdate method for updating only provided movie fields
func (m *Movie) PartialUpdate(title *string, description *string, genres []string, duration *int, rating *string, releaseDate *time.Time, dueDate *time.Time) error {
	// Update title if provided
	if title != nil {
		m.title = *title
	}

	// Update description if provided
	if description != nil {
		m.description = *description
	}

	// Update genres if provided
	if len(genres) > 0 {
		var validatedGenres []Genre
		for _, genre := range genres {
			validatedGenre, err := ParseGenre(genre)
			if err != nil {
				return err
			}
			validatedGenres = append(validatedGenres, validatedGenre)
		}
		m.genres = validatedGenres
	}

	// Update duration if provided
	if duration != nil {
		validatedDuration, err := NewDuration(*duration)
		if err != nil {
			return err
		}
		m.duration = validatedDuration
	}

	// Update rating if provided
	if rating != nil {
		validatedRating, err := ParseRating(*rating)
		if err != nil {
			return err
		}
		m.rating = validatedRating
	}

	// Validate release date and due date together if either is provided
	newReleaseDate := m.releaseDate
	newDueDate := m.dueDate

	if releaseDate != nil {
		newReleaseDate = *releaseDate
	}

	if dueDate != nil {
		newDueDate = *dueDate
	}

	if releaseDate != nil || dueDate != nil {
		if newReleaseDate.After(newDueDate) {
			return ErrReleaseDateAfterDueDate
		}
		m.releaseDate = newReleaseDate
		m.dueDate = newDueDate
	}

	m.updatedAt = time.Now()
	return nil
}