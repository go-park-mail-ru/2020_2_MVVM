package repository

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)

type Dummies struct {
	ID      uuid.UUID
	Vacancy models.Vacancy
	ListVac []models.Vacancy
}

func makeVacRow(vac models.Vacancy) *sqlmock.Rows {
	columns := []string{"vac_id", "empl_id", "comp_id", "title", "salary_min", "salary_max",
		"description", "requirements", "duties", "skills", "sphere", "gender", "employment",
		"area_search", "location", "career_level", "education_level", "experience_month",
		"empl_email", "empl_phone", "path_to_avatar", "date_create"}
	values := []driver.Value{vac.ID, vac.EmpID, vac.CompID, vac.Title, vac.SalaryMin, vac.SalaryMax,
		vac.Description, vac.Requirements, vac.Duties, vac.Skills, vac.Sphere, vac.Gender,
		vac.Employment, vac.AreaSearch, vac.Location, vac.CareerLevel, vac.EducationLevel,
		vac.ExperienceMonth, vac.EmpEmail, vac.EmpPhone, vac.Avatar, vac.DateCreate}
	return sqlmock.NewRows(columns).AddRow(values...)
}

func makeDummies() Dummies {
	ID := uuid.New()
	vac := models.Vacancy{
		ID:     ID,
		EmpID:  ID,
		CompID: ID,
	}
	return Dummies{
		ID:      ID,
		Vacancy: vac,
		ListVac: []models.Vacancy{vac},
	}
}

