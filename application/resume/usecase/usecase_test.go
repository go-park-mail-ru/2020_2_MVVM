package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	resume2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/resume"
	mExperience "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/custom_experience"
	mEducation "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/education"
	mResume "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/resume"
	mUser "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/user"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
	User:   testUser,
}
var testResume = models.Resume{
	ResumeID:  ID,
	CandID:    ID,
	Candidate: candidate,
	Title:     "ID",
}
var educationTest = models.Education{
	EdId:       ID,
	CandID:     ID,
	ResumeId:   ID,
	University: "university",
}
var flag = false
var experienceTest = models.ExperienceCustomComp{
	ID:              ID,
	CandID:          ID,
	ResumeID:        ID,
	NameJob:         "name job",
	ContinueToToday: &flag,
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
	Resume:     testResume,
}

func beforeTest(t *testing.T) (*mResume.Repository, *mUser.UseCase, *mEducation.UseCase, *mExperience.UseCase, ResumeUseCase) {
	infoLogger, _ := logger.New(os.Stdout)
	errorLogger, _ := logger.New(os.Stderr)
	mockRepo := new(mResume.Repository)
	mockEducationUS := new(mEducation.UseCase)
	mockExperienceUS := new(mExperience.UseCase)
	mockUserUS := new(mUser.UseCase)
	usecase := ResumeUseCase{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
		userUseCase: mockUserUS,
		//educationUseCase: mockEducationUS,
		customExpUseCase: mockExperienceUS,
		strg:             mockRepo,
	}
	return mockRepo, mockUserUS, mockEducationUS, mockExperienceUS, usecase
}

func TestResumeGetById(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	mockRepo.On("GetById", ID).Return(&testResume, nil)
	mockRepo.On("GetById", uuid.Nil).Return(nil, assert.AnError)
	answerCorrect, err := usecase.GetById(ID)
	answerNotCorrect, errNotNil := usecase.GetById(uuid.Nil)

	assert.Nil(t, err)
	assert.Equal(t, *answerCorrect, testResume)
	assert.Nil(t, answerNotCorrect)
	assert.Error(t, errNotNil)
}

func TestResumeCreateResume(t *testing.T) {
	mockRepo, _, mockEducationUS, mockExperienceUS, usecase := beforeTest(t)

	testResume.ExperienceCustomComp = []models.ExperienceCustomComp{experienceTest}
	testResume.Education = []models.Education{educationTest}
	testResume.DateCreate = time.Now().String()
	mockRepo.On("Create", mock.Anything).Return(&testResume, nil)
	mockExperienceUS.On("Create", experienceTest).Return(&experienceTest, nil)
	mockEducationUS.On("Create", educationTest).Return(&educationTest, nil)
	answer, err := usecase.Create(testResume)

	assert.Nil(t, err)
	assert.Equal(t, *answer, testResume)
}

func TestResumeUpdateUser(t *testing.T) {
	mockRepo, _, _, expUsecase, usecase := beforeTest(t)

	testResume.ExperienceCustomComp = []models.ExperienceCustomComp{experienceTest}
	testResume.Education = []models.Education{educationTest}
	mockRepo.On("Update", mock.Anything).Return(&testResume, nil).Once()
	expUsecase.On("DropAllFromResume", mock.Anything).Return(nil).Once()
	mockRepo.On("GetById", ID).Return(&testResume, nil).Once()
	answer, err := usecase.Update(testResume)

	assert.Nil(t, err)
	assert.Equal(t, *answer, testResume)

	mockRepo.On("GetById", uuid.Nil).Return(nil, assert.AnError).Once()
	testResume.ResumeID = uuid.Nil
	answerWrong, errNotNil := usecase.Update(testResume)
	assert.Nil(t, answerWrong)
	assert.Error(t, errNotNil)

	testResume.ResumeID = ID
	mockRepo.On("GetById", ID).Return(&testResume, nil).Twice()
	newTestResume := testResume
	newTestResume.CandID = uuid.Nil
	answerWrong2, errNotNil2 := usecase.Update(newTestResume)
	assert.Nil(t, answerWrong2)
	assert.Error(t, errNotNil2)

	expUsecase.On("DropAllFromResume", mock.Anything).Return(assert.AnError)
	answer, err = usecase.Update(testResume)
	assert.Nil(t, answer)
	assert.Error(t, err)
}

