package response

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseResponse interface {
	CreateResponse(models.Response) (*models.Response, error)
	UpdateStatus(models.Response) (*models.Response, error)
	GetAllUserResponses(uuid.UUID) (*models.ResponseWithTitle, error)
}
