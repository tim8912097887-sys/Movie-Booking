package domain

type Rating string

const (
	G    Rating = "G"
	PG   Rating = "PG"
	PG13 Rating = "PG13"
	R    Rating = "R"
	NC17 Rating = "NC17"
)

var validRatings = map[Rating]struct{}{
	G:    {},
	PG:   {},
	PG13: {},
	R:    {},
	NC17: {},
}

func ParseRating(s string) (Rating, error) {
	rating := Rating(s)

	if _, ok := validRatings[rating]; !ok {
		return rating, ErrInvalidRating
	}
	return rating, nil
}

func (r Rating) String() string {
	return string(r)
}