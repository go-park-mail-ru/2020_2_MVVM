package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	resume2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/resume"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocksCommon "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mocksExp "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/custom_experience"
	mocksEduc "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/education"
	mocksResume "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/resume"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"os"
	"testing"
)

const (
	resumeUrlGroup = "/api/v1/resume/"
)

type TestData struct {
	resumeHandler *ResumeHandler
	router        *gin.Engine
	mockUseCase   *mocksResume.UseCase
	mockUSEduc    *mocksEduc.UseCase
	mockUSExp     *mocksExp.UseCase
	mockSB        *mocksCommon.SessionBuilder
	mockSession   *mocksCommon.Session
	httpStatus    []int
	resList       []models.Resume
	resBriefList  []models.BriefResumeInfo
}

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}

func beforeTest() TestData {
	gin.SetMode(gin.TestMode)
	testData := TestData{}
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
	testData.mockSB = new(mocksCommon.SessionBuilder)
	testData.mockSession = new(mocksCommon.Session)
	testData.mockUseCase = new(mocksResume.UseCase)
	testData.mockUSEduc = new(mocksEduc.UseCase)
	testData.mockUSExp = new(mocksExp.UseCase)
	testData.router = gin.Default()
	api := testData.router.Group("/api/v1")
	testData.resumeHandler = NewRest(api.Group("/resume"), testData.mockUseCase,
		//testData.mockUSEduc,
		testData.mockUSExp,
		testData.mockSB,
		func(context *gin.Context) {})
	return testData
}

func getRespStruct(entity interface{}) interface{} {
	switch entity.(type) {
	case resume2.Response:
		resp := entity.(resume2.Response)
		return &resp
	case models.FavoritesForEmpl:
		resp := entity.(models.FavoritesForEmpl)
		return &resp
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
		return models.RespError{Err: err}
	case error:
		err := entity.(error)
		return models.RespError{Err: err.Error()}
	}
	return nil
}

