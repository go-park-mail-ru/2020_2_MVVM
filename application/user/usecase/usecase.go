package usecase

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	logger "github.com/rowdyroad/go-simple-logger"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	iLog   *logger.Logger
	errLog *logger.Logger
	repos  user.RepositoryUser
}

func NewUserUseCase(iLog *logger.Logger, errLog *logger.Logger,
	repos user.RepositoryUser) *UserUseCase {
	return &UserUseCase{
		iLog:   iLog,
		errLog: errLog,
		repos:  repos,
	}
}

func (U *UserUseCase) Login(user models.UserLogin) (*models.User, error) {
	return U.repos.Login(user)
}

func (U *UserUseCase) GetUserByID(id string) (*models.User, error) {
	userById, err := U.repos.GetUserByID(id)
	if err != nil {
		err = fmt.Errorf("error in user get by id func : %w", err)
		return nil, err
	}
	return userById, nil
}

func (U *UserUseCase) CreateUser(user models.User) (*models.User, error) {
	userNew, err := U.repos.CreateUser(user)
	if err != nil {
		if err.Error() != "user already exists" {
			err = fmt.Errorf("error in user get by id func : %w", err)
		}
		return nil, err
	}
	return userNew, nil
}

func (U *UserUseCase) UpdateUser(user_id string, newPassword, oldPassword, nick, name, surname, email, phone,
								areaSearch, socialNetwork string) (*models.User, error) {
	user, err := U.GetUserByID(user_id)
	if err != nil {
		err = fmt.Errorf("error get user with id %s : %w", user_id, err)
		return nil, err
	}

	if nick != "" {
		user.Nickname = nick
	}
	if name != "" {
		user.Name = name
	}
	if surname != "" {
		user.Surname = surname
	}
	if email != "" {
		user.Email = email
	}
	if phone != "" {
		user.Phone = &phone
	}
	if areaSearch != "" {
		user.AreaSearch = &areaSearch
	}
	if socialNetwork != "" {
		user.SocialNetwork = &socialNetwork
	}
	if oldPassword != "" {
		isEqual := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(oldPassword))
		if isEqual != nil {
			return nil, common.ErrInvalidUpdatePassword
		}
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			err = fmt.Errorf("error in crypting password : %W", err)
			return nil, err
		}
		user.PasswordHash = passwordHash
	}

	newUser, err := U.repos.UpdateUser(*user)
	if err != nil {
		err = fmt.Errorf("error in updating user with id = %s : %w", user.ID.String(), err)
		return nil, err
	}

	return newUser, nil
}
