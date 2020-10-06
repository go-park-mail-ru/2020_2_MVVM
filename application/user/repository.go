package user

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type RepositoryUser interface {
	GetUserByID(id string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
}
