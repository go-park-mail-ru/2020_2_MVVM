package vacancy

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type RepositoryVacancy interface {
	CreateVacancy(models.Vacancy) (*models.Vacancy, error)
	UpdateVacancy(models.Vacancy) (*models.Vacancy, error)
	GetVacancyById(uuid.UUID) (*models.Vacancy, error)
	GetVacancyList(uint, uint, uuid.UUID, int) ([]models.Vacancy, error)
	SearchVacancies(models.VacancySearchParams) ([]models.Vacancy, error)
	AddRecommendation(uuid.UUID, int) error
	GetRecommendation(uuid.UUID, int, int) ([]models.Vacancy, error)
}
