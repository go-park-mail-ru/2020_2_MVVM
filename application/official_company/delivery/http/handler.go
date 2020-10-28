package http

import (
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
	Company models.OfficialCompany `json:"official_company"`
}

func NewRest(router *gin.RouterGroup, useCase official_company.IUseCaseOfficialCompany) *CompanyHandler {
	rest := &CompanyHandler{CompUseCase: useCase}
	rest.routes(router)
	return rest
}

func (c *CompanyHandler) routes(router *gin.RouterGroup) {
	router.GET("/by/id/:company_id", c.handlerGetCompanyById)
	router.GET("/page", c.handlerGetCompanyList)
	//router.PUT("/", v.handlerUpdateVacancy)
	router.POST("/", c.handlerCreateCompany)
}

func (c *CompanyHandler) handlerGetCompanyById(ctx *gin.Context) {
	var req struct {
		CompanyID uuid.UUID `json:"company_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	comp, err := c.CompUseCase.GetOfficialCompany(req.CompanyID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{Company: *comp})
}

func (c *CompanyHandler) handlerCreateCompany(ctx *gin.Context) {
	var req struct {
		Name     string   `form:"comp__company-name" binding:"required"`
		Sphere   []string `form:"comp__company-sphere"`
		Location string   `form:"comp__company-location"`
		Link     string   `form:"comp__company-link"`
		VacCount int      `form:"comp__company-vac_count"`
	}

	err := ctx.ShouldBind(&req)
	if errParseForm := ctx.Request.ParseMultipartForm(32 << 15); errParseForm != nil || err != nil {
		if errParseForm != nil {
			err = errParseForm
		}
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	file, err := common.GetImage(ctx, "comp__avatar")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	comp, err := c.CompUseCase.CreateOfficialCompany(models.OfficialCompany{Name: req.Name, Sphere: req.Sphere,
		Location: req.Location, Link: req.Link, VacCount: req.VacCount})
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := common.AddOrUpdateUserImage(*file, comp.ID.String()); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		//return
	}
	ctx.JSON(http.StatusOK, Resp{Company: *comp})
}

func (c *CompanyHandler) handlerGetCompanyList(context *gin.Context) {

}
