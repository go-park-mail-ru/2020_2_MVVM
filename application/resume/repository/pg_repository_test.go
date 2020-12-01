package repository

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

type Dummies struct {
	Resume       models.Resume
	Candidate    models.Candidate
	User         models.User
	FavoriteEmpl models.FavoritesForEmpl
	FavoriteCand models.FavoritesForCand
}

func makeUserRow(user models.User) *sqlmock.Rows {
	columns := []string{"user_id", "user_type", "email", "password_hash",
		"name", "surname", "phone", "social_network"}
	values := []driver.Value{user.ID, user.UserType, user.Email, user.PasswordHash,
		user.Name, user.Surname, user.Phone, user.SocialNetwork}
	return sqlmock.NewRows(columns).AddRow(values...)
}

func makeCandRow(cand models.Candidate) *sqlmock.Rows {
	columns := []string{"cand_id", "user_id"}
	values := []driver.Value{cand.ID, cand.UserID}
	return sqlmock.NewRows(columns).AddRow(values...)
}

func makeResumeRow(resume models.Resume) *sqlmock.Rows {
	columns := []string{"resume_id", "cand_id", "title", "description", "salary_min",
		"salary_max", "gender", "career_level", "education_level", "experience_month",
		"skills", "place", "area_search", "path_to_avatar", "date_create"}
	values := []driver.Value{resume.ResumeID, resume.CandID, resume.Title, resume.Description, resume.SalaryMin,
		resume.SalaryMax, resume.Gender, resume.CareerLevel, resume.EducationLevel, resume.ExperienceMonth,
		resume.Skills, resume.Place, resume.AreaSearch, resume.Avatar, resume.DateCreate}
	return sqlmock.NewRows(columns).AddRow(values...)
}

func makeDummies() Dummies {
	DummyResumeID := uuid.MustParse("9922c7ce-b347-4a26-8413-cb2d307fbbc3")
	DummyUserID := uuid.MustParse("b6dff916-e486-4e42-b68b-000000000000")
	DummyCandID := uuid.MustParse("b6dff916-e486-4e42-b68b-fe90d9a0bfda")
	return Dummies{
		Resume: models.Resume{
			ResumeID:    DummyResumeID,
			CandID:      DummyCandID,
			Title:       "Super title",
			Description: "WILL WORK FOR FOOD",
			Skills:      "NONE",
			Gender:      "female",
			DateCreate:  time.Now(),
			Avatar:      "nowhere.png",
		},
		Candidate: models.Candidate{
			ID:     DummyCandID,
			UserID: DummyUserID,
		},
		User: models.User{
			ID:       DummyUserID,
			UserType: "candidate",
		},
	}
}

func beforeTest(t *testing.T) (resume.Repository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	pg := postgres.Dialector{Config: &postgres.Config{Conn: db}}

	conn, err := gorm.Open(pg, &gorm.Config{
		FullSaveAssociations: false,
	})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}

	mock.MatchExpectationsInOrder(true)
	return NewPgRepository(conn), mock
}