func TestGetResumePageHandler(t *testing.T) {
	var start uint = 0
	var end uint = 2
	td := beforeTest()
	td.router.GET("/page", td.resumeHandler.GetResumePage)

	td.mockUseCase.On("List", start, end).Return(td.resBriefList, nil)
	td.mockUseCase.On("List", start, uint(1000)).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%spage?start=%d&limit=%d", resumeUrlGroup, start, end),
		fmt.Sprintf("%spage", resumeUrlGroup),
		fmt.Sprintf("%spage?start=%d&limit=%d", resumeUrlGroup, start, uint(1000)),
	}

	testExpectedBody := []interface{}{td.resBriefList, common.EmptyFieldErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetResumeSearchHandler(t *testing.T) {
	td := beforeTest()
	td.router.GET("/search", td.resumeHandler.GetResumePage)

	params := resume2.SearchParams{
		KeyWords:        nil,
		SalaryMin:       nil,
		SalaryMax:       nil,
		Gender:          nil,
		EducationLevel:  nil,
		CareerLevel:     nil,
		ExperienceMonth: nil,
		AreaSearch:      nil,
	}
	td.mockUseCase.On("Search", params).Return(td.resBriefList, nil)

	testUrls := []string{
		fmt.Sprintf("%ssearch", resumeUrlGroup),
		fmt.Sprintf("%ssearch", resumeUrlGroup),
	}

	testExpectedBody := []interface{}{td.resBriefList, common.EmptyFieldErr}
	testParamsForPost := []interface{}{params, nil}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, testUrls[i], testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetResumeByID(t *testing.T) {
	td := beforeTest()
	td.router.GET("/by/id/:resume_id", td.resumeHandler.GetResumeByID)

	ID := uuid.New()
	user := models.User{ID: ID}
	cand := models.Candidate{User: user}
	res := models.Resume{
		ResumeID:  ID,
		Candidate: cand,
	}

	favorite := models.FavoritesForEmpl{ FavoriteID: ID }

	td.mockUseCase.On("GetById", ID).Return(&res, nil)
	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetEmplID").Return(ID)
	td.mockUseCase.On("GetFavoriteByResume", ID, ID).Return(&favorite, nil)

	td.mockUseCase.On("GetById", uuid.Nil).Return(nil, assert.AnError)

	res2 := res
	resp := resume2.Response{
		Educations:       res2.Education,
		CustomExperience: res2.ExperienceCustomComp,
	}
	//res2.Candidate = nil
	resp.Resume = res2

	testUrls := []string{
		fmt.Sprintf("%sby/id/%s", resumeUrlGroup, ID.String()),
		fmt.Sprintf("%sby/id/invalidUuid", resumeUrlGroup),
		fmt.Sprintf("%sby/id/%s", resumeUrlGroup, uuid.Nil),
	}

	testExpectedBody := []interface{}{resp, common.EmptyFieldErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetMineResume(t *testing.T) {
	td := beforeTest()
	td.router.GET("/mine", td.resumeHandler.GetMineResume)

	ID := uuid.New()

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetCandID").Return(ID).Once()
	td.mockUseCase.On("GetAllUserResume", ID).Return(td.resBriefList, nil).Once()

	td.mockSession.On("GetCandID").Return(uuid.Nil).Once()

	td.mockSession.On("GetCandID").Return(ID).Once()
	td.mockUseCase.On("GetAllUserResume", ID).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%smine", resumeUrlGroup),
		fmt.Sprintf("%smine", resumeUrlGroup),
		fmt.Sprintf("%smine", resumeUrlGroup),
	}

	testExpectedBody := []interface{}{td.resBriefList, common.AuthRequiredErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCreateResume(t *testing.T) {
	td := beforeTest()
	td.router.POST("/", td.resumeHandler.CreateResume)

	ID := uuid.New()
	user := models.User{ID: ID}
	cand := models.Candidate{User: user}
	res := models.Resume{
		ResumeID:  ID,
		CandID: ID,
		Candidate: cand,
		CandEmail: "a@a.ru",
		CandName: "test",
		CandSurname: "test",
		Gender: "male",
		Skills: "asdasd ad asd ads",
		Title: "!@#DSAD",
		Description: "MYSUPERJOBFPREFSD!",
	}

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetCandID").Return(ID)
	td.mockUseCase.On("Create", mock.AnythingOfType("models.Resume")).Return(&res, nil).Once()
	td.mockUseCase.On("Create", mock.AnythingOfType("models.Resume")).Return(nil, assert.AnError).Once()

	resp := resume2.Response{
		Educations:       res.Education,
		CustomExperience: res.ExperienceCustomComp,
	}
	resp.Resume = res
	testParamsForPost := []interface{}{
		res,
		nil,
		res,
	}
	testExpectedBody := []interface{}{
		resp,
		common.EmptyFieldErr,
		common.DataBaseErr,
	}

	for i := range testParamsForPost {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, resumeUrlGroup, testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}


func TestUpdateResume(t *testing.T) {
	td := beforeTest()
	td.router.PUT("/", td.resumeHandler.UpdateResume)

	ID := uuid.New()
	user := models.User{ID: ID}
	cand := models.Candidate{User: user}
	res := models.Resume{
		ResumeID:  ID,
		CandID: ID,
		Candidate: cand,
		CandEmail: "a@a.ru",
		CandName: "test",
		CandSurname: "test",
		Gender: "male",
		Skills: "asdasd ad asd ads",
		Title: "!@#DSAD",
		Description: "MYSUPERJOBFPREFSD!",
	}

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetCandID").Return(ID)
	td.mockUseCase.On("Update", mock.AnythingOfType("models.Resume")).Return(&res, nil).Once()
	td.mockUseCase.On("Update", mock.AnythingOfType("models.Resume")).Return(nil, assert.AnError)

	res2 := res
	resp := resume2.Response{
		Educations:       res2.Education,
		CustomExperience: res2.ExperienceCustomComp,
	}
	//res2.Candidate = nil
	resp.Resume = res2

	testUrls := []string{
		resumeUrlGroup,
		resumeUrlGroup,
		resumeUrlGroup,
	}
	testParamsForPost := []interface{}{res, nil, res}
	testExpectedBody := []interface{}{resp, common.EmptyFieldErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPut, testUrls[i], testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAddFavorite(t *testing.T) {
	td := beforeTest()
	td.router.POST("/favorite/by/id/:resume_id", td.resumeHandler.AddFavorite)

	ID := uuid.New()
	favoriteForEmpl := models.FavoritesForEmpl{EmplID: ID, ResumeID: ID}
	favoriteId := models.FavoriteID{}

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetEmplID").Return(ID)
	td.mockUseCase.On("AddFavorite", favoriteForEmpl).Return(&favoriteId, nil).Once()

	favoriteForEmplInvalid := models.FavoritesForEmpl{EmplID: ID, ResumeID: uuid.Nil}
	td.mockUseCase.On("AddFavorite", favoriteForEmplInvalid).Return(nil, assert.AnError).Once()


	testUrls := []string{
		fmt.Sprintf("%sfavorite/by/id/%s", resumeUrlGroup, ID),
		fmt.Sprintf("%sfavorite/by/id/invalidID", resumeUrlGroup),
		fmt.Sprintf("%sfavorite/by/id/%s", resumeUrlGroup, uuid.Nil),
	}

	testExpectedBody := []interface{}{favoriteId, common.EmptyFieldErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRemoveFavorite(t *testing.T) {
	td := beforeTest()
	td.router.DELETE("/favorite/by/id/:resume_id", td.resumeHandler.RemoveFavorite)

	ID := uuid.New()
	favoriteForEmpl := models.FavoritesForEmpl{FavoriteID:ID, EmplID: ID}

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetEmplID").Return(ID)
	td.mockUseCase.On("RemoveFavorite", favoriteForEmpl).Return(nil).Once()

	favoriteForEmplInvalid := models.FavoritesForEmpl{FavoriteID:uuid.Nil, EmplID: ID}
	td.mockUseCase.On("RemoveFavorite", favoriteForEmplInvalid).Return(assert.AnError).Once()


	testUrls := []string{
		fmt.Sprintf("%sfavorite/by/id/%s", resumeUrlGroup, ID),
		fmt.Sprintf("%sfavorite/by/id/invalidID", resumeUrlGroup),
		fmt.Sprintf("%sfavorite/by/id/%s", resumeUrlGroup, uuid.Nil),
	}

	testExpectedBody := []interface{}{nil, common.EmptyFieldErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodDelete, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestGetAllFavoritesResume(t *testing.T) {
	td := beforeTest()
	td.router.GET("/myfavorites", td.resumeHandler.GetAllFavoritesResume)

	ID := uuid.New()

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetEmplID").Return(ID).Once()
	td.mockUseCase.On("GetAllEmplFavoriteResume", ID).Return(td.resBriefList, nil).Once()

	td.mockSession.On("GetEmplID").Return(uuid.Nil).Once()

	td.mockSession.On("GetEmplID").Return(ID).Once()
	td.mockUseCase.On("GetAllEmplFavoriteResume", ID).Return(nil, assert.AnError).Once()

	testUrls := []string{
		fmt.Sprintf("%smyfavorites", resumeUrlGroup),
		fmt.Sprintf("%smyfavorites", resumeUrlGroup),
		fmt.Sprintf("%smyfavorites", resumeUrlGroup),
	}

	testExpectedBody := []interface{}{td.resBriefList, common.AuthRequiredErr, common.DataBaseErr}
	for i := range testUrls {
		t.Run("test responses on different urls for getVacancyList handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodGet, testUrls[i], nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, td.httpStatus[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

