package repository

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/uuid"
)

type pgRepository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) vacancy.RepositoryVacancy {
	return &pgRepository{db: db}
}

func (p *pgRepository) CreateVacancy(vac models.Vacancy) (*models.Vacancy, error) {
	var (
		employer models.Employer
		company  models.OfficialCompany
	)
	err := p.db.Model(&employer).Where("empl_id = ?", vac.EmpID).Select()
	if err != nil {
		err = fmt.Errorf("error in FK search for vacancy creation for user with id: %s : error: %w", vac.EmpID, err)
		return nil, err
	}
	if compId := employer.CompanyID; compId != uuid.Nil {
		vac.CompID = compId
		company.ID = compId
	} else {
		return nil, errors.New("error: employer must have company for vacancy creation")
	}
	_, err = p.db.Model(&vac).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting vacancy with title: %s : error: %w", vac.Title, err)
		return nil, err
	}
	_, err = p.db.Model(&company).WherePK().Set("count_vacancy = count_vacancy + 1").Update()
	if err != nil {
		return nil, fmt.Errorf("error in update company with id:  %s vacCount  : error: %w", company.ID.String(), err)
	}
	return &vac, nil
}

func (p *pgRepository) GetVacancyById(id uuid.UUID) (*models.Vacancy, error) {
	var vac models.Vacancy
	if id == uuid.Nil {
		return nil, nil
	}
	vac.ID = id
	err := p.db.Model(&vac).WherePK().Select()
	if err != nil {
		return nil, fmt.Errorf("error in select vacancy with id: %s : error: %w", id.String(), err)
	}
	return &vac, nil
}

func (p *pgRepository) UpdateVacancy(newVac models.Vacancy) (*models.Vacancy, error) {
	oldVac, err := p.GetVacancyById(newVac.ID)
	if err != nil {
		return nil, fmt.Errorf("error in select vacancy with id: %s : error: %w", newVac.ID.String(), err)
	}
	if oldVac.EmpID != newVac.EmpID {
		return nil, fmt.Errorf("this user can't update this vacancy")
	}
	if _, err := p.db.Model(&newVac).WherePK().Returning("*").Update(); err != nil {
		return nil, fmt.Errorf("can't update vacancy with id:%s", newVac.ID)
	}
	return &newVac, nil
}

func (p *pgRepository) GetVacancyList(start uint, limit uint, id uuid.UUID, entityType int) ([]models.Vacancy, error) {
	var (
		vacList []models.Vacancy
		err     error
	)
	if limit <= start {
		return nil, fmt.Errorf("selection with useless positions")
	}
	if entityType == vacancy.ByEmpId {
		err = p.db.Model(&vacList).Where("empl_id= ?", id).Limit(int(limit)).Offset(int(start)).Select()
	} else if entityType == vacancy.ByCompId {
		err = p.db.Model(&vacList).Where("comp_id= ?", id).Limit(int(limit)).Offset(int(start)).Order("sphere").Select()
	} else {
		err = p.db.Model(&vacList).Limit(int(limit)).Offset(int(start)).Order("sphere").Select()
	}
	if err != nil {
		err = fmt.Errorf("error in list selection from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	return vacList, nil
}

func (p *pgRepository) SearchVacancies(params models.VacancySearchParams) ([]models.Vacancy, error) {
	var vacList []models.Vacancy

	err := p.db.Model(&vacList).WhereGroup(func(q *orm.Query) (*orm.Query, error) {
		if params.StartDate != "" {
			q = q.Where("date_create >= (?)", params.StartDate)
		}
		if len(params.Spheres) != 0 {
			q = q.Where("sphere IN (?)", pg.In(params.Spheres))
		}
		if len(params.EducationLevel) != 0 {
			q = q.Where("education_level IN (?)", pg.In(params.EducationLevel))
		}
		if len(params.CareerLevel) != 0 {
			q = q.Where("career_level IN (?)", pg.In(params.CareerLevel))
		}
		if len(params.ExperienceMonth) != 0 {
			q = q.Where("experience_month IN (?)", pg.In(params.ExperienceMonth))
		}
		if len(params.Employment) != 0 {
			q = q.Where("employment IN (?)", pg.In(params.ExperienceMonth))
		}
		if len(params.AreaSearch) != 0 {
			q = q.Where("area_search IN (?)", pg.In(params.AreaSearch))
		}
		if params.SalaryMin != 0 || params.SalaryMax != 0 {
			q = q.Where("salary_min >= ?", params.SalaryMin).
				Where("salary_max <= ?", params.SalaryMax)
		}
		if params.KeyWords != "" {
			q = q.Where("LOWER(title) LIKE ?", "%"+params.KeyWords+"%")
		}
		if params.OrderBy != "" {
			return q.Order(params.OrderBy), nil
		}
		return q, nil
	}).Select()
	if err != nil {
		err = fmt.Errorf("error in vacancies list selection with searchParams: %s", err)
		return nil, err
	}

	return vacList, nil
}
