package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/google/uuid"
)

type CompanyUseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  official_company.OfficialCompanyRepository
}

func (c *CompanyUseCase) GetCompaniesList(start uint, end uint) ([]models.OfficialCompany, error) {
	vacList, err := c.repos.GetCompaniesList(start, end)
	if err != nil {
		err = fmt.Errorf("error in company list creation: %w", err)
		return nil, err
	}
	return vacList, nil
}

func (c *CompanyUseCase) GetMineCompany(empId uuid.UUID) (*models.OfficialCompany, error) {
	comp, err := c.repos.GetMineCompany(empId)
	if err != nil {
		err = fmt.Errorf("error in get by id official company func : %w", err)
		return nil, err
	}
	return comp, nil
}

func NewCompUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos official_company.OfficialCompanyRepository) *CompanyUseCase {
	return &CompanyUseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
	}
}

func (c *CompanyUseCase) CreateOfficialCompany(company models.OfficialCompany, empId uuid.UUID) (*models.OfficialCompany, error) {
	comp, err := c.repos.CreateOfficialCompany(company, empId)
	if err != nil {
		err = fmt.Errorf("error in create official company function: %w", err)
		return nil, err
	}
	return comp, nil
}

func (c *CompanyUseCase) GetOfficialCompany(compId uuid.UUID) (*models.OfficialCompany, error) {
	comp, err := c.repos.GetOfficialCompany(compId)
	if err != nil {
		err = fmt.Errorf("error in get by id official company func : %w", err)
		return nil, err
	}
	return comp, nil
}