func TestCreate(t *testing.T) {
	repo, mock := beforeTest(t)
	r := makeDummies().Resume

	// Ok flow
	queryCandidate := "SELECT .* FROM (.*).\"candidates\" LEFT JOIN (.*).\"users\" (.*) WHERE " +
		"\"candidates\".\"cand_id\" = (.*) ORDER BY \"candidates\".\"cand_id\" LIMIT 1"
	mock.ExpectQuery(queryCandidate).
		WithArgs().
		WillReturnRows(makeCandRow(r.Candidate))
	mock.ExpectQuery("INSERT INTO (.*).\"resume\" .*").
		WithArgs(r.CandID, r.Title, r.SalaryMin, r.SalaryMax, r.Description, r.Skills, r.Gender, r.EducationLevel,
			r.CareerLevel, r.Place, r.ExperienceMonth, r.AreaSearch, r.DateCreate, r.Avatar, r.ResumeID).
		WillReturnRows(sqlmock.NewRows([]string{"resume_id"}).AddRow(r.ResumeID))

	result, err := repo.Create(r)
	assert.Nil(t, err)
	assert.Equal(t, r, *result)

	// Error = no such user
	mock.ExpectQuery(queryCandidate).
		WithArgs().
		WillReturnError(errors.New("TEST ERROR"))
	result, err = repo.Create(r)
	assert.Nil(t, result)
	assert.Error(t, err)

	// Error = failed to create resume
	mock.ExpectQuery(queryCandidate).
		WithArgs().
		WillReturnRows(makeCandRow(r.Candidate))
	mock.ExpectQuery("INSERT INTO (.*).\"resume\" .*").
		WillReturnError(errors.New("TEST ERROR"))

	result, err = repo.Create(r)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestUpdate(t *testing.T) {
	repo, mock := beforeTest(t)
	r := makeDummies().Resume

	// Ok flow
	queryCandidate := "SELECT .* FROM (.*).\"candidates\" LEFT JOIN (.*).\"users\" (.*) WHERE " +
		"\"candidates\".\"cand_id\" = (.*) ORDER BY \"candidates\".\"cand_id\" LIMIT 1"
	mock.ExpectQuery(queryCandidate).
		WithArgs().
		WillReturnRows(makeCandRow(r.Candidate))
	mock.ExpectQuery("INSERT INTO (.*).\"resume\" .*").
		WithArgs(r.CandID, r.Title, r.SalaryMin, r.SalaryMax, r.Description, r.Skills, r.Gender, r.EducationLevel,
			r.CareerLevel, r.Place, r.ExperienceMonth, r.AreaSearch, r.DateCreate, r.Avatar, r.ResumeID).
		WillReturnRows(sqlmock.NewRows([]string{"resume_id"}).AddRow(r.ResumeID))

	result, err := repo.Create(r)
	assert.Nil(t, err)
	assert.Equal(t, r, *result)

	// Error = no such user
	mock.ExpectQuery(queryCandidate).
		WithArgs().
		WillReturnError(errors.New("TEST ERROR"))
	result, err = repo.Update(r)
	assert.Nil(t, result)
	assert.Error(t, err)

	// Error = failed to create resume
	mock.ExpectQuery(queryCandidate).
		WithArgs().
		WillReturnRows(makeCandRow(r.Candidate))
	mock.ExpectQuery("INSERT INTO (.*).\"resume\" .*").
		WillReturnError(errors.New("TEST ERROR"))

	result, err = repo.Update(r)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestSearch(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()
	r := dummies.Resume

	searchParams := resume.SearchParams{
		Gender:          []string{"male"},
		EducationLevel:  []string{"bomz"},
		CareerLevel:     []string{"work for food"},
		ExperienceMonth: []int{12, 24, 100500},
		AreaSearch:      []string{"Moscow", "Paris"},
	}
	kw := "ABRA!"
	searchParams.KeyWords = &kw
	salary := 123
	searchParams.SalaryMin = &salary
	searchParams.SalaryMax = &salary

	r.Candidate = dummies.Candidate
	r.Candidate.User = dummies.User
	// OK flow
	mock.ExpectQuery("SELECT .* FROM (.*).\"resume\" WHERE " +
		"area_search IN (.*) AND " +
		"gender IN (.*) AND " +
		"education_level IN (.*) AND " +
		"career_level IN (.*) AND " +
		"experience_month IN (.*) AND " +
		"salary_min >= .* AND " +
		"salary_max <= .* AND .*").
		WithArgs().
		WillReturnRows(makeResumeRow(r))
	mock.ExpectQuery("SELECT .* FROM (.*).\"candidates\" WHERE \"candidates\".\"cand_id\" = .*").
		WithArgs(r.Candidate.ID).
		WillReturnRows(makeCandRow(r.Candidate))
	mock.ExpectQuery("SELECT .* FROM (.*).\"users\" WHERE \"users\".\"user_id\" = .*").
		WithArgs(r.Candidate.User.ID).
		WillReturnRows(makeUserRow(r.Candidate.User))
	result, err := repo.Search(&searchParams)
	assert.Nil(t, err)
	assert.Equal(t, []models.Resume{r}, result)

	// Error flow, problems with DB
	mock.ExpectQuery("SELECT .* ").WillReturnError(errors.New("TEST ERROR"))
	result, err = repo.Search(&searchParams)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetById(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	r := dummies.Resume
	r.Candidate = dummies.Candidate
	r.Candidate.User = dummies.User
	r.ExperienceCustomComp = []models.ExperienceCustomComp{}

	// OK flow
	mock.ExpectQuery("SELECT .* FROM (.*).\"resume\" WHERE \"resume\".\"resume_id\" = .*").
		WithArgs(r.ResumeID).
		WillReturnRows(makeResumeRow(r))
	mock.ExpectQuery("SELECT .* FROM (.*).\"candidates\" WHERE \"candidates\".\"cand_id\" = .*").
		WithArgs(r.Candidate.ID).
		WillReturnRows(makeCandRow(r.Candidate))
	mock.ExpectQuery("SELECT .* FROM (.*).\"users\" WHERE \"users\".\"user_id\" = .*").
		WithArgs(r.Candidate.User.ID).
		WillReturnRows(makeUserRow(r.Candidate.User))
	mock.ExpectQuery("SELECT .* FROM (.*).\"experience_in_custom_company\" WHERE \"experience_in_custom_company\".\"resume_id\" = .*").
		WithArgs(r.ResumeID).
		WillReturnRows(sqlmock.NewRows([]string{"exp_custom_id", "cand_id", "resume_id", "name_job", "position", "duties", "begin", "finish", "continue_to_today"}))
	result, err := repo.GetById(dummies.Resume.ResumeID)
	assert.Nil(t, err)
	assert.Equal(t, r, *result)

	// Error = database failure
	mock.ExpectQuery(".*").WillReturnError(errors.New("TEST ERROR"))
	result, err = repo.GetById(dummies.Resume.ResumeID)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetAllUserResume(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	r := dummies.Resume
	r.Candidate = dummies.Candidate
	r.Candidate.User = dummies.User

	// OK flow
	mock.ExpectQuery("SELECT .* FROM (.*).\"resume\" WHERE Resume.cand_id = .*").
		WithArgs(r.CandID).
		WillReturnRows(makeResumeRow(r))
	mock.ExpectQuery("SELECT .* FROM (.*).\"candidates\" WHERE \"candidates\".\"cand_id\" = .*").
		WithArgs(r.Candidate.ID).
		WillReturnRows(makeCandRow(r.Candidate))
	mock.ExpectQuery("SELECT .* FROM (.*).\"users\" WHERE \"users\".\"user_id\" = .*").
		WithArgs(r.Candidate.User.ID).
		WillReturnRows(makeUserRow(r.Candidate.User))
	result, err := repo.GetAllUserResume(r.Candidate.ID)
	assert.Nil(t, err)
	assert.Equal(t, []models.Resume{r}, result)

	// Error = database failure
	mock.ExpectQuery(".*").WillReturnError(errors.New("TEST ERROR"))
	result, err = repo.GetAllUserResume(r.Candidate.UserID)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestList(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	r := dummies.Resume
	r.Candidate = dummies.Candidate
	r.Candidate.User = dummies.User

	// OK flow
	mock.ExpectQuery("SELECT .* FROM (.*).\"resume\" LIMIT 100 OFFSET 10").
		WillReturnRows(makeResumeRow(r))
	mock.ExpectQuery("SELECT .* FROM (.*).\"candidates\" WHERE \"candidates\".\"cand_id\" = .*").
		WithArgs(r.Candidate.ID).
		WillReturnRows(makeCandRow(r.Candidate))
	mock.ExpectQuery("SELECT .* FROM (.*).\"users\" WHERE \"users\".\"user_id\" = .*").
		WithArgs(r.Candidate.User.ID).
		WillReturnRows(makeUserRow(r.Candidate.User))
	result, err := repo.List(10, 100)
	assert.Nil(t, err)
	assert.Equal(t, []models.Resume{r}, result)

	// Error = database failure
	mock.ExpectQuery(".*").WillReturnError(errors.New("TEST ERROR"))
	result, err = repo.List(10, 100)
	assert.Nil(t, result)
	assert.Error(t, err)
}

//func TestAddFavorite(t *testing.T) {
//	repo, mock := beforeTest(t)
//	dummies := makeDummies()
//
//	repo.AddFavorite()
//}

//AddFavorite(favoriteForEmpl models.FavoritesForEmpl) (*models.FavoritesForEmpl, error)
//RemoveFavorite(favoriteForEmpl uuid.UUID) error
//GetAllEmplFavoriteResume(emplID uuid.UUID) ([]models.Resume, error)
//GetFavoriteForResume(userID uuid.UUID, resumeID uuid.UUID) (*models.FavoritesForEmpl, error)
//GetFavoriteByID(favoriteID uuid.UUID) (*models.FavoritesForEmpl, error)
//Drop(resume models.Resume) error
