package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api/storage"
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
