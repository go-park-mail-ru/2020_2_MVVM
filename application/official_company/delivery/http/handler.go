package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/official_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mailru/easyjson"
	"net/http"
	"path"
)

type CompanyHandler struct {
	CompUseCase    official_company.IUseCaseOfficialCompany
	SessionBuilder common.SessionBuilder
}

const (
	compCreate = 0
	compUpdate = 1
	compDelete = 2
	compPath   = "company/"
)

func NewRest(router *gin.RouterGroup, useCase official_company.IUseCaseOfficialCompany,
		sessionBuilder common.SessionBuilder, AuthRequired gin.HandlerFunc) *CompanyHandler {
	rest := &CompanyHandler{CompUseCase: useCase, SessionBuilder: sessionBuilder}
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
		router.PUT("/", c.UpdateCompanyHandler)
		//router.DELETE("/", c.DeleteCompanyHandler)
	}
}

func (c *CompanyHandler) GetCompanyHandler(ctx *gin.Context) {
	var req struct {
		CompanyID string `uri:"company_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	compUuid, _ := uuid.Parse(req.CompanyID)
	comp, err := c.CompUseCase.GetOfficialCompany(compUuid)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.Resp{Company: comp}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (c *CompanyHandler) GetUserCompanyHandler(ctx *gin.Context) {
	session := c.SessionBuilder.Build(ctx)
	empId := session.GetEmplID()
	if empId == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.SessionErr)
		return
	}
	comp, err := c.CompUseCase.GetMineCompany(empId)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.Resp{Company: comp}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
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
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	compList, err := c.CompUseCase.GetCompaniesList(req.Start, req.Limit)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.RespList{Companies: compList}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (c *CompanyHandler) SearchCompaniesHandler(ctx *gin.Context) {
	var searchParams models.CompanySearchParams

	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body, &searchParams); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	compList, err := c.CompUseCase.SearchCompanies(searchParams)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.RespList{Companies: compList}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
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
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(nil, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func compHandlerCommon(c *CompanyHandler, ctx *gin.Context, treatmentType int) {
	var (
		req        models.ReqComp
		compNew    *models.OfficialCompany
		err        error
		avatarPath string
	)
	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body, &req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return

	}
	if err := common.ReqValidation(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	file, errImg := common.GetImageFromBase64(req.Avatar)
	if errImg != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, errImg.Error())
		return
	}
	session := c.SessionBuilder.Build(ctx)
	empId := session.GetEmplID()
	if session == nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.SessionErr)
		return
	}
	if empId == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.SessionErr)
		return
	}
	compId := uuid.New()
	avatarName := compPath + compId.String()
	if file != nil {
		avatarPath = common.DOMAIN + path.Join(common.ImgDir, avatarName)
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
			common.WriteErrResponse(ctx, http.StatusConflict, errMsg)
		} else {
			common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		}
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, avatarName); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.Resp{Company: compNew}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func convertSliceToPqArr(slice []int) pq.Int64Array {
	arr := make(pq.Int64Array, len(slice))
	for i, e := range slice {
		arr[i] = int64(e)
	}
	return arr
}
