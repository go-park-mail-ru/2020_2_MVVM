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

func (u *UseCaseChat) CreateChatAndTechMes(response models.Response) (*models.Chat, error) {
	return u.strg.CreateChatAndTechChat(response)
}

func (u *UseCaseChat) CreateTechMesToUpdate(response models.Response) (*models.Chat, error) {
	return u.strg.CreateTechMesToUpdate(response)
}

func (u *UseCaseChat) GetChatHistory(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) (models.ChatHistory, error) {
	history := models.ChatHistory{ChatID: chatID}

	// load messages
	messages, err := u.strg.MessagesForChat(chatID, from, to, offset, limit)
	if err != nil {
		return models.ChatHistory{}, err
	}
	if err = u.strg.MarkMessagesAsRead(chatID, utype, from, to, offset, limit); err != nil {
		return models.ChatHistory{}, err
	}

	// load technical messages
	techMessages, err := u.strg.TechnicalMessagesForChat(chatID, from, to, offset, limit)
	if err != nil {
		return models.ChatHistory{}, err
	}
	if err = u.strg.MarkTechnicalMessagesAsRead(chatID, utype, from, to, offset, limit); err != nil {
		return models.ChatHistory{}, err
	}

	history.TechnicalMessages = *techMessages
	history.Dialog = *messages
	return history, err
}

func (u *UseCaseChat) CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error) {
	mes.DateCreate = time.Now()
	return u.strg.CreateMessage(mes, sender)
}

func (u *UseCaseChat) ListChats(userID uuid.UUID, userType string) ([]models.ChatSummary, error) {
	return u.strg.ListChats(userID, userType)
}