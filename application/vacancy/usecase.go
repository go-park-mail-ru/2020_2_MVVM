package vacancy

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseVacancy interface {
	CreateVacancy(models.Vacancy) (*models.Vacancy, error)
	UpdateVacancy(models.Vacancy) (*models.Vacancy, error)
	GetVacancy(string) (*models.Vacancy, error)
	GetVacancyList(uint, uint, uuid.UUID) ([]models.Vacancy, error)
}
