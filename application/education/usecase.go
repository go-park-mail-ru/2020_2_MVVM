package education

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseEducation interface {
	CreateEducation(resume []models.Education) ([]models.Education, error)
	//UpdateResume(resume models.Resume) (*models.Resume, error)
	GetEducation(id string) (*models.Education, error)
	GetAllResumeEducation(ResumeID uuid.UUID) ([]models.Education, error)
}
