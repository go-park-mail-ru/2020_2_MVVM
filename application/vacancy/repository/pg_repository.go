package repository

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type pgRepository struct {
	db *gorm.DB
}

func NewPgRepository(db *gorm.DB) vacancy.RepositoryVacancy {
	return &pgRepository{db: db}
}

func (p *pgRepository) CreateVacancy(vac models.Vacancy) (*models.Vacancy, error) {
	employer := new(models.Employer)
	company  := new(models.OfficialCompany)

	err := p.db.Table("main.employers").Take(employer,"empl_id = ?", vac.EmpID).Error
	if err != nil {
		err = fmt.Errorf("error in FK search for vacancy creation for user with id: %s : error: %w", vac.EmpID, err)
		return nil, err
	}
	if compId := employer.CompanyID; compId != uuid.Nil {
		vac.ID = uuid.New()
		vac.CompID = compId
		company.ID = compId
	} else {
		return nil, errors.New("error: employer must have company for vacancy creation")
	}
	err = p.db.Table("main.vacancy").Create(&vac).Error
	if err != nil {
		err = fmt.Errorf("error in inserting vacancy with title: %s : error: %w", vac.Title, err)
		return nil, err
	}
	//_, err = p.db.Model(&company).WherePK().Set("count_vacancy = count_vacancy + 1").Update()
	err = p.db.Model(&company).UpdateColumn("count_vacancy", gorm.Expr("count_vacancy + ?", 1)).Error
	if err != nil {
		return nil, fmt.Errorf("error in update company with id:  %s vacCount  : error: %w", company.ID.String(), err)
	}
	return &vac, nil
}

func (p *pgRepository) GetVacancyById(id uuid.UUID) (*models.Vacancy, error) {
	vac :=  new(models.Vacancy)
	if id == uuid.Nil {
		return nil, nil
	}
	err := p.db.Table("main.vacancy").Take(vac, "vac_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("error in select vacancy with id: %s : error: %w", id.String(), err)
	}
	return vac, nil
}

func (p *pgRepository) UpdateVacancy(newVac models.Vacancy) (*models.Vacancy, error) {
	oldVac, err := p.GetVacancyById(newVac.ID)
	if err != nil {
		return nil, fmt.Errorf("error in select vacancy with id: %s : error: %w", newVac.ID.String(), err)
	}
	if oldVac.EmpID != newVac.EmpID {
		return nil, fmt.Errorf("this user can't update this vacancy")
	}
	if err := p.db.Model(&oldVac).Updates(newVac).Error; err != nil {
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
		err = p.db.Table("main.vacancy").Where("empl_id = ?", id).Limit(int(limit)).Offset(int(start)).Order("date_create").Find(&vacList).Error
	} else if entityType == vacancy.ByCompId {
		err = p.db.Table("main.vacancy").Where("comp_id = ?", id).Limit(int(limit)).Offset(int(start)).Order("date_create").Find(&vacList).Error
	} else {
		err = p.db.Table("main.vacancy").Limit(int(limit)).Offset(int(start)).Order("date_create").Find(&vacList).Error
	}
	if err != nil {
		err = fmt.Errorf("error in list selection from %v to %v: error: %w", start, limit, err)
		return nil, err
	}
	if len(vacList) == 0 {
		return nil, nil
	}
	return vacList, nil
}

func (p *pgRepository) SearchVacancies(params models.VacancySearchParams) ([]models.Vacancy, error) {
	var vacList []models.Vacancy

	err := p.db.Table("main.vacancy").Scopes(func(q *gorm.DB) *gorm.DB {
		if params.StartDate != "" {
			q = q.Where("date_create >= (?)", params.StartDate)
		}
		if len(params.Spheres) != 0 {
			q = q.Where("sphere IN (?)", params.Spheres)
		}
		if len(params.EducationLevel) != 0 {
			q = q.Where("education_level IN (?)", params.EducationLevel)
		}
		if len(params.CareerLevel) != 0 {
			q = q.Where("career_level IN (?)", params.CareerLevel)
		}
		if len(params.ExperienceMonth) != 0 {
			q = q.Where("experience_month IN (?)", params.ExperienceMonth)
		}
		if len(params.Employment) != 0 {
			q = q.Where("employment IN (?)", params.ExperienceMonth)
		}
		if len(params.AreaSearch) != 0 {
			q = q.Where("area_search IN (?)", params.AreaSearch)
		}
		if params.SalaryMin != 0 || params.SalaryMax != 0 {
			q = q.Where("salary_min >= ?", params.SalaryMin).
				Where("salary_max <= ?", params.SalaryMax)
		}
		if params.KeyWords != "" {
			q = q.Where("LOWER(title) LIKE ?", "%"+params.KeyWords+"%")
		}
		if params.OrderBy != "" {
			return q.Order(params.OrderBy)
		}
		return q
	}).Find(&vacList).Error

	if err != nil {
		err = fmt.Errorf("error in vacancies list selection with searchParams: %s", err)
		return nil, err
	}
	if len(vacList) == 0 {
		return nil, nil
	}
	return vacList, nil
}
