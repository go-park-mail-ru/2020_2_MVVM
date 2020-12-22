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
		ID: ID,
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
	var  start uint = 0
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
