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
	"time"
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
		router.POST("/", r.handlerCreateResponse)
		router.PUT("/", r.handlerUpdateStatus)
		router.GET("/my", r.handlerGetAllResponses)
	}
}

func (r *ResponseHandler) handlerCreateResponse(ctx *gin.Context) {
	var response models.Response
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := sessions.Default(ctx)
	var userType string
	candIDStr := session.Get("cand_id")
	emplIDStr := session.Get("empl_id")
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
	response.DateCreate = time.Now()
	response.Status = "sent"
	pResponse, err := r.UsecaseResponse.CreateResponse(response)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &pResponse)
}

func (r *ResponseHandler) handlerUpdateStatus(ctx *gin.Context) {
	var response models.Response
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := sessions.Default(ctx)
	var userType string
	candIDStr := session.Get("cand_id")
	emplIDStr := session.Get("empl_id")
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

	pResponse, err := r.UsecaseResponse.UpdateStatus(response)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &pResponse)
}

func (r *ResponseHandler) handlerGetAllResponses(ctx *gin.Context) {
	userID, err := common.HandlerGetCurrentUserID(ctx, "empl_id")
	if err != nil {
		userID, err = common.HandlerGetCurrentUserID(ctx, "cand_id")
		if err != nil {
			ctx.AbortWithError(http.StatusMethodNotAllowed, err)
			return
		}
	}

	pResume, err := r.UsecaseResponse.GetAllUserResponses(userID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, models.Resp{Resume: allResume})
}