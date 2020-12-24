package usecase

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	mocks "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/official_company"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func beforeTest() (*mocks.OfficialCompanyRepository, CompanyUseCase) {
	mockRepo := new(mocks.OfficialCompanyRepository)
	useCase := CompanyUseCase{
		repos: mockRepo,
	}
	return mockRepo, useCase
}

func TestGetOfficialCompany(t *testing.T) {
	mockRepo, useCase := beforeTest()
	compID := uuid.New()
	emptyCompID := uuid.Nil
	comp := models.OfficialCompany{
		ID: compID,
	}
	mockRepo.On("GetOfficialCompany", compID).Return(&comp, nil).Once()
	ansCorrect, errNil1 := useCase.GetOfficialCompany(compID)
	assert.Nil(t, errNil1)
	assert.Equal(t, *ansCorrect, comp)

	mockRepo.On("GetOfficialCompany", emptyCompID).Return(nil, assert.AnError)
	ansNil, errNil2 := useCase.GetOfficialCompany(emptyCompID)
	assert.Error(t, errNil2)
	assert.Nil(t, ansNil)
}

func TestGetMineCompany(t *testing.T) {
	mockRepo, useCase := beforeTest()
	empID := uuid.New()
	compID := uuid.New()
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

func TestCreateOfficialCompany(t *testing.T) {
	mockRepo, useCase := beforeTest()
	empID := uuid.New()
	comp := models.OfficialCompany{
		Name:        "name",
		Description: "description",
		AreaSearch:  "area",
		Link:        "link",
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

func TestGetCompaniesList(t *testing.T) {
	var (
		start uint
		end   uint
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

	paramsErr := models.CompanySearchParams{}
	mockRepo.On("SearchCompanies", params).Return(compList[:2], nil)
	mockRepo.On("SearchCompanies", paramsErr).Return(nil, assert.AnError)
	ansCorrect, errNil := useCase.SearchCompanies(params)
	ansWrong, errNotNil := useCase.SearchCompanies(paramsErr)
	assert.Nil(t, errNil)
	assert.Equal(t, ansCorrect, compList[:2])
	assert.Error(t, errNotNil)
	assert.Nil(t, ansWrong)
}

func TestNewCompUseCase(t *testing.T) {
	useCase := NewCompUseCase(nil, nil, nil)
	assert.Equal(t, useCase, &CompanyUseCase{nil, nil, nil})
}

func TestDeleteOfficialCompany(t *testing.T) {
	mockRepo, useCase := beforeTest()
	ID := uuid.New()

	mockRepo.On("DeleteOfficialCompany", ID, ID).Return(nil).Once()
	errNil := useCase.DeleteOfficialCompany(ID, ID)
	assert.Nil(t, errNil)

	mockRepo.On("DeleteOfficialCompany", ID, ID).Return(assert.AnError)
	errNotNil := useCase.DeleteOfficialCompany(ID, ID)
	assert.Error(t, errNotNil)
}

func TestUpdateOfficialCompany(t *testing.T) {
	mockRepo, useCase := beforeTest()
	ID := uuid.New()
	comp := models.OfficialCompany{
		Name:        "name",
		Description: "description",
		AreaSearch:  "area",
		Link:        "link",
	}
	mockRepo.On("UpdateOfficialCompany", comp, ID).Return(&comp, nil).Once()
	ansCorrect, errNil := useCase.UpdateOfficialCompany(comp, ID)
	assert.Nil(t, errNil)
	assert.Equal(t, *ansCorrect, comp)

	mockRepo.On("UpdateOfficialCompany", comp, uuid.Nil).Return(nil, assert.AnError)
	ansWrong, errNotNil := useCase.UpdateOfficialCompany(comp, uuid.Nil)
	assert.Nil(t, ansWrong)
	assert.Error(t, errNotNil)
}

func TestGetAllCompaniesNames(t *testing.T) {
	mockRepo, useCase := beforeTest()
	compList := []models.BriefCompany{{Name: "test"}}
	mockRepo.On("GetAllCompaniesNames").Return(compList, nil).Once()
	res, err := useCase.GetAllCompaniesNames()
	assert.Equal(t, compList, res)
	assert.Nil(t, err)
}
