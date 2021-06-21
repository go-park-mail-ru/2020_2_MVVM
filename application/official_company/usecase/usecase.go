package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"strings"
)

type CompanyUseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  official_company.OfficialCompanyRepository
}

func (c *CompanyUseCase) DeleteOfficialCompany(compId uuid.UUID, empId uuid.UUID) error {
	err := c.repos.DeleteOfficialCompany(compId, empId)
	if err != nil {
		err = fmt.Errorf("error in delete official company function: %w", err)
	}
	return err
}

func (c *CompanyUseCase) GetCompaniesList(start uint, limit uint) ([]models.OfficialCompany, error) {
	compList, err := c.repos.GetCompaniesList(start, limit)
	if err != nil {
		err = fmt.Errorf("error in company list creation: %w", err)
		return nil, err
	}
	return compList, nil
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
		if err.Error() != common.EmpHaveComp {
			err = fmt.Errorf("error in create official company function: %w", err)
		}
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

func (c *CompanyUseCase) UpdateOfficialCompany(company models.OfficialCompany, empId uuid.UUID) (*models.OfficialCompany, error) {
	comp, err := c.repos.UpdateOfficialCompany(company, empId)
	if err != nil {
		if err.Error() != common.EmpHaveComp {
			err = fmt.Errorf("error in create official company function: %w", err)
		}
		return nil, err
	}
	return comp, nil
}

func (c *CompanyUseCase) SearchCompanies(params models.CompanySearchParams) ([]models.OfficialCompany, error) {
	if params.OrderBy != "" {
		if params.OrderBy == "count_vacancy" {
			if params.ByAsc {
				params.OrderBy += " ASC"
			} else {
				params.OrderBy += " DESC"
			}
		} else {
			params.OrderBy = ""
		}
	}
	params.KeyWords = strings.ToLower(params.KeyWords)
	vacList, err := c.repos.SearchCompanies(params)
	if err != nil {
		return nil, err
	}
	return vacList, nil
}

func (c *CompanyUseCase) GetAllCompaniesNames() ([]models.BriefCompany, error) {
	return c.repos.GetAllCompaniesNames()
}
