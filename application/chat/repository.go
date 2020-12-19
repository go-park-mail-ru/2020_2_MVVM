package chat

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"time"
)
type ChatRepository interface {
	CreateChatAndTechMes(response models.Response) (*models.Chat, error)
	CreateTechMesToUpdate(response models.Response) (*models.Chat, error)
	CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error)

	MessagesForChat(chatID uuid.UUID, from *time.Time, to *time.Time, offset *uint, limit *uint) (*[]models.MessageBrief, error)
	TechnicalMessagesForChat(chatID uuid.UUID, from *time.Time, to *time.Time, offset *uint, limit *uint) (*[]models.TechMessageBrief, error)

	MarkMessagesAsRead(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) error
	MarkTechnicalMessagesAsRead(chatID uuid.UUID, utype string, from *time.Time, to *time.Time, offset *uint, limit *uint) error

	ListChats(userID uuid.UUID, userType string) ([]models.ChatSummary, error)
	GetTotalUnreadMes(userID uuid.UUID, userType string) (*uint, error)

	OnlyUnreadMessagesForChat(chatID uuid.UUID, userType string) (*[]models.MessageBrief, error)
	OnlyUnreadTechnicalMessagesForChat(chatID uuid.UUID, userType string) (*[]models.TechMessageBrief, error)
}
