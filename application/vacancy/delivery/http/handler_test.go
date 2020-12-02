package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/dto/microservises/auth"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/dto/models"
	vacancy2 "github.com/go-park-mail-ru/2020_2_MVVM.git/dto/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocksCommon "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mocksAuth "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/microservices/auth/authmicro"
	mocksVacancy "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/microservices/vacancy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"os"
	"testing"
)

const (
	vacUrlGroup = "/api/v1/vacancy/"
)

var testData struct {
	vacHandler    *VacancyHandler
	router        *gin.Engine
	sessionInfo   auth.SessionInfo
	mockSB        *mocksCommon.SessionBuilder
	mockSession   *mocksCommon.Session
	mockAuth      *mocksAuth.AuthClient
	mockVacClient *mocksVacancy.VacClient
	httpStatus    []int
	vacList       []models.Vacancy
}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	os.Exit(retCode)
}

func setUp() {
	gin.SetMode(gin.TestMode)
	testData.httpStatus = []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusInternalServerError,
	}
	testData.vacList = []models.Vacancy{
		{Title: "title1", Description: "description1", AreaSearch: "area1"},
		{Title: "title2", Description: "description2", AreaSearch: "area2"},
		{Title: "title3", Description: "description3", AreaSearch: "area3"}}
	testData.mockSB = new(mocksCommon.SessionBuilder)
	testData.mockSession = new(mocksCommon.Session)
	testData.mockVacClient = new(mocksVacancy.VacClient)
	testData.mockAuth = new(mocksAuth.AuthClient)
	testData.router = gin.Default()
	testData.vacHandler = NewRest(testData.router.Group(vacUrlGroup), testData.mockSB,
		func(context *gin.Context) {}, testData.mockVacClient, testData.mockAuth)
}

func getRespStruct(entity interface{}) interface{} {
	switch entity.(type) {
	case models.Vacancy:
		vac := entity.(models.Vacancy)
		return vacancy2.Resp{&vac}
	case []models.Vacancy:
		vacList := entity.([]models.Vacancy)
		return vacancy2.RespList{vacList}
	case string:
		err := entity.(string)
		return models.RespError{Err: err}
	case error:
		err := entity.(error)
		return models.RespError{Err: err.Error()}
	}
	return nil
}

