package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type ResumeRepository interface {
	CreateResume(resume models.Resume) (*models.Resume, error)
	UpdateResume(newResume *models.Resume) (*models.Resume, error)
	GetResumeById(id string) (*models.Resume, error)
	GetAllUserResume(userID uuid.UUID) ([]models.Resume, error)
	GetResumeByName(name string) (*models.Resume, error)
	GetResumeArr(start, limit uint) ([]models.Resume, error)
}
