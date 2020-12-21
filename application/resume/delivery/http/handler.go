package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	resume2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/resume"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"path"
	"strings"
)

type ResumeHandler struct {
	UseCaseResume resume.UseCase
	//UseCaseEducation        education.UseCase
	UseCaseCustomExperience custom_experience.UseCase
	SessionBuilder          common.SessionBuilder
}

const resumePath = "resume/"

func NewRest(router *gin.RouterGroup,
	useCaseResume resume.UseCase,
//useCaseEducation education.UseCase,
	useCaseCustomExperience custom_experience.UseCase,
	sessionBuilder common.SessionBuilder,
	AuthRequired gin.HandlerFunc) *ResumeHandler {
	rest := &ResumeHandler{
		UseCaseResume: useCaseResume,
		//UseCaseEducation:        useCaseEducation,
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
		router.DELETE("/resume/:resume_id", r.DeleteResume)

		router.GET("/make/pdf/:resume_id", r.MakePdf)

		router.GET("/favorite/:resume_id", r.GetFavorite)
		router.POST("/favorite/by/id/:resume_id", r.AddFavorite)
		router.DELETE("/favorite/by/id/:favorite_id", r.RemoveFavorite)
		router.GET("/myfavorites", r.GetAllFavoritesResume)
	}
}

