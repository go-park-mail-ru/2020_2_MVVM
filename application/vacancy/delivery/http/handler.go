package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"net/http"
)

type VacancyHandler struct {
	VacUseCase vacancy.IUseCaseVacancy
}

type Resp struct {
	Vacancy *models.Vacancy
}

type RespList struct {
	Vacancies []models.Vacancy `json:"vacancyList"`
}

type vacRequest struct {
	Id              uuid.UUID `form:"sum__vacancy-id"`
	Title           string    `form:"sum__vacancy-name" binding:"required"`
	SalaryMin       int       `form:"salary_min"`
	SalaryMax       int       `form:"salary_max"`
	Description     string    `form:"sum__vacancy-description" binding:"required"`
	Requirements    string    `form:"requirements"`
	Duties          string    `form:"duties"`
	Skills          string    `form:"skills"`
	Spheres         string    `form:"spheres"`
	Employment      string    `form:"employment"`
	WeekWorkHours   int       `form:"week_work_hours"`
	ExperienceMonth string    `form:"experience_work"`
	Location        string    `form:"location"`
	CareerLevel     string    `form:"career_level"`
	EducationLevel  string    `form:"education_level"`
}

const (
	vacCreate = 0
	vacUpdate = 1
)

func NewRest(router *gin.RouterGroup, useCase vacancy.IUseCaseVacancy, AuthRequired gin.HandlerFunc) *VacancyHandler {
	rest := &VacancyHandler{VacUseCase: useCase}
	rest.routes(router, AuthRequired)
	return rest
}

func (v *VacancyHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:vacancy_id", v.handlerGetVacancyById)
	router.GET("/page", v.handlerGetVacancyList)
	router.POST("/search", v.handlerSearchVacancies)
	router.Use(AuthRequired)
	{
		router.GET("/mine", v.handlerGetUserVacancyList)
		router.PUT("/", v.handlerUpdateVacancy)
		router.POST("/", v.handlerCreateVacancy)
	}
}

func (v *VacancyHandler) handlerGetVacancyById(ctx *gin.Context) {
	var req struct {
		VacID uuid.UUID `json:"vacancy_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	vac, err := v.VacUseCase.GetVacancy(req.VacID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vac})
}

func vacHandlerCommon(v *VacancyHandler, ctx *gin.Context, treatmentType int) {
	var req vacRequest

	err := ctx.ShouldBind(&req)
	if errParseForm := ctx.Request.ParseMultipartForm(32 << 15); errParseForm != nil || err != nil {
		if errParseForm != nil {
			err = errParseForm
		}
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	file, errImg := common.GetImage(ctx.Request, "sum__avatar")
	if errImg != nil {
		ctx.AbortWithError(http.StatusBadRequest, errImg)
		return
	}
	vacNew := &models.Vacancy{ID: req.Id, Title: req.Title, SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax,
		Description: req.Description, Requirements: req.Requirements, Duties: req.Duties, Skills: req.Skills, Spheres: req.Spheres,
		Employment: req.Employment, WeekWorkHours: req.WeekWorkHours, ExperienceMonth: req.ExperienceMonth, Location: req.Location,
		CareerLevel: req.CareerLevel, EducationLevel: req.EducationLevel}
	if treatmentType == vacCreate {
		session := sessions.Default(ctx).Get("empl_id")//TODO: session==nil check
		empId, errSession := uuid.Parse(session.(string))
		if errSession != nil {
			ctx.AbortWithError(http.StatusBadRequest, errSession)
			return
		}
		vacNew.EmpID = empId
		vacNew, err = v.VacUseCase.CreateVacancy(*vacNew)
	} else {
		vacNew, err = v.VacUseCase.UpdateVacancy(*vacNew)
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// TODO: fix error code then vacancy successfully loaded and img valid but couldn't be saved on server storage
	if file != nil {
		if err := common.AddOrUpdateUserFile(*file, vacNew.ID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			//return
		}
	}
	ctx.JSON(http.StatusOK, Resp{Vacancy: vacNew})
}

func (v *VacancyHandler) handlerCreateVacancy(ctx *gin.Context) {
	vacHandlerCommon(v, ctx, vacCreate)
}

func (v *VacancyHandler) handlerUpdateVacancy(ctx *gin.Context) {
	vacHandlerCommon(v, ctx, vacUpdate)
}

func (v *VacancyHandler) handlerGetVacancyList(ctx *gin.Context) {
	var req struct {
		Start uint `form:"start"`
		Limit uint `form:"limit"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	vacList, err := v.VacUseCase.GetVacancyList(req.Start, req.Limit, uuid.Nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, RespList{Vacancies: vacList})
}

func (v *VacancyHandler) handlerGetUserVacancyList(ctx *gin.Context) {
	var req struct {
		Start uint `form:"start"`
		End   uint `form:"limit" binding:"required"`
	}

	session := sessions.Default(ctx).Get("empl_id") //TODO: session==nil check
	empId, errSession := uuid.Parse(session.(string))
	if err := ctx.ShouldBindQuery(&req); errSession != nil || err != nil {
		if errSession != nil {
			err = errSession
		}
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userVacList, err := v.VacUseCase.GetVacancyList(req.Start, req.End, empId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, RespList{Vacancies: userVacList})
}

func (v *VacancyHandler) handlerSearchVacancies(ctx *gin.Context) {
	var searchParams models.VacancySearchParams

	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	VacList, err := v.VacUseCase.SearchVacancies(searchParams)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, RespList{Vacancies: VacList})
}
