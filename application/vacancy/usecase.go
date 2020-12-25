package vacancy

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
)

const (
	ByEmpId        = 1
	ByCompId       = 2
	ByVacId        = 3
	Recommendation = 4
	TopDefaultCnt  = 6
	TopAll         = -1
)

type IUseCaseVacancy interface {
	CreateVacancy(models.Vacancy) (*models.Vacancy, error)
	UpdateVacancy(models.Vacancy) (*models.Vacancy, error)
	GetVacancy(uuid.UUID) (*models.Vacancy, error)
	GetVacancyList(uint, uint, uuid.UUID, int) ([]models.Vacancy, error)
	SearchVacancies(models.VacancySearchParams) ([]models.Vacancy, error)
	AddRecommendation(uuid.UUID, int) error
	GetRecommendation(uuid.UUID, int, int) ([]models.Vacancy, error)
	GetVacancyTopSpheres(int32) ([]models.Sphere, *models.VacTopCnt, error)
	DeleteVacancy(id uuid.UUID, empId uuid.UUID) error

	AddFavorite(models.FavoritesForCand) (*models.FavoriteID, error)
	RemoveFavorite(models.FavoritesForCand) error
	GetAllCandFavoriteVacancy(candId uuid.UUID) ([]models.BriefVacancyInfo, error)
}
