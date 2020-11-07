package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateResume(resume models.Resume) (*models.Resume, error)
	UpdateResume(resume models.Resume) (*models.Resume, error)
	SearchResume(searchParams SearchParams) ([]models.BriefResumeInfo, error)

	GetResume(id string) (*models.Resume, error)
	GetResumePage(start, limit uint) ([]models.BriefResumeInfo, error)
	GetAllUserResume(userid uuid.UUID) ([]models.BriefResumeInfo, error)

	AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error)
	GetFavorite(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error)
	RemoveFavorite(favoriteForEmpl uuid.UUID) error
	GetAllEmplFavoriteResume(userid uuid.UUID) ([]models.BriefResumeInfo, error)
}
