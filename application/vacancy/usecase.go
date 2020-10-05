package vacancy

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type IUseCaseVacancy interface {
	CreateVacancy(models.Vacancy) (models.Vacancy, error)
	//UpdateVacancy(id uint) (models.Vacancy, error)
	GetVacancy(id string) (models.Vacancy, error)
	//GetVacancyList(begin, end uint) (models.Vacancy, error)
}