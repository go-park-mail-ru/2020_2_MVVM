package http

import (
	jwt "github.com/appleboy/gin-jwt/v2"
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
	UsecaseResume    resume.IUseCaseResume
	UsecaseEducation education.IUseCaseEducation
	UsecaseCustomCompany custom_company.IUseCaseCustomCompany
	UsecaseCustomExperience custom_experience.IUseCaseCustomExperience
}

func NewRest(router *gin.RouterGroup,
			usecaseResume resume.IUseCaseResume,
			usecaseEducation education.IUseCaseEducation,
			usecaseCustomCompany custom_company.IUseCaseCustomCompany,
			usecaseCustomExperience custom_experience.IUseCaseCustomExperience,
			authMiddleware *jwt.GinJWTMiddleware) *ResumeHandler {
	rest := &ResumeHandler{
		UsecaseResume: usecaseResume,
		UsecaseEducation: usecaseEducation,
		UsecaseCustomCompany: usecaseCustomCompany,
		UsecaseCustomExperience: usecaseCustomExperience,
	}
	rest.routes(router, authMiddleware)
	return rest
}

func (r *ResumeHandler) routes(router *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	router.GET("/by/id/:resume_id", r.handlerGetResumeByID)
	router.GET("/page", r.handlerGetResumeList)

	router.Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/mine", r.handlerGetAllCurrentUserResume)
		router.POST("/add", r.handlerCreateResume)
	}
	//router.PUT("/resume/update", r.handlerUpdateResume)
}

func (r *ResumeHandler) handlerGetAllCurrentUserResume(ctx *gin.Context) {
	identityKey := "myid"
	jwtuser, _ := ctx.Get(identityKey)
	userID := jwtuser.(*models.JWTUserData).ID

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
	identityKey := "myid"
	jwtuser, _ := ctx.Get(identityKey)
	userID := jwtuser.(*models.JWTUserData).ID

	var testResume models.Resume
	if err := ctx.ShouldBindBodyWith(&testResume, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	testResume.UserID = userID
	testResume.DateCreate = time.Now()

	pResume, err := r.UsecaseResume.CreateResume(testResume)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var testEd models.ReqResume
	if err := ctx.ShouldBindBodyWith(&testEd, binding.JSON); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pEducations, err := r.handlerCreateEducation(testEd.Educations.Education, userID, pResume.ID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	pCustomExperience, err := r.handlerCreateCustomExperience(testEd.CustomExperience.ListReqCustomExperience, userID, pResume.ID)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	type Resp struct {
		Resume           models.Resume                 `json:"resume"`
		Educations       []models.Education            `json:"educations"`
		CustomExperience []models.ExperienceCustomComp `json:"experience_custom_company"`
	}

	ctx.JSON(http.StatusOK, Resp{Resume: *pResume, Educations: pEducations, CustomExperience: pCustomExperience})
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

func (r *ResumeHandler) handlerCreateCustomExperience(experiences []models.ReqCustomExperience, userID, resumeID uuid.UUID) ([]models.ExperienceCustomComp, error) {
	var customCompany []models.CustomCompany
	var customExperiense []models.ExperienceCustomComp

	for _, item := range experiences {
		insertInComp := models.CustomCompany{
			Name:     item.CompanyName,
			Location: item.Location,
			Sphere:   item.Sphere,
		}
		customCompany = append(customCompany, insertInComp)

		insertInExp := models.ExperienceCustomComp{
			CandID:      userID,
			ResumeID:    resumeID,
			Position:    item.Position,
			Begin:       time.Now(),
			Finish:      item.Finish,
			Description: item.Description,
		}
		customExperiense = append(customExperiense, insertInExp)
	}

	var pCustomCompany []models.CustomCompany
	for i := range customCompany {
		pCompany, err := r.UsecaseCustomCompany.CreateCustomCompany(customCompany[i])
		if err != nil {
			return nil, err
		}
		pCustomCompany = append(pCustomCompany, *pCompany)
	}

	for i := range customExperiense {
		customExperiense[i].CompanyID = pCustomCompany[i].ID
	}

	pCustomExperience, err := r.UsecaseCustomExperience.CreateCustomExperience(customExperiense)
	if err != nil {
		return nil, err
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
	pEducation, err := r.UsecaseEducation.GetAllResumeEducation(resumeID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	type Resp struct {
		Resume    models.Resume      `json:"resume"`
		Education []models.Education `json:"education"`
	}

	ctx.JSON(http.StatusOK, Resp{Resume: *pResume, Education: pEducation})
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

//func (r *ResumeHandler) handlerUpdateResume(c *gin.Context) {
//	var reqResume struct {
//		ResumeID 	string `json:"resume_id" binding:"required,uuid"`
//		SalaryMin       int       `json:"salary_min"`
//		SalaryMax       int       `json:"salary_max"`
//		Description     string    `json:"description"`
//		Gender          string    `json:"gender"`
//		Level           string    `json:"level"`
//		ExperienceMonth int       `json:"experience_month"`
//		Education       string    `json:"education"`
//	}
//
//	if err := c.ShouldBindJSON(&reqResume); err != nil {
//		c.AbortWithError(http.StatusBadRequest, err)
//		return
//	}
//
//	resumeID, err := uuid.Parse(reqResume.ResumeID)
//	if err != nil {
//		c.AbortWithError(http.StatusBadRequest, err)
//		return
//	}
//
//	resume := models.Resume{
//		ID: resumeID,
//		SalaryMin: reqResume.SalaryMin,
//		SalaryMax: reqResume.SalaryMax,
//		Description: reqResume.Description,
//		Gender: reqResume.Gender,
//		Level: reqResume.Level,
//		ExperienceMonth: reqResume.ExperienceMonth,
//		Education: reqResume.Education,
//	}
//
//	pResume, err := r.Usecase.UpdateResume(resume)
//	if err != nil {
//		c.AbortWithError(http.StatusInternalServerError, err)
//		return
//	}
//
//	type Resp struct {
//		Resume models.Resume `json:"resume"`
//	}
//
//	c.JSON(http.StatusOK, Resp{Resume: *pResume})
//
//}
