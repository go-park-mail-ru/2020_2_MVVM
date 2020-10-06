package repository

import (
	"errors"
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
	var newUser models.User
	err := P.db.Model(&newUser).Where("user_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select user with id: %s : error: %w", id, err)
		return models.User{}, err
	}
	return newUser, nil
}

func (P *pgStorage) CreateUser(user models.User) (models.User, error) {
	_, errInsert := P.db.Model(&user).Returning("*").Insert()
	if errInsert != nil {
		if isExist, err := P.db.Model(&user).Exists(); err != nil {
			errInsert = fmt.Errorf("error in inserting user with name: %s : error: %w", user.Name, err)
		} else if isExist {
			errInsert = errors.New("user already exists")
		}
		return models.User{}, errInsert
	}
	return user, nil
}

func (P *pgStorage) UpdateUser(newUser models.User) (models.User, error) {
	oldUser, err := P.GetUserByID(newUser.ID.String())
	if err != nil {
		return models.User{}, err
	}
	switch {
	case newUser.Nickname != "":
		oldUser.Nickname = newUser.Nickname
	case newUser.Email != "":
		oldUser.Email = newUser.Nickname
	case newUser.Surname != "":
		oldUser.Surname = newUser.Surname
	case newUser.Name != "":
		oldUser.Name = newUser.Name
	case newUser.PasswordHash != nil:
		oldUser.PasswordHash = newUser.PasswordHash
	}
	_, errUpdate := P.db.Model(&oldUser).WherePK().Update()
	if errUpdate != nil {
		if isExist, err := P.db.Model(&oldUser).Exists(); err != nil {
			errUpdate = fmt.Errorf("error in update user with name: %s : error: %w", oldUser.Name, err)
		} else if isExist {
			errUpdate = errors.New("user already exists")
		}
		return models.User{}, errUpdate
	}
	return newUser, nil
}
