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

/*
func TestUserCreateUser(t *testing.T) {
	mockRepo, useCase := beforeTest(t)

	comp := models.comp{
		UserType: "candidate",
		Name:     "name",
		Surname:  "surname",
		Email:    "email@mail.ru",
	}

	mockRepo.On("CreateUser", comp).Return(&comp, nil)
	answer, err := useCase.CreateUser(comp)

	assert.Nil(t, err)
	assert.Equal(t, *answer, comp)
}

func TestUserUpdateUser(t *testing.T) {
	mockRepo, useCase := beforeTest(t)
	compID, _ := uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
	comp := models.comp{
		ID:       compID,
		UserType: "candidate",
		Name:     "newName",
		Surname:  "newSurname",
		Email:    "email@mail.ru",
	}
	mockRepo.On("UpdateUser", comp).Return(&comp, nil)
	mockRepo.On("GetUserByID", compID.String()).Return(&comp, nil)
	answer, err := useCase.UpdateUser(comp.ID.String(), "", "", "newName",
		"newSurname", "", "", "")

	assert.Nil(t, err)
	assert.Equal(t, *answer, comp)
}*/
