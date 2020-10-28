package vacancy

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type RepositoryVacancy interface {
	CreateVacancy(models.Vacancy, uuid.UUID) (*models.Vacancy, error)
	UpdateVacancy(models.Vacancy) (*models.Vacancy, error)
	GetVacancyById(string) (*models.Vacancy, error)
	GetVacancyByName(string) (*models.Vacancy, error)
	GetVacancyList(uint, uint) ([]models.Vacancy, error)
}
