package custom_experience

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type CustomExperienceRepository interface {
	Create(experience models.ExperienceCustomComp) (*models.ExperienceCustomComp, error)
	//Update(id uuid.UUID, updResume *models.Resume) (*models.Resume, error)
	GetById(id string) (*models.ExperienceCustomComp, error)
	GetAllFromResume(experienceID uuid.UUID) ([]models.ExperienceCustomComp, error)
	DeleteAllResumeCustomExperience(experienceID uuid.UUID) error
}
