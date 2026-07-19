package domain

type Genre string

const (
	GenreAction         Genre = "ACTION"
	GenreComedy         Genre = "COMEDY"
	GenreDrama          Genre = "DRAMA"
	GenreThriller       Genre = "THRILLER"
	GenreScienceFiction Genre = "SCIENCE_FICTION"
	GenreAdventure      Genre = "ADVENTURE"
	GenreRomance        Genre = "ROMANCE"
	GenreMystery        Genre = "MYSTERY"
	GenreAnimation      Genre = "ANIMATION"
)

var validateGenres = map[Genre]struct{}{
	GenreAction:         {},
	GenreComedy:         {},
	GenreDrama:          {},
	GenreThriller:       {},
	GenreScienceFiction: {},
	GenreAdventure:      {},
	GenreRomance:        {},
	GenreMystery:        {},
	GenreAnimation:      {},
}

func ParseGenre(g string) (Genre, error) {
	genre := Genre(g)

	if _, ok := validateGenres[genre]; !ok {
		return genre, ErrInvalidGenre
	}
	return genre, nil
}

func (g Genre) String() string {
	return string(g)
}