package repository

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Dummies struct {
	User      models.User
	Candidate models.Candidate
	Employer  models.Employer
}

func makeUserRow(user models.User) ([]string, []driver.Value) {
	columns := []string{"user_id", "user_type", "email", "password_hash",
		"name", "surname", "phone", "social_network"}
	values := []driver.Value{user.ID, user.UserType, user.Email, user.PasswordHash,
		user.Name, user.Surname, user.Phone, user.SocialNetwork}
	return columns, values
}

func makeCandRow(cand models.Candidate) ([]string, []driver.Value) {
	columns := []string{"cand_id", "user_id"}
	values := []driver.Value{cand.ID, cand.UserID}
	return columns, values
}

func makeEmplRow(empl models.Employer) ([]string, []driver.Value) {
	columns := []string{"empl_id", "user_id", "comp_id"}
	values := []driver.Value{empl.ID, empl.UserID, empl.CompanyID}
	return columns, values
}

func makeDummies() Dummies {
	DummyUserID := uuid.New()
	DummyUser := models.User{
		ID:            DummyUserID,
		UserType:      "employer",
		Name:          "ID",
		Surname:       "ID",
		Email:         "ID",
		PasswordHash:  []byte("ASD"),
		Phone:         nil,
		SocialNetwork: nil,
	}

	DummyCandidate := models.Candidate{
		ID:     uuid.New(),
		UserID: DummyUserID,
	}
	DummyEmployer := models.Employer{
		ID:        uuid.New(),
		UserID:    DummyUserID,
		CompanyID: uuid.New(),
	}
	return Dummies{
		User:      DummyUser,
		Candidate: DummyCandidate,
		Employer:  DummyEmployer,
	}
}

func beforeTest(t *testing.T) (user.RepositoryUser, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	pg := postgres.Dialector{Config: &postgres.Config{Conn: db}}

	conn, err := gorm.Open(pg, &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	return NewPgRepository(conn), mock
}

func TestGetUserByID(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	query := "SELECT \\* FROM (.*).\"users\" WHERE user_id = (.*) LIMIT 1"
	cols, row := makeUserRow(dummies.User)
	mock.ExpectQuery(query).
		WithArgs(dummies.User.ID).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(row...))

	fetchedUser, err := repo.GetUserByID(dummies.User.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, *fetchedUser, dummies.User)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(dummies.User.ID).WillReturnError(error)

	fetchedUser, err = repo.GetUserByID(dummies.User.ID.String())
	assert.Nil(t, fetchedUser)
	assert.Error(t, err)
}

func TestCreateUser(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	// 1 - insert user
	query := "INSERT INTO (.*).\"users\" (.*) VALUES (.*) RETURNING \"user_id\""
	duser := dummies.User
	argsUser := []driver.Value{duser.UserType, duser.Name, duser.Surname, duser.Email, duser.PasswordHash, duser.Phone,
		duser.SocialNetwork, duser.ID}
	mock.ExpectQuery(query).
		WithArgs(argsUser...).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(duser.ID))

	// 2 - insert candidate/employer
	mock.ExpectQuery("INSERT INTO (.*).\"(candidates|employers)\" (.*)").
		WithArgs(dummies.Candidate.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"empl_id"}).AddRow(dummies.Employer.ID))

	user, err := repo.CreateUser(dummies.User)
	assert.Nil(t, err)
	assert.Equal(t, *user, dummies.User)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(argsUser).WillReturnError(error)

	user, err = repo.CreateUser(dummies.User)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUpdateUser(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	query := "UPDATE (.*).\"users\" SET (.*) WHERE \"user_id\" = (.*)"
	duser := dummies.User
	argsUser := []driver.Value{duser.UserType, duser.Name, duser.Surname, duser.Email, duser.PasswordHash, duser.Phone,
		duser.SocialNetwork, duser.ID}
	mock.ExpectExec(query).
		WithArgs(argsUser...).
		WillReturnResult(sqlmock.NewResult(1, 1))

	user, err := repo.UpdateUser(dummies.User)
	assert.Nil(t, err)
	assert.Equal(t, *user, dummies.User)

	// Error flow
	error := errors.New("test error")
	mock.ExpectExec(query).
		WithArgs(argsUser...).
		WillReturnError(error)

	user, err = repo.UpdateUser(dummies.User)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	credentials := models.UserLogin{
		Email:    dummies.User.Email,
		Password: "SUPERPASSWORD",
	}
	dummies.User.PasswordHash, _ = bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)

	query := "SELECT \\* FROM (.*).\"users\" WHERE email = (.*) LIMIT 1"
	col, rows := makeUserRow(dummies.User)
	mock.ExpectQuery(query).
		WithArgs(dummies.User.Email).
		WillReturnRows(sqlmock.NewRows(col).AddRow(rows...))

	fetchedUser, err := repo.Login(credentials)
	assert.Nil(t, err)
	assert.Equal(t, *fetchedUser, dummies.User)

	// Error flow = db error
	mock.ExpectQuery(query).
		WithArgs(dummies.User.Email).
		WillReturnError(errors.New("TEST ERROR"))

	fetchedUser, err = repo.Login(credentials)
	assert.Nil(t, fetchedUser)
	assert.Error(t, err)

	// Error flow = credentials error
	credentials.Password = "FAKE PASSWORD"
	mock.ExpectQuery(query).
		WithArgs(dummies.User.Email).
		WillReturnRows(sqlmock.NewRows(col).AddRow(rows...))

	fetchedUser, err = repo.Login(credentials)
	assert.Nil(t, fetchedUser)
	assert.Error(t, err)
}