func (r *ResumeHandler) GetMineResume(ctx *gin.Context) {
	session := r.SessionBuilder.Build(ctx)
	candID := session.GetCandID()
	if candID == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.AuthRequiredErr)
		//ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	result, err := r.UseCaseResume.GetAllUserResume(candID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListBriefResumeInfo(result), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) CreateResume(ctx *gin.Context) {
	var (
		file       io.Reader
		avatarName string
		errImg     *common.Err
	)
	session := r.SessionBuilder.Build(ctx)
	candID := session.GetCandID()
	if candID == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusForbidden, common.AuthRequiredErr)
		return
	}

	template := new(models.Resume)
	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body, template); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := common.ReqValidation(template); err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	template.CandID = candID
	template.ResumeID = uuid.New()
	if !strings.Contains(template.Avatar, "http") {
		file, errImg = common.GetImageFromBase64(template.Avatar)
		if errImg != nil {
			common.WriteErrResponse(ctx, http.StatusBadRequest, errImg.Error())
			return
		}
		avatarName = resumePath + template.ResumeID.String()
		if file != nil {
			template.Avatar = common.DOMAIN + path.Join(common.ImgDir, avatarName)
		}
	}
	result, err := r.UseCaseResume.Create(*template)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if file != nil {
		if err := common.AddOrUpdateUserFile(file, avatarName); err != nil {
			//ctx.AbortWithError(http.StatusInternalServerError, err)
			common.WriteErrResponse(ctx, http.StatusInternalServerError, err.String())
		}
	}

	resp := resume2.Response{
		Educations:       result.Education,
		CustomExperience: result.ExperienceCustomComp,
	}
	result.Education = nil
	result.ExperienceCustomComp = nil
	resp.Resume = *result

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) GetResumeByID(ctx *gin.Context) {
	var request struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resumeID, _ := uuid.Parse(request.ResumeID)

	result, err := r.UseCaseResume.GetById(resumeID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp := resume2.Response{
		Educations:       result.Education,
		CustomExperience: result.ExperienceCustomComp,
	}

	result.Education = nil
	result.ExperienceCustomComp = nil
	resp.Resume = *result

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) GetResumePage(ctx *gin.Context) {
	var request struct {
		Start uint `form:"start"`
		Limit uint `form:"limit" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}

	resumes, err := r.UseCaseResume.List(request.Start, request.Limit)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListBriefResumeInfo(resumes), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) UpdateResume(ctx *gin.Context) {
	var (
		file       io.Reader
		avatarName string
		errImg     *common.Err
		template models.Resume
	)
	session := r.SessionBuilder.Build(ctx)
	candID := session.GetCandID()
	if candID == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusForbidden, common.AuthRequiredErr)
		return
	}
	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body, &template); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := common.ReqValidation(template); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	template.CandID = candID
	if !strings.Contains(template.Avatar, "http") {
		file, errImg = common.GetImageFromBase64(template.Avatar)
		if errImg != nil {
			common.WriteErrResponse(ctx, http.StatusBadRequest, errImg.Error())
			return
		}
		avatarName = resumePath + template.ResumeID.String()
		if file != nil {
			template.Avatar = common.DOMAIN + path.Join(common.ImgDir, avatarName)
		}
	}
	result, err := r.UseCaseResume.Update(template)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if file != nil {
		if err := common.AddOrUpdateUserFile(file, resumePath+result.ResumeID.String()); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	resp := resume2.Response{
		Educations:       result.Education,
		CustomExperience: result.ExperienceCustomComp,
	}
	result.Education = nil
	result.ExperienceCustomComp = nil
	resp.Resume = *result

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) SearchResume(ctx *gin.Context) {
	var request resume2.StartLimit
	if err := ctx.ShouldBindQuery(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}

	var searchParams resume2.SearchParams
	searchParams.StartLimit = request
	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body, &searchParams); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := common.ReqValidation(&searchParams); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	found, err := r.UseCaseResume.Search(searchParams)
	if err != nil {
		//ctx.AbortWithError(http.StatusBadRequest, err)
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListBriefResumeInfo(found), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) AddFavorite(ctx *gin.Context) {
	var request struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resumeID, _ := uuid.Parse(request.ResumeID)

	session := r.SessionBuilder.Build(ctx)
	emplID := session.GetEmplID()
	if emplID == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusForbidden, common.AuthRequiredErr)
		return
	}

	favoriteForEmpl := models.FavoritesForEmpl{EmplID: emplID, ResumeID: resumeID}

	favorite, err := r.UseCaseResume.AddFavorite(favoriteForEmpl)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//type Response struct {
	//	Favorite models.FavoritesForEmpl `json:"favorite_for_empl"`
	//}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(*favorite, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) RemoveFavorite(ctx *gin.Context) {
	var request struct {
		FavoriteID string `uri:"favorite_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	favoriteID, _ := uuid.Parse(request.FavoriteID)

	session := r.SessionBuilder.Build(ctx)
	emplID := session.GetEmplID()
	if emplID == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusForbidden, common.AuthRequiredErr)
		return
	}

	favoriteForEmpl := models.FavoritesForEmpl{FavoriteID: favoriteID, EmplID: emplID}
	err := r.UseCaseResume.RemoveFavorite(favoriteForEmpl)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(nil, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) GetAllFavoritesResume(ctx *gin.Context) {
	session := r.SessionBuilder.Build(ctx)
	emplID := session.GetEmplID()
	if emplID == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.AuthRequiredErr)
		//ctx.AbortWithError(http.StatusBadRequest, errors.Errorf(common.AuthRequiredErr))
		return
	}

	emplFavoriteResume, err := r.UseCaseResume.GetAllEmplFavoriteResume(emplID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListBriefResumeInfo(emplFavoriteResume), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResumeHandler) GetFavorite(ctx *gin.Context) {
	var request struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resumeID, _ := uuid.Parse(request.ResumeID)
	favorite := new(models.FavoritesForEmpl)

	var err error
	session := r.SessionBuilder.Build(ctx)

	emplID := session.GetEmplID()
	if emplID != uuid.Nil {
		favorite, err = r.UseCaseResume.GetFavoriteByResume(emplID, resumeID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	if favorite == nil || favorite.FavoriteID == uuid.Nil {
		if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.FavoriteID{FavoriteID: nil}, ctx.Writer); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	} else {
		if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.FavoriteID{FavoriteID: &favorite.FavoriteID},
			ctx.Writer); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}

func (r *ResumeHandler) DeleteResume(ctx *gin.Context) {
	var req struct {
		ResId string `uri:"resume_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	ResId, _ := uuid.Parse(req.ResId)
	session := r.SessionBuilder.Build(ctx)
	candId := session.GetCandID()
	if candId == uuid.Nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.SessionErr)
		return
	}
	if err := r.UseCaseResume.DeleteResume(ResId, candId); err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
}

func (r *ResumeHandler) MakePdf(ctx *gin.Context) {
	var request struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resumeID, _ := uuid.Parse(request.ResumeID)

	err := r.UseCaseResume.MakePdf(resumeID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	link := models.LinkToPdf{Link: common.DOMAIN + common.PathToPdf + resumeID.String() + ".pdf"}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(link, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
