package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	logger "github.com/rowdyroad/go-simple-logger"
)

type UserUseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  user.RepositoryUser
}

func NewUserUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos user.RepositoryUser) *UserUseCase {
	return &UserUseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
	}
}

func (U *UserUseCase) GetUserByID(id string) (models.User, error) {
	userById, err := U.repos.GetUserByID(id)
	if err != nil {
		err = fmt.Errorf("error in user get by id func : %w", err)
		return models.User{}, err
	}
	return userById, nil
}

func (U *UserUseCase) CreateUser(user models.User) (models.User, error) {
	userNew, err := U.repos.CreateUser(user)
	if err != nil {
		if err.Error() != "user already exists" {
			err = fmt.Errorf("error in user get by id func : %w", err)
		}
		return models.User{}, err
	}
	return userNew, nil
}

func (U *UserUseCase) UpdateUser(userNew models.User) (models.User, error) {
	userNew, err := U.repos.UpdateUser(userNew)
	if err != nil {
		if errMsg := err.Error(); errMsg != "user already exists" && errMsg != "nothing to update" {
			err = fmt.Errorf("error in user update by id func : %w", err)
		}
		return models.User{}, err
	}
	return userNew, nil
}