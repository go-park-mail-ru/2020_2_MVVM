package response

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseResponse interface {
	Create(models.Response) (*models.Response, error)
	UpdateStatus(response models.Response, userType string) (*models.Response, error)
	GetAllCandidateResponses(uuid.UUID) ([]models.ResponseWithTitle, error)
	GetAllEmployerResponses(uuid.UUID) ([]models.ResponseWithTitle, error)
	GetAllResumeWithoutResponse(candID uuid.UUID, vacancyID uuid.UUID) ([]models.BriefResumeInfo, error)
}
