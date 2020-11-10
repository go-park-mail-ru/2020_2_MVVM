package http

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/google/uuid"
	"net/http"
)

type ResponseHandler struct {
	UsecaseResponse response.IUseCaseResponse
}

func NewRest(router *gin.RouterGroup,
	usecaseResponse response.IUseCaseResponse,
	AuthRequired gin.HandlerFunc) *ResponseHandler {
	rest := &ResponseHandler{
		UsecaseResponse: usecaseResponse,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (r *ResponseHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.Use(AuthRequired)
	{
		router.POST("/", r.CreateResponse)
		router.PUT("/", r.UpdateStatus)
		router.GET("/my", r.handlerGetAllResponses)
	}
}

func (r *ResponseHandler) CreateResponse(ctx *gin.Context) {
	var response models.Response
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := sessions.Default(ctx)
	var userType string
	candIDStr := session.Get(common.CandID)
	emplIDStr := session.Get(common.EmplID)
	if candIDStr != nil && emplIDStr == nil{
		userType = "candidate"
	} else if candIDStr == nil && emplIDStr != nil{
		userType = "employer"
	} else {
		err := errors.New("this user cannot respond")
		ctx.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}

	response.Initial = userType
	pResponse, err := r.UsecaseResponse.Create(response)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &pResponse)
}

func (r *ResponseHandler) UpdateStatus(ctx *gin.Context) {
	var response models.Response
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := sessions.Default(ctx)
	var userType string
	candIDStr := session.Get(common.CandID)
	emplIDStr := session.Get(common.EmplID)
	if candIDStr != nil && emplIDStr == nil{
		userType = common.Candidate
	} else if candIDStr == nil && emplIDStr != nil{
		userType = common.Employer
	} else {
		err := errors.New("this user cannot respond")
		ctx.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}

	response.Initial = userType

	pResponse, err := r.UsecaseResponse.UpdateStatus(response, userType)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &pResponse)
}

func (r *ResponseHandler) handlerGetAllResponses(ctx *gin.Context) {
	candID, err := common.HandlerGetCurrentUserID(ctx, common.CandID)
	emplID, err := common.HandlerGetCurrentUserID(ctx, common.EmplID)

	var responses []models.ResponseWithTitle
	if candID != uuid.Nil && emplID == uuid.Nil{
		responses, err = r.UsecaseResponse.GetAllCandidateResponses(candID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else if candID == uuid.Nil && emplID != uuid.Nil{
		responses, err = r.UsecaseResponse.GetAllEmployerResponses(emplID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		err := errors.New("this user cannot have responses")
		ctx.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}

	ctx.JSON(http.StatusOK, responses)
}