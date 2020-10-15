package user

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
)

type IUseCaseUser interface {
	GetUserByID(id string) (*models.User, error)
	CreateUser(user models.User) (*models.User, error)
	UpdateUser(user_id uuid.UUID, newPassword, oldPassword, nick, name, surname, email, phone,
			areaSearch string, socialNetwork []string) (*models.User, error)
}
