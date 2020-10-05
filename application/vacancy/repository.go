package vacancy

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type RepositoryVacancy interface {
	CreateVacancy(vacancy models.Vacancy) (models.Vacancy, error)
	//UpdateVacancy(id uint) (models.Vacancy, error)
	GetVacancyById(id string) (models.Vacancy, error)
	GetVacancyByName(name string) (models.Vacancy, error)
}






















