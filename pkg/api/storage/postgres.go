package storage

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-pg/pg/v9"
)

type pgStorage struct {
	db *pg.DB
}

func (s *pgStorage) NothingFunc() error {
	return nil
}

func (s *pgStorage) GetUserByID(id string) (models.User, error) {
	var user models.User
	err := s.db.Model(&user).Where("user_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select user with id: %s : error: %w", id, err)
		return models.User{}, err
	}
	return user, nil
}

func (s *pgStorage) CreateUser(user models.User) (models.User, error) {
	_, err := s.db.Model(&user).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting user with name: %s : error: %w", user.Name, err)
		return models.User{}, err
	}
	return user, nil
}

func NewPostgresStorage(db *pg.DB) Storage {
	return &pgStorage{db: db}
}
