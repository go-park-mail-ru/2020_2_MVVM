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

func mockDB() (*mocks.DB, pgStorage) {
	db := new(mocks.DB)
	r := pgStorage{db: db}
	return db, r
}

func mockQueryUser(db *mocks.DB) *ormmocks.Query {
	query := new(ormmocks.Query)
	mockCall := db.On("Model", mock.AnythingOfType("*models.User")).Return(query)
	mockCall.RunFn = func(args mock.Arguments) {
		user := args[0].(*models.User)
		*user = testUser
	}
	return query
}

func mockQueryCandidate(db *mocks.DB) *ormmocks.Query  {
	query := new(ormmocks.Query)
	mockCall := db.On("Model", mock.AnythingOfType("*models.Candidate")).Return(query)
	mockCall.RunFn = func(args mock.Arguments) {
		user := args[0].(*models.Candidate)
		*user = cand
	}
	return query
}

func mockQueryEmployer(db *mocks.DB) *ormmocks.Query  {
	query := new(ormmocks.Query)
	mockCall := db.On("Model", mock.AnythingOfType("*models.Employer")).Return(query)
	mockCall.RunFn = func(args mock.Arguments) {
		user := args[0].(*models.Employer)
		*user = empl
	}
	return query
}

func TestGetUserByID(t *testing.T) {
	db, r := mockDB()
	query := mockQueryUser(db)

	query.On("Where", "user_id = ?", ID.String()).Return(query)
	query.On("Select").Return(nil)

	foo, err := r.GetUserByID(ID.String())
	assert.Nil(t, err)
	assert.Equal(t, testUser, *foo)
}

func TestLogin(t *testing.T) {
	db, r := mockDB()
	query := mockQueryUser(db)

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
	db, r := mockDB()
	query := mockQueryUser(db)

	mockResult := MockResult{}
	query.On("WherePK").Return(query)
	query.On("Returning", "*").Return(query)
	query.On("Update").Return(mockResult, nil)

	foo, err := r.UpdateUser(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, *foo)
}

func TestGetCandidateByID(t *testing.T) {
	db, r := mockDB()
	query := mockQueryCandidate(db)

	query.On("Where", "user_id = ?", ID.String()).Return(query)
	query.On("Select").Return(nil)

	foo, err := r.GetCandidateByID(ID.String())
	assert.Nil(t, err)
	assert.Equal(t, cand, *foo)
}

func TestGetEmployerByID(t *testing.T) {
	db, r := mockDB()
	query := mockQueryEmployer(db)
	query.On("Where", "user_id = ?", ID.String()).Return(query)
	query.On("Select").Return(nil)

	foo, err := r.GetEmployerByID(ID.String())
	assert.Nil(t, err)
	assert.Equal(t, empl, *foo)
}

func TestCreate(t *testing.T) {
	db, r := mockDB()
	queryUser := mockQueryUser(db)
	queryEmpl := mockQueryEmployer(db)

	mockResult := MockResult{}
	queryUser.On("Returning", "*").Return(queryUser)
	queryUser.On("Insert").Return(mockResult, nil)

	//mockResult2 := MockResult{}

	queryEmpl.On("Returning", "*").Return(queryEmpl)
	queryEmpl.On("Insert").Return(mockResult, nil)


	foo, err := r.CreateUser(testUser)
	assert.Nil(t, err)
	assert.Equal(t, testUser, *foo)
}
