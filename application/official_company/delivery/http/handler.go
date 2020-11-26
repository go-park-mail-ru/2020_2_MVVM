package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"net/http"
	"os"
	"path"
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
	compCreate = 0
	compUpdate = 1
	compDelete = 2
	compPath   = "company/"
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
	router.DELETE("/", c.DeleteCompanyHandler)
	router.Use(AuthRequired)
	{
		router.GET("/mine", c.GetUserCompanyHandler)
		router.POST("/", c.CreateCompanyHandler)
		router.PUT("/", c.UpdateCompanyHandler)
		//router.DELETE("/", c.DeleteCompanyHandler)
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
	compHandlerCommon(c, ctx, compCreate)
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

func (c *CompanyHandler) UpdateCompanyHandler(ctx *gin.Context) {
	compHandlerCommon(c, ctx, compUpdate)
}

func (c *CompanyHandler) DeleteCompanyHandler(ctx *gin.Context) {
	/*session := sessions.Default(ctx).Get("empl_id")
	empId, errSession := uuid.Parse(session.(string))
	if errSession != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.SessionErr})
		return
	}*/
	empId, err := uuid.Parse("92f68cc8-45e7-41a6-966e-6599d7142ea8")
	err = c.CompUseCase.DeleteOfficialCompany(uuid.Nil, empId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func compHandlerCommon(c *CompanyHandler, ctx *gin.Context, treatmentType int) {
	var (
		req struct {
			Name        string `json:"name" binding:"required" valid:"stringlength(4|30)~название компании должно быть от 4 до 30 символов."`
			Description string `json:"description" binding:"required" valid:"-"`
			Spheres     []int  `json:"spheres" valid:"-"`
			AreaSearch  string `json:"area_search" valid:"stringlength(4|128)~длина названия региона от 4 до 128 смиволов"`
			Link        string `json:"link" valid:"url~неверный формат ссылки"`
			Avatar      string `json:"avatar" valid:"-"`
		}
		compNew    *models.OfficialCompany
		err        error
		avatarPath string
	)
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
	compId := uuid.New()
	avatarName := compPath + compId.String()
	if file != nil {
		fileDir, _ := os.Getwd()
		avatarPath = path.Join(fileDir, common.ImgDir, avatarName)
	}
	if treatmentType == compCreate {
		compNew, err = c.CompUseCase.CreateOfficialCompany(models.OfficialCompany{ID: compId, Name: req.Name, Spheres: convertSliceToPqArr(req.Spheres),
			AreaSearch: req.AreaSearch, Link: req.Link, Description: req.Description, Avatar: avatarPath}, empId)
	} else {
		compNew, err = c.CompUseCase.UpdateOfficialCompany(models.OfficialCompany{ID: compId, Name: req.Name, Spheres: convertSliceToPqArr(req.Spheres),
			AreaSearch: req.AreaSearch, Link: req.Link, Description: req.Description}, empId)
	}
	if err != nil {
		if errMsg := err.Error(); errMsg == common.EmpHaveComp {
			ctx.JSON(http.StatusConflict, common.RespError{Err: errMsg})
		} else {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		}
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, avatarName); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	ctx.JSON(http.StatusOK, Resp{Company: compNew})
}

func convertSliceToPqArr(slice []int) pq.Int64Array {
	arr := make(pq.Int64Array, len(slice))
	for i, e := range slice {
		arr[i] = int64(e)
	}
	return arr
}
