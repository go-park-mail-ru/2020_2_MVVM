package vacancy

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type IUseCaseVacancy interface {
	CreateVacancy(models.Vacancy) (models.Vacancy, error)
	UpdateVacancy(string, models.Vacancy) (models.Vacancy, error)
	GetVacancy(string) (models.Vacancy, error)
	GetVacancyList(uint, uint) ([]models.Vacancy, error)
}
