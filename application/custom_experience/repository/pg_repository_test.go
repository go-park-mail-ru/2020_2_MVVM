package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func beforeTest(t *testing.T) (custom_experience.CustomExperienceRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	pg := postgres.Dialector{Config: &postgres.Config{Conn: db}}

	conn, err := gorm.Open(pg, &gorm.Config{
		FullSaveAssociations: false,
	})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	mock.MatchExpectationsInOrder(true)
	return NewPgRepository(conn), mock
}

func TestNewPgRepository(t *testing.T) {
	res := NewPgRepository(nil)
	assert.Equal(t, res, &pgRepository{db: nil})
}

func TestCreate(t *testing.T) {
	repo, _ := beforeTest(t)
	res, err := repo.Create(models.ExperienceCustomComp{})
	assert.Nil(t, res)
	assert.Nil(t, err)
}

func TestDropAllFromResume(t *testing.T) {
	repo, _ := beforeTest(t)
	id := uuid.New()
	err := repo.DropAllFromResume(id)
	assert.Error(t, err)
}