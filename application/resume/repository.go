package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(resume models.Resume) (*models.Resume, error)
	Drop(resume models.Resume) error
	Update(resume models.Resume) (*models.Resume, error)
	Search(searchParams *SearchParams) ([]models.Resume, error)
	GetById(id uuid.UUID) (*models.Resume, error)
	GetAllUserResume(userID uuid.UUID) ([]models.Resume, error)
	List(start, limit uint) ([]models.Resume, error)
	SelectCandidate (candID uuid.UUID) (*models.Candidate, error)

	AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error)
	RemoveFavorite(favoriteForEmpl uuid.UUID) error
	GetAllEmplFavoriteResume(emplID uuid.UUID) ([]models.Resume, error)
	GetFavoriteForResume(userID uuid.UUID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error)
	GetFavoriteByID(favoriteID uuid.UUID) (*models.FavoritesForEmpl, error)
}
