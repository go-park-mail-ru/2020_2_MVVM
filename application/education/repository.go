package education

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type EducationRepository interface {
	CreateEducation(resume []models.Education) ([]models.Education, error)
	//UpdateEducation(id uuid.UUID, updResume *models.Resume) (*models.Resume, error)
	GetEducationById(id string) (*models.Education, error)
	GetAllResumeEducation(resumeID uuid.UUID) ([]models.Education, error)
	DeleteAllResumeEducation(resumeID uuid.UUID) error
}
