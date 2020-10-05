package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"net/http"
)

type VacancyHandler struct {
	VacUseCase vacancy.IUseCaseVacancy
}

func NewRest(router *gin.RouterGroup, useCase vacancy.IUseCaseVacancy) *VacancyHandler {
	rest := &VacancyHandler{VacUseCase: useCase}
	rest.routes(router)
	return rest
}

func (V *VacancyHandler) routes(router *gin.RouterGroup) {
	router.GET("/vacancy/:resume_id", V.handlerGetVacancyById)
	router.POST("/vacancy/add", V.handlerCreateVacancy)
}

func (V *VacancyHandler) handlerGetVacancyById(ctx *gin.Context) {
	var req struct {
		VacID string `uri:"user_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	vacID, err := uuid.Parse(req.VacID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	vac, err := V.VacUseCase.GetVacancy(vacID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Vacancy models.Vacancy `json:"vacancy"`
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vac})
}

func (V *VacancyHandler) handlerCreateVacancy(ctx *gin.Context) {
	var req struct {
		AuthorID           string `json:"userID" binding:"required"`
		VacancyName        string `json:"vacancy_name" binding:"required"`
		CompanyName        string `json:"company_name" binding:"required"`
		VacancyDescription string `json:"vacancy_description" binding:"required"`
		WorkExperience     string `json:"work_experience" binding:"required"`
		CompanyAddress     string `json:"company_address" binding:"required"`
		Skills             string `json:"skills" binding:"required"`
		Salary             int    `json:"salary" binding:"required"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	vac, err := V.VacUseCase.CreateVacancy(models.Vacancy{FK: userID, VacancyName: req.VacancyName, CompanyName: req.CompanyName,
		VacancyDescription: req.VacancyDescription, WorkExperience: req.WorkExperience, CompanyAddress: req.CompanyAddress,
		Skills: req.Skills, Salary: req.Salary})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Vacancy models.Vacancy `json:"vacancy"`
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vac})
}
