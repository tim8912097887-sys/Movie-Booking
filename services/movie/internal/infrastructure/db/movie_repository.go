package db

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tim8912097887-sys/movie_booking/services/movie/internal/domain"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) Create(ctx context.Context, movie domain.Movie) (domain.Movie, error) {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return domain.Movie{}, err
    }
    defer tx.Rollback(ctx)

    if err := r.insertMovie(ctx, tx, movie); err != nil {
        return domain.Movie{}, err
    }

    if err := r.verifyGenres(ctx, tx, movie.Genres()); err != nil {
        return domain.Movie{}, err
    }

    if err := r.insertMovieGenres(ctx, tx, movie); err != nil {
        return domain.Movie{}, err
    }

    if err := tx.Commit(ctx); err != nil {
        return domain.Movie{}, err
    }

    return r.GetByID(ctx,movie.ID())
}

func (r *MovieRepository) GetByID(ctx context.Context, id string) (domain.Movie, error) {

	query := `
		SELECT m.id, m.title, m.description, ARRAY_AGG(mg.genre_id ORDER BY mg.genre_id) AS genres, m.duration, m.rating,
			m.release_date, m.due_date, m.created_at, m.updated_at, m.deleted_at
		FROM movies AS m
		JOIN movie_genres AS mg ON m.id = mg.movie_id
		WHERE m.id = $1 AND m.deleted_at IS NULL
		GROUP BY m.id, m.title, m.description, m.duration, m.rating,
			m.release_date, m.due_date, m.created_at, m.updated_at, m.deleted_at
	`

	row := r.db.QueryRow(ctx, query, id)
	movie, err := r.scanMovie(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Movie{}, ErrMovieNotFound
		}
		return domain.Movie{}, err
	}

	return movie, nil
}

func (r *MovieRepository) GetAll(ctx context.Context) ([]domain.Movie, error) {

	query := `
		SELECT m.id, m.title, m.description, ARRAY_AGG(mg.genre_id ORDER BY mg.genre_id) AS genres, m.duration, m.rating,
			m.release_date, m.due_date, m.created_at, m.updated_at, m.deleted_at
		FROM movies AS m
		JOIN movie_genres AS mg ON m.id = mg.movie_id
		WHERE m.deleted_at IS NULL
		GROUP BY m.id, m.title, m.description, m.duration, m.rating,
			m.release_date, m.due_date, m.created_at, m.updated_at, m.deleted_at
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	movies := make([]domain.Movie, 0)
	for rows.Next() {
		movie, err := r.scanMovie(rows)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) Update(ctx context.Context, movie domain.Movie) (domain.Movie, error) {

	tx, err := r.db.Begin(ctx)
    if err != nil {
        return domain.Movie{}, err
    }
    defer tx.Rollback(ctx)
	query := `
		UPDATE movies
		SET title = $2,
			description = $3,
			duration = $4,
			rating = $5,
			release_date = $6,
			due_date = $7,
			updated_at = $8
		WHERE id = $1
	`

	row,err := tx.Exec(ctx, query,
		movie.ID(),
		movie.Title(),
		movie.Description(),
		movie.Duration().Minutes(),
		movie.Rating().String(),
		movie.ReleaseDate(),
		movie.DueDate(),
		movie.UpdatedAt(),
	)

	if err != nil {
		return domain.Movie{}, err
	}

	if row.RowsAffected() == 0 {
		return domain.Movie{}, ErrMovieNotFound
	}

	if err := r.verifyGenres(ctx, tx, movie.Genres()); err != nil {
		return domain.Movie{}, err
	}

	query = `
		DELETE FROM movie_genres
		WHERE movie_id = $1
	`

	_, err = tx.Exec(ctx, query, movie.ID())
	if err != nil {
		return domain.Movie{}, err
	}

	if err := r.insertMovieGenres(ctx, tx, movie); err != nil {
		return domain.Movie{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Movie{}, err
	}

	return r.GetByID(ctx,movie.ID())
}

func (r *MovieRepository) Delete(ctx context.Context, id string) error {

	tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

	query := `UPDATE movies SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrMovieNotFound
	}

	query = `DELETE FROM movie_genres WHERE movie_id = $1`
	tag, err = tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrMovieGenreNotFound
	}

	if err := tx.Commit(ctx); err != nil {
        return err
    }

	return nil
}

// Private query functions
func (r *MovieRepository) verifyGenres(
    ctx context.Context,
    tx pgx.Tx,
    genres []domain.Genre,
) error {

    query := `SELECT COUNT(*) FROM genres WHERE name = ANY($1)`

    args := genreStrings(genres)

    var count int
    if err := tx.QueryRow(ctx, query, args).Scan(&count); err != nil {
        return err
    }

    if count != len(genres) {
        return ErrInvalidGenre
    }

    return nil
}

func (r *MovieRepository) insertMovie(
	ctx context.Context,
	tx pgx.Tx,
	movie domain.Movie,
) error {
	query := `
		INSERT INTO movies (id, title, description, duration, rating,
			release_date, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := tx.Exec(ctx, query,
		movie.ID(),
		movie.Title(),
		movie.Description(),
		movie.Duration().Minutes(),
		movie.Rating(),
		movie.ReleaseDate(),
		movie.DueDate(),
		movie.CreatedAt(),
		movie.UpdatedAt(),
	)
	return err
}

func (r *MovieRepository) insertMovieGenres(
	ctx context.Context,
	tx pgx.Tx,
	movie domain.Movie,
) error {
	genres := movie.Genres()
	if len(genres) == 0 {
		return nil
	}

	var (
		query strings.Builder
		args  []any
	)

	query.WriteString(`
		INSERT INTO movie_genres (movie_id, genre_id)
		VALUES
	`)

	for i, genre := range genres {
		if i > 0 {
			query.WriteString(",")
		}

		// ($1,$2), ($3,$4), ...
		query.WriteString(fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))

		args = append(args,
			movie.ID(),
			genre.String(),
		)
	}

	_, err := tx.Exec(ctx, query.String(), args...)
	return err
}

// Utility functions
func (r *MovieRepository) scanMovie(scanner pgxScanner) (domain.Movie, error) {
    var (
        id, title, description, rating string
        genres []string
        duration int
        releaseDate, dueDate, createdAt, updatedAt time.Time
        deletedAt *time.Time
    )

    err := scanner.Scan(
        &id,
        &title,
        &description,
        &genres,        // scan into []string
        &duration,
        &rating,
        &releaseDate,
        &dueDate,
        &createdAt,
        &updatedAt,
        &deletedAt,
    )
    if err != nil {
        return domain.Movie{}, err
    }

    return domain.RehydrateMovie(
        id,
        title,
        description,
        genres,        // pass []string
        duration,
        rating,
        releaseDate,
        dueDate,
        createdAt,
        updatedAt,
        deletedAt,
    )
}


func splitGenres(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.Split(value, ",")
}

func genreStrings(genres []domain.Genre) []string {
	if len(genres) == 0 {
		return nil
	}

	values := make([]string, 0, len(genres))
	for _, genre := range genres {
		values = append(values, genre.String())
	}
	return values
}
