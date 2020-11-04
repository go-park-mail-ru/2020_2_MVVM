package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

type pgReopository struct {
	db *pg.DB
}

func (p *pgReopository) SearchCompanies(params models.CompanySearchParams) ([]models.OfficialCompany, error) {
	var compList []models.OfficialCompany

	err := p.db.Model(&compList).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		if params.VacCount > 0 {
			q = q.Where("count_vacancy >= ?", params.VacCount)
		}
		if len(params.Spheres) != 0 {
			q = q.Where("spheres IN (?)", pg.In(params.Spheres))
		}
		if len(params.Location) != 0 {
			q = q.Where("location IN (?)", pg.In(params.Location))
		}
		if params.KeyWords != "" {
			q = q.Where("LOWER(name) LIKE ?", fmt.Sprintf("%%%s%", params.KeyWords))
		}
		if params.OrderBy != "" {
			return q.Order(params.OrderBy), nil
		}
		return q, nil
	}).Select()
	if err != nil {
		err = fmt.Errorf("error in companies list selection with searchParams: %s", err)
		return nil, err
	}
	return compList, nil
}

func (p *pgReopository) GetCompaniesList(start uint, limit uint) ([]models.OfficialCompany, error) {
	var compList []models.OfficialCompany
	if limit <= start {
		return nil, fmt.Errorf("selection with useless positions")
	}
	err := p.db.Model(&compList).Limit(int(limit)).Offset(int(start)).Select()
	if err != nil {
		err = fmt.Errorf("error in list selection from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	return compList, nil
}

func (p *pgReopository) GetMineCompany(empId string) (*models.OfficialCompany, error) {
	var employer models.Employer
	err := p.db.Model(&employer).Where("empl_id = ?", empId).Select()
	if err != nil {
		err = fmt.Errorf("error in select employer with id: %s : error: %w", empId, err)
		return nil, err
	}
	return p.GetOfficialCompany(employer.CompanyID.String())
}

func NewPgRepository(db *pg.DB) official_company.OfficialCompanyRepository{
	return &pgReopository{db: db}
}

func (p *pgReopository) CreateOfficialCompany(company models.OfficialCompany, empId string) (*models.OfficialCompany, error) {
	if empId == uuid.Nil.String() {
		err := fmt.Errorf("error in inserting official company:empId = nil")
		return nil, err
	}
	//TODO: fix uuid
	s, _ := uuid.Parse(empId)
	employer := models.Employer{ID: s}
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
	_, err = p.db.Model(&employer).WherePK().Column("comp_id").Update()
	if err != nil {
		err = fmt.Errorf("error in update employer(add company) with id:  %s : error: %w", empId, err)
		return nil, err
	}
	return &company, nil
}

func (p *pgReopository) GetOfficialCompany(compId string) (*models.OfficialCompany, error) {
	var company models.OfficialCompany
	if compId == uuid.Nil.String() {
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
