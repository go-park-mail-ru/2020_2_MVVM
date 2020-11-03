package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseResume interface {
	CreateResume(resume models.Resume) (*models.Resume, error)
	UpdateResume(resume models.Resume) (*models.Resume, error)
	SearchResume(searchParams models.SearchResume) ([]models.Resume, error)
	GetResume(id string) (*models.Resume, error)
	GetResumePage(start, limit uint) ([]models.Resume, error)
	GetAllUserResume(userid uuid.UUID) ([]models.Resume, error)

	AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error)
	GetFavoriteForResume(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error)
	RemoveFavorite(favoriteForEmpl uuid.UUID) error
	GetAllEmplFavoriteResume(userid uuid.UUID) ([]models.Resume, error)
}
