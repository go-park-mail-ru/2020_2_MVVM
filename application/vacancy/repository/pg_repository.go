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

func dbSelector(P *pgRepository, pattern string, attribute string) (models.Vacancy, error){
	var vac models.Vacancy
	err := P.db.Model(&vac).Where(pattern, attribute).Select()
	if err != nil {
		err = fmt.Errorf("error in select resume with pattern: %s : error: %w", pattern, err)
		return models.Vacancy{}, err
	}
	return vac, nil
}