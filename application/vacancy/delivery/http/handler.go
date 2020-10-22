package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
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

func (V *VacancyHandler) routes(router *gin.RouterGroup) {
	router.GET("/vacancy/id/:vacancy_id", V.handlerGetVacancyById)
	router.GET("/vacancy/page", V.handlerGetVacancyList)
	router.PUT("/vacancy/update/:vacancy_id", V.handlerUpdateVacancy)
	router.POST("/vacancy/add", V.handlerCreateVacancy)
}

func (V *VacancyHandler) handlerGetVacancyById(ctx *gin.Context) {
	var req struct {
		VacID uuid.UUID `json:"vacancy_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	vac, err := V.VacUseCase.GetVacancy(req.VacID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Vacancy models.Vacancy `json:"vacancy"`
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vac})
}

type Image struct {
	Path        string `validate:"required"`
	Filename    string `validate:"required"`
	Ext         string `validate:"required"`
	ContentType string `validate:"required"`
	Bytes       int64  `validate:"required,gt=0"`
}

func (V *VacancyHandler) handlerCreateVacancy(ctx *gin.Context) {
	var req struct {
		VacancyName        string `form:"vacancy_name" binding:"required"`
		CompanyName        string `form:"company_name" binding:"required"`
		//VacancyDescription string `json:"vacancy_description" binding:"required"`
		//WorkExperience     string `json:"work_experience" binding:"required"`
		//CompanyAddress     string `json:"company_address" binding:"required"`
		//Skills             string `json:"skills" binding:"required"`
		//Salary             int    `json:"salary" binding:"required"`
	}

	file, header, err := ctx.Request.FormFile("Avatar")
	if err != nil {
		if err.Error() == "http: no such file" {
			return
		}
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	img := Image{
		Path: filepath.Dir(header.Filename),
		Filename: filepath.Base(header.Filename),
		Ext: filepath.Ext(header.Filename),
		ContentType: header.Header.Get("Content-Type"),
		Bytes: header.Size,
	}
	if ext := filepath.Ext(header.Filename); ext == ".png" || ext == ".jpeg" {
		if header.Size > 1200000 {
			fmt.Println("too big!")
		} else if err := addOrUpdateUserImage("temp.jpeg", file); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	fmt.Println(img)
	//identityKey := "myid"
	//jwtUser, _ := ctx.Get(identityKey)
	//userID := jwtUser.(*models.JWTUserData).ID
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	/*vac, err := V.VacUseCase.CreateVacancy(models.Vacancy{FK: userID, VacancyName: req.VacancyName, CompanyName: req.CompanyName,
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

func addOrUpdateUserImage(imgPath string, data io.Reader) error {
	path := filepath.Join("static", imgPath)

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, data); err != nil {
		return err
	}
	return nil
}


func (V *VacancyHandler) handlerGetVacancyList(ctx *gin.Context) {
	var req struct {
		Start uint `form:"start"`
		End   uint `form:"end"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	vacList, err := V.VacUseCase.GetVacancyList(req.Start, req.End)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Vacancy []models.Vacancy `json:"vacancyList"`
	}

	ctx.JSON(http.StatusOK, Resp{Vacancy: vacList})
}

func (V *VacancyHandler) handlerUpdateVacancy(ctx *gin.Context) {
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
	vac, err := V.VacUseCase.UpdateVacancy(models.Vacancy{FK: userID, VacancyName: req.VacancyName, CompanyName: req.CompanyName,
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
