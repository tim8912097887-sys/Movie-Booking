package db

type pgxScanner interface {
	Scan(dest ...any) error
}