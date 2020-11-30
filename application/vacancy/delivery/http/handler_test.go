package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/api"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocks2 "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mocks "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/vacancy"
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
	vacHandler  *VacancyHandler
	router      *gin.Engine
	mockUseCase *mocks.IUseCaseVacancy
	sessionInfo api.SessionInfo
	mockSB      *mocks2.SessionBuilder
	mockSession *mocks2.Session
	httpStatus  []int
	vacList     []models.Vacancy
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
	testData.mockUseCase = new(mocks.IUseCaseVacancy)
	//testData.mockAuth = new(mocks2.AuthTest)
	testData.mockSB = new(mocks2.SessionBuilder)
	testData.mockSession = new(mocks2.Session)
	testData.router = gin.Default()
	testData.vacHandler = NewRest(testData.router.Group(vacUrlGroup), testData.mockUseCase, testData.mockSB, nil)
}

func getRespStruct(entity interface{}) interface{} {
	switch entity.(type) {
	case models.Vacancy:
		vac := entity.(models.Vacancy)
		return Resp{&vac}
	case []models.Vacancy:
		vacList := entity.([]models.Vacancy)
		return RespList{vacList}
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
	r, mockUseCase := testData.router, testData.mockUseCase
	vacID := uuid.New()
	vac := models.Vacancy{ID: vacID}
	mockUseCase.On("GetVacancy", vacID).Return(&vac, nil)
	mockUseCase.On("GetVacancy", uuid.Nil).Return(nil, assert.AnError)
	testData.mockSB.On("Build", mock.Anything).Return(testData.mockSession)
	testData.mockSession.On("GetCandID").Return(uuid.New())
	testData.mockSession.On("GetUserID").Return(uuid.New())
	testData.mockUseCase.On("AddRecommendation", mock.Anything, mock.Anything).Return(nil)

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
	r, mockUseCase := testData.router, testData.mockUseCase

	mockUseCase.On("GetVacancyList", start, end, uuid.Nil, vacancy.ByVacId).Return(testData.vacList, nil)
	mockUseCase.On("GetVacancyList", end, end, uuid.Nil, vacancy.ByVacId).Return(nil, assert.AnError)

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
	r, mockUseCase :=testData.router, testData.mockUseCase
	compID := uuid.New()
	mockUseCase.On("GetVacancyList", start, end, compID, vacancy.ByCompId).Return(testData.vacList, nil)
	mockUseCase.On("GetVacancyList", end, end, compID, vacancy.ByCompId).Return(nil, assert.AnError)

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
	r, mockUseCase := testData.router, testData.mockUseCase

	params := models.VacancySearchParams{
		AreaSearch: []string{"area1", "area2"},
	}
	vacList := testData.vacList[:2]
	paramsEmpty := models.VacancySearchParams{}
	mockUseCase.On("SearchVacancies", params).Return(vacList, nil)
	mockUseCase.On("SearchVacancies", paramsEmpty).Return(nil, assert.AnError)

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
