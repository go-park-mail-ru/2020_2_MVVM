package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
	logger "github.com/rowdyroad/go-simple-logger"
)

type UseCaseEducation struct {
	infoLogger  *logger.Logger
	errorLogger *logger.Logger
	strg        education.EducationRepository
}

func NewUsecase(infoLogger *logger.Logger,
				errorLogger *logger.Logger,
				strg education.EducationRepository) *UseCaseEducation {
					usecase := UseCaseEducation {
					infoLogger:  infoLogger,
					errorLogger: errorLogger,
					strg:        strg,
	}
	return &usecase
}

func (u* UseCaseEducation) GetAllResumeEducation(resumeID uuid.UUID) ([]models.Education, error) {
	return u.strg.GetAllResumeEducation(resumeID)
}

func (u *UseCaseEducation) CreateEducation(educations []models.Education) ([]models.Education, error) {
	ed, err := u.strg.CreateEducation(educations)
	if err != nil {
		err = fmt.Errorf("error in create educations function: %w", err)
		return nil, err
	}
	return ed, nil
}

func (u *UseCaseEducation) GetEducation(id string) (*models.Education, error) {
	ed, err := u.strg.GetEducationById(id)
	if err != nil {
		err = fmt.Errorf("error in education get by id func : %w", err)
		return nil, err
	}
	return ed, nil
}

func (u *UseCaseEducation) UpdateEducation(newEducations []models.Education, resumeID uuid.UUID) ([]models.Education, error) {
	err := u.strg.DeleteAllResumeEducation(resumeID)
	if err != nil {
		return nil, err
	}
	ed, err := u.strg.CreateEducation(newEducations)
	if err != nil {
		err = fmt.Errorf("error in update educations function: %w", err)
		return nil, err
	}
	return ed, nil
}
