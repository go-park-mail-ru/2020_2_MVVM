package repository

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type Dummies struct {
	ID      uuid.UUID
	Company models.OfficialCompany
}

func makeCompanyRow(comp models.OfficialCompany) *sqlmock.Rows {
	columns := []string{"comp_id", "name", "spheres", "description",
		"area_search", "link", "count_vacancy", "path_to_avatar"}
	values := []driver.Value{comp.ID, comp.Name, comp.Spheres, comp.Description,
		comp.AreaSearch, comp.Link, comp.VacCount, comp.Avatar}
	return sqlmock.NewRows(columns).AddRow(values...)
}

func makeDummies() Dummies {
	ID := uuid.New()
	return Dummies{
		ID: ID,
		Company: models.OfficialCompany{
			ID:   ID,
			Name: "Name",
		},
	}
}

func beforeTest(t *testing.T) (official_company.OfficialCompanyRepository, sqlmock.Sqlmock) {
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

func TestGetCompaniesList(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	start := 0
	limit := 1

	query := "SELECT \\* FROM (.*).\"official_companies\" ORDER BY name LIMIT 1"
	mock.ExpectQuery(query).
		WillReturnRows(makeCompanyRow(dummies.Company))

	res, err := repo.GetCompaniesList(uint(start), uint(limit))
	assert.Nil(t, err)
	assert.Equal(t, []models.OfficialCompany{dummies.Company}, res)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WillReturnError(error)

	res, err = repo.GetCompaniesList(uint(start), uint(limit))
	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestGetMineCompany(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	query := "SELECT \\* FROM (.*).\"employers\" WHERE empl_id = (.*) LIMIT 1"
	mock.ExpectQuery(query).
		WithArgs(dummies.Company.ID).
		WillReturnRows(makeCompanyRow(dummies.Company))

	query2 := "SELECT \\* FROM \"main\".\"official_companies\" WHERE comp_id = (.*) LIMIT 1"
	mock.ExpectQuery(query2).
		WithArgs(dummies.Company.ID).
		WillReturnRows(makeCompanyRow(dummies.Company))

	res, err := repo.GetMineCompany(dummies.Company.ID)
	assert.Nil(t, err)
	assert.Equal(t, &dummies.Company, res)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WillReturnError(error)

	res, err = repo.GetMineCompany(dummies.Company.ID)
	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestCreateCompany(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	// OK flow
	query := "SELECT \"comp_id\" FROM (.*).\"employers\" WHERE empl_id = (.*) LIMIT 1"
	mock.ExpectQuery(query).
		WithArgs(dummies.Company.ID).
		WillReturnRows(sqlmock.NewRows([]string{"comp_id"}).AddRow(uuid.Nil))

	query2 := "INSERT INTO \"main\".\"official_companies\" (.*) VALUES (.*) RETURNING \"comp_id\""
	mock.ExpectQuery(query2).
		WithArgs(dummies.Company.Name, dummies.Company.Spheres, dummies.Company.Description,
			dummies.Company.AreaSearch, dummies.Company.Link, dummies.Company.VacCount, dummies.Company.Avatar, dummies.Company.ID).
		WillReturnRows(sqlmock.NewRows([]string{"comp_id"}).AddRow(dummies.Company.ID))

	query3 := "UPDATE \"main\".\"employers\" SET \"comp_id\"=(.*) WHERE empl_id = (.*)"
	mock.ExpectExec(query3).
		WithArgs(dummies.Company.ID, dummies.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	res, err := repo.CreateOfficialCompany(dummies.Company, dummies.Company.ID)
	assert.Nil(t, err)
	assert.Equal(t, dummies.Company, *res)

	// Error - select empl
	error := errors.New("test error")
	mock.ExpectQuery(query).WillReturnError(error)
	res, err = repo.GetMineCompany(dummies.Company.ID)
	assert.Nil(t, res)
	assert.Error(t, err)

	// Error - insert company
	mock.ExpectQuery(query).
		WithArgs(dummies.Company.ID).
		WillReturnRows(sqlmock.NewRows([]string{"comp_id"}).AddRow(uuid.Nil))
	mock.ExpectQuery(query2).WillReturnError(error)
	res, err = repo.GetMineCompany(dummies.Company.ID)
	assert.Nil(t, res)
	assert.Error(t, err)

	// Error - update empl
	mock.ExpectQuery(query).
		WithArgs(dummies.Company.ID).
		WillReturnRows(sqlmock.NewRows([]string{"comp_id"}).AddRow(uuid.Nil))
	mock.ExpectQuery(query2).
		WithArgs(dummies.Company.Name, dummies.Company.Spheres, dummies.Company.Description,
			dummies.Company.AreaSearch, dummies.Company.Link, dummies.Company.VacCount, dummies.Company.Avatar, dummies.Company.ID).
		WillReturnRows(sqlmock.NewRows([]string{"comp_id"}).AddRow(dummies.Company.ID))
	mock.ExpectQuery(query3).WillReturnError(error)
	res, err = repo.GetMineCompany(dummies.Company.ID)
	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestSearch(t *testing.T) {
	repo, mock := beforeTest(t)
	dummies := makeDummies()

	searchParams := models.CompanySearchParams{
		KeyWords:   "A",
		AreaSearch: []string{"Moscow", "Paris"},
		Spheres:    []int{1, 2},
		OrderBy:    "name",
		ByAsc:      false,
		VacCount:   2,
	}

	// OK flow
	query := "SELECT .* FROM (.*).\"official_companies\" WHERE " +
		"count_company >= (.*) AND area_search IN (.*) AND " +
		"LOWER(.*) LIKE (.*) AND spheres @> (.*) ORDER BY name"

	mock.ExpectQuery(query).
		WithArgs().
		WillReturnRows(makeCompanyRow(dummies.Company))

	result, err := repo.SearchCompanies(searchParams)
	assert.Nil(t, err)
	assert.Equal(t, []models.OfficialCompany{dummies.Company}, result)

	// Error - select empl
	error := errors.New("test error")
	mock.ExpectQuery(query).WillReturnError(error)
	result, err = repo.SearchCompanies(searchParams)
	assert.Nil(t, result)
	assert.Error(t, err)
}
