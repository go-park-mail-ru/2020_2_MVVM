package http

import (
	"errors"
	"github.com/gin-contrib/sessions"
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
	"time"
)

type ResumeHandler struct {
	UsecaseResume           resume.IUseCaseResume
	UsecaseEducation        education.IUseCaseEducation
	UsecaseCustomCompany    custom_company.IUseCaseCustomCompany
	UsecaseCustomExperience custom_experience.IUseCaseCustomExperience
}

func NewRest(router *gin.RouterGroup,
	usecaseResume resume.IUseCaseResume,
	usecaseEducation education.IUseCaseEducation,
	usecaseCustomCompany custom_company.IUseCaseCustomCompany,
	usecaseCustomExperience custom_experience.IUseCaseCustomExperience,
	AuthRequired gin.HandlerFunc) *ResumeHandler {
	rest := &ResumeHandler{
		UsecaseResume:           usecaseResume,
		UsecaseEducation:        usecaseEducation,
		UsecaseCustomCompany:    usecaseCustomCompany,
		UsecaseCustomExperience: usecaseCustomExperience,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (r *ResumeHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:resume_id", r.handlerGetResumeByID)
	router.GET("/page", r.handlerGetResumeList)
	router.Use(AuthRequired)
	{
		router.GET("/mine", r.handlerGetAllCurrentUserResume)
		router.POST("/", r.handlerCreateResume)
		router.PUT("/", r.handlerUpdateResume)
	}
}

func (r *ResumeHandler) handlerGetAllCurrentUserResume(ctx *gin.Context) {
	//identityKey := "myid"
	//jwtuser, _ := ctx.Get(identityKey)
	//userID := jwtuser.(*models.JWTUserData).ID
	session := sessions.Default(ctx)
	userIDStr := session.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pResume, err := r.UsecaseResume.GetAllUserResume(userID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume []models.Resume `json:"resume"`
	}

	ctx.JSON(http.StatusOK, Resp{Resume: pResume})
}

func (r *ResumeHandler) handlerCreateResume(ctx *gin.Context) {
	// move to constants
	//identityKey := "myid"
	//jwtuser, _ := ctx.Get(identityKey)
	//userID := jwtuser.(*models.JWTUserData).ID

	session := sessions.Default(ctx)
	userIDStr := session.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reqResume models.Resume
	if err := ctx.ShouldBindBodyWith(&reqResume, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	reqResume.UserID = userID
	reqResume.DateCreate = time.Now()

	pResume, err := r.UsecaseResume.CreateResume(reqResume)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type ReqExperienceCustomComp struct {
		NameJob         string  `json:"name_job"`
		Position        *string `json:"position"`
		Begin           string  `json:"begin"`
		Finish          *string `json:"finish"`
		Duties          *string `json:"duties"`
		ContinueToToday bool   `json:"continue_to_today"`
	}

	var AdditionInResume struct {
		Education        []models.Education        `json:"education"`
		CustomExperience []ReqExperienceCustomComp `json:"custom_experience"`
	}

	//var additionParam models.AdditionInResume
	if err := ctx.ShouldBindBodyWith(&AdditionInResume, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pEducations, err := r.handlerCreateEducation(AdditionInResume.Education, userID, pResume.ID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var customExperience []models.ExperienceCustomComp
	for i := range AdditionInResume.CustomExperience {
		item := AdditionInResume.CustomExperience[i]
		dateBedin, err := time.Parse(time.RFC3339, item.Begin + "T00:00:00Z"); if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		var dateFinish time.Time
		if !item.ContinueToToday {
			dateFinish, err = time.Parse(time.RFC3339, *item.Finish + "T00:00:00Z"); if err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
		} else {
			dateFinish = time.Now()
		}
		insertExp := models.ExperienceCustomComp{
			NameJob:         item.NameJob,
			Position:        item.Position,
			Begin:           dateBedin,
			Finish:          &dateFinish,
			Duties:          item.Duties,
			ContinueToToday: &item.ContinueToToday,
		}
		customExperience = append(customExperience, insertExp)
	}

	pCustomExperience, err := r.handlerCreateCustomExperience(customExperience, userID, pResume.ID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	type RespResume struct {
		Resume           models.Resume                 `json:"resume"`
		Educations       []models.Education            `json:"education"`
		CustomExperience []models.ExperienceCustomComp `json:"custom_experience"`
	}

	ctx.JSON(http.StatusOK,
		RespResume{Resume: *pResume, Educations: pEducations, CustomExperience: pCustomExperience})
}

func (r *ResumeHandler) handlerCreateEducation(educations []models.Education, userID, resumeID uuid.UUID) ([]models.Education, error) {
	for i := range educations {
		educations[i].CandId = userID
		educations[i].ResumeId = resumeID
		educations[i].Finish = time.Now()
	}

	pEducations, err := r.UsecaseEducation.CreateEducation(educations)
	if err != nil {
		return nil, err
	}
	return pEducations, nil
}

func (r *ResumeHandler) handlerCreateCustomExperience(experiences []models.ExperienceCustomComp, userID, resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	for i := range experiences {
		experiences[i].CandID = userID
		experiences[i].ResumeID = resumeID
		experiences[i].Begin = time.Now()
	}

	pCustomExperience, err := r.UsecaseCustomExperience.CreateCustomExperience(experiences)
	if err != nil {
		return nil, err
	}

	for i := range pCustomExperience {
		pCustomExperience[i].Begin.Format("2020-10-28")
		//pCustomExperience[i].Finish.Format("2020-10-28")
	}

	return pCustomExperience, nil
}

func (r *ResumeHandler) handlerGetResumeByID(ctx *gin.Context) {
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

	pResume, err := r.UsecaseResume.GetResume(resumeID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pEducations, err := r.UsecaseEducation.GetAllResumeEducation(resumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	pCustomExperience, err := r.UsecaseCustomExperience.GetAllResumeCustomExperience(resumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type RespResume struct {
		Resume           models.Resume                 `json:"resume"`
		Educations       []models.Education            `json:"educations"`
		CustomExperience []models.ExperienceCustomComp `json:"experience_custom_company"`
	}

	ctx.JSON(http.StatusOK,
		RespResume{Resume: *pResume, Educations: pEducations, CustomExperience: pCustomExperience})
}

func (r *ResumeHandler) handlerGetResumeList(ctx *gin.Context) {
	var reqResume struct {
		Start uint `form:"start"`
		Limit uint `form:"limit" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&reqResume); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	rList, err := r.UsecaseResume.GetResumePage(reqResume.Start, reqResume.Limit)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Resume []models.Resume `json:"resume"`
	}

	ctx.JSON(http.StatusOK, Resp{Resume: rList})
}

func (r *ResumeHandler) handlerUpdateResume(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userIDStr := session.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var reqResume models.Resume
	if err := ctx.ShouldBindBodyWith(&reqResume, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if userID != reqResume.UserID {
		errMsg := "this user has not update this resume"
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: errMsg})
	}

	reqResume.DateCreate = time.Now()

	pResume, err := r.UsecaseResume.UpdateResume(reqResume)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var additionParam models.AdditionInResume
	if err := ctx.ShouldBindBodyWith(&additionParam, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	pEducations, err := r.handlerUpdateEducation(additionParam.Education, userID, pResume.ID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pCustomExperience, err := r.handlerUpdateCustomExperience(additionParam.CustomExperience, userID, pResume.ID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	type RespResume struct {
		Resume           models.Resume                 `json:"resume"`
		Educations       []models.Education            `json:"education"`
		CustomExperience []models.ExperienceCustomComp `json:"custom_experience"`
	}

	ctx.JSON(http.StatusOK,
		RespResume{Resume: *pResume, Educations: pEducations, CustomExperience: pCustomExperience})

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

	pEducations, err := r.UsecaseEducation.UpdateEducation(educations, resumeID)
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

	pEducations, err := r.UsecaseCustomExperience.UpdateCustomExperience(experience, resumeID)
	if err != nil {
		return nil, err
	}
	return pEducations, nil
}
