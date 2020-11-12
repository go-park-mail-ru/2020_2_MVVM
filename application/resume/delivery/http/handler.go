package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"net/http"
)

type ResumeHandler struct {
	UseCaseResume           resume.UseCase
	UseCaseEducation        education.UseCase
	UseCaseCustomExperience custom_experience.UseCase
}

const resumePath = "resume/"

func NewRest(router *gin.RouterGroup,
	useCaseResume resume.UseCase,
	useCaseEducation education.UseCase,
	useCaseCustomExperience custom_experience.UseCase,
	AuthRequired gin.HandlerFunc) *ResumeHandler {
	rest := &ResumeHandler{
		UseCaseResume:           useCaseResume,
		UseCaseEducation:        useCaseEducation,
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
	if err := common.ReqValidation(template); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}
	template.CandID = candID

	file, errImg := common.GetImageFromBase64(template.Avatar)
	if errImg != nil {
		ctx.JSON(http.StatusBadRequest, errImg)
		return
	}

	result, err := r.UseCaseResume.Create(*template)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if file != nil {
		if err := common.AddOrUpdateUserFile(file, resumePath+result.ResumeID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	resp := resume.Response{
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

	resp := resume.Response{
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}

	resumes, err := r.UseCaseResume.List(request.Start, request.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
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
	if err := common.ReqValidation(template); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}
	template.CandID = candID

	file, errImg := common.GetImageFromBase64(template.Avatar)
	if errImg != nil {
		ctx.JSON(http.StatusBadRequest, errImg)
		return
	}

	result, err := r.UseCaseResume.Update(template)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if file != nil {
		if err := common.AddOrUpdateUserFile(file, resumePath+result.ResumeID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	resp := resume.Response{
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err := common.ReqValidation(&searchParams); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}
	found, err := r.UseCaseResume.Search(searchParams)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.DataBaseErr})
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
