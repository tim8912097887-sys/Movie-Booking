-- +goose Up
CREATE TYPE movie_rating AS ENUM ('G', 'PG', 'PG13', 'R', 'NC17');
CREATE TABLE movies (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    duration INT NOT NULL,       
    rating movie_rating NOT NULL, 
    release_date TIMESTAMP WITH TIME ZONE NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_movies_rating 
ON movies (rating) 
WHERE deleted_at IS NULL;
-- +goose Down
DROP TABLE movies;
