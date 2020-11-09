package repository

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"log"
	"testing"
)

func SetupDB() (sqlmock.Sqlmock, *pg.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("cant create mock: %s", err)
	}
	repo := &pgStorage{
		db: db,
	}

	db.
	db1 := pg.Connect(db)

	DB, erro := pg.Conn("postgres", db)
	if erro != nil {
		log.Fatalf("Got an unexpected error: %s", err)

	}
	return mock, DB
}

func TestPostgresForSerials_GetSeriesByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repo := &pgStorage{
		db: db,
	}



	mock, DB := SetupDB()
	defer DB.Close()
}