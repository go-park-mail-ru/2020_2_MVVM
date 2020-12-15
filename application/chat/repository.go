package chat

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
)

type ChatRepository interface {
	Create(response models.Response) (*models.Chat, error)
	GetById(chatID uuid.UUID, start uint, limit uint) ([]models.Message, error)
	CreateMessage(mes models.Message, sender uuid.UUID) (*models.Message, error)
}
