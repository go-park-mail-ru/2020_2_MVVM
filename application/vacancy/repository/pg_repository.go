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

// TODO:
//у пользователя мб несколько вакансий, которые привязаны к одному user_id (FK)

func (P *pgRepository) UpdateVacancy(newVac models.Vacancy) (models.Vacancy, error) {
	/*oldVac, err := P.GetVacancyById(newVac.FK.String())
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
	_, err = P.db.Model(&oldVac).WherePK().Update()
	if err != nil {
		err = fmt.Errorf("error in update resume with id: %s : error: %w", newVac.ID, err)
		return models.Vacancy{}, err
	}
	return oldVac, nil
	*/
	return newVac, nil
}

func (P *pgRepository) GetVacancyList(start uint, end uint) ([]models.Vacancy, error) {
	if end <= start {
		return nil, fmt.Errorf("selection with useless positions")
	}
	var vacList []models.Vacancy
	err := P.db.Model(&vacList).Where(fmt.Sprintf("vacancy_idx >= %v", start)).Limit(int(end)).Select()
	if err != nil {
		err = fmt.Errorf("error in list selection from %v to %v: error: %w", start, end, err)
		return nil, err
	}
	return vacList, nil
}
