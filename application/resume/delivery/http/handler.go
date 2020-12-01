package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
	"path"
)

type ResumeHandler struct {
	UseCaseResume           resume.UseCase
	UseCaseEducation        education.UseCase
	UseCaseCustomExperience custom_experience.UseCase
	SessionBuilder          common.SessionBuilder
}

const resumePath = "resume/"

func NewRest(router *gin.RouterGroup,
	useCaseResume resume.UseCase,
	useCaseEducation education.UseCase,
	useCaseCustomExperience custom_experience.UseCase,
	sessionBuilder common.SessionBuilder,
	AuthRequired gin.HandlerFunc) *ResumeHandler {
	rest := &ResumeHandler{
		UseCaseResume:           useCaseResume,
		UseCaseEducation:        useCaseEducation,
		UseCaseCustomExperience: useCaseCustomExperience,
		SessionBuilder:          sessionBuilder,
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
	session := r.SessionBuilder.Build(ctx)
	candID := session.GetCandID()
	if candID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err:  common.AuthRequiredErr})
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := r.UseCaseResume.GetAllUserResume(candID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err:  common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (r *ResumeHandler) CreateResume(ctx *gin.Context) {
	session := r.SessionBuilder.Build(ctx)
	candID := session.GetCandID()
	if candID == uuid.Nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var template models.Resume
	if err := ctx.ShouldBindJSON(&template); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := common.ReqValidation(template); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}
	template.CandID = candID
	template.ResumeID = uuid.New()

	file, errImg := common.GetImageFromBase64(template.Avatar)
	if errImg != nil {
		ctx.JSON(http.StatusBadRequest, errImg)
		return
	}
	avatarName := resumePath + template.ResumeID.String()
	if file != nil {
		template.Avatar = path.Join(common.DOMAIN, common.ImgDir, avatarName)
	}
	result, err := r.UseCaseResume.Create(template)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err:  common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if file != nil {
		if err := common.AddOrUpdateUserFile(file, avatarName); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	resp := resume.Response{
		User:             result.Candidate.User,
		Educations:       result.Education,
		CustomExperience: result.ExperienceCustomComp,
		IsFavorite:       nil,
	}

	//result.Candidate = nil
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
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
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err:  common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var isFavorite *uuid.UUID = nil
	session := r.SessionBuilder.Build(ctx)
	if session != nil {
		emplID := session.GetEmplID()
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
	}

	resp := resume.Response{
		User:             result.Candidate.User,
		Educations:       result.Education,
		CustomExperience: result.ExperienceCustomComp,
		IsFavorite:       isFavorite,
	}

	//result.Candidate = nil
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
	session := r.SessionBuilder.Build(ctx)
	candID := session.GetCandID()
	if candID == uuid.Nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var template models.Resume
	if err := ctx.ShouldBindJSON(&template); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
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
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err:  common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if file != nil {
		if err := common.AddOrUpdateUserFile(file, resumePath+result.ResumeID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	resp := resume.Response{
		User:             result.Candidate.User,
		Educations:       result.Education,
		CustomExperience: result.ExperienceCustomComp,
		IsFavorite:       nil,
	}

	//result.Candidate = nil
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resumeID, err := uuid.Parse(request.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	session := r.SessionBuilder.Build(ctx)
	emplID := session.GetEmplID()
	if emplID == uuid.Nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	favoriteForEmpl := models.FavoritesForEmpl{EmplID: emplID, ResumeID: resumeID}

	favorite, err := r.UseCaseResume.AddFavorite(favoriteForEmpl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err:  common.DataBaseErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//type Response struct {
	//	Favorite models.FavoritesForEmpl `json:"favorite_for_empl"`
	//}

	ctx.JSON(http.StatusOK, *favorite)
}

func (r *ResumeHandler) RemoveFavorite(ctx *gin.Context) {
	var request struct {
		FavoriteID string `uri:"favorite_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	favoriteID, err := uuid.Parse(request.FavoriteID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	session := r.SessionBuilder.Build(ctx)
	emplID := session.GetEmplID()
	if emplID == uuid.Nil {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	favoriteForEmpl := models.FavoritesForEmpl{FavoriteID: favoriteID, EmplID: emplID}
	err = r.UseCaseResume.RemoveFavorite(favoriteForEmpl)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err:  common.DataBaseErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (r *ResumeHandler) GetAllFavoritesResume(ctx *gin.Context) {
	session := r.SessionBuilder.Build(ctx)
	emplID := session.GetEmplID()
	if emplID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err:  common.AuthRequiredErr})
		ctx.AbortWithError(http.StatusBadRequest, errors.Errorf(common.AuthRequiredErr))
		return
	}

	emplFavoriteResume, err := r.UseCaseResume.GetAllEmplFavoriteResume(emplID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err:  common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, emplFavoriteResume)
}
