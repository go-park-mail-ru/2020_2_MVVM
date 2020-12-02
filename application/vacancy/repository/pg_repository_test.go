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
	return Dummies{
		ID: ID,
		Vacancy: models.Vacancy{
			ID:     ID,
			EmpID:  ID,
			CompID: ID,
		},
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
	query := "SELECT \\* FROM \"main\".\"vacancy\" WHERE vac_id = (.*) LIMIT 1"
	mock.ExpectQuery(query).
		WithArgs(dummies.Vacancy.ID).
		WillReturnRows(makeVacRow(dummies.Vacancy))

	query2 := "UPDATE \"main\".\"vacancy\" SET (.*) WHERE \"vac_id\" = (.*)"
	mock.ExpectExec(query2).
		WithArgs(dummies.ID, dummies.ID, dummies.ID, dummies.Vacancy.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := repo.UpdateVacancy(dummies.Vacancy)
	assert.Nil(t, err)
	assert.Equal(t, *result, dummies.Vacancy)

	// Error flow
	//error := errors.New("test error")
	//mock.ExpectQuery(query).
	//	WithArgs(dummies.Vacancy.ID).
	//	WillReturnError(error)
	//
	//result, err = repo.GetVacancyById(dummies.Vacancy.ID)
	//assert.Nil(t, result)
	//assert.Error(t, err)
}
