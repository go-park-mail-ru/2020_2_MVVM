package repository

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	vacancy2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/vacancy"
	mocks "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/vacancy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

var ID = uuid.New()
//var testResponse = models.Response{
//	ID:        ID,
//	ResumeID:  ID,
//	VacancyID: ID,
//}
//var testResume = models.Resume{
//	ResumeID: ID,
//	CandID:   ID,
//	Title:    "ID",
//}

type Dummies struct {
	Response    models.Response
	ListResp    []models.Response
	RespColumns []string
	RespRow     []driver.Value

	Resume        models.Resume
	ListResume    []models.Resume
	ResumeColumns []string
	ResumeRow     []driver.Value

	Vacancy        models.Vacancy
	ListVacancy    []models.Vacancy
	VacancyColumns []string
	VacancyRow     []driver.Value
}

func makeDummies() Dummies {
	DummyRespID := uuid.New()
	DummyResp := models.Response{
		ID:        DummyRespID,
		ResumeID:  DummyRespID,
		VacancyID: DummyRespID,
		Initial:   common.Candidate,
		Status:    "sent",
	}
	var DummyResume = models.Resume{
		ResumeID: ID,
		CandID:   ID,
		Title:    "ID",
	}
	var DummyVac = models.Vacancy{
		ID:     ID,
		EmpID:  ID,
		CompID: ID,
		Title:  "ID",
	}

	RespTableColumns := []string{"response_id", "resume_id", "vacancy_id", "status", "unread",
		"initial", "date_create"}
	DummyRespRow := []driver.Value{DummyResp.ID.String(), DummyResp.ResumeID.String(),
		DummyResp.VacancyID.String(), DummyResp.Status, false,
		DummyResp.Initial, DummyResp.DateCreate}

	ResumeColumns := []string{"resume_id", "cand_id", "title", "description", "salary_min", "salary_max",
		"gender", "career_level", "education_level", "experience_month", "skills", "place", "area_search",
		"path_to_avatar", "date_create"}
	ResumeRow := []driver.Value{DummyResume.ResumeID, DummyResume.CandID, DummyResume.Title, DummyResume.Description,
		DummyResume.SalaryMin, DummyResume.SalaryMax, DummyResume.Gender, DummyResume.CareerLevel, DummyResume.EducationLevel,
		DummyResume.ExperienceMonth, DummyResume.Skills, DummyResume.Place, DummyResume.AreaSearch, DummyResume.Avatar,
		DummyResume.DateCreate}

	VacancyColumns := []string{"vac_id", "empl_id", "comp_id", "title", "salary_min", "salary_max", "description",
		"requirements", "duties", "skills", "sphere", "gender", "employment", "area_search", "location",
		"career_level", "education_level", "experience_month", "empl_email", "empl_phone",
		"path_to_avatar", "date_create"}
	VacancyRow := []driver.Value{DummyVac.ID, DummyVac.EmpID, DummyVac.CompID, DummyVac.Title, DummyVac.SalaryMin,
		DummyVac.SalaryMax, DummyVac.Description, DummyVac.Requirements, DummyVac.Duties, DummyVac.Skills, DummyVac.Sphere,
		DummyVac.Gender, DummyVac.Employment, DummyVac.AreaSearch, DummyVac.Location, DummyVac.CareerLevel, DummyVac.EducationLevel,
		DummyVac.ExperienceMonth, DummyVac.EmpEmail, DummyVac.EmpPhone, DummyVac.Avatar, DummyVac.DateCreate}

	return Dummies{
		Response:       DummyResp,
		ListResp:       []models.Response{DummyResp},
		Resume:         DummyResume,
		ListResume:     []models.Resume{DummyResume},
		Vacancy:        DummyVac,
		ListVacancy:    []models.Vacancy{DummyVac},
		RespColumns:    RespTableColumns,
		RespRow:        DummyRespRow,
		ResumeColumns:  ResumeColumns,
		ResumeRow:      ResumeRow,
		VacancyColumns: VacancyColumns,
		VacancyRow:     VacancyRow,
	}
}

