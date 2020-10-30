package usecase

import (
	"fmt"
	logger "github.com/apsdehal/go-logger"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
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

func (u *UserUseCase) UpdateEmployer(employerNew models.Employer) (*models.Employer, error) {
	newEmpl, err := u.repos.UpdateEmployer(employerNew)
	if err != nil {
		err = fmt.Errorf("error in updating empl with id = %s : %w", newEmpl.ID.String(), err)
		return nil, err
	}

	return newEmpl, nil
}

func (u *UserUseCase) CreateUser(user models.User) (*models.User, error) {
	userNew, err := u.repos.CreateUser(user)
	if err != nil {
		if err.Error() != "user already exists" {
			err = fmt.Errorf("error in user get by id func : %w", err)
		}
		return nil, err
	}
	return userNew, nil
}

func (u *UserUseCase) UpdateUser(user_id string, newPassword, oldPassword, nick, name, surname, email, phone,
	socialNetwork string) (*models.User, error) {
	user, err := u.GetUserByID(user_id)
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

	newUser, err := u.repos.UpdateUser(*user)
	if err != nil {
		err = fmt.Errorf("error in updating user with id = %s : %w", user.ID.String(), err)
		return nil, err
	}

	return newUser, nil
}
