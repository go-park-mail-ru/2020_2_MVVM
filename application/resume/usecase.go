package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type UseCase interface {
	Create(template models.Resume) (*models.Resume, error)
	Update(resume models.Resume) (*models.Resume, error)
	Search(searchParams SearchParams) ([]models.BriefResumeInfo, error)

	GetById(id uuid.UUID) (*models.Resume, error)
	List(start, limit uint) ([]models.BriefResumeInfo, error)
	GetAllUserResume(userid uuid.UUID) ([]models.BriefResumeInfo, error)

	AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error)
	GetFavoriteByResume(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error)
	GetFavoriteByID(favoriteID uuid.UUID) (*models.FavoritesForEmpl, error)
	RemoveFavorite(favorite models.FavoritesForEmpl) error
	GetAllEmplFavoriteResume(userID uuid.UUID) ([]models.BriefResumeInfo, error)
}
