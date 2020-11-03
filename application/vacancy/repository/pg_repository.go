package repository

import (
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
		company models.OfficialCompany
	)
	err := p.db.Model(&employer).Where("empl_id = ?", vac.EmpID).Select()
	if err != nil {
		err = fmt.Errorf("error in FK search for vacancy creation for user with id: %s : error: %w", vac.EmpID, err)
		return nil, err
	}
	vac.CompID = employer.CompanyID
	company.ID = employer.CompanyID
	_, err = p.db.Model(&vac).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting vacancy with title: %s : error: %w", vac.Title, err)
		return nil, err
	}
	//TODO: update vac_count in companies

	return &vac, nil
}

func (p *pgRepository) GetVacancyById(id string) (*models.Vacancy, error) {
	return dbSelector(p, "vacancy_id = ?", id)
}

func (p *pgRepository) GetVacancyByName(name string) (*models.Vacancy, error) {
	return dbSelector(p, "title = ?", name)
}

func dbSelector(p *pgRepository, pattern string, attribute string) (*models.Vacancy, error) {
	var vac models.Vacancy
	err := p.db.Model(&vac).Where(pattern, attribute).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume with pattern: %s : error: %w", pattern, err)
		return nil, err
	}
	return &vac, nil
}

func (p *pgRepository) UpdateVacancy(newVac models.Vacancy) (*models.Vacancy, error) {
	/*oldVac, err := p.GetVacancyById(newVac.FK.String())
	if err != nil {
		return models.Vacancy{}, err
	}
	switch {
	case newVac.VacancyName != "":
		oldVac.VacancyName = newVac.VacancyName
		fallthrough
	case newVac.CompanyName != "":
		oldVac.CompanyName = newVac.CompanyName
		fallthrough
	case newVac.VacancyDescription != "":
		oldVac.VacancyDescription = newVac.VacancyDescription
		fallthrough
	case newVac.CompanyAddress != "":
		oldVac.CompanyAddress = newVac.CompanyAddress
		fallthrough
	case newVac.WorkExperience != "":
		oldVac.WorkExperience = newVac.WorkExperience
		fallthrough
	case newVac.Skills != "":
		oldVac.Skills = newVac.Skills
		fallthrough
	case newVac.Salary != 0:
		oldVac.Salary = newVac.Salary
	}
	_, err = p.db.Model(&oldVac).WherePK().Update()
	if err != nil {
		err = fmt.Errorf("error in update resume with id: %s : error: %w", newVac.ID, err)
		return models.Vacancy{}, err
	}
	return oldVac, nil
	*/
	return &newVac, nil
}

func (p *pgRepository) GetVacancyList(start uint, limit uint, empId uuid.UUID) ([]models.Vacancy, error) {
	var (
		vacList []models.Vacancy
		err     error
	)
	if limit <= start {
		return nil, fmt.Errorf("selection with useless positions")
	}
	if empId != uuid.Nil {
		err = p.db.Model(&vacList).Where("empl_id= ?", empId).Limit(int(limit)).Offset(int(start)).Select()
	} else {
		err = p.db.Model(&vacList).Limit(int(limit)).Offset(int(start)).Select()
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
			q = q.Where("date_create >= ?", params.StartDate)
		}
		if len(params.Spheres) != 0 {
			q = q.Where("spheres IN (?)", pg.In(params.Spheres))
		}
		if len(params.WeekWorkHours) != 0 {
			q = q.Where("week_work_hours IN (?)", pg.In(params.WeekWorkHours))
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
		if params.SalaryMin != 0 || params.SalaryMax != 0 {
			q = q.Where("salary_min >= ?", params.SalaryMin).
				Where("salary_max <= ?", params.SalaryMax)
		}
		if params.KeyWords != "" {
			q = q.Where("LOWER(title) LIKE ?", fmt.Sprintf("%%%s%", params.KeyWords)).
				WhereOr("LOWER(location) LIKE ?", fmt.Sprintf("%%%s%", params.KeyWords))
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
