package repository

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/go-pg/pg/v9"
	"golang.org/x/crypto/bcrypt"
)

func NewPgRepository(db *pg.DB) user.RepositoryUser {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}

func (P *pgStorage) Login(user models.UserLogin) (*models.User, error) {
	var userDB models.User
	err := P.db.Model(&userDB).
		Where("email = ?", user.Email).
		Where("nickname = ?", user.Nickname).
		Select()
	if err != nil {
		return nil, err
	}
	// compare password with the hashed one
	err = bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password))
	if err != nil {
		return nil, err
	}
	return &userDB, nil
}

func (P *pgStorage) GetUserByID(id string) (*models.User, error) {
	var newUser models.User
	err := P.db.Model(&newUser).Where("cand_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select user with id: %s : error: %w", id, err)
		return nil, err
	}
	return &newUser, nil
}

func (P *pgStorage) CreateUser(user models.User) (*models.User, error) {
	_, errInsert := P.db.Model(&user).Returning("*").Insert()
	if errInsert != nil {
		if isExist, err := P.db.Model(&user).Exists(); err != nil {
			errInsert = fmt.Errorf("error in inserting user with name: %s : error: %w", user.Name, err)
		} else if isExist {
			errInsert = errors.New("user already exists")
		}
		return nil, errInsert
	}
	return &user, nil
}

/*
func (P *pgStorage) UpdateUser(newUser models.User) (models.User, error) {
	oldUser, err := P.GetUserByID(newUser.ID.String())
	if err != nil {
		return models.User{}, err
	}
	switch {
	case newUser.Nickname != "":
		oldUser.Nickname = newUser.Nickname
		fallthrough
	case newUser.Email != "":
		oldUser.Email = newUser.Nickname
		fallthrough
	case newUser.Surname != "":
		oldUser.Surname = newUser.Surname
		fallthrough
	case newUser.Name != "":
		oldUser.Name = newUser.Name
		fallthrough
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
*/
func (P *pgStorage) UpdateUser(userNew models.User) (*models.User, error) {
	_, err := P.db.Model(&userNew).WherePK().Returning("*").Update()
	if err != nil {
		err = fmt.Errorf("error in updating user with id %s, : %w", userNew.ID.String(), err)
		return nil, err
	}
	return &userNew, nil
}
