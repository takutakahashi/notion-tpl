package store

import (
	"time"
)

type Store struct {
	path string
}

func New() Store {
	return Store{}
}

func (s Store) LastUpdated() time.Time {
	return time.Now()
}
