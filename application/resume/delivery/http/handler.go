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
	//router.GET("/resume/list", r.handlerGetResumeList)
	router.PUT("/resume/update", r.handlerUpdateResume)

}

func (r *ResumeHandler) handlerCreateResume(c *gin.Context) {
	var reqResume struct {
		UserID uuid.UUID `json:"userID" binding:"required"`
		Title    string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqResume); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resume := models.Resume{
		FK:    reqResume.UserID,
		Title: reqResume.Title,
	}

	pResume, err := r.Usecase.CreateResume(resume)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume models.Resume `json:"resume"`
	}

	c.JSON(http.StatusOK, Resp{Resume: *pResume})
}

func (r *ResumeHandler) handlerGetResumeByID(c *gin.Context) {
	var reqResume struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&reqResume); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resumeID, err := uuid.Parse(reqResume.ResumeID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pResume, err := r.Usecase.GetResume(resumeID.String())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume models.Resume `json:"resume"`
	}

	c.JSON(http.StatusOK, Resp{Resume: *pResume})
}

func (r *ResumeHandler) handlerGetResumeList(c *gin.Context) {
	var reqResume struct {
		Begin uint `json:"begin"`
		End   uint `json:"end"`
	}

	if err := c.ShouldBindJSON(&reqResume); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	rList, err := r.Usecase.GetResumeList(reqResume.Begin, reqResume.End)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume []models.Resume `json:"resume"`
	}

	c.JSON(http.StatusOK, Resp{Resume: rList})
}

func (r *ResumeHandler) handlerUpdateResume(c *gin.Context) {
	var reqResume struct {
		ResumeID 	string `json:"resume_id" binding:"required,uuid"`
		Title       string    `json:"title"`
		Salary      int       `json:"salary"`
		Description string    `json:"description"`
		Skills      string    `json:"skills"`
		Views       int       `json:"views"`
	}

	if err := c.ShouldBindJSON(&reqResume); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resumeID, err := uuid.Parse(reqResume.ResumeID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resume := models.Resume{
		ID: resumeID,
		Title: reqResume.Title,
		Salary: reqResume.Salary,
		Description: reqResume.Description,
		Skills: reqResume.Skills,
		Views: reqResume.Views,
	}

	pResume, err := r.Usecase.UpdateResume(resume)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type Resp struct {
		Resume models.Resume `json:"resume"`
	}

	c.JSON(http.StatusOK, Resp{Resume: *pResume})

}