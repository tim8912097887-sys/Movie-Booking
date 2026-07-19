-- +goose Up
CREATE TABLE movie_genres (
    movie_id UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    genre_id VARCHAR(50) NOT NULL REFERENCES genres(name) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (movie_id, genre_id) 
);
CREATE INDEX idx_movie_genres_genre_id_movie_id 
ON movie_genres (genre_id, movie_id);
-- +goose Down
DROP TABLE movie_genres;
