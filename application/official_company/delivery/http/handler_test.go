package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/microservises/auth"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocksCommon "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mocks "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/official_company"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"os"
	"testing"
)

const (
	compUrlGroup = "/api/v1/company/"
)

var testData struct {
	compHandler *CompanyHandler
	router      *gin.Engine
	sessionInfo auth.SessionInfo
	mockSB      *mocksCommon.SessionBuilder
	mockSession *mocksCommon.Session
	mockUseCase *mocks.IUseCaseOfficialCompany
	httpStatus  []int
	compList    []models.OfficialCompany
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
	testData.compList = []models.OfficialCompany{
		{Name: "name1", Description: "description1", Link: "link1", AreaSearch: "area1"},
		{Name: "name2", Description: "description2", Link: "link2", AreaSearch: "area2"},
		{Name: "name3", Description: "description3", Link: "link3", AreaSearch: "area3"},
	}
	testData.mockUseCase = new(mocks.IUseCaseOfficialCompany)
	testData.mockSB = new(mocksCommon.SessionBuilder)
	testData.mockSession = new(mocksCommon.Session)
	testData.router = gin.Default()
	api := testData.router.Group("api/v1")
	testData.compHandler = NewRest(api.Group("/company"), testData.mockUseCase, testData.mockSB, func(context *gin.Context) {})
}

func getRespStruct(entity interface{}) interface{} {
	switch entity.(type) {
	case models.OfficialCompany:
		comp := entity.(models.OfficialCompany)
		return models.Resp{&comp}
	case []models.OfficialCompany:
		compList := entity.([]models.OfficialCompany)
		return models.RespList{compList}
	case string:
		err := entity.(string)
		return models.RespError{Err: err}
	case error:
		err := entity.(error)
		return models.RespError{Err: err.Error()}
	}
	return nil
}

func TestGetCompanyHandler(t *testing.T) {
	c, r, mockUseCase := testData.compHandler, testData.router, testData.mockUseCase
	r.GET("/by/id/:company_id", c.GetCompanyHandler)
	compID := uuid.New()
	comp := models.OfficialCompany{ID: compID}
	mockUseCase.On("GetOfficialCompany", compID).Return(&comp, nil)
	mockUseCase.On("GetOfficialCompany", uuid.Nil).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%sby/id/%s", compUrlGroup, compID),
		fmt.Sprintf("%sby/id/invalidUuid", compUrlGroup),
		fmt.Sprintf("%sby/id/%s", compUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{comp, common.EmptyFieldErr, common.DataBaseErr}

	for i := range testUrls {
		t.Run("test responses on different urls for getCompany handler", func(t *testing.T) {
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

func TestGetCompanyListHandler(t *testing.T) {
	var start uint = 0
	var end uint = 3
	c, r, mockUseCase := testData.compHandler, testData.router, testData.mockUseCase
	r.GET("/page", c.GetCompanyListHandler)

	mockUseCase.On("GetCompaniesList", start, end).Return(testData.compList, nil)
	mockUseCase.On("GetCompaniesList", end, end).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%spage?start=%d&limit=%d", compUrlGroup, start, end),
		fmt.Sprintf("%spage", compUrlGroup),
		fmt.Sprintf("%spage?start=%d&limit=%d", compUrlGroup, end, end),
	}

	testExpectedBody := []interface{}{testData.compList, common.EmptyFieldErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getCompanyList handler", func(t *testing.T) {
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

func TestSearchCompaniesHandler(t *testing.T) {
	c, r, mockUseCase := testData.compHandler, testData.router, testData.mockUseCase
	r.POST("/search", c.SearchCompaniesHandler)

	params := models.CompanySearchParams{
		AreaSearch: []string{"area1", "area2"},
	}
	compList := testData.compList[:2]
	paramsEmpty := models.CompanySearchParams{}
	mockUseCase.On("SearchCompanies", params).Return(compList, nil)
	mockUseCase.On("SearchCompanies", paramsEmpty).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%ssearch", compUrlGroup),
		fmt.Sprintf("%ssearch", compUrlGroup),
		fmt.Sprintf("%ssearch", compUrlGroup),
	}

	testExpectedBody := []interface{}{compList, common.EmptyFieldErr, common.DataBaseErr}
	testParamsForPost := []interface{}{params, nil, paramsEmpty}
	for i := range testUrls {
		t.Run("test responses on different urls for SearchCompanies handler", func(t *testing.T) {
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

func TestGetUserCompanyHandler(t *testing.T) {
	c, r, mockUseCase := testData.compHandler, testData.router, testData.mockUseCase
	r.GET("/mine", c.GetUserCompanyHandler)
	empID := uuid.New()
	comp := models.OfficialCompany{ID: empID}
	mockUseCase.On("GetMineCompany", empID).Return(&comp, nil)
	mockUseCase.On("GetMineCompany", uuid.Nil).Return(nil, assert.AnError).Once()
	testData.mockSB.On("Build", mock.Anything).Return(testData.mockSession)
	testData.mockSession.On("GetEmplID").Return(empID).Once()
	testData.mockSession.On("GetEmplID").Return(uuid.Nil).Once()
	testUrl := fmt.Sprintf("%smine", compUrlGroup)
	testExpectedBody := []interface{}{comp, common.SessionErr}

	for i := range testExpectedBody {
		t.Run("test responses on different urls for getUserCompany handler", func(t *testing.T) {
			w, err := general.PerformRequest(r, http.MethodGet, testUrl, nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}