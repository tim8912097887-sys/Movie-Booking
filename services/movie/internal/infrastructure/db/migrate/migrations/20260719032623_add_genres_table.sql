-- +goose Up
CREATE TABLE genres (
    name VARCHAR(50) PRIMARY KEY
);
INSERT INTO genres (name) VALUES
    ('ACTION'),
    ('COMEDY'),
    ('DRAMA'),
    ('THRILLER'),
    ('SCIENCE_FICTION'),
    ('ADVENTURE'),
    ('ROMANCE'),
    ('MYSTERY'),
    ('ANIMATION');
-- +goose Down
DROP TABLE genres;
