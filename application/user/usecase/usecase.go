package usecase

import (
	"errors"
	"fmt"
	logger "github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
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

func (u *UserUseCase) Login(user models.UserLogin) (*models.User, error) {
	return u.repos.Login(user)
}

func (u *UserUseCase) GetUserByID(id string) (*models.User, error) {
	userById, err := u.repos.GetUserByID(id)
	if err != nil {
		err = fmt.Errorf("error in user get by id func : %w", err)
		return nil, err
	}
	return userById, nil
}

func (u *UserUseCase) GetCandByID(id string) (*models.User, error) {
	return u.repos.GetCandByID(id)
}

func (u *UserUseCase) GetEmplByID(id string) (*models.User, error) {
	return u.repos.GetEmplByID(id)
}

func (u *UserUseCase) GetCandidateByID(id string) (*models.Candidate, error) {
	candById, err := u.repos.GetCandidateByID(id)
	if err != nil {
		err = fmt.Errorf("error in cand get by id func : %w", err)
		return nil, err
	}
	return candById, nil
}

func (u *UserUseCase) GetEmployerByID(id string) (*models.Employer, error) {
	emplById, err := u.repos.GetEmployerByID(id)
	if err != nil {
		err = fmt.Errorf("error in empl get by id func : %w", err)
		return nil, err
	}
	return emplById, nil
}

func (u *UserUseCase) CreateUser(user models.User) (*models.User, error) {
	userNew, err := u.repos.CreateUser(user)
	if err != nil {
		if err.Error() != common.UserExistErr {
			err = fmt.Errorf("error in user get by id func : %w", err)
		}
		return nil, err
	}
	return userNew, nil
}

func (u *UserUseCase) UpdateUser(userNew models.User) (*models.User, error) {
	userOld, err := u.GetUserByID(userNew.ID.String())
	if err != nil {
		err = fmt.Errorf("error get user with id %s : %w", userNew.ID, err)
		return nil, err
	}
	if userNew.Name != "" {
		userOld.Name = userNew.Name
	}
	if userNew.Surname != "" {
		userOld.Surname = userNew.Surname
	}
	if userNew.Email != "" {
		userOld.Email = userNew.Email
	}
	if userNew.Phone != nil && *userNew.Phone != "" {
		userOld.Phone = userNew.Phone
	}
	if userNew.SocialNetwork != nil && *userNew.SocialNetwork != "" {
		userOld.SocialNetwork = userNew.SocialNetwork
	}
	if userNew.PasswordHash != nil {
		isEqual := bcrypt.CompareHashAndPassword(userOld.PasswordHash, userNew.PasswordHash)
		if isEqual != nil {
			return nil, errors.New(common.WrongPasswd)
		}
		userOld.PasswordHash = userNew.PasswordHash
	}
	newUser, err := u.repos.UpdateUser(*userOld)
	if err != nil {
		if err.Error() != common.UserExistErr {
			err = fmt.Errorf("error in updating user with id = %s : %w", userOld.ID.String(), err)
		}
		return nil, err
	}
	return newUser, nil
}
