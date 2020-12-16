package chat

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
)

type IUseCaseChat interface {
	Create(response models.Response) (*models.Chat, error)
	GetByID(chatID uuid.UUID, start uint, limit uint) ([]models.Message, error)
	CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error)
	ListChats(userID uuid.UUID, userType string) ([]models.BriefChat, error)
}
