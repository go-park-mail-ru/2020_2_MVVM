package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/testing/general"
	mocks2 "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/common"
	mocks "github.com/go-park-mail-ru/2020_2_MVVM.git/testing/mocks/application/user"
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

var testData struct {
	userHandler *UserHandler
	router      *gin.Engine
	mockUseCase *mocks.UseCase
	mockAuth    *mocks2.AuthTest
	httpStatus  []int
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
	testData.mockUseCase = new(mocks.UseCase)
	//testData.mockAuth = new(mocks2.AuthTest)
	testData.router = gin.Default()
	api := testData.router.Group("api/v1")
	//testData.mockAuth.On("AuthRequired").Return(nil)
	testData.userHandler = NewRest(api.Group("/users"), testData.mockUseCase, nil)
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

func TestGetUserByIdHandler(t *testing.T) {
	u, r, mockUseCase := testData.userHandler, testData.router, testData.mockUseCase
	r.GET("/by/id/:user_id", u.GetUserByIdHandler)
	candID := uuid.New()
	user := models.User{ID: candID}
	mockUseCase.On("GetUserByID", candID.String()).Return(&user, nil)
	mockUseCase.On("GetUserByID", uuid.Nil.String()).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%sby/id/%s", userUrlGroup, candID),
		fmt.Sprintf("%sby/id/invalidUuid", userUrlGroup),
		fmt.Sprintf("%sby/id/%s", userUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{user, common.EmptyFieldErr, common.DataBaseErr}

	for i := range testUrls {
		t.Run("test responses on different urls for getUser handler", func(t *testing.T) {
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

func TestGetCandByIdHandler(t *testing.T) {
	u, r, mockUseCase := testData.userHandler, testData.router, testData.mockUseCase
	r.GET("cand/by/id/:cand_id", u.GetCandByIdHandler)
	candID := uuid.New()
	user := models.User{}
	mockUseCase.On("GetCandByID", candID.String()).Return(&user, nil)
	mockUseCase.On("GetCandByID", uuid.Nil.String()).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%scand/by/id/%s", userUrlGroup, candID),
		fmt.Sprintf("%scand/by/id/invalidUuid", userUrlGroup),
		fmt.Sprintf("%scand/by/id/%s", userUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{user, common.EmptyFieldErr, common.DataBaseErr}

	for i := range testUrls {
		t.Run("test responses on different urls for GetCandByID handler", func(t *testing.T) {
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

func TestGetEmplByIdHandler(t *testing.T) {
	u, r, mockUseCase := testData.userHandler, testData.router, testData.mockUseCase
	r.GET("empl/by/id/:empl_id", u.GetEmplByIdHandler)
	empID := uuid.New()
	user := models.User{}
	mockUseCase.On("GetEmplByID", empID.String()).Return(&user, nil)
	mockUseCase.On("GetEmplByID", uuid.Nil.String()).Return(nil, assert.AnError)

	testUrls := []string{
		fmt.Sprintf("%sempl/by/id/%s", userUrlGroup, empID),
		fmt.Sprintf("%sempl/by/id/invalidUuid", userUrlGroup),
		fmt.Sprintf("%sempl/by/id/%s", userUrlGroup, uuid.Nil),
	}
	testExpectedBody := []interface{}{user, common.EmptyFieldErr, common.DataBaseErr}

	for i := range testUrls {
		t.Run("test responses on different urls for GetEmplByID handler", func(t *testing.T) {
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

type userReq struct {
	UserType      string `json:"user_type"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	SocialNetwork string `json:"social_network"`
}

func TestCreateUserHandler(t *testing.T) {
	u, r, mockUseCase := testData.userHandler, testData.router, testData.mockUseCase
	r.POST("/", u.CreateUserHandler)
	req := userReq{UserType: "employer", Password: "password", Name: "name", Surname: "surname", Email: "email", Phone: "", SocialNetwork: ""}

	userEmpty := models.User{}
	mockUseCase.On("CreateUser", mock.Anything).Return(&userEmpty, nil)
	//mockUseCase.On("CreateUser", mock.Call{}).Return(nil, assert.AnError)

	testUrls := []string{
		userUrlGroup,
		userUrlGroup,
		//userUrlGroup,
	}
	testExpectedBody := []interface{}{userEmpty, common.EmptyFieldErr, common.DataBaseErr}
	testParamsForPost := []interface{}{req, nil, req}
	for i := range testUrls {
		t.Run("test responses on different urls for CreateUser handler", func(t *testing.T) {
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
