package resume

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseResume interface {
	CreateResume(resume models.Resume) (*models.Resume, error)
	UpdateResume(resume models.Resume) (*models.Resume, error)
	GetResume(id string) (*models.Resume, error)
	GetResumePage(start, limit uint) ([]models.Resume, error)
	GetAllUserResume(userid uuid.UUID) ([]models.Resume, error)
}
