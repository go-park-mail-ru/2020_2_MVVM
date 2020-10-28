package user

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
)

type IUseCaseUser interface {
	//Login(user models.User) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetEmployerByID(id string) (*models.Employer, error)
	UpdateEmployer(employer models.Employer) (*models.Employer, error)
	GetCandidateByID(id string) (*models.Candidate, error)
	CreateUser(user models.User) (*models.User, error)
	UpdateUser(user_id string, newPassword, oldPassword, nick, name, surname, email, phone,
			   socialNetwork string) (*models.User, error)
	Login(user models.UserLogin) (*models.User, error)
}
