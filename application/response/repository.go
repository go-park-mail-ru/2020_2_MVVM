package response

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type ResponseRepository interface {
	CreateResponse(models.Response) (*models.Response, error)
	UpdateStatus(models.Response) (*models.Response, error)
}
