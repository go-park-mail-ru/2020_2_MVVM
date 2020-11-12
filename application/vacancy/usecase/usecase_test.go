package usecase

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
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
	params := models.VacancySearchParams{AreaSearch: []string{"area1"}, SalaryMax: math.MaxInt64, DaysFromNow: 7, OrderBy: "salary_max"}
	paramsErr := models.VacancySearchParams{SalaryMax: -1}
	params.OrderBy += " DESC"
	params.StartDate = time.Now().AddDate(0, 0, -params.DaysFromNow).Format("2006-01-02")
	mockRepo.On("SearchVacancies", params).Return(vacList[:2], nil)
	mockRepo.On("SearchVacancies", paramsErr).Return(nil, assert.AnError)
	params.SalaryMax = 0
	params.StartDate = ""
	params.OrderBy = "salary_max"
	ansCorrect, errNil := useCase.SearchVacancies(params)
	ansWrong, errNotNil := useCase.SearchVacancies(paramsErr)
	assert.Nil(t, errNil)
	assert.Equal(t, ansCorrect, vacList[:2])
	assert.Error(t, errNotNil)
	assert.Nil(t, ansWrong)
}