func TestGetCandByID(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	candQuery := "SELECT \\* FROM (.*).\"candidates\" WHERE cand_id = (.*) LIMIT 1"
	candCol, candRow := makeCandRow(dummies.Candidate)
	mock.ExpectQuery(candQuery).
		WithArgs(dummies.Candidate.ID).
		WillReturnRows(sqlmock.NewRows(candCol).AddRow(candRow...))

	userQuery := "SELECT \\* FROM (.*).\"users\" WHERE user_id = (.*) LIMIT 1"
	userCol, userRow := makeUserRow(dummies.User)
	mock.ExpectQuery(userQuery).
		WithArgs(dummies.User.ID).
		WillReturnRows(sqlmock.NewRows(userCol).AddRow(userRow...))

	fetchedUser, err := repo.GetCandByID(dummies.Candidate.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, *fetchedUser, dummies.User)

	// Error flow = user is not a candidate
	mock.ExpectQuery(candQuery).
		WithArgs(dummies.Candidate.ID).
		WillReturnError(errors.New("TEST ERROR"))
	fetchedUser, err = repo.GetCandByID(dummies.Candidate.ID.String())
	assert.Nil(t, fetchedUser)
	assert.Error(t, err)

	// Error flow = user fetch error
	mock.ExpectQuery(candQuery).
		WithArgs(dummies.Candidate.ID).
		WillReturnRows(sqlmock.NewRows(candCol).AddRow(candRow...))
	mock.ExpectQuery(userQuery).
		WithArgs(dummies.User.ID).
		WillReturnError(errors.New("TEST ERROR"))
	fetchedUser, err = repo.GetCandByID(dummies.Candidate.ID.String())
	assert.Nil(t, fetchedUser)
	assert.Error(t, err)
}

func TestGetEmplByID(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	candQuery := "SELECT \\* FROM (.*).\"employers\" WHERE empl_id = (.*) LIMIT 1"
	candCol, candRow := makeEmplRow(dummies.Employer)
	mock.ExpectQuery(candQuery).
		WithArgs(dummies.Employer.ID).
		WillReturnRows(sqlmock.NewRows(candCol).AddRow(candRow...))

	userQuery := "SELECT \\* FROM (.*).\"users\" WHERE user_id = (.*) LIMIT 1"
	userCol, userRow := makeUserRow(dummies.User)
	mock.ExpectQuery(userQuery).
		WithArgs(dummies.User.ID).
		WillReturnRows(sqlmock.NewRows(userCol).AddRow(userRow...))

	fetchedUser, err := repo.GetEmplByID(dummies.Employer.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, *fetchedUser, dummies.User)

	// Error flow = user is not a candidate
	mock.ExpectQuery(candQuery).
		WithArgs(dummies.Employer.ID).
		WillReturnError(errors.New("TEST ERROR"))
	fetchedUser, err = repo.GetEmplByID(dummies.Employer.ID.String())
	assert.Nil(t, fetchedUser)
	assert.Error(t, err)

	// Error flow = user fetch error
	mock.ExpectQuery(candQuery).
		WithArgs(dummies.Employer.ID).
		WillReturnRows(sqlmock.NewRows(candCol).AddRow(candRow...))
	mock.ExpectQuery(userQuery).
		WithArgs(dummies.User.ID).
		WillReturnError(errors.New("TEST ERROR"))
	fetchedUser, err = repo.GetEmplByID(dummies.Employer.ID.String())
	assert.Nil(t, fetchedUser)
	assert.Error(t, err)
}

func TestGetEmployerByID(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	query := "SELECT \\* FROM (.*).\"employers\" WHERE user_id = (.*) LIMIT 1"
	cols, row := makeEmplRow(dummies.Employer)
	mock.ExpectQuery(query).
		WithArgs(dummies.Employer.UserID).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(row...))

	fetchedEmpl, err := repo.GetEmployerByID(dummies.Employer.UserID.String())
	assert.Nil(t, err)
	assert.Equal(t, *fetchedEmpl, dummies.Employer)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(dummies.Employer.UserID).WillReturnError(error)

	fetchedEmpl, err = repo.GetEmployerByID(dummies.Employer.UserID.String())
	assert.Nil(t, fetchedEmpl)
	assert.Error(t, err)
}

func TestGetCandidateByID(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	query := "SELECT \\* FROM (.*).\"candidates\" WHERE user_id = (.*) LIMIT 1"
	cols, row := makeCandRow(dummies.Candidate)
	mock.ExpectQuery(query).
		WithArgs(dummies.Candidate.UserID).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(row...))

	fetchedCand, err := repo.GetCandidateByID(dummies.Candidate.UserID.String())
	assert.Nil(t, err)
	assert.Equal(t, *fetchedCand, dummies.Candidate)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(dummies.Candidate.UserID).WillReturnError(error)

	fetchedCand, err = repo.GetCandidateByID(dummies.Candidate.UserID.String())
	assert.Nil(t, fetchedCand)
	assert.Error(t, err)
}

//GetEmployerByID(id string) (*models.Employer, error)
//GetCandidateByID(id string) (*models.Candidate, error)
//UpdateEmployer(employer models.Employer) (*models.Employer, error)
