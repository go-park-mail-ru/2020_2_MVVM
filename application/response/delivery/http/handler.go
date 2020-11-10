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
		router.GET("/free/resumes/:entity_id", r.handlerGetAllResumeWithoutResponse) // vacancy_id
		router.GET("/free/vacancies/:entity_id", r.handlerGetAllVacancyWithoutResponse) // resume_id
	}
}

func (r *ResponseHandler) CreateResponse(ctx *gin.Context) {
	var response models.Response
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
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
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: err.Error()})
		return
	}

	response.Initial = userType
	pResponse, err := r.UsecaseResponse.Create(response)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &pResponse)
}

func (r *ResponseHandler) UpdateStatus(ctx *gin.Context) {
	var response models.Response
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
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
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: err.Error()})
		return
	}

	response.Initial = userType

	pResponse, err := r.UsecaseResponse.UpdateStatus(response, userType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: err.Error()})
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
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: err.Error()})
			return
		}
	} else if candID == uuid.Nil && emplID != uuid.Nil{
		responses, err = r.UsecaseResponse.GetAllEmployerResponses(emplID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: err.Error()})
			return
		}
	} else {
		err := errors.New("this user cannot have responses")
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

func (r *ResponseHandler) handlerGetAllResumeWithoutResponse(ctx *gin.Context) {
	candID, vacancyID, err := r.handlerGetAllEntityWithoutResponse(ctx, common.CandID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: err.Error()})
		return
	}

	resumes, err := r.UsecaseResponse.GetAllResumeWithoutResponse(candID, vacancyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resumes)
}

func (r *ResponseHandler) handlerGetAllVacancyWithoutResponse(ctx *gin.Context) {
	emplID, resumeID, err := r.handlerGetAllEntityWithoutResponse(ctx, common.EmplID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: err.Error()})
		return
	}
	vacancies, err := r.UsecaseResponse.GetAllVacancyWithoutResponse(emplID, resumeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, vacancies)
}

func (r *ResponseHandler) handlerGetAllEntityWithoutResponse(ctx *gin.Context, userType string) (uuid.UUID, uuid.UUID, error) {
	var req struct {
		EntityID string `uri:"entity_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	userID, err := common.GetCurrentUserId(ctx, userType)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return userID, entityID, nil
}