package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	mocks "github.com/go-park-mail-ru/2020_2_MVVM.git/mocks/application/official_company"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func beforeTest() (*mocks.OfficialCompanyRepository, CompanyUseCase) {
	infoLogger, _ := logger.New("test", 1, os.Stdout)
	errorLogger, _ := logger.New("test", 2, os.Stderr)
	mockRepo := new(mocks.OfficialCompanyRepository)
	useCase := CompanyUseCase{
		iLog:   infoLogger,
		errLog: errorLogger,
		repos:  mockRepo,
	}
	return mockRepo, useCase
}

func TestGetOfficialCompany(t *testing.T) {
	mockRepo, useCase := beforeTest()
	compID, _ := uuid.Parse("28b49c0d-5ded-4f52-afa5-51948c05e0f5")
	emptyCompID := uuid.Nil
	comp := models.OfficialCompany{
		ID: compID,
	}
	mockRepo.On("GetOfficialCompany", compID).Return(&comp, nil)
	mockRepo.On("GetOfficialCompany", emptyCompID).Return(nil, nil)
	ansCorrect, errNil1 := useCase.GetOfficialCompany(compID)
	ansNil, errNil2 := useCase.GetOfficialCompany(emptyCompID)
	assert.Nil(t, errNil1)
	assert.Nil(t, errNil2)
	assert.Nil(t, ansNil)
	assert.Equal(t, *ansCorrect, comp)
}

func TestGetMineCompany(t *testing.T) {
	mockRepo, useCase := beforeTest()
	empID, _ := uuid.Parse("28b49c0d-5ded-4f52-afa5-51948c05e0f5")
	compID, _ := uuid.Parse("1dbd3178-84ad-4d0d-b61d-6116fda279ef")
	employer := models.Employer{
		ID:        empID,
		CompanyID: compID,
	}
	comp := models.OfficialCompany{
		ID: employer.CompanyID,
	}
	mockRepo.On("GetMineCompany", empID).Return(&comp, nil)
	mockRepo.On("GetMineCompany", compID).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.GetMineCompany(empID)
	ansWrong, errNotNil := useCase.GetMineCompany(compID)
	assert.Nil(t, errNil)
	assert.Equal(t, *ansCorrect, comp)
	assert.Nil(t, ansWrong)
	assert.Error(t, errNotNil)
}

func TestCreateOfficialCompany(t *testing.T)  {
	mockRepo, useCase := beforeTest()
	empID, _ := uuid.Parse("28b49c0d-5ded-4f52-afa5-51948c05e0f5")
	comp := models.OfficialCompany{
		Name: "name",
		Description: "description",
		AreaSearch: "area",
		Link: "link",
	}
	mockRepo.On("CreateOfficialCompany", comp, empID).Return(&comp, nil)
	mockRepo.On("CreateOfficialCompany", comp, uuid.Nil).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.CreateOfficialCompany(comp, empID)
	ansWrong, errNotNil := useCase.CreateOfficialCompany(comp, uuid.Nil)
	assert.Nil(t, errNil)
	assert.Equal(t, *ansCorrect, comp)
	assert.Nil(t, ansWrong)
	assert.Error(t, errNotNil)
}

func TestGetCompaniesList(t *testing.T)  {
	var (
		start uint
		end uint
	)
	start = 0
	end = 3
	mockRepo, useCase := beforeTest()
	compList := []models.OfficialCompany{
		{Name: "name1", Description: "description1", Link: "link1", AreaSearch: "area1"},
		{Name: "name2", Description: "description2", Link: "link2", AreaSearch: "area2"},
		{Name: "name3", Description: "description3", Link: "link3", AreaSearch: "area3"},
	}

	mockRepo.On("GetCompaniesList", start, end).Return(compList, nil)
	mockRepo.On("GetCompaniesList", start, start).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.GetCompaniesList(start, end)
	ansWrong, errNotNil := useCase.GetCompaniesList(start, start)
	assert.Nil(t, errNil)
	assert.Equal(t, ansCorrect, compList)
	assert.Nil(t, ansWrong)
	assert.Error(t, errNotNil)
}

func TestSearchCompanies(t *testing.T) {
	compList := []models.OfficialCompany{
		{Name: "name1", Description: "description1", Link: "link1", AreaSearch: "area1"},
		{Name: "Name", Description: "description2", Link: "link2", AreaSearch: "area1"},
		{Name: "NAME", Description: "description3", Link: "link3", AreaSearch: "area2"},
	}

	mockRepo, useCase := beforeTest()
	params := models.CompanySearchParams{
		AreaSearch: []string{"area1"},
	}
	mockRepo.On("SearchCompanies", params).Return(compList[1:], nil)
	ansCorrect, errNil := useCase.SearchCompanies(params)
	assert.Nil(t, errNil)
	assert.Equal(t, ansCorrect, compList[1:])
}