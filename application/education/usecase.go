package education

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseEducation interface {
	CreateEducation(educations []models.Education) ([]models.Education, error)
	UpdateEducation(education []models.Education, resumeID uuid.UUID) ([]models.Education, error)
	GetEducation(id string) (*models.Education, error)
	GetAllResumeEducation(resumeID uuid.UUID) ([]models.Education, error)
}
