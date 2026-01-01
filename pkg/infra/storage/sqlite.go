package storage

import (
	"github.com/PhilippReinke/kino-berlin/pkg/domain"
)

type SQLite struct {
}

var _ domain.Storage = &SQLite{}

func NewSQLite() *SQLite {
	return &SQLite{}
}

func (s *SQLite) Upsert(screening domain.Screening) error {
	panic("unimplemented")

	return nil
}

func (s *SQLite) Fetch(filter ...domain.Filter) ([]domain.Screening, error) {
	panic("unimplemented")

	return []domain.Screening{}, nil
}
