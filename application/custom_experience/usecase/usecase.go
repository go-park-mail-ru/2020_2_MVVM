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

func (u *UseCase) GetAllResumeCustomExperienceWithCompanies(resumeID uuid.UUID) ([]models.CustomExperienceWithCompanies, error) {
	exp, err := u.customExperienceRepository.GetAllResumeCustomExperience(resumeID)
	if err != nil {
		err = fmt.Errorf("error in get custom experience function: %w", err)
		return nil, err
	}

	var expWithComp []models.CustomExperienceWithCompanies

	for _, item := range exp {
		company, err := u.customCompanyRepository.GetCustomCompanyById(item.CompanyID.String())
		if err != nil {
			err = fmt.Errorf("error in get custom company function: %w", err)
			return nil, err
		}
		insert := models.CustomExperienceWithCompanies{
			CompanyName: company.Name,
			Location:    company.Location,
			Sphere:      company.Sphere,
			Position:    item.Position,
			Begin:       item.Begin,
			Finish:      item.Finish,
			Description: item.Description,
		}
		expWithComp = append(expWithComp, insert)
	}
	return expWithComp, nil
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
