package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
)

type pgRepository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) vacancy.RepositoryVacancy {
	return &pgRepository{db: db}
}

func (p *pgRepository) CreateVacancy(vac models.Vacancy) (*models.Vacancy, error) {
	var employer models.Employer
	err := p.db.Model(&employer).Where("empl_id = ?", vac.EmpID).Select()
	if err != nil {
		err = fmt.Errorf("error in FK search for vacancy creation for user with id: %s : error: %w", vac.EmpID, err)
		return nil, err
	}
	vac.CompID = employer.CompanyID
	_, err = p.db.Model(&vac).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting vacancy with title: %s : error: %w", vac.Title, err)
		return nil, err
	}
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
