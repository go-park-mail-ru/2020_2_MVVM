package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"net/http"
)

type ResumeResponse struct {
	Resume           models.Resume                  `json:"resume"`
	User             models.User                    `json:"user"`
	Educations       []*models.Education            `json:"education"`
	CustomExperience []*models.ExperienceCustomComp `json:"custom_experience"`
	IsFavorite       *uuid.UUID                     `json:"is_favorite"`
}

type ResumeHandler struct {
	UseCaseResume           resume.UseCase
	UseCaseEducation        education.UseCase
	UseCaseCustomCompany    custom_company.UseCase
	UseCaseCustomExperience custom_experience.UseCase
}

const resumePath = "resume/"

func NewRest(router *gin.RouterGroup,
	useCaseResume resume.UseCase,
	useCaseEducation education.UseCase,
	useCaseCustomCompany custom_company.UseCase,
	useCaseCustomExperience custom_experience.UseCase,
	AuthRequired gin.HandlerFunc) *ResumeHandler {
	rest := &ResumeHandler{
		UseCaseResume:           useCaseResume,
		UseCaseEducation:        useCaseEducation,
		UseCaseCustomCompany:    useCaseCustomCompany,
		UseCaseCustomExperience: useCaseCustomExperience,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (r *ResumeHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:resume_id", r.GetResumeByID)
	router.GET("/page", r.GetResumePage)
	router.POST("/search", r.SearchResume)
	router.Use(AuthRequired)
	{
		router.GET("/mine", r.GetMineResume)
		router.POST("/", r.CreateResume)
		router.PUT("/", r.UpdateResume)

		router.POST("/favorite/by/id/:resume_id", r.AddFavorite)
		router.DELETE("/favorite/by/id/:favorite_id", r.RemoveFavorite)
		router.GET("/myfavorites", r.GetAllFavoritesResume)
	}
}

func (r *ResumeHandler) GetMineResume(ctx *gin.Context) {
	candID, err := common.GetCurrentUserId(ctx, "cand_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := r.UseCaseResume.GetAllUserResume(candID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (r *ResumeHandler) CreateResume(ctx *gin.Context) {
	candID, err := common.GetCurrentUserId(ctx, "cand_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var template *models.Resume
	if err := ctx.ShouldBindBodyWith(&template, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	template.CandID = candID

	file, errImg := common.GetImageFromBase64(template.Avatar)
	if errImg != nil {
		ctx.JSON(http.StatusBadRequest, errImg)
		return
	}

	resume, err := r.UseCaseResume.Create(*template)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if file != nil {
		if err := common.AddOrUpdateUserFile(file, resumePath+resume.ResumeID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	resp := ResumeResponse{
		User:             *resume.Candidate.User,
		Educations:       resume.Education,
		CustomExperience: resume.ExperienceCustomComp,
		IsFavorite:       nil,
	}

	resume.Candidate = nil
	resume.Education = nil
	resume.ExperienceCustomComp = nil
	resp.Resume = *resume

	ctx.JSON(http.StatusOK, resp)
}

func (r *ResumeHandler) GetResumeByID(ctx *gin.Context) {
	var request struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resumeID, err := uuid.Parse(request.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	result, err := r.UseCaseResume.GetById(resumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var isFavorite *uuid.UUID = nil
	emplID, err := common.GetCurrentUserId(ctx, "empl_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if emplID != uuid.Nil {
		favorite, err := r.UseCaseResume.GetFavoriteByResume(emplID, result.ResumeID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if favorite != nil {
			isFavorite = &favorite.FavoriteID
		}
	}

	resp := ResumeResponse{
		User:             *result.Candidate.User,
		Educations:       result.Education,
		CustomExperience: result.ExperienceCustomComp,
		IsFavorite:       isFavorite,
	}

	result.Candidate = nil
	result.Education = nil
	result.ExperienceCustomComp = nil
	resp.Resume = *result

	ctx.JSON(http.StatusOK, resp)
}

func (r *ResumeHandler) GetResumePage(ctx *gin.Context) {
	var request struct {
		Start uint `form:"start"`
		Limit uint `form:"limit" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resumes, err := r.UseCaseResume.List(request.Start, request.Limit)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, resumes)
}

func (r *ResumeHandler) UpdateResume(ctx *gin.Context) {
	candID, err := common.GetCurrentUserId(ctx, "cand_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if candID == uuid.Nil {
		ctx.AbortWithError(http.StatusForbidden, err)
		return
	}

	var template models.Resume
	if err := ctx.ShouldBindBodyWith(&template, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	template.CandID = candID



	result, err := r.UseCaseResume.Update(template)

	//for i := range additionParam.CustomExperience {
	//	item := additionParam.CustomExperience[i]
	//	dateBegin, err := time.Parse(time.RFC3339, item.Begin+"T00:00:00Z")
	//	if err != nil {
	//		ctx.AbortWithError(http.StatusBadRequest, err)
	//		return
	//	}
	//	var dateFinish time.Time
	//	if !item.ContinueToToday {
	//		dateFinish, err = time.Parse(time.RFC3339, *item.Finish+"T00:00:00Z")
	//		if err != nil {
	//			ctx.AbortWithError(http.StatusBadRequest, err)
	//			return
	//		}
	//	} else {
	//		dateFinish = time.Now()
	//	}
	//	//dateBegin := time.Now()
	//	//dateFinish := time.Now()
	//	insertExp := models.ExperienceCustomComp{
	//		NameJob:         item.NameJob,
	//		Position:        item.Position,
	//		Begin:           dateBegin,
	//		Finish:          &dateFinish,
	//		Duties:          item.Duties,
	//		ContinueToToday: &item.ContinueToToday,
	//	}
	//	customExperience = append(customExperience, insertExp)
	//}
	//
	//pCustomExperience, err := r.handlerUpdateCustomExperience(customExperience, candID, pResume.ResumeID)
	//if err != nil {
	//	ctx.AbortWithError(http.StatusBadRequest, err)
	//	return
	//}

	resp := ResumeResponse{
		User:             *result.Candidate.User,
		Educations:       result.Education,
		CustomExperience: result.ExperienceCustomComp,
		IsFavorite:       nil,
	}

	result.Candidate = nil
	result.Education = nil
	result.ExperienceCustomComp = nil
	resp.Resume = *result

	ctx.JSON(http.StatusOK, resp)

}

func (r *ResumeHandler) SearchResume(ctx *gin.Context) {
	var searchParams resume.SearchParams
	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	found, err := r.UseCaseResume.Search(searchParams)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, found)
}

func (r *ResumeHandler) AddFavorite(ctx *gin.Context) {
	var request struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resumeID, err := uuid.Parse(request.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	emplID, err := common.GetCurrentUserId(ctx, "empl_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if emplID == uuid.Nil {
		ctx.AbortWithError(http.StatusForbidden, err)
		return
	}

	favoriteForEmpl := models.FavoritesForEmpl{EmplID: emplID, ResumeID: resumeID}

	favorite, err := r.UseCaseResume.AddFavorite(favoriteForEmpl)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	type Response struct {
		Favorite models.FavoritesForEmpl `json:"favorite_for_empl"`
	}

	ctx.JSON(http.StatusOK, Response{Favorite: *favorite})
}

func (r *ResumeHandler) RemoveFavorite(ctx *gin.Context) {
	var request struct {
		FavoriteID string `uri:"favorite_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	favoriteID, err := uuid.Parse(request.FavoriteID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	emplID, err := common.GetCurrentUserId(ctx, "empl_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if emplID == uuid.Nil {
		ctx.AbortWithError(http.StatusForbidden, err)
		return
	}

	favoriteForEmpl := models.FavoritesForEmpl{FavoriteID: favoriteID, EmplID: emplID}
	err = r.UseCaseResume.RemoveFavorite(favoriteForEmpl)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *ResumeHandler) GetAllFavoritesResume(ctx *gin.Context) {
	emplID, err := common.GetCurrentUserId(ctx, "empl_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	emplFavoriteResume, err := r.UseCaseResume.GetAllEmplFavoriteResume(emplID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, emplFavoriteResume)
}
