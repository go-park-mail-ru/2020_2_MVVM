package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
	logger "github.com/rowdyroad/go-simple-logger"
)

type UseCase struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        custom_experience.CustomExperienceRepository
}

func NewUsecase(infoLogger *logger.Logger,
				errorLogger *logger.Logger,
				strg custom_experience.CustomExperienceRepository) *UseCase {
					usecase := UseCase {
					infoLogger:  infoLogger,
					errorLogger: errorLogger,
					strg:        strg,
	}
	return &usecase
}

func (u* UseCase) GetAllResumeCustomExperience(resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	return u.strg.GetAllResumeCustomExperience(resumeID)
}

func (u *UseCase) CreateCustomExperience(experience []models.ExperienceCustomComp) ([]models.ExperienceCustomComp, error) {
	ed, err := u.strg.CreateCustomExperience(experience)
	if err != nil {
		err = fmt.Errorf("error in create custom experience function: %w", err)
		return nil, err
	}
	return ed, nil
}

func (u *UseCase) GetCustomExperience(id string) (*models.ExperienceCustomComp, error) {
	experience, err := u.strg.GetCustomExperienceById(id)
	if err != nil {
		err = fmt.Errorf("error in get by id custom experience func : %w", err)
		return nil, err
	}
	return experience, nil
}
