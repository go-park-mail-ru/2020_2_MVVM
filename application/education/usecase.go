package education

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type UseCase interface {
	Create(educations []*models.Education) ([]models.Education, error)
	Update(education []models.Education, resumeID uuid.UUID) ([]models.Education, error)
	GetById(id string) (*models.Education, error)
	GetAllFromResume(resumeID uuid.UUID) ([]models.Education, error)
}
