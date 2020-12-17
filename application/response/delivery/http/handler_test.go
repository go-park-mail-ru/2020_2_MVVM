package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	response2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/response"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocksCommon "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mResponse "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/response"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"os"
	"testing"
)

const (
	responseUrlGroup = "/api/v1/response/"
)

type TestData struct {
	responseHandler *ResponseHandler
	router          *gin.Engine
	mockUseCase     *mResponse.IUseCaseResponse
	mockSB          *mocksCommon.SessionBuilder
	mockSession     *mocksCommon.Session
	httpStatus      []int
}

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}

func beforeTest() TestData {
	testData := TestData{}
	gin.SetMode(gin.TestMode)
	testData.httpStatus = []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusInternalServerError,
	}
	testData.mockUseCase = new(mResponse.IUseCaseResponse)
	testData.mockSB = new(mocksCommon.SessionBuilder)
	testData.mockSession = new(mocksCommon.Session)
	testData.router = gin.Default()
	api := testData.router.Group("api/v1")
	testData.responseHandler = NewRest(api.Group("/response"), testData.mockUseCase, testData.mockSB, func(context *gin.Context) {})
	return testData
}

func getRespStruct(entity interface{}) interface{} {
	switch entity.(type) {
	case response2.RespNotifications:
		notify := entity.(response2.RespNotifications)
		return &notify
	case models.Response:
		response := entity.(models.Response)
		return &response
	case []models.ResponseWithTitle:
		responseList := entity.([]models.ResponseWithTitle)
		return &responseList
	case []models.BriefResumeInfo:
		resumeList := entity.([]models.BriefResumeInfo)
		return &resumeList
	case []models.Vacancy:
		vacancyList := entity.([]models.Vacancy)
		return &vacancyList
	case string:
		err := entity.(string)
		return models.RespError{Err: err}
	case error:
		err := entity.(error)
		return models.RespError{Err: err.Error()}
	}
	return nil
}

