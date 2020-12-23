package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/authmicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/vacancyMicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	vacancy2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/vacancy"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
)

type VacancyHandler struct {
	vacancyClient  vacancyMicro.VacClient
	SessionBuilder common.SessionBuilder
	authClient     authmicro.AuthClient
}

const (
	vacCreate = 0
	vacUpdate = 1
	vacPath   = "vacancy/"
)

func NewRest(router *gin.RouterGroup,
	sessionBuilder common.SessionBuilder,
	AuthRequired gin.HandlerFunc,
	vacancyClient vacancyMicro.VacClient,
	authClient authmicro.AuthClient) *VacancyHandler {
	rest := &VacancyHandler{
		SessionBuilder: sessionBuilder,
		vacancyClient:  vacancyClient,
		authClient:     authClient,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (v *VacancyHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:vacancy_id", v.GetVacancyByIdHandler)
	router.GET("/comp", v.GetCompVacancyListHandler)
	router.GET("/top/:top_spheres_cnt", v.GetVacancyTopSpheres)
	router.GET("/top", v.GetVacancyTopSpheresAll)
	router.GET("/page", v.GetVacancyListHandler)
	router.POST("/search", v.SearchVacanciesHandler)
	router.Use(AuthRequired)
	{
		router.GET("/mine", v.GetUserVacancyListHandler)
		router.PUT("/", v.UpdateVacancyHandler)
		router.DELETE("/:vacancy_id", v.DeleteVacancyHandler)
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
	vac, err := v.vacancyClient.GetVacancy(vacId)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	sessionID, _ := ctx.Cookie("session")
	if sessionID != "" {
		if session, _ := v.authClient.Check(sessionID); session != nil {
			if session.GetCandID() != uuid.Nil && vac.Sphere != -1 {
				err := v.vacancyClient.AddRecommendation(session.GetUserID(), vac.Sphere)
				if err != nil {
					_ = ctx.Error(err)
				}
			}
		}
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy2.Resp{Vacancy: vac}, ctx.Writer); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
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
		req     vacancy2.VacListRequest
		err     error
		vacList []models.Vacancy
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	session := v.SessionBuilder.Build(ctx)
	if session.GetCandID() != uuid.Nil {
		vacList, err = v.vacancyClient.GetRecommendation(session.GetUserID(), int(req.Start), int(req.Limit))
		if err != nil {
			if err.Error() == common.NoRecommendation {
				common.WriteErrResponse(ctx, http.StatusOK, common.NoRecommendation)
				return
			}
			common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
			return
		}
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy2.RespList{Vacancies: vacList}, ctx.Writer); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (v *VacancyHandler) SearchVacanciesHandler(ctx *gin.Context) {
	var searchParams models.VacancySearchParams

	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body, &searchParams); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	if err := common.ReqValidation(&searchParams); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	vacList, err := v.vacancyClient.SearchVacancies(searchParams)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy2.RespList{Vacancies: vacList}, ctx.Writer); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (v *VacancyHandler) GetVacancyTopSpheresAll(ctx *gin.Context) {
	topSphereHandlerCommon(v, ctx, vacancy.TopAll)
}

func (v *VacancyHandler) GetVacancyTopSpheres(ctx *gin.Context) {
	var (
		req    vacancy2.TopSpheres
		topCnt int32 = vacancy.TopDefaultCnt
	)

	if err := ctx.ShouldBindUri(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	if *req.TopSpheresCnt > 0 {
		topCnt = *req.TopSpheresCnt
	}
	topSphereHandlerCommon(v, ctx, topCnt)
}

func (v *VacancyHandler) DeleteVacancyHandler(ctx *gin.Context) {
	var req struct {
		VacID string `uri:"vacancy_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	vacId, _ := uuid.Parse(req.VacID)
	session := v.SessionBuilder.Build(ctx)
	empId := session.GetEmplID()
	if empId == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.SessionErr)
		return
	}
	if err := v.vacancyClient.DeleteVacancy(vacId, empId); err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(nil, ctx.Writer); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func topSphereHandlerCommon(v *VacancyHandler, ctx *gin.Context, topCnt int32) {
	spheresInfo, vacInfo, err := v.vacancyClient.GetVacancyTopSpheres(topCnt)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy2.RespTop{TopSpheres: spheresInfo, NewVacCnt: vacInfo.NewVacCnt, AllVacCnt: vacInfo.AllVacCnt}, ctx.Writer); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func vacHandlerCommon(v *VacancyHandler, ctx *gin.Context, treatmentType int) {
	var (
		req vacancy2.VacRequest
		err error
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
	if req.Sphere == nil {
		noSphere := -1
		req.Sphere = &noSphere
	}
	vacNew := &models.Vacancy{Title: req.Title, SalaryMin: req.SalaryMin, SalaryMax: req.SalaryMax, AreaSearch: req.AreaSearch,
		Description: req.Description, Requirements: req.Requirements, Duties: req.Duties, Skills: req.Skills, Sphere: *req.Sphere,
		Employment: req.Employment, ExperienceMonth: req.ExperienceMonth, Location: req.Location, CareerLevel: req.CareerLevel,
		EducationLevel: req.EducationLevel, EmpPhone: req.EmpPhone, EmpEmail: req.EmpEmail, Gender: req.Gender}
	session := v.SessionBuilder.Build(ctx)
	vacNew.EmpID = session.GetEmplID()
	if treatmentType == vacCreate {
		if session == nil {
			common.WriteErrResponse(ctx, http.StatusForbidden, common.SessionErr)
			return
		}
		vacNew, err = v.vacancyClient.CreateVacancy(*vacNew)
	} else if vacNew.ID, _ = uuid.Parse(req.Id); vacNew.ID != uuid.Nil {
		vacNew, err = v.vacancyClient.UpdateVacancy(*vacNew)
	} else {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, vacPath+vacNew.ID.String()); err != nil {
			common.WriteErrResponse(ctx, http.StatusInternalServerError, err.Error())
		}
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy2.Resp{Vacancy: vacNew}, ctx.Writer); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func vacListHandlerCommon(v *VacancyHandler, ctx *gin.Context, entityType int) {
	var (
		req     vacancy2.VacListRequest
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
	vacList, err = v.vacancyClient.GetVacancyList(req.Start, req.Limit, id, entityType)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(vacancy2.RespList{Vacancies: vacList}, ctx.Writer); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