func beforeTest2(t *testing.T) (response.ResponseRepository, sqlmock.Sqlmock, *mocks.RepositoryVacancy) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	pg := postgres.Dialector{Config: &postgres.Config{Conn: db}}
	conn, err := gorm.Open(pg, &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening gorm database", err)
	}
	vacMock := new(mocks.RepositoryVacancy)
	return NewPgRepository(conn, vacMock), mock, vacMock
}

func TestResponseGetUserByID(t *testing.T) {
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()

	query := "SELECT \\* FROM \"main\".\"response\" WHERE \"response\".\"response_id\" = (.*) ORDER BY \"response\".\"response_id\" LIMIT 1"
	mock.ExpectQuery(query).
		WithArgs(dummies.Response.ID).
		WillReturnRows(sqlmock.NewRows(dummies.RespColumns).AddRow(dummies.RespRow...))

	resp, err := repo.GetByID(dummies.Response.ID)
	assert.Nil(t, err)
	assert.Equal(t, *resp, dummies.Response)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(dummies.Response.ID).WillReturnError(error)

	resp, err = repo.GetByID(dummies.Response.ID)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestCreateResponse(t *testing.T) {
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()

	mock.ExpectQuery("INSERT INTO (.*).\"response\" (.*)").
		WithArgs(dummies.Response.ResumeID, dummies.Response.VacancyID, dummies.Response.Initial,
			dummies.Response.Status, dummies.Response.DateCreate, dummies.Response.ID).
		WillReturnRows(sqlmock.NewRows([]string{"response_id"}).AddRow(dummies.Response.ID))

	resp, err := repo.Create(dummies.Response)
	assert.Nil(t, err)
	assert.Equal(t, *resp, dummies.Response)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery("INSERT INTO (.*).\"response\" (.*)").
		WithArgs(dummies.Response.ResumeID, dummies.Response.VacancyID, dummies.Response.Initial,
			dummies.Response.Status, dummies.Response.DateCreate, dummies.Response.ID).
		WillReturnError(error)

	resp, err = repo.GetByID(dummies.Response.ID)
	assert.Nil(t, resp)
	assert.Error(t, err)

}

func TestUpdateResponse(t *testing.T) {
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()

	mock.ExpectExec("UPDATE \"main\".\"response\" SET \"status\"=(.*) WHERE \"response_id\" = (.*)").
		WithArgs(dummies.Response.Status, dummies.Response.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	resp, err := repo.UpdateStatus(dummies.Response)
	assert.Nil(t, err)
	assert.Equal(t, *resp, dummies.Response)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery("UPDATE \"main\".\"response\" SET \"status\"=(.*) WHERE \"response_id\" = (.*)").
		WithArgs(dummies.Response.ResumeID).WillReturnError(error)

	resp, err = repo.UpdateStatus(dummies.Response)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestGetResumeAllResponse(t *testing.T) {
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()

	query := "SELECT \\* FROM \"main\".\"response\" WHERE resume_id = (.*)"
	mock.ExpectQuery(query).
		WithArgs(dummies.Response.ResumeID).
		WillReturnRows(sqlmock.NewRows(dummies.RespColumns).AddRow(dummies.RespRow...))

	resp, err := repo.GetResumeAllResponse(dummies.Response.ID)
	assert.Nil(t, err)
	assert.Equal(t, resp, dummies.ListResp)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(dummies.Response.ResumeID).WillReturnError(error)

	resp, err = repo.GetResumeAllResponse(dummies.Response.ID)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestGetVacancyAllResponse(t *testing.T) {
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()

	query := "SELECT \\* FROM \"main\".\"response\" WHERE vacancy_id = (.*)"
	mock.ExpectQuery(query).
		WithArgs(dummies.Response.VacancyID).
		WillReturnRows(sqlmock.NewRows(dummies.RespColumns).AddRow(dummies.RespRow...))

	resp, err := repo.GetVacancyAllResponse(dummies.Response.ID)
	assert.Nil(t, err)
	assert.Equal(t, resp, dummies.ListResp)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(dummies.Response.VacancyID).WillReturnError(error)

	resp, err = repo.GetVacancyAllResponse(dummies.Response.ID)
	assert.Nil(t, resp)
	assert.Error(t, err)
}

func TestGetAllResumeWithoutResponse(t *testing.T) {
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()
	query := "select main.resume.* from main.resume left join main.response on main.response.resume_id = main.resume.resume_id where cand_id = (.*) group by main.resume.resume_id having sum"
	mock.ExpectQuery(query).
		WithArgs(dummies.Response.ID, dummies.Response.VacancyID).
		WillReturnRows(sqlmock.NewRows(dummies.ResumeColumns).AddRow(dummies.ResumeRow...))

	resume, err := repo.GetAllResumeWithoutResponse(dummies.Response.ID, dummies.Response.VacancyID)
	assert.Nil(t, err)
	assert.Equal(t, resume, dummies.ListResume)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(dummies.Response.VacancyID).WillReturnError(error)

	resume, err = repo.GetAllResumeWithoutResponse(dummies.Response.ID, dummies.Response.VacancyID)
	assert.Nil(t, resume)
	assert.Error(t, err)
}

func TestGetAllVacancyWithoutResponse(t *testing.T) {
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()
	query := "select main.vacancy.* from main.vacancy left join main.response on main.response.vacancy_id = main.vacancy.vac_id where empl_id = (.*) group by main.vacancy.vac_id having sum"
	mock.ExpectQuery(query).
		WithArgs(dummies.Response.ID, dummies.Response.ResumeID).
		WillReturnRows(sqlmock.NewRows(dummies.VacancyColumns).AddRow(dummies.VacancyRow...))

	vac, err := repo.GetAllVacancyWithoutResponse(dummies.Response.ID, dummies.Response.ResumeID)
	assert.Nil(t, err)
	assert.Equal(t, vac, dummies.ListVacancy)

	// Error flow
	error := errors.New("test error")
	mock.ExpectQuery(query).WithArgs(dummies.Response.VacancyID).WillReturnError(error)

	vac, err = repo.GetAllVacancyWithoutResponse(dummies.Response.ID, dummies.Response.ResumeID)
	assert.Nil(t, vac)
	assert.Error(t, err)
}

func TestGetResponsesCnt(t *testing.T) {
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()
	query := "select count(.*) from main.resume join main.response using(.*) .*"
	var a uint = 0
	mock.ExpectQuery(query).
		WithArgs(dummies.Response.ID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(a))

	vac, err := repo.GetResponsesCnt(dummies.Response.ID, common.Candidate)
	assert.Nil(t, err)
	assert.Equal(t, vac, a)

	query2 := "select count(.*) from main.vacancy join main.response on vacancy_id=vac_id .*"
	mock.ExpectQuery(query2).
		WithArgs(dummies.Response.ID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(a))
	vac, err = repo.GetResponsesCnt(dummies.Response.ID, common.Employer)
	assert.Nil(t, err)
	assert.Equal(t, vac, a)
}

func TestGetRespNotifications(t *testing.T) {
	//for resume
	repo, mock, _ := beforeTest2(t)
	dummies := makeDummies()
	query := "UPDATE \"main\".\"response\" SET \"unread\"=(.*) WHERE resume_id = (.*) and response_id IN .*"
	ID2 := uuid.New()
	mock.ExpectExec(query).
		WithArgs(false, ID.String(), ID.String(), ID2.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	query2 := "SELECT \\* FROM \"main\".\"response\" WHERE resume_id = (.*) and unread = (.*)"
	mock.ExpectQuery(query2).
		WithArgs(ID, true).
		WillReturnRows(sqlmock.NewRows(dummies.RespColumns).AddRow(dummies.RespRow...))

	vac, err := repo.GetRespNotifications([]uuid.UUID{ID, ID2}, ID, common.Resume)
	assert.Nil(t, err)
	assert.Equal(t, vac, dummies.ListResp)

	//for vacancy
	query = "UPDATE \"main\".\"response\" SET \"unread\"=(.*) WHERE vacancy_id = (.*) and response_id IN .*"
	mock.ExpectExec(query).
		WithArgs(false, ID.String(), ID.String(), ID2.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	query2 = "SELECT \\* FROM \"main\".\"response\" WHERE vacancy_id = (.*) and unread = (.*)"
	mock.ExpectQuery(query2).
		WithArgs(ID, true).
		WillReturnRows(sqlmock.NewRows(dummies.RespColumns).AddRow(dummies.RespRow...))

	vac, err = repo.GetRespNotifications([]uuid.UUID{ID, ID2}, ID, common.Vacancy)
	assert.Nil(t, err)
	assert.Equal(t, vac, dummies.ListResp)

	//error flow
	error := errors.New("test error")
	mock.ExpectExec(query).
		WithArgs(false, ID.String(), ID.String(), ID2.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery(query).WithArgs(false, ID.String(), ID.String(), ID2.String()).WillReturnError(error)

	vac, err = repo.GetRespNotifications([]uuid.UUID{ID, ID2}, ID, common.Vacancy)
	assert.Nil(t, vac)
	assert.Error(t, err)

	mock.ExpectExec(query).WithArgs(false, ID.String(), ID.String(), ID2.String()).WillReturnError(error)

	vac, err = repo.GetRespNotifications([]uuid.UUID{ID, ID2}, ID, common.Vacancy)
	assert.Nil(t, vac)
	assert.Error(t, err)
}

func TestGetRecommendedVacCnt(t *testing.T) {
	repo, mock, vacMock := beforeTest2(t)
	pair := vacancy2.Pair{
		SphereInd: 1,
		Score:     2,
	}
	pairs := []vacancy2.Pair{pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair,
		pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair}
	vacMock.On("GetPreferredSpheres", ID).Return(pairs, nil).Twice()
	query := "select count(.*) from main.vacancy where date(.*) >= (.*) and sphere in .*"
	var count uint = 0
	startDate := "20-20-2020"
	step := 2
	curSphere := 0
	for curSphere < vacancy2.CountSpheres {
		mock.ExpectQuery(query).
			WithArgs(startDate, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
		curSphere += step
	}
	vac, err := repo.GetRecommendedVacCnt(ID, startDate)

	assert.Nil(t, err)
	assert.Equal(t, vac, count)

	//error flow
	//error := errors.New("test error")
	//mock.ExpectExec(query).WithArgs(startDate, 1, 1).WillReturnError(error)
	//
	//vac, err = repo.GetRecommendedVacCnt(ID, startDate)
	//assert.Nil(t, vac)
	//assert.Error(t, err)

	vacMock.On("GetPreferredSpheres", ID).Return(nil, assert.AnError)
	vac, err = repo.GetRecommendedVacCnt(ID, startDate)
	assert.Error(t, err)
	assert.Equal(t, vac, uint(0))
}

func TestGetRecommendedVacancies(t *testing.T) {
	repo, mock, vacMock := beforeTest2(t)
	dummies := makeDummies()

	start := 0
	limit := 1
	pair := vacancy2.Pair{
		SphereInd: 1,
		Score:     2,
	}
	pairs := []vacancy2.Pair{pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair,
		pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair, pair}
	vacMock.On("GetPreferredSpheres", ID).Return(pairs, nil).Once()

	query := "SELECT \\* FROM \"main\".\"vacancy\" WHERE date(.*) >= (.*) and sphere in (.*) LIMIT (.*)"

	startDate := "20-20-2020"
	step := 2
	curSphere := 0
	for curSphere < vacancy2.CountSpheres {
		mock.ExpectQuery(query).
			WithArgs(startDate, 1, 1).
			WillReturnRows(sqlmock.NewRows(dummies.VacancyColumns).AddRow(dummies.VacancyRow...))
		curSphere += step
	}
	vac, err := repo.GetRecommendedVacancies(ID, start, limit, startDate)

	assert.Nil(t, err)
	assert.Equal(t, vac, dummies.ListVacancy)

	//error flow
	vacMock.On("GetPreferredSpheres", ID).Return(nil, assert.AnError)
	vac, err = repo.GetRecommendedVacancies(ID, start, limit, startDate)
	assert.Error(t, err)
	assert.Nil(t, vac)
}
