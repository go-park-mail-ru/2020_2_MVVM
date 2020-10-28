package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
)

type CompanyUseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  official_company.OfficialCompanyRepository
}

func NewCompUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos official_company.OfficialCompanyRepository) *CompanyUseCase {
	return &CompanyUseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
	}
}

func (u *CompanyUseCase) CreateOfficialCompany(company models.OfficialCompany) (*models.OfficialCompany, error) {
	comp, err := u.repos.CreateOfficialCompany(company)
	if err != nil {
		err = fmt.Errorf("error in create official company function: %w", err)
		return nil, err
	}
	return comp, nil
}

func (u *CompanyUseCase) GetOfficialCompany(id string) (*models.OfficialCompany, error) {
	comp, err := u.repos.GetOfficialCompanyById(id)
	if err != nil {
		err = fmt.Errorf("error in get by id official company func : %w", err)
		return nil, err
	}
	return comp, nil
}
