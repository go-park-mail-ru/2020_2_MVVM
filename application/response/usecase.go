package response

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
)

type IUseCaseResponse interface {
	Create(models.Response) (*models.Response, error)
	UpdateStatus(response models.Response, userType string) (*models.Response, error)
	GetAllCandidateResponses(candId uuid.UUID, respIds []uuid.UUID) ([]models.ResponseWithTitle, error)
	GetAllEmployerResponses(emplId uuid.UUID, respIds []uuid.UUID) ([]models.ResponseWithTitle, error)
	GetAllResumeWithoutResponse(candID uuid.UUID, vacancyID uuid.UUID) ([]models.BriefResumeInfo, error)
	GetAllVacancyWithoutResponse(emplID uuid.UUID, resumeID uuid.UUID) ([]models.Vacancy, error)
	GetResponsesCnt(userId uuid.UUID, userType string) (uint, error)
	GetRecommendedVacCnt(candId uuid.UUID, daysFromNow int) (uint, error)
	GetRecommendedVacancies(candId uuid.UUID, start uint, limit uint, daysFromNow int) ([]models.Vacancy, error)
}
