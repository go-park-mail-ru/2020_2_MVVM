package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(resume models.Resume) (*models.Resume, error)
	Update(newResume *models.Resume) (*models.Resume, error)
	Search(searchParams *SearchParams) ([]models.ResumeWithCandidate, error)
	GetById(id uuid.UUID) (*models.Resume, error)
	GetAllUserResume(userID uuid.UUID) ([]models.ResumeWithCandidate, error)
	List(start, limit uint) ([]models.Resume, error)

	AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error)
	RemoveFavorite(favoriteForEmpl uuid.UUID) error
	GetAllEmplFavoriteResume(empl_id uuid.UUID) ([]models.FavoritesForEmplWithResume, error)
	GetFavoriteForResume(userID uuid.UUID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error)
}