func beforeTest(t *testing.T) (vacancy.RepositoryVacancy, sqlmock.Sqlmock) {
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

func TestGetVacancyById(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	query := "SELECT \\* FROM \"main\".\"vacancy\" WHERE vac_id = (.*) LIMIT 1"
	mock.ExpectQuery(query).
		WithArgs(dummies.Vacancy.ID).
		WillReturnRows(makeVacRow(dummies.Vacancy))

	result, err := repo.GetVacancyById(dummies.Vacancy.ID)
	assert.Nil(t, err)
	assert.Equal(t, *result, dummies.Vacancy)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).
		WithArgs(dummies.Vacancy.ID).
		WillReturnError(error)

	result, err = repo.GetVacancyById(dummies.Vacancy.ID)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestUpdateVacancy(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	query := "select empl_id from main.vacancy where vac_id = (.*)"
	mock.ExpectQuery(query).
		WithArgs(dummies.Vacancy.ID).
		WillReturnRows(sqlmock.NewRows([]string{"empl_id"}).AddRow(&dummies.Vacancy.EmpID))

	query2 := "UPDATE \"main\".\"vacancy\" SET (.*) WHERE \"vac_id\" = (.*)"
	mock.ExpectExec(query2).
		WithArgs(dummies.ID, dummies.ID, dummies.ID, dummies.Vacancy.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	query3 := "SELECT \\* FROM \"main\".\"vacancy\" WHERE vac_id = (.*) LIMIT 1"
	mock.ExpectQuery(query3).
		WithArgs(dummies.Vacancy.ID).
		WillReturnRows(makeVacRow(dummies.Vacancy))

	result, err := repo.UpdateVacancy(dummies.Vacancy)
	assert.Nil(t, err)
	assert.Equal(t, *result, dummies.Vacancy)

	// Error flow
	testErr := errors.New("test testErr")
	mock.ExpectQuery(query).
		WithArgs(dummies.Vacancy.ID).
		WillReturnError(testErr)

	result, err = repo.UpdateVacancy(dummies.Vacancy)
	assert.Nil(t, result)
	assert.Error(t, err)

	mock.ExpectQuery(query).
		WithArgs(dummies.Vacancy.ID).
		WillReturnRows(makeVacRow(dummies.Vacancy))
	mock.ExpectQuery(query2).
		WithArgs(dummies.ID, dummies.ID, dummies.ID, dummies.Vacancy.ID).
		WillReturnError(testErr)

	result, err = repo.UpdateVacancy(dummies.Vacancy)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestGetVacancyList(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()
	var start uint = 0
	var limit uint = 2

	// all
	query := "SELECT \\* FROM \"main\".\"vacancy\" ORDER BY date_create desc LIMIT 2"
	mock.ExpectQuery(query).
		WithArgs().
		WillReturnRows(makeVacRow(dummies.Vacancy))

	result, err := repo.GetVacancyList(start, limit, dummies.Vacancy.ID, vacancy.TopAll)
	assert.Nil(t, err)
	assert.Equal(t, result, dummies.ListVac)

	//employer
	query2 := "SELECT \\* FROM \"main\".\"vacancy\" WHERE empl_id = (.*) ORDER BY date_create desc LIMIT 2"
	mock.ExpectQuery(query2).
		WithArgs(dummies.ID).
		WillReturnRows(makeVacRow(dummies.Vacancy))

	result, err = repo.GetVacancyList(start, limit, dummies.Vacancy.ID, vacancy.ByEmpId)
	assert.Nil(t, err)
	assert.Equal(t, result, dummies.ListVac)

	//company
	query3 := "SELECT \\* FROM \"main\".\"vacancy\" WHERE comp_id = (.*) ORDER BY date_create desc LIMIT 2"
	mock.ExpectQuery(query3).
		WithArgs(dummies.ID).
		WillReturnRows(makeVacRow(dummies.Vacancy))

	result, err = repo.GetVacancyList(start, limit, dummies.Vacancy.ID, 2)
	assert.Nil(t, err)
	assert.Equal(t, result, dummies.ListVac)

	// Error flow
	testErr := errors.New("test error")
	mock.ExpectQuery(query).
		WithArgs().
		WillReturnError(testErr)

	result, err = repo.GetVacancyList(start, limit, dummies.Vacancy.ID, -1)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestDeleteVacancy(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	query := "DELETE FROM \"main\".\"vacancy\" WHERE vac_id = (.*) AND empl_id = (.*)"
	dvacancy := dummies.Vacancy
	mock.ExpectQuery(query).
		WithArgs(nil).
		WillReturnError(errors.New("TEST ERROR"))
	err2 := repo.DeleteVacancy(dvacancy.ID, dvacancy.ID)
	assert.Error(t, err2)
}

func TestGetRecommendation(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	query := "SELECT \\* FROM \"main\".\"vacancy\" WHERE sphere IN (.*) ORDER BY date_create desc LIMIT 1 OFFSET 1"
	mock.ExpectQuery(query).
		WithArgs(0).
		WillReturnRows(makeVacRow(dummies.Vacancy))
	res, err1 := repo.GetRecommendation(1, 1, 0, []int{0})
	assert.Nil(t, err1)
	assert.Equal(t, []models.Vacancy{dummies.Vacancy}, res)
	mock.ExpectQuery(query).
		WithArgs(nil).
		WillReturnError(errors.New("TEST ERROR"))
	res, err2 := repo.GetRecommendation(1, -1, 0, []int{0})
	assert.Nil(t, res)
	assert.Error(t, err2)
}

func TestGetPreferredSalary(t *testing.T) {
	repo, mock := beforeTest(t)

	query := "select avg(.*) as avg from main.resume " +
		"join main.candidates on resume.cand_id = candidates.cand_id where user_id = (.*) and salary_min>0 and salary_max>0"
	mock.ExpectQuery(query).
		WithArgs(uuid.Nil).
		WillReturnError(errors.New("TEST ERROR"))
	res, err2 := repo.GetPreferredSalary(uuid.Nil)
	assert.Nil(t, res)
	assert.Nil(t, err2)
}

func TestGetPreferredSpheres(t *testing.T) {
	repo, mock := beforeTest(t)

	rec := models.Recommendation{}
	id := uuid.New()
	query := "SELECT \\* FROM \"main\".\"recommendation\" WHERE user_id = (.*)"
	mock.ExpectQuery(query).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"rec_id"}).AddRow(rec.ID))
	res, err1 := repo.GetPreferredSpheres(id)
	//assert.Equal(t, rec, res)
	assert.Nil(t, err1)
	mock.ExpectQuery(query).
		WithArgs(uuid.Nil).
		WillReturnError(errors.New("TEST ERROR"))
	res, err2 := repo.GetPreferredSpheres(uuid.Nil)
	assert.Nil(t, res)
	assert.Error(t, err2)
}

func TestAddRecommendation(t *testing.T) {
	repo, mock := beforeTest(t)
	query := "insert into main.recommendation (user_id, sphere0) values (.*) ON CONFLICT (.*) DO UPDATE SET (.*)"
	mock.ExpectQuery(query).
		WithArgs(uuid.Nil).
		WillReturnError(errors.New("TEST ERROR"))
	err2 := repo.AddRecommendation(uuid.Nil, 0)
	assert.Error(t, err2)
}

func TestSearchVacancies(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()
	params := models.VacancySearchParams{Sphere: []int{0, 1}, StartDate: time.Now().String()}
	query := "SELECT \\* FROM \"main\".\"vacancy\" WHERE date(.*) >= (.*) AND sphere IN (.*) ORDER BY date_create desc"
	mock.ExpectQuery(query).
		WithArgs().
		WillReturnRows(makeVacRow(dummies.Vacancy))
	res, err1 := repo.SearchVacancies(params)
	assert.Equal(t, []models.Vacancy{dummies.Vacancy}, res)
	assert.Nil(t, err1)
	mock.ExpectQuery(query).
		WithArgs().
		WillReturnError(errors.New("TEST ERROR"))
	res, err2 := repo.SearchVacancies(params)
	assert.Nil(t, res)
	assert.Error(t, err2)
}

func TestCreateVacancy(t *testing.T) {
	repo, mock := beforeTest(t)
	vac := models.Vacancy{Title: "test"}
	query := "SELECT \\* FROM \"main\".\"employers\" WHERE empl_id = (.*) LIMIT 1"
	mock.ExpectQuery(query).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"empl_id", "comp_id"}).AddRow(vac.EmpID, vac.CompID))
	res, err1 := repo.CreateVacancy(vac)
	assert.Nil(t, res)
	assert.Equal(t, err1, errors.New("error: employer must have company for vacancy creation"))
	vac.CompID = uuid.New()
	mock.ExpectQuery(query).
		WithArgs().
		WillReturnRows(sqlmock.NewRows([]string{"empl_id", "comp_id"}).AddRow(vac.EmpID, vac.CompID))
	res, err2 := repo.CreateVacancy(vac)
	assert.Nil(t, res)
	assert.Error(t, err2)
	mock.ExpectQuery(query).
		WithArgs().
		WillReturnError(errors.New("TEST ERROR"))
	res, err3 := repo.CreateVacancy(vac)
	assert.Nil(t, res)
	assert.Error(t, err3)

}