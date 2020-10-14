package http

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"net/http"
)

type ResumeHandler struct {
	Usecase resume.IUseCaseResume
}

func NewRest(router *gin.RouterGroup, usecase resume.IUseCaseResume, authMiddleware *jwt.GinJWTMiddleware) *ResumeHandler {
	rest := &ResumeHandler{Usecase: usecase}
	rest.routes(router, authMiddleware)
	return rest
}

func (r *ResumeHandler) routes(router *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	router.GET("/by/id/:resume_id", r.handlerGetResumeByID)
	router.GET("/page", r.handlerGetResumeList)

	router.Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/mine", r.handlerGetAllCurrentUserResume)
		router.POST("/add", r.handlerCreateResume)
	}
	//router.PUT("/resume/update", r.handlerUpdateResume)
}

func (r *ResumeHandler) handlerGetAllCurrentUserResume(c *gin.Context) {
	identityKey := "myid"
	jwtuser, _ := c.Get(identityKey)
	userID := jwtuser.(*models.JWTUserData).ID

	pResume, err := r.Usecase.GetAllUserResume(userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume []models.Resume `json:"resume"`
	}

	c.JSON(http.StatusOK, Resp{Resume: pResume})
}

func (r *ResumeHandler) handlerCreateResume(c *gin.Context) {
	// move to constants
	identityKey := "myid"
	jwtuser, _ := c.Get(identityKey)
	userID := jwtuser.(*models.JWTUserData).ID

	var reqResume struct {
		SalaryMin       *int       `json:"salary_min"`
		SalaryMax       *int       `json:"salary_max"`
		Description     *string    `json:"description"`
		Gender          *string    `json:"gender"`
		Level           *string    `json:"level"`
		ExperienceMonth *int       `json:"experience_month"`
		Education       *string    `json:"education"`
	}
	if err := c.ShouldBindJSON(&reqResume); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resume := models.Resume{
		UserID: userID,
		SalaryMin: reqResume.SalaryMin,
		SalaryMax: reqResume.SalaryMax,
		Description: reqResume.Description,
		Gender: reqResume.Gender,
		Level: reqResume.Level,
		ExperienceMonth: reqResume.ExperienceMonth,
		Education: reqResume.Education,
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
		Start uint `form:"start"`
		Limit uint `form:"limit" binding:"required"`
	}

	if err := c.ShouldBindQuery(&reqResume); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	rList, err := r.Usecase.GetResumePage(reqResume.Start, reqResume.Limit)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume []models.Resume `json:"resume"`
	}

	c.JSON(http.StatusOK, Resp{Resume: rList})
}

//func (r *ResumeHandler) handlerUpdateResume(c *gin.Context) {
//	var reqResume struct {
//		ResumeID 	string `json:"resume_id" binding:"required,uuid"`
//		SalaryMin       int       `json:"salary_min"`
//		SalaryMax       int       `json:"salary_max"`
//		Description     string    `json:"description"`
//		Gender          string    `json:"gender"`
//		Level           string    `json:"level"`
//		ExperienceMonth int       `json:"experience_month"`
//		Education       string    `json:"education"`
//	}
//
//	if err := c.ShouldBindJSON(&reqResume); err != nil {
//		c.AbortWithError(http.StatusBadRequest, err)
//		return
//	}
//
//	resumeID, err := uuid.Parse(reqResume.ResumeID)
//	if err != nil {
//		c.AbortWithError(http.StatusBadRequest, err)
//		return
//	}
//
//	resume := models.Resume{
//		ID: resumeID,
//		SalaryMin: reqResume.SalaryMin,
//		SalaryMax: reqResume.SalaryMax,
//		Description: reqResume.Description,
//		Gender: reqResume.Gender,
//		Level: reqResume.Level,
//		ExperienceMonth: reqResume.ExperienceMonth,
//		Education: reqResume.Education,
//	}
//
//	pResume, err := r.Usecase.UpdateResume(resume)
//	if err != nil {
//		c.AbortWithError(http.StatusInternalServerError, err)
//		return
//	}
//
//	type Resp struct {
//		Resume models.Resume `json:"resume"`
//	}
//
//	c.JSON(http.StatusOK, Resp{Resume: *pResume})
//
//}