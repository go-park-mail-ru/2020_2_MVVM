package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocksCommon "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mUser "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"os"
	"testing"
)

const (
	userUrlGroup = "/api/v1/users/"
)

type TestData struct {
	userHandler *UserHandler
	router      *gin.Engine
	mockUseCase *mUser.UseCase
	mockAuth    *mocksCommon.AuthTest
	mockSB      *mocksCommon.SessionBuilder
	mockSession *mocksCommon.Session
	httpStatus  []int
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
	testData.mockUseCase = new(mUser.UseCase)
	testData.mockSB = new(mocksCommon.SessionBuilder)
	testData.mockSession = new(mocksCommon.Session)
	//testData.mockAuth = new(mocksCommon.AuthTest)
	testData.router = gin.Default()
	api := testData.router.Group("api/v1")
	//testData.mockAuth.On("AuthRequired").Return(nil)
	testData.userHandler = NewRest(api.Group("/users"), testData.mockUseCase, testData.mockSB, func(context *gin.Context) {})
	return testData
}

func getRespStruct(entity interface{}) interface{} {
	switch entity.(type) {
	case models.User:
		user := entity.(models.User)
		return Resp{&user}
	case string:
		err := entity.(string)
		return common.RespError{Err: err}
	case error:
		err := entity.(error)
		return common.RespError{Err: err.Error()}
	}
	return nil
}

