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
	Companies []models.OfficialCompany `json:"companyList"`
}

const (
	compPath = "company/"
)

func NewRest(router *gin.RouterGroup, useCase official_company.IUseCaseOfficialCompany, AuthRequired gin.HandlerFunc) *CompanyHandler {
	rest := &CompanyHandler{CompUseCase: useCase}
	rest.Routes(router, AuthRequired)
	return rest
}

func (c *CompanyHandler) Routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:company_id", c.GetCompanyHandler)
	router.GET("/page", c.GetCompanyListHandler)
	router.POST("/search", c.SearchCompaniesHandler)
	router.Use(AuthRequired)
	{
		router.GET("/mine", c.GetUserCompanyHandler)
		router.POST("/", c.CreateCompanyHandler)
		//router.PUT("/", c.handlerUpdateCompany)
	}
}

func (c *CompanyHandler) GetCompanyHandler(ctx *gin.Context) {
	var req struct {
		CompanyID string `uri:"company_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	compUuid, _ := uuid.Parse(req.CompanyID)
	comp, err := c.CompUseCase.GetOfficialCompany(compUuid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, Resp{Company: comp})
}

func (c *CompanyHandler) GetUserCompanyHandler(ctx *gin.Context) {
	session := sessions.Default(ctx).Get("empl_id")
	empId, errSession := uuid.Parse(session.(string))
	if errSession != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.SessionErr})
		return
	}
	comp, err := c.CompUseCase.GetMineCompany(empId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, Resp{Company: comp})
}

func (c *CompanyHandler) CreateCompanyHandler(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required" valid:"utfletternum~название компании может содержать только буквы и цифры.,stringlength(4|15)~название компании должно быть от 2 до 15 символов."`
		Description string `json:"description" binding:"required" valid:"-"`
		Spheres     []int  `json:"spheres" valid:"-"`
		AreaSearch  string `json:"area_search" valid:"utfletter~неверный регион,stringlength(4|128)~длина названия региона от 4 до 128 смиволов"`
		Link        string `json:"link" valid:"url~неверный формат ссылки"`
		Avatar      string `json:"avatar" valid:"-"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err := common.ReqValidation(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}

	file, errImg := common.GetImageFromBase64(req.Avatar)
	if errImg != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: errImg.Error()})
		return
	}
	session := sessions.Default(ctx).Get("empl_id")
	if session == nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.SessionErr})
		return
	}
	empId, errSession := uuid.Parse(session.(string))
	if errSession != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.SessionErr})
		return
	}
	compNew, err := c.CompUseCase.CreateOfficialCompany(models.OfficialCompany{Name: req.Name, Spheres: req.Spheres,
		AreaSearch: req.AreaSearch, Link: req.Link, Description: req.Description}, empId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.DataBaseErr})
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, compPath+compNew.ID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	ctx.JSON(http.StatusOK, Resp{Company: compNew})
}

func (c *CompanyHandler) GetCompanyListHandler(ctx *gin.Context) {
	var req struct {
		Start uint `form:"start"`
		Limit uint `form:"limit" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	compList, err := c.CompUseCase.GetCompaniesList(req.Start, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, RespList{Companies: compList})
}

func (c *CompanyHandler) SearchCompaniesHandler(ctx *gin.Context) {
	var searchParams models.CompanySearchParams

	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	compList, err := c.CompUseCase.SearchCompanies(searchParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}
	ctx.JSON(http.StatusOK, RespList{Companies: compList})
}
