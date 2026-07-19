package domain

type Duration struct {
	minutes int
}

func NewDuration(minutes int) (Duration, error) {
	if minutes <= 0 || minutes > 240 {
		return Duration{}, ErrInvalidDuration
	}
	return Duration{minutes: minutes}, nil
}

func (d Duration) Minutes() int {
	return d.minutes
}
