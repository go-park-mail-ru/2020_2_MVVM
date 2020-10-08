package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/go-pg/pg/v9"
)

func NewPgRepository(db *pg.DB) user.RepositoryUser {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}

func (P *pgStorage) GetUserByID(id string) (models.User, error) {
	var user models.User
	err := P.db.Model(&user).Where("user_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select user with id: %s : error: %w", id, err)
		return models.User{}, err
	}
	return user, nil
}

func (P *pgStorage) CreateUser(user models.User) (models.User, error) {
	_, err := P.db.Model(&user).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting user with name: %s : error: %w", user.Name, err)
		return models.User{}, err
	}
	return user, nil
}

func (P *pgStorage) UpdateUser(userNew models.User) (models.User, error) {
	_, err := P.db.Model(&userNew).WherePK().Returning("*").Update()
	if err != nil {
		err = fmt.Errorf("error in updating user with id %s, : %w", userNew.ID.String(), err)
		return models.User{}, err
	}
	return userNew, nil
}
