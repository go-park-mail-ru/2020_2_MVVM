package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/google/uuid"
	"net/http"
)

type CompanyHandler struct {
	CompUseCase official_company.IUseCaseOfficialCompany
}

type Resp struct {
	Company *models.OfficialCompany `json:"company"`
}

type RespList struct {
	Companies []models.OfficialCompany `json:"companies_list"`
}

const (
	compPath = "company/"
)

func NewRest(router *gin.RouterGroup, useCase official_company.IUseCaseOfficialCompany, AuthRequired gin.HandlerFunc) *CompanyHandler {
	rest := &CompanyHandler{CompUseCase: useCase}
	rest.routes(router, AuthRequired)
	return rest
}

func (c *CompanyHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:company_id", c.handlerGetCompany)
	router.GET("/page", c.handlerGetCompanyList)
	router.POST("/search", c.handlerSearchCompanies)
	router.Use(AuthRequired)
	{
		router.GET("/mine", c.handlerGetUserCompany)
		router.POST("/", c.handlerCreateCompany)
		//router.PUT("/", v.handlerUpdateVacancy)
	}
}

func (c *CompanyHandler) handlerGetCompany(ctx *gin.Context) {
	var req struct {
		CompanyID string `uri:"company_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	compUuid, _ := uuid.Parse(req.CompanyID)
	comp, err := c.CompUseCase.GetOfficialCompany(compUuid)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{Company: comp})
}

func (c *CompanyHandler) handlerGetUserCompany(ctx *gin.Context) {
	session := sessions.Default(ctx).Get("empl_id")
	empId, errSession := uuid.Parse(session.(string))
	if errSession != nil {
		ctx.AbortWithError(http.StatusBadRequest, errSession)
		return
	}
	comp, err := c.CompUseCase.GetMineCompany(empId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{Company: comp})
}

func (c *CompanyHandler) handlerCreateCompany(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Spheres     []int  `json:"spheres"`
		Location    string `json:"location" binding:"required"`
		Link        string `json:"link"`
		Avatar      string
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "empty required fields")
		return
	}
	file, errImg := common.GetImageFromBase64(req.Avatar)
	if errImg != nil {
		ctx.JSON(http.StatusBadRequest, errImg)
		return
	}
	session := sessions.Default(ctx).Get("empl_id")
	if session == nil {
		ctx.JSON(http.StatusInternalServerError, "session error")
		return
	}
	empId, errSession := uuid.Parse(session.(string))
	if errSession != nil {
		ctx.JSON(http.StatusInternalServerError, "session error")
		return
	}
	compNew, err := c.CompUseCase.CreateOfficialCompany(models.OfficialCompany{Name: req.Name, Spheres: req.Spheres,
		Location: req.Location, Link: req.Link, Description: req.Description}, empId)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, compPath+compNew.ID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	ctx.JSON(http.StatusOK, Resp{Company: compNew})
}

func (c *CompanyHandler) handlerGetCompanyList(ctx *gin.Context) {
	var req struct {
		Start uint `form:"start"`
		Limit uint `form:"limit" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	compList, err := c.CompUseCase.GetCompaniesList(req.Start, req.Limit)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, RespList{Companies: compList})
}

func (c *CompanyHandler) handlerSearchCompanies(ctx *gin.Context) {
	var searchParams models.CompanySearchParams

	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	compList, err := c.CompUseCase.SearchCompanies(searchParams)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, RespList{Companies: compList})
}
