package repository

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/go-pg/pg/v9"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewPgRepository(db *pg.DB) user.RepositoryUser {
	return &pgStorage{db: db}
}

type pgStorage struct {
	db *pg.DB
}


func (p *pgStorage) Login(user models.UserLogin) (*models.User, error) {
	var userDB models.User
	err := p.db.Model(&userDB).
		Where("email = ?", user.Email).
		Where("nickname = ?", user.Nickname).
		Select()
	if err != nil {
		return nil, err
	}
	// compare password with the hashed one
	err = bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password))
	if err != nil {
		return nil, err
	}
	return &userDB, nil
}

func (p *pgStorage) GetEmployerByID(id string) (*models.Employer, error) {
	var employer models.Employer
	err := p.db.Model(&employer).Where("user_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select employer with id: %s : error: %w", id, err)
		return nil, err
	}
	return &employer, nil
}

func (p *pgStorage) UpdateEmployer(employerNew models.Employer) (*models.Employer, error) {
	if employerNew.CompanyID != uuid.Nil {
		company := models.OfficialCompany{ID: employerNew.CompanyID}
		if isExist, err := p.db.Model(&company).Exists(); err != nil {
			return nil, err
		} else if !isExist {
			return nil, errors.New("new company doesn't exist")
		}
		_, err := p.db.Model(&employerNew).WherePK().Returning("*").Update()
		if err != nil {
			err = fmt.Errorf("error in updating employer with id %s, : %w", employerNew.ID.String(), err)
			return nil, err
		}
	}
	return &employerNew, nil
}
func (p *pgStorage) GetCandidateByID(id string) (*models.Candidate, error) {
	var candidate models.Candidate
	err := p.db.Model(&candidate).Where("user_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select user with id: %s : error: %w", id, err)
		return nil, err
	}
	return &candidate, nil
}

func (p *pgStorage) GetUserByID(id string) (*models.User, error) {
	var newUser models.User
	err := p.db.Model(&newUser).Where("user_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select user with id: %s : error: %w", id, err)
		return nil, err
	}
	return &newUser, nil
}

func (p *pgStorage) GetCandByID(id string) (*models.User, error) {
	var cand models.Candidate
	err := p.db.Model(&cand).Where("cand_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select candidate with id: %s : error: %w", id, err)
		return nil, err
	}
	return p.GetUserByID(cand.UserID.String())
}

func (p *pgStorage) GetEmplByID(id string) (*models.User, error) {
	var empl models.Employer
	err := p.db.Model(&empl).Where("empl_id = ?", id).Select()
	if err != nil {
		err = fmt.Errorf("error in select employer with id: %s : error: %w", id, err)
		return nil, err
	}
	return p.GetUserByID(empl.UserID.String())
}

func (p *pgStorage) CreateUser(user models.User) (*models.User, error) {
	_, errInsert := p.db.Model(&user).Returning("*").Insert()
	if errInsert != nil {
		if isExist, err := p.db.Model(&user).Exists(); err != nil {
			errInsert = fmt.Errorf("error in inserting user with name: %s : error: %w", user.Name, err)
		} else if isExist {
			errInsert = errors.New("user already exists")
		}
		return nil, errInsert
	}
	if user.UserType == "employer" {
		newEmpl := models.Employer{UserID: user.ID}
		_, errInsert = p.db.Model(&newEmpl).Returning("*").Insert()
	} else if user.UserType == "candidate" {
		newCand := models.Candidate{UserID: user.ID}
		_, errInsert = p.db.Model(&newCand).Returning("*").Insert()
	}
	if errInsert != nil {
		return nil, errInsert
	}
	return &user, nil
}

func (p *pgStorage) UpdateUser(userNew models.User) (*models.User, error) {
	_, err := p.db.Model(&userNew).WherePK().Returning("*").Update()
	if err != nil {
		err = fmt.Errorf("error in updating user with id %s, : %w", userNew.ID.String(), err)
		return nil, err
	}
	return &userNew, nil
}
