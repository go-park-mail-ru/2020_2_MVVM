package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
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

var ID, _ = uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
var testUser = models.User{
	ID:       ID,
	UserType: "candidate",
	Name:     "ID",
	Surname:  "ID",
	Email:    "ID",
}
var candidate = models.Candidate{
	ID:     ID,
	UserID: ID,
	User:   &testUser,
}
var testResume = models.Resume{
	ResumeID:  ID,
	CandID:    ID,
	Candidate: &candidate,
	Title:     "ID",
}
var briefResume = models.BriefResumeInfo{
	ResumeID: ID,
	CandID:   ID,
	UserID:   ID,
	Title:    "ID",
	Name:     "ID",
	Surname:  "ID",
	Email:    "ID",
}

var favorite = models.FavoritesForEmpl{
	FavoriteID: ID,
	EmplID:     ID,
	ResumeID:   ID,
	Resume:     &testResume,
}

func beforeTest(t *testing.T) (*mResume.Repository, *mUser.UseCase, *mEducation.UseCase, *mExperience.UseCase, ResumeUseCase) {
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

	mockRepo.On("GetById", ID).Return(&testResume, nil)
	answer, err := usecase.GetById(ID)

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
	mockRepo.On("GetById", ID).Return(&testResume, nil)
	mockEducationUS.On("DropAllFromResume", ID).Return(nil)
	mockExperienceUS.On("DropAllFromResume", ID).Return(nil)
	answer, err := usecase.Update(testResume)

	assert.Nil(t, err)
	assert.Equal(t, *answer, testResume)
}

func TestResumeSearch(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	word := "word"
	params := resume.SearchParams{
		KeyWords: &word,
	}

	var listResume = []models.Resume{testResume}
	var listBrief = []models.BriefResumeInfo{briefResume}
	mockRepo.On("Search", &params).Return(listResume, nil)
	answer, err := usecase.Search(params)

	assert.Nil(t, err)
	assert.Equal(t, answer, listBrief)
}

//List(start, limit uint) ([]models.BriefResumeInfo, error)
func TestResumeList(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	var listResume = []models.Resume{testResume}
	var listBrief = []models.BriefResumeInfo{briefResume}
	mockRepo.On("List", uint(1), uint(1)).Return(listResume, nil)
	answer, err := usecase.List(uint(1), uint(1))

	assert.Nil(t, err)
	assert.Equal(t, answer, listBrief)
}

//GetAllUserResume(userid uuid.UUID) ([]models.BriefResumeInfo, error)
func TestResumeGetAllUserResume(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	var listResume = []models.Resume{testResume}
	var listBrief = []models.BriefResumeInfo{briefResume}
	mockRepo.On("GetAllUserResume", ID).Return(listResume, nil)
	answer, err := usecase.GetAllUserResume(ID)

	assert.Nil(t, err)
	assert.Equal(t, answer, listBrief)
}

//AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error)
func TestResumeAddFavorite(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	mockRepo.On("AddFavorite", favorite).Return(&favorite, nil)
	answer, err := usecase.AddFavorite(favorite)

	assert.Nil(t, err)
	assert.Equal(t, *answer, favorite)
}

//GetFavoriteByResume(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error)
func TestResumeGetFavoriteByResume(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)


	mockRepo.On("GetFavoriteByResume", ID, ID).Return(&favorite, nil)
	answer, err := usecase.GetFavoriteByResume(ID, ID)

	assert.Nil(t, err)
	assert.Equal(t, *answer, favorite)
}