func TestResumeSearch(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	word := "word"
	params := resume2.SearchParams{
		KeyWords: &word,
	}

	var listResume = []models.Resume{testResume}
	var listBrief = []models.BriefResumeInfo{briefResume}
	mockRepo.On("Search", &params).Return(listResume, nil)
	answer, err := usecase.Search(params)

	assert.Nil(t, err)
	assert.Equal(t, answer, listBrief)
}

func TestResumeList(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	var listResume = []models.Resume{testResume}
	var listBrief = []models.BriefResumeInfo{briefResume}
	mockRepo.On("List", uint(1), uint(1)).Return(listResume, nil).Once()
	answer, err := usecase.List(uint(1), uint(1))
	assert.Nil(t, err)
	assert.Equal(t, answer, listBrief)

	mockRepo.On("List", uint(1), uint(1)).Return(nil, errors.Errorf("")).Once()
	answerWrong, errNotNil := usecase.List(uint(1), uint(1))
	assert.Nil(t, answerWrong)
	assert.Error(t, errNotNil)

	answerWrong2, errNotNil2 := usecase.List(uint(1), uint(1000))
	assert.Nil(t, answerWrong2)
	assert.Error(t, errNotNil2)
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

	favoriteId := models.FavoriteID{}
	mockRepo.On("AddFavorite", favorite).Return(&favoriteId, nil)
	answer, err := usecase.AddFavorite(favorite)

	assert.Nil(t, err)
	assert.Equal(t, *answer, favoriteId)
}

//GetFavoriteByResume(userID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error)
func TestResumeGetFavoriteByResume(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	mockRepo.On("GetFavoriteForResume", ID, ID).Return(&favorite, nil)
	answer, err := usecase.GetFavoriteByResume(ID, ID)

	assert.Nil(t, err)
	assert.Equal(t, *answer, favorite)
}

//GetFavoriteByID(favoriteID uuid.UUID) (*models.FavoritesForEmpl, error)
func TestResumeGetFavoriteByID(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	mockRepo.On("GetFavoriteByID", ID).Return(&favorite, nil)
	answer, err := usecase.GetFavoriteByID(ID)

	assert.Nil(t, err)
	assert.Equal(t, *answer, favorite)
}

//RemoveFavorite(favorite models.FavoritesForEmpl) error
func TestResumeRemoveFavorite(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	mockRepo.On("GetFavoriteByID", ID).Return(&favorite, nil)
	mockRepo.On("RemoveFavorite", ID).Return(nil)
	err := usecase.RemoveFavorite(favorite)

	assert.Nil(t, err)
}

func TestResumeGetAllEmplFavoriteResume(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)

	var listResume = []models.Resume{testResume}
	var listBrief = []models.BriefResumeInfo{briefResume}
	mockRepo.On("GetAllEmplFavoriteResume", ID).Return(listResume, nil)
	answer, err := usecase.GetAllEmplFavoriteResume(ID)

	assert.Nil(t, err)
	assert.Equal(t, answer, listBrief)
}

func TestDeleteResume(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)
	id := uuid.New()
	mockRepo.On("Delete", id, id).Return(nil)
	err := usecase.DeleteResume(id, id)
	assert.Nil(t, err)
}

func TestMakePdf(t *testing.T) {
	mockRepo, _, _, _, usecase := beforeTest(t)
	id := uuid.New()
	res := models.Resume{}

	mockRepo.On("GetByIdWithCand", id).Return(&res, nil).Once()
	mockRepo.On("GetByIdWithCand", id).Return(nil, assert.AnError).Once()
	err1 := usecase.MakePdf(id)
	err2 := usecase.MakePdf(id)
	assert.Error(t, err1)
	assert.Error(t, err2)
}

func TestNewUseCase(t *testing.T) {
	useCase := NewUseCase(nil, nil, nil, nil, nil)
	assert.Equal(t, useCase, &ResumeUseCase{nil, nil, nil, nil, nil})
}
