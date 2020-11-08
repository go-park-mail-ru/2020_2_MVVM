package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type UseCaseEducation struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        education.Repository
}

func NewUsecase(infoLogger *logger.Logger,
				errorLogger *logger.Logger,
				strg education.Repository) *UseCaseEducation {
					usecase := UseCaseEducation {
					infoLogger:  infoLogger,
					errorLogger: errorLogger,
					strg:        strg,
	}
	return &usecase
}

func (u* UseCaseEducation) GetAllFromResume(resumeID uuid.UUID) ([]models.Education, error) {
	return u.strg.GetAllFromResume(resumeID)
}

func (u *UseCaseEducation) Create(educations models.Education) (*models.Education, error) {
	ed, err := u.strg.Create(educations)
	if err != nil {
		err = fmt.Errorf("error in create educations function: %w", err)
		return nil, err
	}
	return ed, nil
}

func (u *UseCaseEducation) GetById(id string) (*models.Education, error) {
	ed, err := u.strg.GetById(id)
	if err != nil {
		err = fmt.Errorf("error in education get by id func : %w", err)
		return nil, err
	}
	return ed, nil
}

func (u *UseCaseEducation) DropAllFromResume(resumeID uuid.UUID) error {
	return u.strg.DropAllFromResume(resumeID)
}
