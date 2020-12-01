package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
)

type VacancyHandler struct {
	VacUseCase     vacancy.IUseCaseVacancy
	SessionBuilder common.SessionBuilder
}

const (
	vacCreate = 0
	vacUpdate = 1
	vacPath   = "vacancy/"
)

func NewRest(router *gin.RouterGroup,
	useCase vacancy.IUseCaseVacancy,
	sessionBuilder common.SessionBuilder,
	AuthRequired gin.HandlerFunc) *VacancyHandler {
	rest := &VacancyHandler{
		VacUseCase:     useCase,
		SessionBuilder: sessionBuilder,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (v *VacancyHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:vacancy_id", v.GetVacancyByIdHandler)
	router.GET("/comp", v.GetCompVacancyListHandler)
	router.GET("/page", v.GetVacancyListHandler)
	router.POST("/search", v.SearchVacanciesHandler)
	router.Use(AuthRequired)
	{
		router.GET("/mine", v.GetUserVacancyListHandler)
		router.PUT("/", v.UpdateVacancyHandler)
		router.POST("/", v.CreateVacancyHandler)
		router.GET("/recommendation", v.GetRecommendationUserVacancy)
	}
}

func (v *VacancyHandler) GetVacancyByIdHandler(ctx *gin.Context) {
	var req struct {
		VacID string `uri:"vacancy_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	vacId, _ := uuid.Parse(req.VacID)
	vac, err := v.VacUseCase.GetVacancy(vacId)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	session := v.SessionBuilder.Build(ctx)
	candID := session.GetCandID()
	userID := session.GetUserID()

	if candID != uuid.Nil && vac.Sphere != -1 {
		err := v.VacUseCase.AddRecommendation(userID, vac.Sphere)
		if err != nil {
			ctx.Error(err)
		}
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy.Resp{Vacancy: vac}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (v *VacancyHandler) CreateVacancyHandler(ctx *gin.Context) {
	vacHandlerCommon(v, ctx, vacCreate)
}

func (v *VacancyHandler) UpdateVacancyHandler(ctx *gin.Context) {
	vacHandlerCommon(v, ctx, vacUpdate)
}

func (v *VacancyHandler) GetVacancyListHandler(ctx *gin.Context) {
	vacListHandlerCommon(v, ctx, vacancy.ByVacId)
}

func (v *VacancyHandler) GetUserVacancyListHandler(ctx *gin.Context) {
	vacListHandlerCommon(v, ctx, vacancy.ByEmpId)
}

func (v *VacancyHandler) GetCompVacancyListHandler(ctx *gin.Context) {
	vacListHandlerCommon(v, ctx, vacancy.ByCompId)
}

func (v *VacancyHandler) GetRecommendationUserVacancy(ctx *gin.Context) {
	var (
		req     vacancy.VacListRequest
		err     error
		vacList []models.Vacancy
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := v.SessionBuilder.Build(ctx)
	candID := session.GetCandID()
	userID := session.GetUserID()

	if err == nil && candID != uuid.Nil {
		vacList, err = v.VacUseCase.GetRecommendation(userID, int(req.Start), int(req.Limit))
		if err != nil {
			if err.Error() == common.NoRecommendation {
				common.WriteErrResponse(ctx, http.StatusOK, common.NoRecommendation)
				return
			}
			common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
			//ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy.RespList{Vacancies: vacList}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (v *VacancyHandler) SearchVacanciesHandler(ctx *gin.Context) {
	var searchParams models.VacancySearchParams

	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body,  &searchParams); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}

	if err := common.ReqValidation(&searchParams); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	vacList, err := v.VacUseCase.SearchVacancies(searchParams)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy.RespList{Vacancies: vacList}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func vacHandlerCommon(v *VacancyHandler, ctx *gin.Context, treatmentType int) {
	var (
		req vacancy.VacRequest
		err error
	)

	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body,  &req); err != nil {
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
	if req.Sphere == nil {
		noSphere := -1
		req.Sphere = &noSphere
	}
	vacNew := &models.Vacancy{Title: req.Title, SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax, AreaSearch: req.AreaSearch,
		Description: req.Description, Requirements: req.Requirements, Duties: req.Duties, Skills: req.Skills, Sphere: *req.Sphere,
		Employment: req.Employment, ExperienceMonth: req.ExperienceMonth, Location: req.Location, CareerLevel: req.CareerLevel,
		EducationLevel: req.EducationLevel, EmpPhone: req.EmpPhone, EmpEmail: req.EmpEmail, Gender: req.Gender}
	if treatmentType == vacCreate {
		session := v.SessionBuilder.Build(ctx)
		emplID := session.GetEmplID()
		if session == nil {
			ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.SessionErr})
			return
		}
		vacNew.EmpID = emplID
		vacNew, err = v.VacUseCase.CreateVacancy(*vacNew)
	} else {
		vacNew.ID, _ = uuid.Parse(req.Id)
		vacNew, err = v.VacUseCase.UpdateVacancy(*vacNew)
	}
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.DataBaseErr)
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, vacPath+vacNew.ID.String()); err != nil {
			common.WriteErrResponse(ctx, http.StatusInternalServerError, err.Error())
			//ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy.Resp{Vacancy: vacNew}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func vacListHandlerCommon(v *VacancyHandler, ctx *gin.Context, entityType int) {
	var (
		req     vacancy.VacListRequest
		err     error
		vacList []models.Vacancy
		id      = uuid.Nil
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	if entityType == vacancy.ByEmpId {
		session := v.SessionBuilder.Build(ctx)
		emplID := session.GetEmplID()
		if id, err = uuid.Parse(emplID.String()); err != nil {
			common.WriteErrResponse(ctx, http.StatusBadRequest, common.SessionErr)
			return
		}
	} else if entityType == vacancy.ByCompId {
		id, _ = uuid.Parse(req.CompId)
	}
	vacList, err = v.VacUseCase.GetVacancyList(req.Start, req.Limit, id, entityType)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy.RespList{Vacancies: vacList}, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
