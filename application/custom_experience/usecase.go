package custom_experience

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type UseCase interface {
	CreateCustomExperience(experience []models.ExperienceCustomComp) ([]models.ExperienceCustomComp, error)
	GetCustomExperience(id string) (*models.ExperienceCustomComp, error)
	UpdateCustomExperience(experience []models.ExperienceCustomComp, resumeID uuid.UUID) ([]models.ExperienceCustomComp, error)
	GetAllResumeCustomExperience(ResumeID uuid.UUID) ([]models.ExperienceCustomComp, error)
}
