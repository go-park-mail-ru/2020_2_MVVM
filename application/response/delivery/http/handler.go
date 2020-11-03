package http

import (
	"github.com/gin-gonic/gin"
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
		router.POST("/", r.handlerCreateResponse)
		router.GET("/by/id/:resume_id", r.handlerGetResumeByID)
		router.PUT("/", r.handlerUpdateResume)
		router.GET("/mine", r.handlerGetAllCurrentUserResume)
	}
}

func (r *ResponseHandler) handlerCreateResponse(ctx *gin.Context) {
	var req struct {
		ResumeID  string `uri:"resume_id" binding:"required,uuid"`
		VacancyID string `uri:"vacancy_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	

	user, err := u.UserUseCase.GetUserByID(req.UserID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: *user})
}
