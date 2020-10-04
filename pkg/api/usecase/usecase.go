package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api/storage"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/models"
	logger "github.com/rowdyroad/go-simple-logger"
)

type Usecase struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        storage.Storage
}

func NewUsecase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	strg storage.Storage) *Usecase {
	usecase := Usecase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		strg:        strg,
	}
	return &usecase
}

func (u *Usecase) DoNothing() error {
	err := u.strg.NothingFunc()
	if err != nil {
		err = fmt.Errorf("error in nothing func : %w", err)
		return err
	}
	return nil
}

func (u *Usecase) GetUserByID(id string) (models.User, error) {
	user, err := u.strg.GetUserByID(id)
	if err != nil {
		err = fmt.Errorf("error in user get by id func : %w", err)
		return models.User{}, err
	}
	return user, nil
}

func (u *Usecase) CreateUser(user models.User) (models.User, error) {
	user, err := u.strg.CreateUser(user)
	if err != nil {
		err = fmt.Errorf("error in user get by id func : %w", err)
		return models.User{}, err
	}
	return user, nil
}
