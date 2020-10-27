package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"net/http"
	"path/filepath"
)

type VacancyHandler struct {
	VacUseCase vacancy.IUseCaseVacancy
}

func NewRest(router *gin.RouterGroup, useCase vacancy.IUseCaseVacancy) *VacancyHandler {
	rest := &VacancyHandler{VacUseCase: useCase}
	rest.routes(router)
	return rest
}

func (v *VacancyHandler) routes(router *gin.RouterGroup) {
	router.GET("/vacancy/id/:vacancy_id", v.handlerGetVacancyById)
	router.GET("/vacancy/page", v.handlerGetVacancyList)
	router.PUT("/vacancy/update/:vacancy_id", v.handlerUpdateVacancy)
	router.POST("/vacancy/add", v.handlerCreateVacancy)
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
	type Resp struct {
		Vacancy models.Vacancy `json:"vacancy"`
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vac})
}

func (v *VacancyHandler) handlerCreateVacancy(ctx *gin.Context) {
	var req struct {
		//VacancyName string `form:"sum__company-vacancy_name" binding:"required"`
		//CompanyName string `form:"sum__company-name" binding:"required"`
		//VacancyDescription string `form:"sum__company-vacancy_description" binding:"required"`
		//WorkExperience     string `json:"work_experience" binding:"required"`
		//CompanyAddress     string `form:"sum__company-address" binding:"required"`
		//Skills             string `json:"skills" binding:"required"`
		//Salary             int    `json:"salary" binding:"required"`
	}
	err := ctx.ShouldBind(&req)
	if errParseForm := ctx.Request.ParseMultipartForm(32 << 15); errParseForm != nil || err != nil {
		if errParseForm != nil {
			err = errParseForm
		}
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	file, header, err := ctx.Request.FormFile("sum__avatar")
	if err == nil {
		if err := common.FileValidation(header, file, []string{"image/jpeg", "image/png"}, common.MaxImgSize); err.Code() == common.FileValid {
			if err := common.AddOrUpdateUserImage(file, fmt.Sprintf("temp%s", filepath.Ext(header.Filename))); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			ctx.JSON(http.StatusOK, err.String())
			return
		}
	} else if err.Error() != "http: no such file" {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	/*vac, err := v.VacUseCase.CreateVacancy(models.Vacancy{FK: userID, VacancyName: req.VacancyName, CompanyName: req.CompanyName,
	VacancyDescription: req.VacancyDescription, WorkExperience: req.WorkExperience, CompanyAddress: req.CompanyAddress,
	Skills: req.Skills, Salary: req.Salary})*/
	/*if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}*/
	type Resp struct {
		Vacancy models.Vacancy `json:"vacancyUser"`
	}

	ctx.JSON(http.StatusOK, Resp{})
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
	identityKey := "myid"
	jwtUser, _ := ctx.Get(identityKey)
	userID := jwtUser.(*models.JWTUserData).ID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
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
}
