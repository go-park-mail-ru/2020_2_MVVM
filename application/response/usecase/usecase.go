package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
)

type UseCaseResponse struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        response.ResponseRepository
}

func NewUsecase(infoLogger *logger.Logger,
				errorLogger *logger.Logger,
				strg response.ResponseRepository) *UseCaseResponse {
					usecase := UseCaseResponse {
					infoLogger:  infoLogger,
					errorLogger: errorLogger,
					strg:        strg,
	}
	return &usecase
}

func (u* UseCaseResponse) CreateResponse(response models.Response) (*models.Response, error) {
	return u.strg.CreateResponse(response)
}

func (u* UseCaseResponse) UpdateStatus(response models.Response) (*models.Response, error) {
	if response.Status == "sent" {
		return nil, nil
	}
	return u.strg.UpdateStatus(response)
}
