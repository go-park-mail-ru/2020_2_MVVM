package response

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type ResponseRepository interface {
	GetByID(responseID uuid.UUID) (*models.Response, error)
	Create(models.Response) (*models.Response, error)
	UpdateStatus(models.Response) (*models.Response, error)
	GetRespNotifications(respIds map[uuid.UUID]bool) ([]models.Response, error)
	GetResumeAllResponse(uuid uuid.UUID) ([]models.Response, error)
	GetVacancyAllResponse(uuid uuid.UUID) ([]models.Response, error)
	GetAllResumeWithoutResponse(candID uuid.UUID, vacancyID uuid.UUID) ([]models.Resume, error)
	GetAllVacancyWithoutResponse(emplID uuid.UUID, resumeID uuid.UUID) ([]models.Vacancy, error)
}
