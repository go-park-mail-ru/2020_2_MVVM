package usecase

import (
	"context"
	"github.com/apsdehal/go-logger"
)

type Usecase struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
}

func NewUsecase(infoLogger, errorLogger *logger.Logger) Usecase {
	return Usecase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}
}

func (u *Usecase) WS(ctx context.Context, lat, lng float64) {

}
