package usecase

import (
	"fmt"
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	vacancy2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/vacancy"
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

func (v VacancyUseCase) GetVacancy(id uuid.UUID) (*models.Vacancy, error) {
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

func (v VacancyUseCase) GetVacancyList(start uint, limit uint, id uuid.UUID, typeDb int) ([]models.Vacancy, error) {
	vacList, err := v.repos.GetVacancyList(start, limit, id, typeDb)
	if err != nil {
		err = fmt.Errorf("error in vacancy list creation: %w", err)
		return nil, err
	}
	return vacList, nil
}

func (v VacancyUseCase) SearchVacancies(params models.VacancySearchParams) ([]models.Vacancy, error) {
	if params.SalaryMax == 0 {
		params.SalaryMax = math.MaxInt32
	}
	if params.DaysFromNow > 0 {
		params.StartDate = time.Now().AddDate(0, 0, -params.DaysFromNow).Format("2006-01-02")
	}
	if params.OrderBy != "" {
		if s := params.OrderBy; s == "salary_min" || s == "salary_max" || s == "experience_month" || s == "date_create" {
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
	params.KeywordsGeo = strings.ToLower(params.KeywordsGeo)
	vacList, err := v.repos.SearchVacancies(params)
	if err != nil {
		return nil, err
	}
	return vacList, nil
}

func (v VacancyUseCase) AddRecommendation(userID uuid.UUID, sphere int) error {
	return v.repos.AddRecommendation(userID, sphere)
}

func (v VacancyUseCase) GetRecommendation(userID uuid.UUID, start int, limit int) ([]models.Vacancy, error) {
	preferredSphere, err := v.repos.GetPreferredSpheres(userID)
	if err != nil {
		if err.Error() == common.NoRecommendation {
			return nil, err
		}
		return nil, fmt.Errorf("error in GetUserRecommendations: %w", err)
	}
	step := 2
	curSphere := 0
	//preferredSalary, err := v.repos.GetPreferredSalary(userID)
	if err != nil {
		return nil, fmt.Errorf("error in GetPreferredSalary: %w", err)
	}

	var vacList []models.Vacancy

	for len(vacList) < limit && curSphere < vacancy2.CountSpheres {
		arr := []int{preferredSphere[curSphere].SphereInd, preferredSphere[curSphere+1].SphereInd}
		list, err := v.repos.GetRecommendation(start, limit, arr)
		vacList = append(vacList, list...)
		if err != nil {
			err = fmt.Errorf("error in GetRecommendation: %w", err)
			return nil, err
		}
		curSphere += step
		start = 0
	}
	end := limit
	if limit > len(vacList) {
		end = len(vacList)
	}
	return vacList[0:end], err
}

func (v VacancyUseCase) GetVacancyTopSpheres(topSpheresCnt int32) ([]models.Sphere, *models.VacTopCnt, error) {
	return v.repos.GetVacancyTopSpheres(topSpheresCnt)
}

func (v VacancyUseCase) DeleteVacancy(id uuid.UUID, empId uuid.UUID) error {
	return v.repos.DeleteVacancy(id, empId)
}



func DoBriefRespVacancy(vacancyList []models.Vacancy) ([]models.BriefVacancyInfo, error) {
	var briefRespResumes []models.BriefVacancyInfo
	for i := range vacancyList {
		brief, err := vacancyList[i].Brief()
		if err != nil {
			return nil, err
		}
		briefRespResumes = append(briefRespResumes, *brief)
	}
	return briefRespResumes, nil
}

func (v VacancyUseCase) GetFavoriteByVacancy(candId uuid.UUID, vacId uuid.UUID) (*models.FavoriteID, error) {
	return v.repos.GetFavoriteByVacancy(candId, vacId)
}

func (v VacancyUseCase) AddFavorite(favorite models.FavoritesForCand) (*models.FavoriteID, error) {
	return v.repos.AddFavorite(favorite)
}
func (v VacancyUseCase) RemoveFavorite(favorite models.FavoritesForCand) error {
	return v.repos.RemoveFavorite(favorite)
}
func (v VacancyUseCase) GetAllCandFavoriteVacancy(candId uuid.UUID) ([]models.BriefVacancyInfo, error) {
	vac, err := v.repos.GetAllCandFavoriteVacancy(candId)
	if err != nil {
		err = fmt.Errorf("error in get list favorite vacancy %w", err)
		return nil, err
	}
	return DoBriefRespVacancy(vac)
}
