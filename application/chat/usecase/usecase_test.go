package usecase

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	mocks "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/chat"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func beforeTest() (*mocks.ChatRepository, UseCaseChat) {
	mockRepo := new(mocks.ChatRepository)
	useCase := UseCaseChat{
		strg: mockRepo,
	}
	return mockRepo, useCase
}

func TestCreateChatAndTechMes(t *testing.T) {
	repo, usecase := beforeTest()
	resp := models.Response{}
	repo.On("CreateChatAndTechMes", resp).Return(nil, nil)
	res, err := usecase.CreateChatAndTechMes(resp)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestCreateTechMesToUpdate(t *testing.T) {
	repo, usecase := beforeTest()
	resp := models.Response{}
	repo.On("CreateTechMesToUpdate", resp).Return(nil, nil)
	res, err := usecase.CreateTechMesToUpdate(resp)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestCreateMessage(t *testing.T) {
	repo, usecase := beforeTest()
	mes := models.Message{}
	id := uuid.New()
	repo.On("CreateMessage", mock.Anything, id).Return(nil, nil)
	res, err := usecase.CreateMessage(mes, id)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestListChats(t *testing.T) {
	repo, usecase := beforeTest()
	id := uuid.New()
	repo.On("ListChats", id, "test").Return(nil, nil)
	res, err := usecase.ListChats(id, "test")
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestGetTotalUnreadMes(t *testing.T) {
	repo, usecase := beforeTest()
	id := uuid.New()
	repo.On("GetTotalUnreadMes", id, "test").Return(nil, nil)
	res, err := usecase.GetTotalUnreadMes(id, "test")
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestGetChatHistory(t *testing.T) {
	repo, usecase := beforeTest()
	id := uuid.New()
	utype := common.Candidate
	start := time.Now()
	end := start
	var offset uint =  1
	var limit uint = 1
	var messages []models.MessageBrief
	var techMess []models.TechMessageBrief
	repo.On("MessagesForChat", id, &start, &end, &offset, &limit).Return(&messages, nil)
	repo.On("MarkMessagesAsRead", id, utype, &start, &end, &offset, &limit).Return(nil)
	repo.On("TechnicalMessagesForChat", id, &start, &end, &offset, &limit).Return(&techMess, nil)
	repo.On("MarkTechnicalMessagesAsRead", id, utype, &start, &end, &offset, &limit).Return(nil)
	res, err := usecase.GetChatHistory(id, utype, &start, &end, &offset, &limit)
	assert.Equal(t, models.ChatHistory{ChatID:id}, res)
	assert.Nil(t, err)
}
