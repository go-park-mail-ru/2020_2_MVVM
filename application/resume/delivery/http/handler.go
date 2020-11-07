package http

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_company"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/custom_experience"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/education"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/resume"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type ResumeHandler struct {
	UseCaseResume           resume.UseCase
	UseCaseEducation        education.UseCase
	UseCaseCustomCompany    custom_company.UseCase
	UseCaseCustomExperience custom_experience.UseCase
}

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

		// TODO почему favorite тут а не в employer?
		router.POST("/favorite/by/id/:resume_id", r.AddFavorite)
		router.DELETE("/favorite/by/id/:favorite_id", r.RemoveFavorite)
		router.GET("/myfavorites", r.GetAllFavoritesResume)
	}
}

// TODO: связанная с auth функциональность должна быть в auth
func (r *ResumeHandler) GetCurrentUserId(ctx *gin.Context, user_type string) (id uuid.UUID, err error) {
	session := sessions.Default(ctx)
	userIDStr := session.Get(user_type)
	if userIDStr == nil {
		return uuid.Nil, nil
	}
	userID, err := uuid.Parse(userIDStr.(string))

	if err != nil {
		return uuid.Nil, err
	}
	return userID, nil
}

func (r *ResumeHandler) GetMineResume(ctx *gin.Context) {
	candID, err := r.GetCurrentUserId(ctx, "cand_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pResume, err := r.UseCaseResume.GetAllUserResume(candID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, pResume)
}

func (r *ResumeHandler) CreateResume(ctx *gin.Context) {
	candID, err := r.GetCurrentUserId(ctx, "cand_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reqResume models.Resume
	if err := ctx.ShouldBindBodyWith(&reqResume, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reqResume.CandID = candID
	reqResume.DateCreate = time.Now()

	pResume, err := r.UseCaseResume.CreateResume(reqResume)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	user := *pResume.Candidate.User
	pResume.Candidate = nil

	var additionParam models.AdditionInResume
	if err := ctx.ShouldBindBodyWith(&additionParam, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pEducations, err := r.createEducation(additionParam.Education, candID, pResume.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var customExperience []models.ExperienceCustomComp
	for i := range additionParam.CustomExperience {
		item := additionParam.CustomExperience[i]
		dateBegin, err := time.Parse(time.RFC3339, item.Begin+"T00:00:00Z")
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var dateFinish time.Time
		if !item.ContinueToToday {
			dateFinish, err = time.Parse(time.RFC3339, *item.Finish+"T00:00:00Z")
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			dateFinish = time.Now()
		}
		//dateBegin := time.Now()
		//dateFinish := time.Now()

		insertExp := models.ExperienceCustomComp{
			NameJob:         item.NameJob,
			Position:        item.Position,
			Begin:           dateBegin,
			Finish:          &dateFinish,
			Duties:          item.Duties,
			ContinueToToday: &item.ContinueToToday,
		}
		customExperience = append(customExperience, insertExp)
	}

	pCustomExperience, err := r.createCustomExperience(customExperience, candID, pResume.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK,
		models.RespResume{Resume: *pResume, User: user, Educations: pEducations, CustomExperience: pCustomExperience})
}

func (r *ResumeHandler) createEducation(educations []models.Education, userID, resumeID uuid.UUID) ([]models.Education, error) {
	for i := range educations {
		educations[i].CandId = userID
		educations[i].ResumeId = resumeID
		educations[i].Finish = time.Now()
	}

	pEducations, err := r.UseCaseEducation.CreateEducation(educations)
	if err != nil {
		return nil, err
	}
	return pEducations, nil
}

func (r *ResumeHandler) createCustomExperience(experiences []models.ExperienceCustomComp, userID, resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	for i := range experiences {
		experiences[i].CandID = userID
		experiences[i].ResumeID = resumeID
	}

	pCustomExperience, err := r.UseCaseCustomExperience.CreateCustomExperience(experiences)
	if err != nil {
		return nil, err
	}

	return pCustomExperience, nil
}

func (r *ResumeHandler) GetResumeByID(ctx *gin.Context) {
	var reqResume struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&reqResume); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resumeID, err := uuid.Parse(reqResume.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pResume, err := r.UseCaseResume.GetResume(resumeID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user := *pResume.Candidate.User
	pResume.Candidate = nil

	pEducations, err := r.UseCaseEducation.GetAllResumeEducation(resumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pCustomExperience, err := r.UseCaseCustomExperience.GetAllResumeCustomExperience(resumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var isFavorite *uuid.UUID = nil
	emplID, err := r.GetCurrentUserId(ctx, "empl_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if emplID != uuid.Nil {
		favorite, err := r.UseCaseResume.GetFavorite(emplID, pResume.ResumeID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if favorite != nil {
			isFavorite = &favorite.ID
		}
	}

	ctx.JSON(http.StatusOK,
		models.RespResume{Resume: *pResume,
			User:             user,
			Educations:       pEducations,
			CustomExperience: pCustomExperience,
			IsFavorite:       isFavorite})
}

func (r *ResumeHandler) GetResumePage(ctx *gin.Context) {
	var reqResume struct {
		Start uint `form:"start"`
		Limit uint `form:"limit" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&reqResume); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	respResumes, err := r.UseCaseResume.GetResumePage(reqResume.Start, reqResume.Limit)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, respResumes)
}

func (r *ResumeHandler) UpdateResume(ctx *gin.Context) {
	candID, err := r.GetCurrentUserId(ctx, "cand_id")
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if candID == uuid.Nil {
		ctx.AbortWithError(http.StatusForbidden, err)
		return
	}

	var reqResume models.Resume
	if err := ctx.ShouldBindBodyWith(&reqResume, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reqResume.CandID = candID
	reqResume.DateCreate = time.Now()

	pResume, err := r.UseCaseResume.UpdateResume(reqResume)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	user := *pResume.Candidate.User
	pResume.Candidate = nil

	var additionParam models.AdditionInResume
	if err := ctx.ShouldBindBodyWith(&additionParam, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	pEducations, err := r.handlerUpdateEducation(additionParam.Education, candID, pResume.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var customExperience []models.ExperienceCustomComp
	for i := range additionParam.CustomExperience {
		item := additionParam.CustomExperience[i]
		dateBegin, err := time.Parse(time.RFC3339, item.Begin+"T00:00:00Z")
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var dateFinish time.Time
		if !item.ContinueToToday {
			dateFinish, err = time.Parse(time.RFC3339, *item.Finish+"T00:00:00Z")
			if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			dateFinish = time.Now()
		}
		//dateBegin := time.Now()
		//dateFinish := time.Now()
		insertExp := models.ExperienceCustomComp{
			NameJob:         item.NameJob,
			Position:        item.Position,
			Begin:           dateBegin,
			Finish:          &dateFinish,
			Duties:          item.Duties,
			ContinueToToday: &item.ContinueToToday,
		}
		customExperience = append(customExperience, insertExp)
	}

	pCustomExperience, err := r.handlerUpdateCustomExperience(customExperience, candID, pResume.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK,
		models.RespResume{Resume: *pResume, User: user, Educations: pEducations, CustomExperience: pCustomExperience})

}

func (r *ResumeHandler) handlerUpdateEducation(educations []models.Education, userID, resumeID uuid.UUID) ([]models.Education, error) {
	for i := range educations {
		if educations[i].CandId == uuid.Nil && educations[i].ResumeId == uuid.Nil {
			educations[i].CandId = userID
			educations[i].ResumeId = resumeID
		} else if educations[i].CandId != userID && educations[i].ResumeId != resumeID {
			return nil, errors.New("this user has not update this resume")
		}
	}

	pEducations, err := r.UseCaseEducation.UpdateEducation(educations, resumeID)
	if err != nil {
		return nil, err
	}
	return pEducations, nil
}

func (r *ResumeHandler) handlerUpdateCustomExperience(experience []models.ExperienceCustomComp, userID, resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	for i := range experience {
		if experience[i].CandID == uuid.Nil && experience[i].ResumeID == uuid.Nil {
			experience[i].CandID = userID
			experience[i].ResumeID = resumeID
		} else if experience[i].CandID != userID && experience[i].ResumeID != resumeID {
			return nil, errors.New("this user has not update this resume")
		}
	}

	pEducations, err := r.UseCaseCustomExperience.UpdateCustomExperience(experience, resumeID)
	if err != nil {
		return nil, err
	}
	return pEducations, nil
}

func (r *ResumeHandler) SearchResume(ctx *gin.Context) {
	var searchParams resume.SearchParams
	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	found, err := r.UseCaseResume.SearchResume(searchParams)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, found)
}

func (r *ResumeHandler) AddFavorite(ctx *gin.Context) {
	var reqFavorite struct {
		ResumeID string `uri:"resume_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&reqFavorite); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resumeID, err := uuid.Parse(reqFavorite.ResumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	emplID, err := r.GetCurrentUserId(ctx, "empl_id")
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

	type RespFavorite struct {
		Favorite models.FavoritesForEmpl `json:"favorite_for_empl"`
	}

	ctx.JSON(http.StatusOK, RespFavorite{Favorite: *favorite})
}

func (r *ResumeHandler) RemoveFavorite(ctx *gin.Context) {
	var reqFavorite struct {
		FavoriteID string `uri:"favorite_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&reqFavorite); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	favoriteID, err := uuid.Parse(reqFavorite.FavoriteID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = r.UseCaseResume.RemoveFavorite(favoriteID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *ResumeHandler) GetAllFavoritesResume(ctx *gin.Context) {
	emplID, err := r.GetCurrentUserId(ctx, "empl_id")
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

func (r *ResumeHandler) handlerGetAllForListResume(resume []models.Resume) ([]models.RespResume, error) {
	var allResume []models.RespResume
	for i := range resume {
		exp, err := r.UseCaseCustomExperience.GetAllResumeCustomExperience(resume[i].ResumeID)
		if err != nil {
			return nil, err
		}
		educ, err := r.UseCaseEducation.GetAllResumeEducation(resume[i].ResumeID)
		if err != nil {
			return nil, err
		}
		wholeResume := models.RespResume{
			Resume:           resume[i],
			Educations:       educ,
			CustomExperience: exp,
		}

		allResume = append(allResume, wholeResume)
	}
	return allResume, nil
}
