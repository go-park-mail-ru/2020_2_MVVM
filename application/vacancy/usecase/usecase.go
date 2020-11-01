package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
)

type VacancyUseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  vacancy.RepositoryVacancy
}

func NewVacUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos vacancy.RepositoryVacancy) *VacancyUseCase {
	return &VacancyUseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
	}
}

func (v VacancyUseCase) CreateVacancy(vacancy models.Vacancy, userId uuid.UUID) (*models.Vacancy, error) {
	vac, err := v.repos.CreateVacancy(vacancy, userId)
	if err != nil {
		err = fmt.Errorf("error in vacancy creation: %w", err)
		return nil, err
	}
	return vac, nil
}

func (v VacancyUseCase) GetVacancy(id string) (*models.Vacancy, error) {
	vac, err := v.repos.GetVacancyById(id)
	if err != nil {
		err = fmt.Errorf("error in vacancy selection get: %w", err)
		return nil, err
	}
	return vac, nil
}

func (v VacancyUseCase) UpdateVacancy(newVac models.Vacancy) (*models.Vacancy, error) {
	vac, err := v.repos.UpdateVacancy(newVac)
	if err != nil {
		err = fmt.Errorf("error in vacancy update: %w", err)
		return nil, err
	}
	return vac, nil
}

func (v VacancyUseCase) GetVacancyList(start, end uint, empId uuid.UUID) ([]models.Vacancy, error) {
	vacList, err := v.repos.GetVacancyList(start, end, empId)
	if err != nil {
		err = fmt.Errorf("error in vacancy list creation: %w", err)
		return nil, err
	}
	return vacList, nil
}
