package usecase

import (
	"github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/chat"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"time"
)

type UseCaseChat struct {
	infoLogger     *logger.Logger
	errorLogger    *logger.Logger
	strg           chat.ChatRepository
}

func NewUsecase(infoLogger *logger.Logger,
	errorLogger *logger.Logger,
	strg chat.ChatRepository,
) *UseCaseChat {
	usecase := UseCaseChat{
		infoLogger:     infoLogger,
		errorLogger:    errorLogger,
		strg:           strg,
	}
	return &usecase
}

func (u *UseCaseChat) Create(response models.Response) (*models.Chat, error) {
	return u.strg.Create(response)
}

func (u *UseCaseChat) GetByID(chatID uuid.UUID, start uint, limit uint) ([]models.Message, error) {
	if limit == 0 {
		limit = 20
	}
	return u.strg.GetById(chatID, start, limit)
}

func (u *UseCaseChat) CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error) {
	mes.DateCreate = time.Now()
	return u.strg.CreateMessage(mes, sender)
}

func (u *UseCaseChat) ListChats(userID uuid.UUID, userType string) ([]models.BriefChat, error) {
	return u.strg.ListChats(userID, userType)
}