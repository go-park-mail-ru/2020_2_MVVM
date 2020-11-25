package http

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/google/uuid"
	"net/http"
)

type RespNotifications struct {
	UnreadResponses      []models.ResponseWithTitle `json:"unread_responses"`
	RecommendedVacancies []models.Vacancy           `json:"new_vacancies"`
}

type ResponseHandler struct {
	UsecaseResponse response.IUseCaseResponse
	SessionBuilder  common.SessionBuilder
}

func NewRest(router *gin.RouterGroup,
	usecaseResponse response.IUseCaseResponse,
	sessionBuilder common.SessionBuilder,
	AuthRequired gin.HandlerFunc) *ResponseHandler {
	rest := &ResponseHandler{
		UsecaseResponse: usecaseResponse,
		SessionBuilder:  sessionBuilder,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (r *ResponseHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.Use(AuthRequired)
	{
		router.POST("/notify", r.handlerGetAllNotifications)
		router.POST("/", r.CreateResponse)
		router.POST("/update", r.UpdateStatus)
		router.GET("/my", r.handlerGetAllResponses)
		router.GET("/free/resumes/:entity_id", r.handlerGetAllResumeWithoutResponse)    // vacancy_id
		router.GET("/free/vacancies/:entity_id", r.handlerGetAllVacancyWithoutResponse) // resume_id
	}
}

func (r *ResponseHandler) CreateResponse(ctx *gin.Context) {
	var response models.Response
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := r.SessionBuilder.Build(ctx)
	var userType string
	candIDStr := session.Get(common.CandID)
	emplIDStr := session.Get(common.EmplID)
	if candIDStr != nil && emplIDStr == nil {
		userType = "candidate"
	} else if candIDStr == nil && emplIDStr != nil {
		userType = "employer"
	} else {
		err := errors.New("this user cannot respond")
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: common.AuthRequiredErr})
		ctx.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}

	response.Initial = userType
	pResponse, err := r.UsecaseResponse.Create(response)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &pResponse)
}

func (r *ResponseHandler) UpdateStatus(ctx *gin.Context) {
	var response models.Response
	if err := ctx.ShouldBindJSON(&response); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := r.SessionBuilder.Build(ctx)
	var userType string
	candIDStr := session.Get(common.CandID)
	emplIDStr := session.Get(common.EmplID)
	if candIDStr != nil && emplIDStr == nil {
		userType = common.Candidate
	} else if candIDStr == nil && emplIDStr != nil {
		userType = common.Employer
	} else {
		err := errors.New("this user cannot respond")
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: common.AuthRequiredErr})
		ctx.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}

	response.Initial = userType

	pResponse, err := r.UsecaseResponse.UpdateStatus(response, userType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, &pResponse)
}

func (r *ResponseHandler) handlerGetAllResponses(ctx *gin.Context) {
	session := r.SessionBuilder.Build(ctx)
	emplID, err := common.GetCurrentUserId(session, common.EmplID)
	candID, err := common.GetCurrentUserId(session, common.CandID)

	var responses []models.ResponseWithTitle
	if candID != uuid.Nil && emplID == uuid.Nil {
		responses, err = r.UsecaseResponse.GetAllCandidateResponses(candID, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else if candID == uuid.Nil && emplID != uuid.Nil {
		responses, err = r.UsecaseResponse.GetAllEmployerResponses(emplID, nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		err := errors.New("this user cannot have responses")
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, responses)
}

func (r *ResponseHandler) handlerGetAllResumeWithoutResponse(ctx *gin.Context) {
	candID, vacancyID, err := r.handlerGetAllEntityWithoutResponse(ctx, common.CandID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resumes, err := r.UsecaseResponse.GetAllResumeWithoutResponse(candID, vacancyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, resumes)
}

func (r *ResponseHandler) handlerGetAllVacancyWithoutResponse(ctx *gin.Context) {
	emplID, resumeID, err := r.handlerGetAllEntityWithoutResponse(ctx, common.EmplID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	vacancies, err := r.UsecaseResponse.GetAllVacancyWithoutResponse(emplID, resumeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, vacancies)
}

func (r *ResponseHandler) handlerGetAllEntityWithoutResponse(ctx *gin.Context, userType string) (uuid.UUID, uuid.UUID, error) {
	var req struct {
		EntityID string `uri:"entity_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindUri(&req); err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	session := r.SessionBuilder.Build(ctx)
	userID, err := common.GetCurrentUserId(session, userType)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return userID, entityID, nil
}

func getResponses(r *ResponseHandler, session sessions.Session, respIds map[uuid.UUID]bool) ([]models.ResponseWithTitle, int, error) {
	emplID, err := common.GetCurrentUserId(session, common.EmplID)
	candID, err := common.GetCurrentUserId(session, common.CandID)

	var responses []models.ResponseWithTitle
	if candID != uuid.Nil && emplID == uuid.Nil {
		responses, err = r.UsecaseResponse.GetAllCandidateResponses(candID, respIds)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New(common.DataBaseErr)
		}
	} else if candID == uuid.Nil && emplID != uuid.Nil {
		responses, err = r.UsecaseResponse.GetAllEmployerResponses(emplID, respIds)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New(common.DataBaseErr)
		}
	} else {
		return nil, http.StatusMethodNotAllowed, errors.New("this user cannot have responses")
	}
	return responses, http.StatusOK, nil
}

func (r *ResponseHandler) handlerGetAllNotifications(ctx *gin.Context) {
	var (
		responses     []models.ResponseWithTitle
		err           error
		status        int
		req           struct {
			NewVacNotifications  map[uuid.UUID]bool     `json:"new_vacancies"` // notifications about new vacancies, nil means all
			NewRespNotifications map[uuid.UUID]bool     `json:"unread_responses"`
			Chat                 map[uuid.UUID][]string `json:"unread_messages"`
		}
	)

	session := r.SessionBuilder.Build(ctx)

	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	/*if req.NewVacNotifications != nil {
		for key, isActive := range req.NewVacNotifications {
			if !isActive {
				delete(req.NewVacNotifications, key)
			}
		}
	} else {
		// достаем все вакансии или резюме, подобранные для пользователя
	}*/
	if req.NewRespNotifications != nil {
		responses, status, err = getResponses(r, session, req.NewRespNotifications)
	} else {
		responses, status, err = getResponses(r, session, nil)
	}
	if err != nil {
		ctx.JSON(status, common.RespError{Err: err.Error()})
	}
	ctx.JSON(http.StatusOK, RespNotifications{UnreadResponses: responses})
}