func TestCreateResponse(t *testing.T) {
	td := beforeTest()
	td.router.POST("/", td.responseHandler.CreateResponse)
	ID := uuid.New()
	response := models.Response{
		ID:        ID,
		ResumeID:  ID,
		VacancyID: ID,
		Initial:   common.Employer,
	}
	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetCandID").Return(uuid.Nil).Once()
	td.mockSession.On("GetEmplID").Return(ID).Once()
	td.mockUseCase.On("CreateChatAndTechChat", response).Return(&response, nil).Once()

	response2 := response
	response2.Initial = common.Candidate
	td.mockSession.On("GetCandID").Return(ID).Once()
	td.mockSession.On("GetEmplID").Return(uuid.Nil).Once()
	td.mockUseCase.On("CreateChatAndTechChat", response2).Return(&response2, nil).Once()

	td.mockSession.On("GetCandID").Return(uuid.Nil).Once()
	td.mockSession.On("GetEmplID").Return(uuid.Nil).Once()

	td.mockSession.On("GetCandID").Return(ID).Once()
	td.mockSession.On("GetEmplID").Return(uuid.Nil).Once()
	td.mockUseCase.On("CreateChatAndTechChat", response2).Return(nil, assert.AnError).Once()

	testUrls := []string{
		responseUrlGroup,
		responseUrlGroup,
		responseUrlGroup,
		responseUrlGroup,
		responseUrlGroup,
	}
	httpStatus := []int{
		http.StatusOK,
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusMethodNotAllowed,
		http.StatusInternalServerError,
	}

	testExpectedBody := []interface{}{response, response2, common.EmptyFieldErr, common.AuthRequiredErr, common.DataBaseErr}
	testParamsForPost := []interface{}{response, response2, nil, response, response2}

	for i := range testUrls {
		t.Run("test responses on different urls for create response handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, testUrls[i], testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateStatusResponse(t *testing.T) {
	td := beforeTest()
	td.router.POST("/update", td.responseHandler.UpdateStatus)
	ID := uuid.New()
	response := models.Response{
		ID:        ID,
		ResumeID:  ID,
		VacancyID: ID,
		Initial:   common.Employer,
	}
	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetCandID").Return(uuid.Nil).Once()
	td.mockSession.On("GetEmplID").Return(ID).Once()
	td.mockUseCase.On("UpdateStatus", response, common.Employer).Return(&response, nil).Once()

	response2 := response
	response2.Initial = common.Candidate
	td.mockSession.On("GetCandID").Return(ID).Once()
	td.mockSession.On("GetEmplID").Return(uuid.Nil).Once()
	td.mockUseCase.On("UpdateStatus", response2, common.Candidate).Return(&response2, nil).Once()

	td.mockSession.On("GetCandID").Return(uuid.Nil).Once()
	td.mockSession.On("GetEmplID").Return(uuid.Nil).Once()

	td.mockSession.On("GetCandID").Return(ID).Once()
	td.mockSession.On("GetEmplID").Return(uuid.Nil).Once()
	td.mockUseCase.On("UpdateStatus", response2, common.Candidate).Return(nil, assert.AnError).Once()

	testUrls := [5]string{
		fmt.Sprintf("%supdate", responseUrlGroup),
		fmt.Sprintf("%supdate", responseUrlGroup),
		fmt.Sprintf("%supdate", responseUrlGroup),
		fmt.Sprintf("%supdate", responseUrlGroup),
		fmt.Sprintf("%supdate", responseUrlGroup),
	}
	httpStatus := []int{
		http.StatusOK,
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusMethodNotAllowed,
		http.StatusInternalServerError,
	}

	testExpectedBody := []interface{}{response, response2, common.EmptyFieldErr, common.AuthRequiredErr, common.DataBaseErr}
	testParamsForPost := []interface{}{response, response2, nil, response, response2}

	for i := range testUrls {
		t.Run("test responses on different urls for create response handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, testUrls[i], testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetAllResponses(t *testing.T) {
	td := beforeTest()
	td.router.GET("/my", td.responseHandler.handlerGetAllResponses)
	ID := uuid.New()

	response := models.ResponseWithTitle{ResponseID: ID}
	listResp := []models.ResponseWithTitle{response}

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetCandID").Return(uuid.Nil).Twice()
	td.mockSession.On("GetEmplID").Return(ID).Twice()
	td.mockUseCase.On("GetAllEmployerResponses", ID, []uuid.UUID(nil)).Return(listResp, nil).Once()
	td.mockUseCase.On("GetAllEmployerResponses", ID, []uuid.UUID(nil)).Return(nil, assert.AnError).Once()

	td.mockSession.On("GetCandID").Return(ID).Twice()
	td.mockSession.On("GetEmplID").Return(uuid.Nil).Twice()
	td.mockUseCase.On("GetAllCandidateResponses", ID, []uuid.UUID(nil)).Return(listResp, nil).Once()
	td.mockUseCase.On("GetAllCandidateResponses", ID, []uuid.UUID(nil)).Return(nil, assert.AnError).Once()

	td.mockSession.On("GetCandID").Return(uuid.Nil).Once()
	td.mockSession.On("GetEmplID").Return(uuid.Nil).Once()

	testUrls := []string{
		fmt.Sprintf("%smy", responseUrlGroup),
		fmt.Sprintf("%smy", responseUrlGroup),
		fmt.Sprintf("%smy", responseUrlGroup),
		fmt.Sprintf("%smy", responseUrlGroup),
		fmt.Sprintf("%smy", responseUrlGroup),
	}
	httpStatus := []int{
		http.StatusOK,
		http.StatusInternalServerError,
		http.StatusOK,
		http.StatusInternalServerError,
		http.StatusMethodNotAllowed,
	}
	testExpectedBody := []interface{}{
		listResp,
		common.DataBaseErr,
		listResp,
		common.DataBaseErr,
		"this user cannot have responses",
	}

	for i := range testUrls {
		t.Run("test responses on different urls for create response handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetAllResumeWithoutResponse(t *testing.T) {
	td := beforeTest()
	td.router.GET("/free/resumes/:entity_id", td.responseHandler.handlerGetAllResumeWithoutResponse)
	ID := uuid.New()
	resume := models.BriefResumeInfo{
		ResumeID: ID,
		CandID:   ID,
		UserID:   ID,
	}
	listResume := []models.BriefResumeInfo{resume}

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetCandID").Return(ID).Once()
	td.mockUseCase.On("GetAllResumeWithoutResponse", ID, ID).Return(listResume, nil).Once()

	td.mockSession.On("GetCandID").Return(ID).Once()
	td.mockUseCase.On("GetAllResumeWithoutResponse", ID, uuid.Nil).Return(nil, assert.AnError).Once()

	testUrls := []string{
		fmt.Sprintf("%sfree/resumes/%s", responseUrlGroup, ID),
		fmt.Sprintf("%sfree/resumes/invalidid", responseUrlGroup),
		fmt.Sprintf("%sfree/resumes/%s", responseUrlGroup, uuid.Nil),
	}
	httpStatus := []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusInternalServerError,
		//http.StatusInternalServerError,
	}
	testExpectedBody := []interface{}{
		listResume,
		common.EmptyFieldErr,
		common.DataBaseErr,
	}

	for i := range testUrls {
		t.Run("test responses on different urls for create response handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetAllVacancyWithoutResponse(t *testing.T) {
	td := beforeTest()
	td.router.GET("/free/vacancies/:entity_id", td.responseHandler.handlerGetAllResumeWithoutResponse)
	ID := uuid.New()
	vacancy := models.Vacancy{ID: ID}
	listVacancy := []models.Vacancy{vacancy}

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetEmplID").Return(ID).Once()
	td.mockUseCase.On("GetAllVacancyWithoutResponse", ID, ID).Return(listVacancy, nil).Once()

	td.mockSession.On("GetEmplID").Return(ID).Once()
	td.mockUseCase.On("GetAllVacancyWithoutResponse", ID, ID).Return(nil, assert.AnError).Once()

	testUrls := []string{
		fmt.Sprintf("%sfree/vacancies/%s", responseUrlGroup, ID),
		fmt.Sprintf("%sfree/vacancies/invalidID", responseUrlGroup),
		fmt.Sprintf("%sfree/vacancies/%s", responseUrlGroup, ID),
	}
	httpStatus := []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusInternalServerError,
		//http.StatusMethodNotAllowed,
	}
	testExpectedBody := []interface{}{
		listVacancy,
		common.EmptyFieldErr,
		common.DataBaseErr,
	}

	for i := range testUrls {
		t.Run("test responses on different urls for create response handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetAllNotifications(t *testing.T) {
	td := beforeTest()
	td.router.POST("/notify", td.responseHandler.handlerGetAllNotifications)
	ID := uuid.New()
	var count uint = 2
	var start uint = 0
	var end uint = 1

	req1 := response2.ReqNotify{
		VacInLastNDays:       nil,
		OnlyVacCnt:           true,
		ListStart:            0,
		ListEnd:              1,
		NewRespNotifications: nil,
		OnlyRespCnt:          true,
	}

	resp1 := response2.RespNotifications{
		UnreadResp:        nil,
		UnreadRespCnt:     count,
		RecommendedVac:    nil,
		RecommendedVacCnt: 0,
	}

	//count flow
	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetEmplID").Return(ID)
	td.mockSession.On("GetCandID").Return(nil)
	td.mockSession.On("GetUserID").Return(ID)
	td.mockUseCase.On("GetResponsesCnt", ID, common.EmplID).Return(count, nil).Once()
	td.mockUseCase.On("GetRecommendedVacCnt", ID, common.Week).Return(count, nil).Once()

	//vacancies flow
	vacancy := models.Vacancy{ID: ID}
	listVacancy := []models.Vacancy{vacancy, vacancy}
	response := models.ResponseWithTitle{ResponseID: ID}
	listResp := []models.ResponseWithTitle{response, response}
	var days int = 1
	req2 := response2.ReqNotify{
		VacInLastNDays:       &days,
		OnlyVacCnt:           false,
		ListStart:            start,
		ListEnd:              end,
		NewRespNotifications: []uuid.UUID{ID},
		OnlyRespCnt:          false,
	}

	resp2 := response2.RespNotifications{
		UnreadResp:        listResp,
		UnreadRespCnt:     count,
		RecommendedVac:    listVacancy,
		RecommendedVacCnt: count,
	}

	td.mockUseCase.On("GetAllEmployerResponses", ID, []uuid.UUID{ID}).Return(listResp, nil).Once()
	td.mockUseCase.On("GetRecommendedVacancies", ID, start, end, 1).Return(listVacancy, nil).Once()

	testUrls := []string{
		fmt.Sprintf("%snotify", responseUrlGroup),
		fmt.Sprintf("%snotify", responseUrlGroup),
	}
	httpStatus := []int{
		http.StatusOK,
		http.StatusOK,
	}
	testExpectedBody := []interface{}{resp1, resp2}
	testParamsForPost := []interface{}{req1, req2}

	for i := range testUrls {
		t.Run("test responses on different urls for create response handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, testUrls[i], testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetAllNotificationsErrors(t *testing.T) {
	td := beforeTest()
	td.router.POST("/notify", td.responseHandler.handlerGetAllNotifications)
	ID := uuid.New()

	req1 := response2.ReqNotify{
		VacInLastNDays:       nil,
		OnlyVacCnt:           true,
		ListStart:            0,
		ListEnd:              1,
		NewRespNotifications: nil,
		OnlyRespCnt:          true,
	}

	//not allowed
	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetEmplID").Return(nil)
	td.mockSession.On("GetCandID").Return(nil).Once()

	//nil body
	td.mockSession.On("GetCandID").Return(ID)

	testUrls := []string{
		fmt.Sprintf("%snotify", responseUrlGroup),
		fmt.Sprintf("%snotify", responseUrlGroup),
	}
	httpStatus := []int{
		http.StatusMethodNotAllowed,
		http.StatusBadRequest,
	}
	testExpectedBody := []interface{}{
		common.AuthRequiredErr,
		common.EmptyFieldErr,
	}

	testParamsForPost := []interface{}{req1, nil}

	for i := range testUrls {
		t.Run("test responses on different urls for create response handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, testUrls[i], testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}