func TestGetVacancyByIdHandler(t *testing.T) {
	r, mockVacClient := testData.router, testData.mockVacClient
	vacID := uuid.New()
	vac := models.Vacancy{ID: vacID}
	mockVacClient.On("GetVacancy", vacID).Return(&vac, nil).Once()
	mockVacClient.On("GetVacancy", uuid.Nil).Return(nil, assert.AnError).Once()
	testData.mockSB.On("Build", mock.Anything).Return(testData.mockSession)
	testData.mockSession.On("GetCandID").Return(uuid.New()).Once()
	testData.mockVacClient.On("AddRecommendation", mock.Anything, mock.Anything).Return(nil).Once()

	testUrls := []string{
		fmt.Sprintf("%sby/id/%s", vacUrlGroup, vacID),
		fmt.Sprintf("%sby/id/invalidUuid", vacUrlGroup),
		fmt.Sprintf("%sby/id/%s", vacUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{vac, common.EmptyFieldErr, common.DataBaseErr}

	for i := range testUrls {
		t.Run("test responses on different urls for getVacancy handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetVacancyListHandler(t *testing.T) {
	var start uint = 0
	var end uint = 3
	r, mockVacClient := testData.router, testData.mockVacClient

	mockVacClient.On("GetVacancyList", start, end, uuid.Nil, vacancy.ByVacId).Return(testData.vacList, nil).Once()
	mockVacClient.On("GetVacancyList", end, end, uuid.Nil, vacancy.ByVacId).Return(nil, assert.AnError).Once()

	testUrls := []string{
		fmt.Sprintf("%spage?start=%d&limit=%d", vacUrlGroup, start, end),
		fmt.Sprintf("%spage", vacUrlGroup),
		fmt.Sprintf("%spage?start=%d&limit=%d", vacUrlGroup, end, end),
	}

	testExpectedBody := []interface{}{testData.vacList, common.EmptyFieldErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetCompVacancyListHandler(t *testing.T) {
	var start uint = 0
	var end uint = 3
	r, mockVacClient := testData.router, testData.mockVacClient
	compID := uuid.New()
	mockVacClient.On("GetVacancyList", start, end, compID, vacancy.ByCompId).Return(testData.vacList, nil).Once()
	mockVacClient.On("GetVacancyList", end, end, compID, vacancy.ByCompId).Return(nil, assert.AnError).Once()

	testUrls := []string{
		fmt.Sprintf("%scomp?start=%d&limit=%d&comp_id=%s", vacUrlGroup, start, end, compID),
		fmt.Sprintf("%scomp", vacUrlGroup),
		fmt.Sprintf("%scomp?start=%d&limit=%d&comp_id=%s", vacUrlGroup, end, end, compID),
	}

	testExpectedBody := []interface{}{testData.vacList, common.EmptyFieldErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getCompVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSearchVacanciesHandler(t *testing.T) {
	r, mockVacClient := testData.router, testData.mockVacClient

	params := models.VacancySearchParams{
		AreaSearch: []string{"area1", "area2"},
	}
	vacList := testData.vacList[:2]
	paramsEmpty := models.VacancySearchParams{}
	mockVacClient.On("SearchVacancies", params).Return(vacList, nil).Once()
	mockVacClient.On("SearchVacancies", paramsEmpty).Return(nil, assert.AnError).Once()

	testUrls := []string{
		fmt.Sprintf("%ssearch", vacUrlGroup),
		fmt.Sprintf("%ssearch", vacUrlGroup),
		fmt.Sprintf("%ssearch", vacUrlGroup),
	}

	testExpectedBody := []interface{}{vacList, common.EmptyFieldErr, common.DataBaseErr}
	testParamsForPost := []interface{}{params, nil, paramsEmpty}
	for i := range testUrls {
		t.Run("test responses on different urls for SearchVacancies handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodPost, testUrls[i], testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateVacancyHandler(t *testing.T) {
	r, mockVacClient := testData.router, testData.mockVacClient
	vacNew := models.Vacancy{Title: "abcde", Description: "b"}

	mockVacClient.On("CreateVacancy", mock.Anything).Return(&vacNew, nil).Once()
	mockVacClient.On("CreateVacancy", mock.Anything).Return(nil, assert.AnError).Once()
	testData.mockSB.On("Build", mock.Anything).Return(testData.mockSession)
	testData.mockSession.On("GetEmplID").Return(uuid.New()).Twice()

	testExpectedBody := []interface{}{vacNew, common.EmptyFieldErr, common.DataBaseErr}
	testParamsForPost := []interface{}{vacNew, nil, models.Vacancy{}}
	for i := range testExpectedBody {
		t.Run("test responses on different urls for SearchVacancies handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodPost, vacUrlGroup, testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateVacancyHandler(t *testing.T) {
	r, mockVacClient := testData.router, testData.mockVacClient
	vacNew := models.Vacancy{ID: uuid.New(), Title: "abcde", Description: "b"}

	mockVacClient.On("UpdateVacancy", mock.Anything).Return(&vacNew, nil).Once()
	mockVacClient.On("UpdateVacancy", mock.Anything).Return(nil, assert.AnError).Once()
	testData.mockSB.On("Build", mock.Anything).Return(testData.mockSession)
	testData.mockSession.On("GetEmplID").Return(uuid.New()).Twice()

	testExpectedBody := []interface{}{vacNew, common.EmptyFieldErr, common.DataBaseErr}
	testParamsForPut := []interface{}{vacNew, nil, models.Vacancy{ID: uuid.New()}}
	for i := range testExpectedBody {
		t.Run("test responses on different urls for SearchVacancies handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodPut, vacUrlGroup, testParamsForPut[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetUserVacancyListHandler(t *testing.T) {
	r, mockVacClient := testData.router, testData.mockVacClient
	var start uint = 0
	var end uint = 3

	empID := uuid.New()
	mockVacClient.On("GetVacancyList", start, end, empID, vacancy.ByEmpId).Return(testData.vacList, nil).Once()
	mockVacClient.On("GetVacancyList", end, end, empID, vacancy.ByEmpId).Return(nil, assert.AnError).Once()
	testData.mockSB.On("Build", mock.Anything).Return(testData.mockSession)
	testData.mockSession.On("GetEmplID").Return(empID).Twice()

	testExpectedBody := []interface{}{testData.vacList, common.EmptyFieldErr, common.DataBaseErr}
	testUrls := []string{
		fmt.Sprintf("%smine?start=%d&limit=%d", vacUrlGroup, start, end),
		fmt.Sprintf("%smine", vacUrlGroup),
		fmt.Sprintf("%smine?start=%d&limit=%d", vacUrlGroup, end, end),
	}
	for i := range testExpectedBody {
		t.Run("test responses on different urls for getVacancy handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetRecommendationUserVacancy(t *testing.T) {
	r, mockVacClient := testData.router, testData.mockVacClient
	var start = 0
	var end = 3

	userID := uuid.New()
	mockVacClient.On("GetRecommendation", userID, start, end).Return(testData.vacList, nil).Once()
	mockVacClient.On("GetRecommendation", userID, end, end).Return(nil, assert.AnError).Once()
	testData.mockSB.On("Build", mock.Anything).Return(testData.mockSession)
	testData.mockSession.On("GetCandID").Return(uuid.New()).Twice()
	testData.mockSession.On("GetUserID").Return(userID).Twice()

	testExpectedBody := []interface{}{testData.vacList, common.EmptyFieldErr, common.DataBaseErr}
	testUrls := []string{
		fmt.Sprintf("%srecommendation?start=%d&limit=%d", vacUrlGroup, start, end),
		fmt.Sprintf("%srecommendation", vacUrlGroup),
		fmt.Sprintf("%srecommendation?start=%d&limit=%d", vacUrlGroup, end, end),
	}
	for i := range testExpectedBody {
		t.Run("test responses on different urls for getVacancy handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}
