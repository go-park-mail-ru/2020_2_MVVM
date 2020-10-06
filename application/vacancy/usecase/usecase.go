package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	logger "github.com/rowdyroad/go-simple-logger"
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

func (V VacancyUseCase) CreateVacancy(vacancy models.Vacancy) (models.Vacancy, error) {
	vac, err := V.repos.CreateVacancy(vacancy)
	if err != nil {
		err = fmt.Errorf("error in vacancy creation get by id func : %w", err)
		return models.Vacancy{}, err
	}
	return vac, nil
}

func (V VacancyUseCase) GetVacancy(id string) (models.Vacancy, error) {
	vac, err := V.repos.GetVacancyById(id)
	if err != nil {
		err = fmt.Errorf("error in vacancy selection get by id func : %w", err)
		return models.Vacancy{}, err
	}
	return vac, nil
}
