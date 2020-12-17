package chat

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"time"
)
type ChatRepository interface {
	CreateChatAndTechChat(response models.Response) (*models.Chat, error)
	CreateTechMesToUpdate(response models.Response) (*models.Chat, error)
	MessagesForChat(chatID uuid.UUID, from *time.Time, to *time.Time, offset *uint, limit *uint) (*[]models.MessageBrief, error)
	MarkMessagesAsRead(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) error
	TechnicalMessagesForChat(chatID uuid.UUID, from *time.Time, to *time.Time, offset *uint, limit *uint) (*[]models.TechMessageBrief, error)
	MarkTechnicalMessagesAsRead(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) error
	CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error)
	ListChats(userID uuid.UUID, userType string) ([]models.ChatSummary, error)
}
