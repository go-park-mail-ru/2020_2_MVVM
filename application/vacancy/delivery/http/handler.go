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
	Vacancy models.Vacancy `json:"vacancyUser"`
}

func NewRest(router *gin.RouterGroup, useCase vacancy.IUseCaseVacancy) *VacancyHandler {
	rest := &VacancyHandler{VacUseCase: useCase}
	rest.routes(router)
	return rest
}

func (v *VacancyHandler) routes(router *gin.RouterGroup) {
	router.GET("/by/id/:vacancy_id", v.handlerGetVacancyById)
	router.GET("/page", v.handlerGetVacancyList)
	//router.PUT("/", v.handlerUpdateVacancy)
	router.POST("/", v.handlerCreateVacancy)
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

	ctx.JSON(http.StatusOK, Resp{Vacancy: *vac})
}

func (v *VacancyHandler) handlerCreateVacancy(ctx *gin.Context) {
	var req struct {
		Title           string `form:"sum__vacancy-name" binding:"required"`
		SalaryMin       int    `form:"salary_min"`
		SalaryMax       int    `form:"salary_max"`
		Description     string `form:"sum__vacancy-description" binding:"required"`
		Requirements    string `form:"requirements"`
		Duties          string `form:"duties"`
		Skills          string `form:"skills"`
		Spheres         string `form:"spheres"`
		Employment      string `form:"employment"`
		WeekWorkHours   int    `form:"week_work_hours"`
		ExperienceMonth string `form:"experience_work"`
		Location        string `form:"location"`
		CareerLevel     string `form:"career_level"`
		EducationLevel  string `form:"education_level"`
	}
	session := sessions.Default(ctx)
	userIDStr := session.Get("user_id")
	userId, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	errBind := ctx.ShouldBind(&req)
	if errParseForm := ctx.Request.ParseMultipartForm(32 << 15); errParseForm != nil || errBind != nil {
		if errParseForm != nil {
			errBind = errParseForm
		}
		ctx.AbortWithError(http.StatusBadRequest, errBind)
		return
	}
	file, errImg := common.GetImage(ctx, "sum__avatar")
	if errImg != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	vac, err := v.VacUseCase.CreateVacancy(models.Vacancy{Title: req.Title, SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax,
		Description: req.Description, Requirements: req.Requirements, Duties: req.Duties, Skills: req.Skills, Spheres: req.Spheres,
		Employment: req.Employment, WeekWorkHours: req.WeekWorkHours, ExperienceMonth: req.ExperienceMonth, Location: req.Location,
		CareerLevel: req.CareerLevel, EducationLevel: req.EducationLevel}, userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	// TODO: fix error code then vacancy successfully loaded and img valid but couldn't be saved on server storage
	if err := common.AddOrUpdateUserImage(*file, vac.ID.String()); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		//return
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: *vac})
}

func (v *VacancyHandler) handlerGetVacancyList(ctx *gin.Context) {
	var req struct {
		Start uint `form:"start"`
		End   uint `form:"end"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	vacList, err := v.VacUseCase.GetVacancyList(req.Start, req.End)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Vacancy []models.Vacancy `json:"vacancyList"`
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vacList})
}

/*
func (v *VacancyHandler) handlerUpdateVacancy(ctx *gin.Context) {
	var req struct {
		VacancyName        string `json:"vacancy_name" binding:"required"`
		CompanyName        string `json:"company_name" binding:"required"`
		VacancyDescription string `jsnewon:"vacancy_description" binding:"required"`
		WorkExperience     string `json:"work_experience" binding:"required"`
		CompanyAddress     string `json:"company_address" binding:"required"`
		Skills             string `json:"skills" binding:"required"`
		Salary             int    `json:"salary" binding:"required"`
	}
	vac, err := v.VacUseCase.UpdateVacancy(models.Vacancy{FK: userID, VacancyName: req.VacancyName, CompanyName: req.CompanyName,
		VacancyDescription: req.VacancyDescription, WorkExperience: req.WorkExperience, CompanyAddress: req.CompanyAddress,
		Skills: req.Skills, Salary: req.Salary})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Vacancy models.Vacancy `json:"vacancyUser"`
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vac})
}*/
