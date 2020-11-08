package custom_experience

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type UseCase interface {
	Create(experience models.ExperienceCustomComp) (*models.ExperienceCustomComp, error)
	GetById(id string) (*models.ExperienceCustomComp, error)
	DropAllFromResume(resumeID uuid.UUID) error
	GetAllFromResume(ResumeID uuid.UUID) ([]models.ExperienceCustomComp, error)
}
