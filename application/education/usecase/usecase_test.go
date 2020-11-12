package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/education"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func beforeTest() (*mocks.Repository, UseCaseEducation) {
	infoLogger, _ := logger.New("test", 1, os.Stdout)
	errorLogger, _ := logger.New("test", 2, os.Stderr)
	mockRepo := new(mocks.Repository)
	usecase := UseCaseEducation{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		strg:        mockRepo,
	}
	return mockRepo, usecase
}

var ID, _ = uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
var testEducation = models.Education{
	EdId:        ID,
	CandID:      ID,
	ResumeId:    ID,
	University:  "",
}

func TestEducationCreate(t *testing.T) {
	mockRepo, usecase := beforeTest()

	mockRepo.On("Create", testEducation).Return(&testEducation, nil)
	answer, err := usecase.Create(testEducation)
	assert.Nil(t, err)
	assert.Equal(t, *answer, testEducation)
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