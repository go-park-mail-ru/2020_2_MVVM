package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/dto/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/custom_experience"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func beforeTest() (*mocks.CustomExperienceRepository, UseCase) {
	infoLogger, _ := logger.New("test", 1, os.Stdout)
	errorLogger, _ := logger.New("test", 2, os.Stderr)
	mockRepo := new(mocks.CustomExperienceRepository)
	usecase := UseCase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		customExperienceRepository:        mockRepo,
	}
	return mockRepo, usecase
}

var ID, _ = uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
var testExperience = models.ExperienceCustomComp{
	ID:              ID,
	CandID:          ID,
	ResumeID:        ID,
	NameJob:         "",
}

func TestEducationCreate(t *testing.T) {
	mockRepo, usecase := beforeTest()

	mockRepo.On("Create", testExperience).Return(&testExperience, nil)
	answer, err := usecase.Create(testExperience)
	assert.Nil(t, err)
	assert.Equal(t, *answer, testExperience)
}

func TestEducationDropAll(t *testing.T) {
	mockRepo, usecase := beforeTest()

	mockRepo.On("DropAllFromResume", ID).Return(nil)
	err := usecase.DropAllFromResume(ID)
	assert.Nil(t, err)

	mockRepo.On("DropAllFromResume", uuid.Nil).Return(assert.AnError)
	err2 := usecase.DropAllFromResume(uuid.Nil)
	assert.Error(t, err2)
}
