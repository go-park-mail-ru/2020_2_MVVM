package usecase

import (
	RepositoryVacancy "github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy/repository"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/vacancy"
	mocks "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/vacancy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func beforeTest() (*mocks.RepositoryVacancy, VacancyUseCase) {
	mockRepo := new(mocks.RepositoryVacancy)
	useCase := VacancyUseCase{
		repos: mockRepo,
	}
	return mockRepo, useCase
}

func TestGetVacancy(t *testing.T) {
	mockRepo, useCase := beforeTest()
	vacID := uuid.New()

	vacancy := models.Vacancy{
		ID: vacID,
	}
	vacEmptyId := uuid.Nil
	mockRepo.On("GetVacancyById", vacID).Return(&vacancy, nil)
	mockRepo.On("GetVacancyById", vacEmptyId).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.GetVacancy(vacID)
	ansWrong, errNotNil := useCase.GetVacancy(vacEmptyId)
	assert.Nil(t, errNil)
	assert.Equal(t, *ansCorrect, vacancy)
	assert.Error(t, errNotNil)
	assert.Nil(t, ansWrong)
}

func TestCreateVacancy(t *testing.T) {
	mockRepo, useCase := beforeTest()

	vacancy := models.Vacancy{
		Title:       "title",
		Description: "description",
		AreaSearch:  "area",
	}
	vacEmpty := models.Vacancy{}
	mockRepo.On("CreateVacancy", vacancy).Return(&vacancy, nil)
	mockRepo.On("CreateVacancy", vacEmpty).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.CreateVacancy(vacancy)
	ansWrong, errNotNil := useCase.CreateVacancy(vacEmpty)
	assert.Nil(t, errNil)
	assert.Equal(t, *ansCorrect, vacancy)
	assert.Error(t, errNotNil)
	assert.Nil(t, ansWrong)
}

func TestUpdateVacancy(t *testing.T) {
	mockRepo, useCase := beforeTest()

	vacancyNew := models.Vacancy{
		Title:       "title",
		Description: "description",
		AreaSearch:  "area",
	}
	vacEmpty := models.Vacancy{}
	mockRepo.On("UpdateVacancy", vacancyNew).Return(&vacancyNew, nil)
	mockRepo.On("UpdateVacancy", vacEmpty).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.UpdateVacancy(vacancyNew)
	ansWrong, errNotNil := useCase.UpdateVacancy(vacEmpty)
	assert.Nil(t, errNil)
	assert.Equal(t, *ansCorrect, vacancyNew)
	assert.Error(t, errNotNil)
	assert.Nil(t, ansWrong)
}

func TestGetVacancyList(t *testing.T) {

	var start uint = 0
	var end uint = 3
	id := uuid.New()
	typeDb := 0
	mockRepo, useCase := beforeTest()

	vacList := []models.Vacancy{
		{Title: "title1", Description: "description1", AreaSearch: "area1"},
		{Title: "title2", Description: "description2", AreaSearch: "area2"},
		{Title: "title3", Description: "description3", AreaSearch: "area3"}}

	mockRepo.On("GetVacancyList", start, end, id, typeDb).Return(vacList, nil)
	mockRepo.On("GetVacancyList", start, start, id, typeDb).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.GetVacancyList(start, end, id, typeDb)
	ansWrong, errNotNil := useCase.GetVacancyList(start, start, id, typeDb)
	assert.Nil(t, errNil)
	assert.Equal(t, ansCorrect, vacList)
	assert.Error(t, errNotNil)
	assert.Nil(t, ansWrong)
}

func TestSearchVacancies(t *testing.T) {

	mockRepo, useCase := beforeTest()

	vacList := []models.Vacancy{
		{Title: "title1", Description: "description1", AreaSearch: "area1", SalaryMax: 15},
		{Title: "title2", Description: "description2", AreaSearch: "area1", SalaryMax: 12},
		{Title: "title3", Description: "description3", AreaSearch: "area3", SalaryMax: 8}}

	params := models.VacancySearchParams{AreaSearch: []string{"area1"}, SalaryMax: 0, DaysFromNow: 7, OrderBy: "salary_max", ByAsc: false}
	finalParams := params
	finalParams.OrderBy += " DESC"
	finalParams.StartDate = time.Now().AddDate(0, 0, -params.DaysFromNow).Format("2006-01-02")
	finalParams.SalaryMax = math.MaxInt32

	mockRepo.On("SearchVacancies", finalParams).Return(vacList[:2], nil)
	ansCorrect, errNil := useCase.SearchVacancies(params)
	assert.Equal(t, ansCorrect, vacList[:2])
	assert.Nil(t, errNil)

	//paramsErr := models.VacancySearchParams{SalaryMax: -1}
	//mockRepo.On("SearchVacancies", paramsErr).Return(nil, assert.AnError)
	//ansWrong, errNotNil := useCase.SearchVacancies(paramsErr)
	//assert.Error(t, errNotNil)
	//assert.Nil(t, ansWrong)
}

func TestAddRecommendation(t *testing.T) {
	mockRepo, useCase := beforeTest()
	userID := uuid.New()
	sphere := 0

	emptyId := uuid.Nil
	mockRepo.On("AddRecommendation", userID, sphere).Return(nil)
	mockRepo.On("AddRecommendation", emptyId, sphere).Return(assert.AnError)
	errNil := useCase.AddRecommendation(userID, 0)
	errNotNil := useCase.AddRecommendation(emptyId, sphere)
	assert.Nil(t, errNil)
	assert.Error(t, errNotNil)
}

func TestGetRecommendation(t *testing.T) {
	mockRepo, useCase := beforeTest()
	userID := uuid.New()
	start := 0
	vacList := []models.Vacancy{
		{Title: "title1", Description: "description1", AreaSearch: "area1", SalaryMax: 15},
		{Title: "title2", Description: "description2", AreaSearch: "area1", SalaryMax: 12},
		{Title: "title3", Description: "description3", AreaSearch: "area3", SalaryMax: 8}}
	vacPair := []vacancy.Pair{{Score: 5, SphereInd: 2}, {Score: 3, SphereInd: 1}}
	var salary float64 = 0

	length := len(vacList)
	emptyId := uuid.Nil
	mockRepo.On("GetPreferredSpheres", userID).Return(vacPair, nil)
	mockRepo.On("GetPreferredSpheres", emptyId).Return(nil, assert.AnError)
	mockRepo.On("GetPreferredSalary", userID).Return(&salary, nil)
	mockRepo.On("GetPreferredSalary", emptyId).Return(nil, assert.AnError)
	mockRepo.On("GetRecommendation", start, length, salary, []int{2, 1}).Return(vacList, nil).Once()
	mockRepo.On("GetRecommendation", start, length, salary, []int{2, 1}).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.GetRecommendation(userID, 0, length)
	ansWrong, errNotNil := useCase.GetRecommendation(emptyId, 0, length)
	assert.Equal(t, ansCorrect, vacList)
	assert.Nil(t, ansWrong)
	assert.Nil(t, errNil)
	assert.Error(t, errNotNil)
}

func TestNewVacUseCase(t *testing.T) {
	vacancyRep := RepositoryVacancy.NewPgRepository(nil)
	vac := NewVacUseCase(nil, nil, vacancyRep)
	assert.Equal(t, vac, &VacancyUseCase{iLog: nil, errLog: nil, repos: vacancyRep})
}
