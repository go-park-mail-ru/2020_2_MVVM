package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
	logger "github.com/rowdyroad/go-simple-logger"
)

type UseCase struct {
	infoLogger                 *logger.Logger
	errorLogger                *logger.Logger
	customExperienceRepository custom_experience.CustomExperienceRepository
	customCompanyRepository    custom_company.CustomCompanyRepository
}

func NewUsecase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	customExperienceRepository custom_experience.CustomExperienceRepository,
	customCompanyRepository custom_company.CustomCompanyRepository) *UseCase {
	usecase := UseCase{
		infoLogger:                 infoLogger,
		errorLogger:                errorLogger,
		customExperienceRepository: customExperienceRepository,
		customCompanyRepository:    customCompanyRepository,
	}
	return &usecase
}

func (u *UseCase) GetAllResumeCustomExperience(resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	return u.customExperienceRepository.GetAllResumeCustomExperience(resumeID)
}

func (u *UseCase) CreateCustomExperience(experience []models.ExperienceCustomComp) ([]models.ExperienceCustomComp, error) {
	ed, err := u.customExperienceRepository.CreateCustomExperience(experience)
	if err != nil {
		err = fmt.Errorf("error in create custom experience function: %w", err)
		return nil, err
	}
	return ed, nil
}

func (u *UseCase) GetCustomExperience(id string) (*models.ExperienceCustomComp, error) {
	experience, err := u.customExperienceRepository.GetCustomExperienceById(id)
	if err != nil {
		err = fmt.Errorf("error in get by id custom experience func : %w", err)
		return nil, err
	}
	return experience, nil
}

func (u *UseCase) UpdateCustomExperience(newExperience []models.ExperienceCustomComp, resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	err := u.customExperienceRepository.DeleteAllResumeCustomExperience(resumeID)
	if err != nil {
		return nil, err
	}
	exp, err := u.customExperienceRepository.CreateCustomExperience(newExperience)
	if err != nil {
		err = fmt.Errorf("error in update custom experience function: %w", err)
		return nil, err
	}
	return exp, nil
}
