package custom_experience

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseCustomExperience interface {
	CreateCustomExperience(experience []models.ExperienceCustomComp) ([]models.ExperienceCustomComp, error)
	//UpdateResume(resume models.Resume) (*models.Resume, error)
	GetCustomExperience(id string) (*models.ExperienceCustomComp, error)
	GetAllResumeCustomExperience(ResumeID uuid.UUID) ([]models.ExperienceCustomComp, error)
	GetAllResumeCustomExperienceWithCompanies(ResumeID uuid.UUID) ([]models.CustomExperienceWithCompanies, error)
}
