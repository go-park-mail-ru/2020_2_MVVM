package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"net/http"
)

type ResumeHandler struct {
	Usecase resume.IUseCaseResume
}

func NewRest(router *gin.RouterGroup, usecase resume.IUseCaseResume) *ResumeHandler {
	rest := &ResumeHandler{Usecase: usecase}
	rest.routes(router)
	return rest
}

func (r *ResumeHandler) routes(router *gin.RouterGroup) {
	router.GET("/resume/:resume_id", r.handlerGetResumeByID)
	router.POST("/resume/add", r.handlerCreateResume)
}

//func (r *ResumeHandler) routes(router *gin.RouterGroup) {
//	router.GET("/resume/:resume_id", r.handlerGetResumeByID)
//	router.POST("/resume/add", r.handlerCreateResume)
//}
//
//func (r *ResumeHandler) handlerGetNothing(c *gin.Context) {
//	err := r.Usecase.DoNothing()
//	if err != nil {
//		c.AbortWithError(http.StatusInternalServerError, err)
//		return
//	}
//	type Resp struct {
//		Status string
//	}
//
//	c.JSON(http.StatusOK, Resp{Status: "ok"})
//}

func (r *ResumeHandler) handlerGetResumeByID(c *gin.Context) {
	var req struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resumeID, err := uuid.Parse(req.ResumeID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resume, err := r.Usecase.GetResume(resumeID.String())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume models.Resume `json:"resume"`
	}

	c.JSON(http.StatusOK, Resp{Resume: resume})
}

func (r *ResumeHandler) handlerCreateResume(c *gin.Context) {
	var reqPayload struct {
		UserID uuid.UUID `json:"userID" binding:"required"`
		Title    string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resume := models.Resume{
		FK: reqPayload.UserID,
		Title: reqPayload.Title,
	}

	resume, err := r.Usecase.CreateResume(resume)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume models.Resume `json:"resume"`
	}

	c.JSON(http.StatusOK, Resp{Resume: resume})
}
