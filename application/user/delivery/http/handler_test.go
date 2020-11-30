package http

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocksCommon "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mocksAuth "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/microservices/auth/authmicro"
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
	mockSB      *mocksCommon.SessionBuilder
	mockSession *mocksCommon.Session
	mockAuth    *mocksAuth.AuthClient
	sessionInfo common.BasicSession
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
	testData.mockAuth = new(mocksAuth.AuthClient)
	testData.router = gin.Default()
	api := testData.router.Group("api/v1")
	testData.userHandler = NewRest(api.Group("/users"), testData.mockUseCase, testData.mockAuth,
		common.AuthCookieConfig{}, testData.mockSB, func(context *gin.Context) {})
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
	req := userReq{
		UserType: common.Candidate,
		Password: "password",
		Name: "name",
		Surname: "surname",
		Email: "email@email.ru",
		Phone: "",
		SocialNetwork: "",
	}

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

	td.mockAuth.On("Login", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&td.sessionInfo, nil)
	td.mockUseCase.On("CreateUser", mock.Anything).Return(&userCand, nil).Once()
	td.mockUseCase.On("CreateUser", mock.Anything).Return(nil, assert.AnError)

	//reqUser := models.UserLogin{
	//	Email:    "email@email.ru",
	//	Password: "password",
	//}

	userEmpl := userCand
	IDEmpl := uuid.New()
	userEmpl.ID = IDEmpl
	userEmpl.UserType = common.Employer

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

	testExpectedBody := []interface{}{userCand, common.EmptyFieldErr, common.DataBaseErr}
	testParamsForPost := []interface{}{req, nil, req}
	for i := range testExpectedBody {
		t.Run("test responses on different urls for CreateUser handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, userUrlGroup, testParamsForPost[i])
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
	td.mockSession.On("GetUserID").Return(userID).Once()
	td.mockSession.On("GetUserID").Return(uuid.Nil)

	td.mockUseCase.On("GetUserByID", userID.String()).Return(&user, nil).Once()
	td.mockUseCase.On("GetUserByID", uuid.Nil.String()).Return(nil, assert.AnError)

	testStatuses := []int{
		http.StatusOK,
		http.StatusInternalServerError,
	}
	testExpectedBody := []interface{}{user, common.DataBaseErr}

	for i := range testStatuses {
		t.Run("test responses on different urls for get current User handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodGet, fmt.Sprintf("%sme", userUrlGroup), nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testStatuses[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	td := beforeTest()
	sessionID := uuid.New()

	td.router.POST("/logout", td.userHandler.LogoutHandler)
	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetSessionID").Return(sessionID.String())
	td.mockAuth.On("Logout", sessionID.String()).Return(nil)

	testExpectedBody := []interface{}{
		nil,
		//common.SessionErr,
	}
	testStatuses := []int{
		http.StatusOK,
		//http.StatusInternalServerError,
	}

	for i := range testExpectedBody {
		t.Run("test responses on different urls for logoutUser handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, fmt.Sprintf("%slogout", userUrlGroup), nil)
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testStatuses[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUpdateUserHandler(t *testing.T) {
	td := beforeTest()
	td.router.PUT("/", td.userHandler.UpdateUserHandler)
	id := uuid.New()

	td.mockSB.On("Build", mock.AnythingOfType("*gin.Context")).Return(td.mockSession)
	td.mockSession.On("GetUserID").Return(uuid.Nil).Once()
	td.mockSession.On("GetUserID").Return(uuid.Nil).Once()
	td.mockSession.On("GetUserID").Return(id).Once()

	req := userReq{Name: "name", Surname: "surname", Email: "email@mail.ru"}

	userNew := models.User{
		Name:    req.Name,
		Surname: req.Surname,
		Email:   req.Email,
	}
	td.mockUseCase.On("UpdateUser", mock.Anything).Return(nil, assert.AnError).Once()
	td.mockUseCase.On("UpdateUser", mock.Anything).Return(nil, errors.New(common.WrongPasswd)).Once()
	td.mockUseCase.On("UpdateUser", mock.Anything).Return(&userNew, nil).Once()

	testStatuses := []int{
		http.StatusBadRequest,
		http.StatusBadRequest,
		http.StatusInternalServerError,
		http.StatusConflict,
		http.StatusOK,
	}
	testExpectedBody := []interface{}{
		common.EmptyFieldErr,
		errors.New("Неправильные значения полей: имя должно содержать только буквы"),
		common.DataBaseErr,
		common.WrongPasswd,
		userNew,
	}
	testParamsForPost := []interface{}{nil, userReq{Name: ")"}, req, req, req}
	for i := range testStatuses {
		t.Run("test responses on different urls for UpdateUser handler", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPut, userUrlGroup, testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testStatuses[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRegister(t *testing.T) {
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
	td.mockAuth.On("Login", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&td.sessionInfo, nil)

	td.mockUseCase.On("GetCandidateByID", ID).Return(&cand, nil).Once()

	empl := models.Employer{
		ID:     ID,
		UserID: ID,
	}
	td.mockUseCase.On("GetEmployerByID", ID.String()).Return(&empl, nil).Once()

	testStatuses := []int{
		http.StatusOK,
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusBadRequest,
	}
	testExpectedBody := []interface{}{
		nil,//userCand,
		nil,//userEmpl,
		common.EmptyFieldErr,
		errors.New("Неправильные значения полей: длина пароля должна быть от 5 до 25 символов."),
	}
	testParamsForPost := []interface{}{
		reqUser,
		reqEmpl,
		nil,
		models.UserLogin{Email: "e@e.ru", Password: "err"},
	}

	for i := range testStatuses {
		t.Run("test responses on different urls for register func", func(t *testing.T) {
			w, err := general.PerformRequest(td.router, http.MethodPost, fmt.Sprintf("%slogin", userUrlGroup), testParamsForPost[i])
			if err != nil {
				t.Fatalf("Couldn't create request: %v\n", err)
			}
			if err := general.ResponseComparator(*w, testStatuses[i], getRespStruct(testExpectedBody[i])); err != nil {
				t.Fatal(err)
			}
		})
	}
}
