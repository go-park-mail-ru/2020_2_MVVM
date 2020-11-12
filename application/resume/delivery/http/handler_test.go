package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocksCommon "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mocksExp "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/custom_experience"
	mocksEduc "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/education"
	mocksResume "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/resume"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

const (
	resumeUrlGroup = "/api/v1/resume/"
)

var testData struct {
	resumeHandler *ResumeHandler
	router        *gin.Engine
	mockUseCase   *mocksResume.UseCase
	mockUSEduc    *mocksEduc.UseCase
	mockUSExp     *mocksExp.UseCase
	mockAuth      *mocksCommon.AuthTest
	httpStatus    []int
	resList       []models.Resume
	resBriefList  []models.BriefResumeInfo
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
	testData.resList = []models.Resume{
		{Title: "title1", Description: "description1"},
		{Title: "title2", Description: "description2"},
		{Title: "title3", Description: "description3"}}
	testData.resBriefList = []models.BriefResumeInfo{
		{Title: "title1", Description: "description1"},
		{Title: "title2", Description: "description2"},
		{Title: "title3", Description: "description3"}}
	testData.mockUseCase = new(mocksResume.UseCase)
	testData.mockUSEduc = new(mocksEduc.UseCase)
	testData.mockUSExp = new(mocksExp.UseCase)
	testData.router = gin.Default()
	api := testData.router.Group("/api/v1")
	testData.resumeHandler = NewRest(api.Group("/resume"), testData.mockUseCase, testData.mockUSEduc, testData.mockUSExp, nil)
}

func getRespStruct(entity interface{}) interface{} {
	switch entity.(type) {
	case models.Resume:
		resume := entity.(models.Resume)
		return &resume
	case []models.Resume:
		resumeList := entity.([]models.Resume)
		return resumeList
	case []models.BriefResumeInfo:
		resumeList := entity.([]models.BriefResumeInfo)
		return resumeList
	case string:
		err := entity.(string)
		return common.RespError{Err: err}
	case error:
		err := entity.(error)
		return common.RespError{Err: err.Error()}
	}
	return nil
}

func TestGetResumePageHandler(t *testing.T) {
	var start uint = 0
	var end uint = 2
	v, r, mockUseCase := testData.resumeHandler, testData.router, testData.mockUseCase
	r.GET("/page", v.GetResumePage)

	mockUseCase.On("List", start, end).Return(testData.resBriefList, nil)
	mockUseCase.On("List", start, uint(1000)).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%spage?start=%d&limit=%d", resumeUrlGroup, start, end),
		fmt.Sprintf("%spage", resumeUrlGroup),
		fmt.Sprintf("%spage?start=%d&limit=%d", resumeUrlGroup, start, uint(1000)),
	}

	testExpectedBody := []interface{}{testData.resBriefList, common.EmptyFieldErr, common.DataBaseErr}
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

func TestGetResumeSearchHandler(t *testing.T) {
	v, r, mockUseCase := testData.resumeHandler, testData.router, testData.mockUseCase
	r.GET("/search", v.GetResumePage)

	params := resume.SearchParams{
		KeyWords:        nil,
		SalaryMin:       nil,
		SalaryMax:       nil,
		Gender:          nil,
		EducationLevel:  nil,
		CareerLevel:     nil,
		ExperienceMonth: nil,
		AreaSearch:      nil,
	}
	mockUseCase.On("Search", params).Return(testData.resBriefList, nil)

	testUrls := []string{
		fmt.Sprintf("%ssearch", resumeUrlGroup),
		fmt.Sprintf("%ssearch", resumeUrlGroup),
	}

	testExpectedBody := []interface{}{testData.resBriefList, common.EmptyFieldErr}
	testParamsForPost := []interface{}{params, nil}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
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
