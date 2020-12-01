package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/authmicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/vacancy/vacancyMicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/vacancy"
	"github.com/google/uuid"
	"net/http"
)

type VacancyHandler struct {
	vacancyClient  vacancyMicro.VacClient
	SessionBuilder common.SessionBuilder
	autClient      authmicro.AuthClient
}

type Resp struct {
	Vacancy *models.Vacancy `json:"vacancy"`
}

type RespList struct {
	Vacancies []models.Vacancy `json:"vacancyList"`
}

type vacRequest struct {
	Id              string `json:"vac_id,uuid" valid:"-"`
	Avatar          string `json:"avatar" valid:"-"`
	Title           string `json:"title" binding:"required" valid:"stringlength(4|128)~название вакансии должно быть от 4 до 128 символов в длину."`
	Gender          string `json:"gender" valid:"-"`
	SalaryMin       int    `json:"salary_min" valid:"-"`
	SalaryMax       int    `json:"salary_max" valid:"-"`
	Description     string `json:"description" binding:"required" valid:"-"`
	Requirements    string `json:"requirements" valid:"-"`
	Duties          string `json:"duties" valid:"-"`
	Skills          string `json:"skills" valid:"-"`
	Sphere          *int   `json:"sphere" valid:"numeric~сфера деятельности должна содержать только код"`
	Employment      string `json:"employment" valid:"-"`
	ExperienceMonth int    `json:"experience_month" valid:"-"`
	Location        string `json:"location" valid:"stringlength(4|512)~длина адреса от 4 до 512 смиволов"`
	AreaSearch      string `json:"area_search" valid:"stringlength(4|128)~длина названия региона от 4 до 128 смиволов"`
	CareerLevel     string `json:"career_level" valid:"-"`
	EducationLevel  string `json:"education_level" valid:"-"`
	EmpEmail        string `json:"email" valid:"email"`
	EmpPhone        string `json:"phone" valid:"numeric~номер телефона должен состоять только из цифр.,stringlength(4|18)~номер телефона от 4 до 18 цифр"`
}

type vacListRequest struct {
	Start  uint   `form:"start"`
	Limit  uint   `form:"limit" binding:"required"`
	CompId string `form:"comp_id,uuid"`
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
	autClient authmicro.AuthClient) *VacancyHandler {
	rest := &VacancyHandler{
		SessionBuilder: sessionBuilder,
		vacancyClient:  vacancyClient,
		autClient:      autClient,
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	vacId, _ := uuid.Parse(req.VacID)
	vac, err := v.vacancyClient.GetVacancy(vacId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}
	sessionID, _ := ctx.Cookie("session")
	if sessionID != "" {
		if session, _ := v.autClient.Check(sessionID); session != nil {
			if session.GetCandID() != uuid.Nil && vac.Sphere != -1 {
				err := v.vacancyClient.AddRecommendation(session.GetUserID(), vac.Sphere)
				if err != nil {
					_ = ctx.Error(err)
				}
			}
		}
	}
	ctx.JSON(http.StatusOK, Resp{Vacancy: vac})
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
		req     vacListRequest
		err     error
		vacList []models.Vacancy
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	session := v.SessionBuilder.Build(ctx)
	if session.GetCandID() != uuid.Nil {
		vacList, err = v.vacancyClient.GetRecommendation(session.GetUserID(), int(req.Start), int(req.Limit))
		if err != nil {
			if err.Error() == common.NoRecommendation {
				ctx.JSON(http.StatusOK, common.RespError{Err: common.NoRecommendation})
				return
			}
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
			return
		}
	}
	ctx.JSON(http.StatusOK, RespList{Vacancies: vacList})

}

func (v *VacancyHandler) SearchVacanciesHandler(ctx *gin.Context) {
	var searchParams models.VacancySearchParams

	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}

	if err := common.ReqValidation(&searchParams); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}
	VacList, err := v.vacancyClient.SearchVacancies(searchParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}
	ctx.JSON(http.StatusOK, RespList{Vacancies: VacList})
}

func vacHandlerCommon(v *VacancyHandler, ctx *gin.Context, treatmentType int) {
	var (
		req vacRequest
		err error
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
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.SessionErr})
			return
		}
		vacNew, err = v.vacancyClient.CreateVacancy(*vacNew)
	} else if vacNew.ID, _ = uuid.Parse(req.Id); vacNew.ID != uuid.Nil {
		vacNew, err = v.vacancyClient.UpdateVacancy(*vacNew)
	} else {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, vacPath+vacNew.ID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
	ctx.JSON(http.StatusOK, Resp{Vacancy: vacNew})
}

func vacListHandlerCommon(v *VacancyHandler, ctx *gin.Context, entityType int) {
	var (
		req     vacListRequest
		err     error
		vacList []models.Vacancy
		id      = uuid.Nil
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if entityType == vacancy.ByEmpId {
		session := v.SessionBuilder.Build(ctx)
		emplID := session.GetEmplID()
		if id, err = uuid.Parse(emplID.String()); err != nil {
			ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.SessionErr})
			return
		}
	} else if entityType == vacancy.ByCompId {
		id, _ = uuid.Parse(req.CompId)
	}
	vacList, err = v.vacancyClient.GetVacancyList(req.Start, req.Limit, id, entityType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}
	ctx.JSON(http.StatusOK, RespList{Vacancies: vacList})
}
