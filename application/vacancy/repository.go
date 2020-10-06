package vacancy

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type RepositoryVacancy interface {
	CreateVacancy(models.Vacancy) (models.Vacancy, error)
	UpdateVacancy(string, models.Vacancy) (models.Vacancy, error)
	GetVacancyById(string) (models.Vacancy, error)
	GetVacancyByName(string) (models.Vacancy, error)
	GetVacancyList(uint, uint) ([]models.Vacancy, error)
}
