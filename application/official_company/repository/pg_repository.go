package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
)

type pgReopository struct {
	db *pg.DB
}

func (p *pgReopository) GetCompaniesList(start uint, end uint) ([]models.OfficialCompany, error) {
	var compList []models.OfficialCompany
	if end <= start {
		return nil, fmt.Errorf("selection with useless positions")
	}
	err := p.db.Model(&compList).Limit(int(end)).Offset(int(start)).Select()
	if err != nil {
		err = fmt.Errorf("error in list selection from %v to %v: error: %w", start, end, err)
		return nil, err
	}
	return compList, nil
}

func (p *pgReopository) GetMineCompany(empId uuid.UUID) (*models.OfficialCompany, error) {
	var employer models.Employer
	err := p.db.Model(&employer).Where("empl_id = ?", empId).Select()
	if err != nil {
		err = fmt.Errorf("error in select employer with id: %s : error: %w", empId, err)
		return nil, err
	}
	return p.GetOfficialCompany(employer.CompanyID)
}

func NewPgRepository(db *pg.DB) official_company.OfficialCompanyRepository{
	return &pgReopository{db: db}
}

func (p *pgReopository) CreateOfficialCompany(company models.OfficialCompany, empId uuid.UUID) (*models.OfficialCompany, error) {
	if empId == uuid.Nil {
		err := fmt.Errorf("error in inserting official company:empId = nil")
		return nil, err
	}
	employer := models.Employer{ID: empId}
	err := p.db.Model(&employer).WherePK().Select()
	if err != nil || employer.CompanyID != uuid.Nil {
		err = fmt.Errorf("error employer with id = %d doesn't exist or already have company", empId)
	}
	_, err = p.db.Model(&company).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting official company: error: %w", err)
		return nil, err
	}
	employer.CompanyID = company.ID
	_, err = p.db.Model(&employer).WherePK().UpdateNotZero()
	if err != nil {
		err = fmt.Errorf("error in update employer(add company) with id:  %s : error: %w", empId, err)
		return nil, err
	}
	return &company, nil
}

func (p *pgReopository) GetOfficialCompany(compId uuid.UUID) (*models.OfficialCompany, error) {
	var company models.OfficialCompany
	if compId == uuid.Nil {
		return nil, nil
	}
	err := p.db.Model(&company).Where("comp_id = ?", compId).Select()
	if err != nil {
		/*if err.Error() == "pg: no rows in result set" {
			return nil, nil
		}*/
		err = fmt.Errorf("error in select official company with id: %s : error: %w", compId, err)
		return nil, err
	}
	return &company, nil
}
