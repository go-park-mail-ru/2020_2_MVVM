package user

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
)

type RepositoryUser interface {
	GetUserByID(id string) (*models.User, error)
	GetCandByID(id string) (*models.User, error)
	GetEmplByID(id string) (*models.User, error)
	GetEmployerByID(id string) (*models.Employer, error)
	GetCandidateByID(id string) (*models.Candidate, error)
	CreateUser(user models.User, companyID *uuid.UUID) (*models.User, error)
	UpdateUser(user models.User) (*models.User, error)
	DeleteUser(id uuid.UUID) error
	Login(user models.UserLogin) (*models.User, error)
}
