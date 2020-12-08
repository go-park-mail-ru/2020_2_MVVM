package repository

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type pgRepository struct {
	db *gorm.DB
}

func NewPgRepository(db *gorm.DB) official_company.OfficialCompanyRepository {
	return &pgRepository{db: db}
}

func (p *pgRepository) UpdateOfficialCompany(newComp models.OfficialCompany, empId uuid.UUID) (*models.OfficialCompany, error) {
	/*if _, err := p.db.Model(&newComp).WherePK().Returning("*").Update(); err != nil {
		return nil, fmt.Errorf("can't update company with id: %s", newComp.ID)
	}
	return &newComp, nil*/
	return nil, nil

}

func (p *pgRepository) DeleteOfficialCompany(compId uuid.UUID, empId uuid.UUID) error {
	/*var (
		comp *models.OfficialCompany
		err  error
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
	}*/
	return nil
}

func (p *pgRepository) SearchCompanies(params models.CompanySearchParams) ([]models.OfficialCompany, error) {
	var compList []models.OfficialCompany

	err := p.db.Table("main.official_companies").Scopes(func(q *gorm.DB) *gorm.DB {
		if params.VacCount > 0 {
			q = q.Where("count_company >= (?)", params.VacCount)
		}
		if len(params.AreaSearch) != 0 {
			q = q.Where("area_search IN (?)", params.AreaSearch)
		}
		if params.KeyWords != "" {
			q = q.Where("LOWER(name) LIKE (?)", "%"+params.KeyWords+"%")
		}
		if len(params.Sphere) != 0 {
			q = q.Where("spheres @> (?)", pq.Array(params.Sphere))
		}
		if params.OrderBy != "" {
			return q.Order(params.OrderBy)
		}
		return q
	}).Find(&compList).Error
	if err != nil {
		return nil, fmt.Errorf("error in companies list selection with searchParams: %s", err)
	}
	if len(compList) == 0 {
		return nil, nil
	}
	return compList, nil
}

func (p *pgRepository) GetCompaniesList(start uint, limit uint) ([]models.OfficialCompany, error) {
	var compList []models.OfficialCompany
	if limit <= start {
		return nil, fmt.Errorf("selection with useless positions")
	}
	err := p.db.Table("main.official_companies").Limit(int(limit)).Offset(int(start)).Order("name").Find(&compList).Error
	if err != nil {
		err = fmt.Errorf("error in list selection from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	if len(compList) == 0 {
		return nil, nil
	}
	return compList, nil
}

func (p *pgRepository) GetMineCompany(empId uuid.UUID) (*models.OfficialCompany, error) {
	var employer models.Employer
	err := p.db.Table("main.employers").Take(&employer, "empl_id = ?", empId).Error
	if err != nil {
		err = fmt.Errorf("error in select employer with id: %s : error: %w", empId, err)
		return nil, err
	}
	return p.GetOfficialCompany(employer.CompanyID)
}

func (p *pgRepository) CreateOfficialCompany(comp models.OfficialCompany, empId uuid.UUID) (*models.OfficialCompany, error) {
	if empId == uuid.Nil {
		return nil, fmt.Errorf("error in inserting official company:empId = nil")
	}
	employer := models.Employer{}
	err := p.db.Table("main.employers").Select("comp_id").Take(&employer, "empl_id = ?", empId).Error
	if err != nil || employer.CompanyID != uuid.Nil {
		return nil, errors.New(common.EmpHaveComp)
	}
	err = p.db.Table("main.official_companies").Create(&comp).Error
	if err != nil {
		return nil, fmt.Errorf("error in inserting official company: error: %w", err)
	}
	err = p.db.Table("main.employers").Where("empl_id = ?", empId).UpdateColumn("comp_id", comp.ID).Error
	if err != nil {
		return nil, fmt.Errorf("error in update employer(add company) with id:  %s : error: %w", empId.String(), err)
	}
	return &comp, nil
}

func (p *pgRepository) GetOfficialCompany(compId uuid.UUID) (*models.OfficialCompany, error) {
	var comp models.OfficialCompany

	if compId == uuid.Nil {
		return nil, nil
	}
	err := p.db.Table("main.official_companies").Take(&comp, "comp_id = ?", compId).Error
	if err != nil {
		return nil, fmt.Errorf("error in select official company with id: %s : error: %w", compId, err)
	}

	return &comp, nil
}

func (p *pgRepository) GetAllCompaniesNames() ([]models.BriefCompany, error) {
	var listComp []models.BriefCompany
	err := p.db.Table("main.official_companies").Find(&listComp).Error
	if err != nil {
		return nil, fmt.Errorf("error in select official company names: error: %w", err)
	}
	return listComp, nil
}