type userReq struct {
	UserType      string `json:"user_type"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	SocialNetwork string `json:"social_network"`
}

func TestGetUserByIdHandler(t *testing.T) {
	td := beforeTest()
	td.router.GET("/by/id/:user_id", td.userHandler.GetUserByIdHandler)
	candID := uuid.New()
	user := models.User{ID: candID}
	td.mockUseCase.On("GetUserByID", candID.String()).Return(&user, nil)
	td.mockUseCase.On("GetUserByID", uuid.Nil.String()).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%sby/id/%s", userUrlGroup, candID),
		fmt.Sprintf("%sby/id/invalidUuid", userUrlGroup),
		fmt.Sprintf("%sby/id/%s", userUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{user, common.EmptyFieldErr, common.DataBaseErr}

	for i := range testUrls {
		t.Run("test responses on different urls for getUser handler", func(t *testing.T) {
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

func TestGetCandByIdHandler(t *testing.T) {
	td := beforeTest()
	td.router.GET("cand/by/id/:cand_id", td.userHandler.GetCandByIdHandler)
	candID := uuid.New()
	user := models.User{}
	td.mockUseCase.On("GetCandByID", candID.String()).Return(&user, nil)
	td.mockUseCase.On("GetCandByID", uuid.Nil.String()).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%scand/by/id/%s", userUrlGroup, candID),
		fmt.Sprintf("%scand/by/id/invalidUuid", userUrlGroup),
		fmt.Sprintf("%scand/by/id/%s", userUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{user, common.EmptyFieldErr, common.DataBaseErr}

	for i := range testUrls {
		t.Run("test responses on different urls for GetCandByID handler", func(t *testing.T) {
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

func TestGetEmplByIdHandler(t *testing.T) {
	td := beforeTest()
	td.router.GET("empl/by/id/:empl_id", td.userHandler.GetEmplByIdHandler)
	empID := uuid.New()
	user := models.User{}
	td.mockUseCase.On("GetEmplByID", empID.String()).Return(&user, nil)
	td.mockUseCase.On("GetEmplByID", uuid.Nil.String()).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%sempl/by/id/%s", userUrlGroup, empID),
		fmt.Sprintf("%sempl/by/id/invalidUuid", userUrlGroup),
		fmt.Sprintf("%sempl/by/id/%s", userUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{user, common.EmptyFieldErr, common.DataBaseErr}

	for i := range testUrls {
		t.Run("test responses on different urls for GetEmplByID handler", func(t *testing.T) {
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

func TestCreateUserHandler(t *testing.T) {
	td := beforeTest()
	td.router.POST("/", td.userHandler.CreateUserHandler)
	req := userReq{UserType: common.Candidate, Password: "password", Name: "name", Surname: "surname", Email: "email@email.ru", Phone: "", SocialNetwork: ""}

	IDCand := uuid.New()
	userCand := models.User{
		ID:            IDCand,
		UserType:      common.Candidate,
		Name:          "name",
		Surname:       "surname",
		Email:         "email@email.ru",
		PasswordHash:  []byte{1, 2, 3},
		Phone:         nil,
		SocialNetwork: nil,
	}

	td.mockUseCase.On("CreateUser", mock.Anything).Return(&userCand, nil)

	reqUser := models.UserLogin{
		Email:    "email@email.ru",
		Password: "password",
	}

	userEmpl := userCand
	IDEmpl := uuid.New()
	userEmpl.ID = IDEmpl
	userEmpl.UserType = common.Employer

	td.mockUseCase.On("Login", reqUser).Return(&userCand, nil)
	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	cand := models.Candidate{
		ID:     IDCand,
		UserID: IDCand,
	}
	td.mockUseCase.On("GetCandidateByID", IDCand.String()).Return(&cand, nil)
	td.mockSession.On("Set", common.CandID, IDCand.String())
	td.mockSession.On("Set", common.EmplID, nil)
	td.mockSession.On("Set", common.UserID, IDCand.String())
	td.mockSession.On("Save").Return(nil)

	testUrls := []string{
		userUrlGroup,
	}
	testExpectedBody := []interface{}{userCand, common.EmptyFieldErr, common.DataBaseErr}
	testParamsForPost := []interface{}{req, nil, req}
	for i := range testUrls {
		t.Run("test responses on different urls for CreateUser handler", func(t *testing.T) {
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

func TestGetCurrentUserHandler(t *testing.T) {
	td := beforeTest()
	td.router.GET("/me", td.userHandler.GetCurrentUserHandler)

	userID := uuid.New()
	user := models.User{ID: userID}

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("Get", common.UserID).Return(userID.String())
	//td.mockSession.On("Get", common.UserID).Return(uuid.Nil)

	td.mockUseCase.On("GetUserByID", userID.String()).Return(&user, nil)
	//mockUseCase.On("GetUserByID", uuid.Nil.String()).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%sme", userUrlGroup),
		//fmt.Sprintf("%sby/id/invalidUuid", userUrlGroup),
		//fmt.Sprintf("%sby/id/%s", userUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{user}

	for i := range testUrls {
		t.Run("test responses on different urls for get current User handler", func(t *testing.T) {
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

func TestLogoutHandler(t *testing.T) {
	td := beforeTest()
	td.router.POST("/logout", td.userHandler.LogoutHandler)

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("Clear")
	td.mockSession.On("Options", mock.Anything)
	td.mockSession.On("Save").Return(nil)

	testUrls := []string{
		fmt.Sprintf("%slogout", userUrlGroup),
	}
	testExpectedBody := []interface{}{nil}

	for i := range testUrls {
		t.Run("test responses on different urls for logoutUser handler", func(t *testing.T) {
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

func TestUpdateUserHandler(t *testing.T) {
	td := beforeTest()
	td.router.PUT("/", td.userHandler.UpdateUserHandler)

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("Get", common.UserID).Return(mock.Anything)
	//mockSession.On("Get", common.UserID).Return(assert.AnError)

	req := userReq{Name: "name", Surname: "surname", Email: "email", Phone: "", SocialNetwork: ""}

	userEmpty := models.User{
		UserType:      "",
		Name:          "name",
		Surname:       "surname",
		Email:         "email",
		PasswordHash:  nil,
		Phone:         nil,
		SocialNetwork: nil,
	}
	td.mockUseCase.On("UpdateUser", mock.Anything, "", "", "name", "surname", "email", "", "").Return(&userEmpty, nil)
	//mockUseCase.On("UpdateUser", mock.Call{}).Return(nil, assert.AnError)

	testUrls := []string{
		userUrlGroup,
	}
	testExpectedBody := []interface{}{userEmpty}
	testParamsForPost := []interface{}{req}
	for i := range testUrls {
		t.Run("test responses on different urls for UpdateUser handler", func(t *testing.T) {
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

func TestLoginHandler(t *testing.T) {
	td := beforeTest()
	td.router.POST("/login", td.userHandler.LoginHandler)

	reqUser := models.UserLogin{
		Email:    "email@email.ru",
		Password: "password",
	}
	reqEmpl := reqUser
	reqEmpl.Email = "email1edfs@email1.ru"

	ID := uuid.New()
	userCand := models.User{
		ID:            ID,
		UserType:      common.Candidate,
		Name:          "name",
		Surname:       "surname",
		Email:         "email678",
		PasswordHash:  []byte{1, 2, 3, 14},
		Phone:         nil,
		SocialNetwork: nil,
	}
	userEmpl := userCand
	userEmpl.UserType = common.Employer

	cand := models.Candidate{
		ID:     ID,
		UserID: ID,
	}
	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockUseCase.On("Login", reqUser).Return(&userCand, nil)
	td.mockUseCase.On("GetCandidateByID", ID.String()).Return(&cand, nil)
	td.mockSession.On("Set", common.CandID, ID.String())
	td.mockSession.On("Set", common.EmplID, nil)
	td.mockSession.On("Set", common.UserID, ID.String())

	empl := models.Employer{
		ID:     ID,
		UserID: ID,
	}
	td.mockUseCase.On("Login", reqEmpl).Return(&userEmpl, nil)
	td.mockUseCase.On("GetEmployerByID", ID.String()).Return(&empl, nil)
	td.mockSession.On("Set", common.CandID, nil)
	td.mockSession.On("Set", common.EmplID, ID.String())
	td.mockSession.On("Set", common.UserID, ID.String())
	td.mockSession.On("Save").Return(nil)

	testUrls := []string{
		fmt.Sprintf("%slogin", userUrlGroup),
		fmt.Sprintf("%slogin", userUrlGroup),
	}
	testExpectedBody := []interface{}{userCand, userEmpl}
	testParamsForPost := []interface{}{reqUser, reqEmpl}

	for i := range testUrls {
		t.Run("test responses on different urls for loginUser handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, testUrls[i], testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, 200, getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}
