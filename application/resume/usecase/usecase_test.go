package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	mExperience "github.com/go-park-mail-ru/2020_2_MVVM.git/mocks/application/custom_experience"
	mEducation "github.com/go-park-mail-ru/2020_2_MVVM.git/mocks/application/education"
	mResume "github.com/go-park-mail-ru/2020_2_MVVM.git/mocks/application/resume"
	mUser "github.com/go-park-mail-ru/2020_2_MVVM.git/mocks/application/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
	"time"
)

var resumeID, _ = uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
var testResume = models.Resume{
	ResumeID:             resumeID,
}


func beforeTest(t *testing.T) (*mResume.Repository, *mUser.UseCase, *mEducation.UseCase, *mExperience.UseCase,  ResumeUseCase) {
	infoLogger, _ := logger.New(os.Stdout)
	errorLogger, _ := logger.New(os.Stderr)
	mockRepo := new(mResume.Repository)
	mockEducationUS := new(mEducation.UseCase)
	mockExperienceUS := new(mExperience.UseCase)
	mockUserUS := new(mUser.UseCase)
	usecase := ResumeUseCase{
		infoLogger:       infoLogger,
		errorLogger:      errorLogger,
		userUseCase:      mockUserUS,
		educationUseCase: mockEducationUS,
		customExpUseCase: mockExperienceUS,
		strg:             mockRepo,
	}
	return mockRepo, mockUserUS, mockEducationUS, mockExperienceUS, usecase
}

func TestResumeGetById(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	mockRepo.On("GetById", resumeID).Return(&testResume, nil)
	answer, err := usecase.GetById(resumeID)

	assert.Nil(t, err)
	assert.Equal(t, *answer, testResume)
}

func TestResumeCreateResume(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	testResume.DateCreate = time.Now()
	mockRepo.On("Create", mock.Anything).Return(&testResume, nil)
	answer, err := usecase.Create(testResume)

	assert.Nil(t, err)
	assert.Equal(t, *answer, testResume)
}

func TestResumeUpdateUser(t *testing.T) {
	mockRepo, _, mockEducationUS, mockExperienceUS, usecase := beforeTest(t)

	mockRepo.On("Update", testResume).Return(&testResume, nil)
	mockRepo.On("GetById", resumeID).Return(&testResume, nil)
	mockEducationUS.On("DropAllFromResume", resumeID).Return(nil)
	mockExperienceUS.On("DropAllFromResume", resumeID).Return(nil)
	answer, err := usecase.Update(testResume)

	assert.Nil(t, err)
	assert.Equal(t, *answer, testResume)
}

//Search(searchParams SearchParams) ([]models.BriefResumeInfo, error)
func TestResumeSearch(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	mockRepo.On("GetById", resumeID).Return(&testResume, nil)
	answer, err := usecase.GetById(resumeID)

	assert.Nil(t, err)
	assert.Equal(t, *answer, testResume)
}
