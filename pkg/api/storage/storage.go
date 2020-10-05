package storage

import "github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"

type Storage interface {
	NothingFunc() error
	GetUserByID(id string) (models.User, error)
	CreateUser(user models.User) (models.User, error)
}
