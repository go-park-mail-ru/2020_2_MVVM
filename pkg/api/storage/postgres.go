package storage

import "github.com/go-pg/pg/v9"

type pgStorage struct {
	db *pg.DB
}

func (s *pgStorage) NothingFunc() error {
	return nil
}

func NewPostgresStorage(db *pg.DB) Storage {
	return &pgStorage{db: db}
}
