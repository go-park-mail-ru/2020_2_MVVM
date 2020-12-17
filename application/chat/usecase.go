package chat

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"time"
)

type IUseCaseChat interface {
	CreateChatAndTechMes(response models.Response) (*models.Chat, error)
	CreateTechMesToUpdate(response models.Response) (*models.Chat, error)
	GetChatHistory(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) (models.ChatHistory, error)
	CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error)
	ListChats(userID uuid.UUID, userType string) ([]models.ChatSummary, error)
}
