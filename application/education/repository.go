package education

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(resume models.Education) (*models.Education, error)
	//Drop(id uuid.UUID, updResume *models.Resume) (*models.Resume, error)
	//GetById(id string) (*models.Education, error)
	//GetAllFromResume(resumeID uuid.UUID) ([]models.Education, error)
	DropAllFromResume(resumeID uuid.UUID) error
}
