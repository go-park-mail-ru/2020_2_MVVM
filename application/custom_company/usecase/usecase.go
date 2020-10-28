package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type UseCase struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        custom_company.CustomCompanyRepository
}

func NewUseCase(infoLogger *logger.Logger,
				errorLogger *logger.Logger,
				strg custom_company.CustomCompanyRepository) *UseCase {
					usecase := UseCase {
					infoLogger:  infoLogger,
					errorLogger: errorLogger,
					strg:        strg,
	}
	return &usecase
}

func (u *UseCase) CreateCustomCompany(company models.CustomCompany) (*models.CustomCompany, error) {
	ed, err := u.strg.CreateCustomCompany(company)
	if err != nil {
		err = fmt.Errorf("error in create custom company function: %w", err)
		return nil, err
	}
	return ed, nil
}

func (u *UseCase) GetCustomCompany(id string) (*models.CustomCompany, error) {
	ed, err := u.strg.GetCustomCompanyById(id)
	if err != nil {
		err = fmt.Errorf("error in get by id custom company func : %w", err)
		return nil, err
	}
	return ed, nil
}
