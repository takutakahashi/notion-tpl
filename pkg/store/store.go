package store

import (
	"time"
)

type Store struct {
	path string
}

func setTime(path string, t time.Time) error {
	return nil
}

func getTime(path string) (time.Time, error) {
	return time.Now(), nil
}

func New(path string) Store {
	return Store{path: path}
}

func (s Store) LastUpdated() time.Time {
	return time.Now().AddDate(0, 0, -10)
}
