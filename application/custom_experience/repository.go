package custom_experience

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type CustomExperienceRepository interface {
	CreateCustomExperience(experience []models.ExperienceCustomComp) ([]models.ExperienceCustomComp, error)
	//UpdateEducation(id uuid.UUID, updResume *models.Resume) (*models.Resume, error)
	GetCustomExperienceById(id string) (*models.ExperienceCustomComp, error)
	GetAllResumeCustomExperience(experienceID uuid.UUID) ([]models.ExperienceCustomComp, error)
}
