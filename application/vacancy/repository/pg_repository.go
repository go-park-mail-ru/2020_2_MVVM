package repository

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-pg/pg/v9"
)

type pgRepository struct {
	db *pg.DB
}

func NewPgRepository(db *pg.DB) vacancy.RepositoryVacancy {
	return &pgRepository{db: db}
}

func (P *pgRepository) CreateVacancy(vac models.Vacancy) (models.Vacancy, error) {
	_, err := P.db.Model(&vac).Returning("*").Insert()
	if err != nil {
		err = fmt.Errorf("error in inserting vacancy with title: %s : error: %w", vac.VacancyName, err)
		return models.Vacancy{}, err
	}
	return vac, nil
}

func (P *pgRepository) GetVacancyById(id string) (models.Vacancy, error) {
	return dbSelector(P, "vacancy_id = ?", id)
}

func (P *pgRepository) GetVacancyByName(name string) (models.Vacancy, error) {
	return dbSelector(P, "vacancy_name = ?", name)
}

func dbSelector(P *pgRepository, pattern string, attribute string) (models.Vacancy, error) {
	var vac models.Vacancy
	err := P.db.Model(&vac).Where(pattern, attribute).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume with pattern: %s : error: %w", pattern, err)
		return models.Vacancy{}, err
	}
	return vac, nil
}

func (P *pgRepository) UpdateVacancy(id string, newVac models.Vacancy) (models.Vacancy, error) {
	oldVac, err := P.GetVacancyById(id)
	if err != nil {
		return models.Vacancy{}, err
	}
	switch {
	case newVac.VacancyName != "":
		oldVac.VacancyName = newVac.VacancyName
	case newVac.CompanyName != "":
		oldVac.CompanyName = newVac.CompanyName
	case newVac.VacancyDescription != "":
		oldVac.VacancyDescription = newVac.VacancyDescription
	case newVac.CompanyAddress != "":
		oldVac.CompanyAddress = newVac.CompanyAddress
	case newVac.WorkExperience != "":
		oldVac.WorkExperience = newVac.WorkExperience
	case newVac.Skills != "":
		oldVac.Skills = newVac.Skills
	case newVac.Salary != 0:
		oldVac.Salary = newVac.Salary
	}
	_, err = P.db.Model(&oldVac).WherePK().Update()
	if err != nil {
		err = fmt.Errorf("error in update resume with id: %s : error: %w", id, err)
		return models.Vacancy{}, err
	}
	return oldVac, nil
}

func (P *pgRepository) GetVacancyList(start uint, end uint) ([]models.Vacancy, error) {
	if end <= start {
		return nil, fmt.Errorf("selection with useless positions")
	}
	var vacList []models.Vacancy
	err := P.db.Model(&vacList).Where(fmt.Sprintf("vacancy_idx > %v", end)).Limit(int(start)).Select()
	if err != nil {
		err = fmt.Errorf("error in list selection from %v to %v: error: %w", start, end, err)
		return nil, err
	}
	return vacList, nil
}
