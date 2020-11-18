package repository

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

type pgRepository struct {
	db *pg.DB
}

func (p *pgRepository) UpdateOfficialCompany(newComp models.OfficialCompany, empId uuid.UUID) (*models.OfficialCompany, error) {
	if _, err := p.db.Model(&newComp).WherePK().Returning("*").Update(); err != nil {
		return nil, fmt.Errorf("can't update company with id: %s", newComp.ID)
	}
	return &newComp, nil

}

func (p *pgRepository) DeleteOfficialCompany(compId uuid.UUID, empId uuid.UUID) error {
	var (
		comp *models.OfficialCompany
		err error
	)
	if comp, err = p.GetMineCompany(empId); err != nil {
		return fmt.Errorf("can't delete employer company with EmpId: %s", empId)
	}

	var employer = models.Employer{ID: empId}
	_, err = p.db.Model(&employer).Column("comp_id").WherePK().Update()
	if err != nil {
		err = fmt.Errorf("error in select employer with id: %s : error: %w", empId, err)
		return err
	}
	if _, err := p.db.Model(comp).WherePK().Delete(); err != nil {
		return fmt.Errorf("can't delete company with id: %s", comp.ID)
	}
	return nil
}

func NewPgRepository(db *pg.DB) official_company.OfficialCompanyRepository {
	return &pgRepository{db: db}
}

func (p *pgRepository) SearchCompanies(params models.CompanySearchParams) ([]models.OfficialCompany, error) {
	var compList []models.OfficialCompany

	err := p.db.Model(&compList).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		if params.VacCount > 0 {
			q = q.Where("count_company >= (?)", params.VacCount)
		}
		if len(params.AreaSearch) != 0 {
			q = q.Where("area_search IN (?)", pg.In(params.AreaSearch))
		}
		if params.KeyWords != "" {
			q = q.Where("LOWER(name) LIKE (?)", "%"+params.KeyWords+"%")
		}
		if len(params.Spheres) != 0 {
			q = q.Where("spheres @> (?)", pg.Array(params.Spheres))
		}
		if params.OrderBy != "" {
			return q.Order(params.OrderBy), nil
		}
		return q, nil
	}).Select()
	if err != nil {
		return nil, fmt.Errorf("error in companies list selection with searchParams: %s", err)
	}
	return compList, nil
}

func (p *pgRepository) GetCompaniesList(start uint, limit uint) ([]models.OfficialCompany, error) {
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

func (p *pgRepository) GetMineCompany(empId uuid.UUID) (*models.OfficialCompany, error) {
	var employer models.Employer
	err := p.db.Model(&employer).Where("empl_id = ?", empId).Select()
	if err != nil {
		err = fmt.Errorf("error in select employer with id: %s : error: %w", empId, err)
		return nil, err
	}
	return p.GetOfficialCompany(employer.CompanyID)
}

func (p *pgRepository) CreateOfficialCompany(company models.OfficialCompany, empId uuid.UUID) (*models.OfficialCompany, error) {
	if empId == uuid.Nil {
		return nil, fmt.Errorf("error in inserting official company:empId = nil")
	}
	employer := models.Employer{ID: empId}
	err := p.db.Model(&employer).WherePK().Select()
	if err != nil || employer.CompanyID != uuid.Nil {
		return nil, errors.New(common.EmpHaveComp)
		//fmt.Errorf("error employer with id = %s doesn't exist or already have company", empId.String())
	}
	//_, err = p.db.Model(&company).WherePK().Drop()
	_, err = p.db.Model(&company).Returning("*").Insert()
	if err != nil {
		return nil, fmt.Errorf("error in inserting official company: error: %w", err)
	}
	employer.CompanyID = company.ID
	_, err = p.db.Model(&employer).WherePK().Column("comp_id").Update()
	if err != nil {
		return nil, fmt.Errorf("error in update employer(add company) with id:  %s : error: %w", empId.String(), err)
	}
	return &company, nil
}

func (p *pgRepository) GetOfficialCompany(compId uuid.UUID) (*models.OfficialCompany, error) {
	var company models.OfficialCompany
	if compId == uuid.Nil {
		return nil, nil
	}
	err := p.db.Model(&company).Where("comp_id = ?", compId).Select()
	if err != nil {
		/*if err.Error() == "pg: no rows in result set" {
			return nil, nil
		}*/
		return nil, fmt.Errorf("error in select official company with id: %s : error: %w", compId, err)
	}
	return &company, nil
}
