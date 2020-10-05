package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type ResumeRepository interface {
	CreateResume(resume models.Resume) (*models.Resume, error)
	UpdateResume(id uuid.UUID, updResume *models.Resume) (*models.Resume, error)
	GetResumeById(id string) (*models.Resume, error)
	GetResumeByName(name string) (*models.Resume, error)
	GetResumeArr(begin, end uint) ([]models.Resume, error)
}
