package custom_experience

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/dto/models"
	"github.com/google/uuid"
)

type CustomExperienceRepository interface {
	Create(experience models.ExperienceCustomComp) (*models.ExperienceCustomComp, error)
	DropAllFromResume(resumeID uuid.UUID) error
	//GetById(id string) (*models.ExperienceCustomComp, error)
	//GetAllFromResume(experienceID uuid.UUID) ([]models.ExperienceCustomComp, error)
}
