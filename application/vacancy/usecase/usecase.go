package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"math"
	"strings"
	"time"
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

func (v VacancyUseCase) CreateVacancy(vacancy models.Vacancy) (*models.Vacancy, error) {
	vac, err := v.repos.CreateVacancy(vacancy)
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

func (v VacancyUseCase) GetVacancyList(start uint, limit uint, empId uuid.UUID) ([]models.Vacancy, error) {
	vacList, err := v.repos.GetVacancyList(start, limit, empId)
	if err != nil {
		err = fmt.Errorf("error in vacancy list creation: %w", err)
		return nil, err
	}
	return vacList, nil
}

func (v VacancyUseCase) SearchVacancies(params models.VacancySearchParams) ([]models.Vacancy, error) {
	if params.SalaryMax == 0 {
		params.SalaryMax = math.MaxInt64
	}
	if params.DaysFromNow >= 0 {
		params.StartDate = time.Now().AddDate(0, 0, -params.DaysFromNow).Format("2006-01-02")
	}
	if params.OrderBy != "" {
		if s := params.OrderBy; s == "salary_min" || s == "salary_max" || s == "week_work_hours" || s == "date_create" {
			if params.ByAsc {
				params.OrderBy += " ASC"
			} else {
				params.OrderBy += " DESC"
			}
		} else {
			params.OrderBy = ""
		}
	}
	params.KeyWords = strings.ToLower(params.KeyWords)
	vacList, err := v.repos.SearchVacancies(params)
	if err != nil {
		return nil, err
	}
	return vacList, nil
}
