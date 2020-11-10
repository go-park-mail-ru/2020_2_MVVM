package response

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type ResponseRepository interface {
	Create(models.Response) (*models.Response, error)
	UpdateStatus(models.Response) (*models.Response, error)
	GetResumeAllResponse(uuid uuid.UUID) ([]models.Response, error)
}
