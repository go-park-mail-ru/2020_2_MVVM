package http

/*
import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocksCommon "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mocksExp "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/custom_experience"
	mocksEduc "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/education"
	mocksResume "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/resume"
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
	vacList       []models.Resume
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
	testData.vacList = []models.Resume{
		{Title: "title1", Description: "description1"},
		{Title: "title2", Description: "description2"},
		{Title: "title3", Description: "description3"}}
	testData.mockUseCase = new(mocksResume.UseCase)
	testData.mockUSEduc = new(mocksEduc.UseCase)
	testData.mockUSExp = new(mocksExp.UseCase)
	testData.router = gin.Default()
	api := testData.router.Group("api/v1")
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
	case string:
		err := entity.(string)
		return common.RespError{Err: err}
	case error:
		err := entity.(error)
		return common.RespError{Err: err.Error()}
	}
	return nil
}

func TestGetResumeListHandler(t *testing.T) {
	var start uint = 0
	var end uint = 3
	v, r, mockUseCase := testData.resumeHandler, testData.router, testData.mockUseCase
	r.GET("/page", v.GetResumePage)

	mockUseCase.On("GetResumePage", start, end).Return(testData.vacList, nil)
	//mockUseCase.On("GetVacancyList", end, end, uuid.Nil, vacancy.ByVacId).Return(nil, assert.AnError)

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
*/
//func TestGetResumeByIdHandler(t *testing.T) {
//	h, r, mockUseCase := testData.resumeHandler, testData.router, testData.mockUseCase
//	r.GET("/by/id/:resume_id", h.GetResumeByID)
//	ID := uuid.New()
//	resume := models.Resume{ResumeID: ID}
//	mockUseCase.On("GetById", ID).Return(&resume, nil)
//	//mockUseCase.On("GetVacancy", uuid.Nil).Return(nil, assert.AnError)
//
//	testUrls := []string{
//		fmt.Sprintf("%sby/id/%s", resumeUrlGroup, ID),
//		//fmt.Sprintf("%sby/id/invalidUuid", resumeUrlGroup),
//		//fmt.Sprintf("%sby/id/%s", resumeUrlGroup, uuid.Nil),
//	}
//	testExpectedBody := []interface{}{resume, common.EmptyFieldErr, common.DataBaseErr}
//
//	for i := range testUrls {
//		t.Run("test responses on different urls for getVacancy handler", func(t *testing.T) {
//			w, err := general.PerformRequest(r, http.MethodGet, testUrls[i], nil)
//			if err != nil {
//				t.Fatalf("Couldn't create request: %v\n", err)
//			}
//			if err := general.ResponseComparator(*w, testData.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
//				t.Fatal(err)
//			}
//		})
//	}
//}
