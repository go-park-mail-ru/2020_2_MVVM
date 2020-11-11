package repository

import (
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"testing"

	model "github.com/go-pg/pg/v9/orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/slax0rr/go-pg-wrapper/mocks"
	ormmocks "gitlab.com/slax0rr/go-pg-wrapper/mocks/orm"
)

type MockResult struct {
}

func (r MockResult) Model() model.Model {
	panic("implement me!")
}
func (r MockResult) RowsAffected() int {
	panic("implement me!")
}
func (r MockResult) RowsReturned() int {
	panic("implement me!")
}

var ID, _ = uuid.Parse("ab055db8-a31e-4927-9e9c-55441a98c429")
var passwordHash, _ = bcrypt.GenerateFromPassword([]byte("ID"), bcrypt.DefaultCost)
var testUser = models.User{
	ID:            ID,
	UserType:      "employer",
	Name:          "ID",
	Surname:       "ID",
	Email:         "ID",
	PasswordHash:  passwordHash,
	Phone:         nil,
	SocialNetwork: nil,
}

var cand = models.Candidate{
	ID:     ID,
	UserID: ID,
	User:   nil,
}

var empl = models.Employer{
	ID:        ID,
	UserID:    ID,
	CompanyID: ID,
}

func mockQueryUser() (pgStorage, *ormmocks.Query) {
	db := new(mocks.DB)
	query := new(ormmocks.Query)

	r := pgStorage{db: db}

	mockCall := db.On("Model", mock.Anything).Return(query)
	mockCall.RunFn = func(args mock.Arguments) {
		user := args[0].(*models.User)
		*user = testUser
	}
	return r, query
}

func mockQueryCandidate() (pgStorage, *ormmocks.Query) {
	db := new(mocks.DB)
	query := new(ormmocks.Query)
	r := pgStorage{db: db}

	mockCall := db.On("Model", mock.Anything).Return(query)
	mockCall.RunFn = func(args mock.Arguments) {
		user := args[0].(*models.Candidate)
		*user = cand
	}
	return r, query
}

func mockQueryEmployer() (pgStorage, *ormmocks.Query) {
	db := new(mocks.DB)
	query := new(ormmocks.Query)
	r := pgStorage{db: db}

	mockCall := db.On("Model", mock.Anything).Return(query)
	mockCall.RunFn = func(args mock.Arguments) {
		user := args[0].(*models.Employer)
		*user = empl
	}
	return r, query
}

func TestGetUserByID(t *testing.T) {
	r, query := mockQueryUser()

	query.On("Where", "user_id = ?", ID.String()).Return(query)
	query.On("Select").Return(nil)

	foo, err := r.GetUserByID(ID.String())
	assert.Nil(t, err)
	assert.Equal(t, testUser, *foo)
}

func TestLogin(t *testing.T) {
	r, query := mockQueryUser()

	var userLogin = models.UserLogin{
		Email:    "ID",
		Password: "ID",
	}
	query.On("Where", "email = ?", userLogin.Email).Return(query)
	query.On("Select").Return(nil)

	foo, err := r.Login(userLogin)
	assert.Nil(t, err)
	assert.Equal(t, testUser, *foo)
}

func TestUpdate(t *testing.T) {
	r, query := mockQueryUser()

	mockResult := MockResult{}
	query.On("WherePK").Return(query)
	query.On("Returning", "*").Return(query)
	query.On("Update").Return(mockResult, nil)

	foo, err := r.UpdateUser(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, *foo)
}

func TestGetCandidateByID(t *testing.T) {
	r, query := mockQueryCandidate()

	query.On("Where", "user_id = ?", ID.String()).Return(query)
	query.On("Select").Return(nil)

	foo, err := r.GetCandidateByID(ID.String())
	assert.Nil(t, err)
	assert.Equal(t, cand, *foo)
}

func TestGetEmployerByID(t *testing.T) {
	r, query := mockQueryEmployer()
	query.On("Where", "user_id = ?", ID.String()).Return(query)
	query.On("Select").Return(nil)

	foo, err := r.GetEmployerByID(ID.String())
	assert.Nil(t, err)
	assert.Equal(t, empl, *foo)
}

//func TestCreate(t *testing.T) {
//	repoUser, queryUser := mockQueryUser()
//
//	mockResult := MockResult{}
//	queryUser.On("Returning", "*").Return(queryUser)
//	queryUser.On("Insert").Return(mockResult, nil)
//
//	mockResult2 := MockResult{}
//	_, queryEmpl := mockQueryEmployer()
//	queryEmpl.On("Returning", "*").Return(queryEmpl)
//	queryEmpl.On("Insert").Return(mockResult2, nil)
//
//
//	foo, err := repoUser.CreateUser(testUser)
//	assert.Nil(t, err)
//	assert.Equal(t, testUser, *foo)
//}

//func (p *pgStorage) CreateUser(user models.User) (*models.User, error) {
//	_, errInsert := p.db.Model(&user).Returning("*").Insert()
//	if errInsert != nil {
//		if isExist, err := p.db.Model(&user).Exists(); err != nil {
//			errInsert = fmt.Errorf("error in inserting user with name: %s : error: %w", user.Name, err)
//		} else if isExist {
//			errInsert = errors.New("user already exists")
//		}
//		return nil, errInsert
//	}
//	if user.UserType == "employer" {
//		newEmpl := models.Employer{UserID: user.ID}
//		_, errInsert = p.db.Model(&newEmpl).Returning("*").Insert()
//	} else if user.UserType == "candidate" {
//		newCand := models.Candidate{UserID: user.ID}
//		_, errInsert = p.db.Model(&newCand).Returning("*").Insert()
//	}
//	if errInsert != nil {
//		return nil, errInsert
//	}
//	return &user, nil
//}
