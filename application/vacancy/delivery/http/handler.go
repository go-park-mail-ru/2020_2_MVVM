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
	Vacancy *models.Vacancy `json:"vacancy"`
}

type RespErr struct {
	Error string `json:"error"`
}

type RespList struct {
	Vacancies []models.Vacancy `json:"vacancyList"`
}

type vacRequest struct {
	Id              string `json:"vac_id,uuid"`
	Avatar          string `json:"avatar"`
	Title           string `json:"title" binding:"required"`
	SalaryMin       int    `json:"salary_min"`
	SalaryMax       int    `json:"salary_max"`
	Description     string `json:"description" binding:"required"`
	Requirements    string `json:"requirements"`
	Duties          string `json:"duties"`
	Skills          string `json:"skills"`
	Sphere          int    `json:"sphere"`
	Employment      string `json:"employment"`
	ExperienceMonth string `json:"experience_work"`
	Location        string `json:"location"`
	CareerLevel     string `json:"career_level"`
	EducationLevel  string `json:"education_level"`
	EmpEmail        string `json:"email"`
	EmpPhone        string `json:"phone"`
}

type vacListRequest struct {
	Start  uint   `form:"start"`
	Limit  uint   `form:"limit" binding:"required"`
	CompId string `form:"comp_id,uuid"`
}

const (
	vacCreate = 0
	vacUpdate = 1
	vacPath   = "vacancy/"
	emptyFieldErr = "empty required fields"
	sessionErr = "session error"
)

func NewRest(router *gin.RouterGroup, useCase vacancy.IUseCaseVacancy, AuthRequired gin.HandlerFunc) *VacancyHandler {
	rest := &VacancyHandler{VacUseCase: useCase}
	rest.routes(router, AuthRequired)
	return rest
}

func (v *VacancyHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:vacancy_id", v.handlerGetVacancyById)
	router.GET("/page/comp", v.handlerGetCompVacancyList)
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
		VacID string `uri:"vacancy_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, RespErr{err.Error()})
		return
	}
	vacId, _ := uuid.Parse(req.VacID)
	vac, err := v.VacUseCase.GetVacancy(vacId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, RespErr{err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vac})
}

func (v *VacancyHandler) handlerCreateVacancy(ctx *gin.Context) {
	vacHandlerCommon(v, ctx, vacCreate)
}

func (v *VacancyHandler) handlerUpdateVacancy(ctx *gin.Context) {
	vacHandlerCommon(v, ctx, vacUpdate)
}

func (v *VacancyHandler) handlerGetVacancyList(ctx *gin.Context) {
	vacListHandlerCommon(v, ctx, vacancy.ByVacId)
}

func (v *VacancyHandler) handlerGetUserVacancyList(ctx *gin.Context) {
	vacListHandlerCommon(v, ctx, vacancy.ByEmpId)
}

func (v *VacancyHandler) handlerGetCompVacancyList(ctx *gin.Context) {
	vacListHandlerCommon(v, ctx, vacancy.ByCompId)
}

func (v *VacancyHandler) handlerSearchVacancies(ctx *gin.Context) {
	var searchParams models.VacancySearchParams

	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.JSON(http.StatusInternalServerError, RespErr{emptyFieldErr})
		return
	}
	VacList, err := v.VacUseCase.SearchVacancies(searchParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, RespErr{err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, RespList{Vacancies: VacList})
}

func vacHandlerCommon(v *VacancyHandler, ctx *gin.Context, treatmentType int) {
	var (
		req vacRequest
		err error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, RespErr{emptyFieldErr})
		return
	}
	file, errImg := common.GetImageFromBase64(req.Avatar)
	if errImg != nil {
		ctx.JSON(http.StatusBadRequest, RespErr{errImg.Error()})
		return
	}
	vacNew := &models.Vacancy{Title: req.Title, SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax,
		Description: req.Description, Requirements: req.Requirements, Duties: req.Duties, Skills: req.Skills,
		Sphere: req.Sphere, Employment: req.Employment, ExperienceMonth: req.ExperienceMonth, Location: req.Location,
		CareerLevel: req.CareerLevel, EducationLevel: req.EducationLevel, EmpPhone: req.EmpPhone, EmpEmail: req.EmpEmail}
	if treatmentType == vacCreate {
		session := sessions.Default(ctx).Get("empl_id")
		if session == nil {
			ctx.JSON(http.StatusInternalServerError, RespErr{sessionErr})
			return
		}
		empId, errSession := uuid.Parse(session.(string))
		if errSession != nil {
			ctx.JSON(http.StatusInternalServerError, RespErr{sessionErr})
			return
		}
		vacNew.EmpID = empId
		vacNew, err = v.VacUseCase.CreateVacancy(*vacNew)
	} else {
		vacNew.ID, _ = uuid.Parse(req.Id)
		vacNew, err = v.VacUseCase.UpdateVacancy(*vacNew)
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, RespErr{err.Error()})
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, vacPath+vacNew.ID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	ctx.JSON(http.StatusOK, Resp{Vacancy: vacNew})
}

func vacListHandlerCommon(v *VacancyHandler, ctx *gin.Context, entityType int) {
	var (
		req     vacListRequest
		err     error
		vacList []models.Vacancy
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, RespErr{emptyFieldErr})
		return
	}
	if entityType == vacancy.ByEmpId {
		session := sessions.Default(ctx).Get("empl_id")
		empId, errSession := uuid.Parse(session.(string))
		if errSession != nil {
			ctx.JSON(http.StatusBadRequest, RespErr{errSession.Error()})
			return
		}
		vacList, err = v.VacUseCase.GetVacancyList(req.Start, req.Limit, empId, vacancy.ByEmpId)
	} else if entityType == vacancy.ByCompId {
		compId, _ := uuid.Parse(req.CompId)
		if compId != uuid.Nil {
			vacList, err = v.VacUseCase.GetVacancyList(req.Start, req.Limit, compId, vacancy.ByCompId)
		}
	} else {
		vacList, err = v.VacUseCase.GetVacancyList(req.Start, req.Limit, uuid.Nil, vacancy.ByVacId)
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, RespErr{err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, RespList{Vacancies: vacList})
}
