package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func beforeTest(t *testing.T) (*mocks.RepositoryUser, UserUseCase) {
	infoLogger, _ := logger.New("test", 1, os.Stdout)
	errorLogger, _ := logger.New("test", 2, os.Stderr)
	mockRepo := new(mocks.RepositoryUser)
	usecase := UserUseCase{
		iLog:   infoLogger,
		errLog: errorLogger,
		repos:  mockRepo,
	}
	return mockRepo, usecase
}

func TestUserGetById(t *testing.T) {
	mockRepo, usecase := beforeTest(t)
	userID, _ := uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")

	user := models.User{
		ID: userID,
	}

	mockRepo.On("GetUserByID", userID.String()).Return(&user, nil)
	answer, err := usecase.GetUserByID(userID.String())

	assert.Nil(t, err)
	assert.Equal(t, *answer, user)
}

func TestUserCreateUser(t *testing.T) {
	mockRepo, usecase := beforeTest(t)

	user := models.User{
		UserType: "candidate",
		Name:     "name",
		Surname:  "surname",
		Email:    "email@mail.ru",
	}

	mockRepo.On("CreateUser", user).Return(&user, nil)
	answer, err := usecase.CreateUser(user)

	assert.Nil(t, err)
	assert.Equal(t, *answer, user)
}

func TestUserUpdateUser(t *testing.T) {
	mockRepo, usecase := beforeTest(t)
	userID, _ := uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
	user := models.User{
		ID:       userID,
		UserType: "candidate",
		Name:     "newName",
		Surname:  "newSurname",
		Email:    "email@mail.ru",
	}
	mockRepo.On("UpdateUser", user).Return(&user, nil)
	mockRepo.On("GetUserByID", userID.String()).Return(&user, nil)
	answer, err := usecase.UpdateUser(user.ID.String(), "", "", "newName",
		"newSurname", "", "", "")

	assert.Nil(t, err)
	assert.Equal(t, *answer, user)
}

func TestUserLogin(t *testing.T) {
	mockRepo, usecase := beforeTest(t)
	userID, _ := uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
	user := models.User{
		ID:       userID,
		UserType: "candidate",
		Name:     "newName",
		Surname:  "newSurname",
		Email:    "email@mail.ru",
	}

	userLogin := models.UserLogin{
		Email:    "email@mail.ru",
		Password: "123",
	}
	mockRepo.On("Login", userLogin).Return(&user, nil)
	answer, err := usecase.Login(userLogin)

	assert.Nil(t, err)
	assert.Equal(t, *answer, user)
}

func TestUserGetEmployerByID(t *testing.T) {
	mockRepo, usecase := beforeTest(t)
	userID, _ := uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
	user := models.Employer{
		UserID:    userID,
		Favorites: nil,
	}
	mockRepo.On("GetEmployerByID", userID.String()).Return(&user, nil)
	answer, err := usecase.GetEmployerByID(userID.String())

	assert.Nil(t, err)
	assert.Equal(t, *answer, user)
}

func TestUserGetCandidateByID(t *testing.T) {
	mockRepo, usecase := beforeTest(t)
	userID, _ := uuid.Parse("77b2e989-6be6-4db5-a657-f25487638af9")
	user := models.Candidate{
		UserID:    userID,
	}
	mockRepo.On("GetCandidateByID", userID.String()).Return(&user, nil)
	answer, err := usecase.GetCandidateByID(userID.String())

	assert.Nil(t, err)
	assert.Equal(t, *answer, user)
}