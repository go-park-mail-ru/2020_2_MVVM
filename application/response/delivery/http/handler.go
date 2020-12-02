package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/response"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/dto/models"
	response2 "github.com/go-park-mail-ru/2020_2_MVVM.git/dto/response"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
)

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
	response := new(models.Response)

	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body,  response); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := r.SessionBuilder.Build(ctx)
	var userType string
	candID := session.GetCandID()
	emplID := session.GetEmplID()
	if candID != uuid.Nil && emplID == uuid.Nil {
		userType = "candidate"
	} else if candID == uuid.Nil && emplID != uuid.Nil {
		userType = "employer"
	} else {
		common.WriteErrResponse(ctx, http.StatusMethodNotAllowed, common.AuthRequiredErr)
		//ctx.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}

	response.Initial = userType
	pResponse, err := r.UsecaseResponse.Create(*response)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(pResponse, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResponseHandler) UpdateStatus(ctx *gin.Context) {
	response := new(models.Response)
	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body,  response); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := r.SessionBuilder.Build(ctx)
	var userType string
	candID := session.GetCandID()
	emplID := session.GetEmplID()
	if candID != uuid.Nil && emplID == uuid.Nil {
		userType = common.Candidate
	} else if candID == uuid.Nil && emplID != uuid.Nil {
		userType = common.Employer
	} else {
		common.WriteErrResponse(ctx, http.StatusMethodNotAllowed, common.AuthRequiredErr)
		//ctx.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}

	response.Initial = userType
	pResponse, err := r.UsecaseResponse.UpdateStatus(*response, userType)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(pResponse, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResponseHandler) handlerGetAllResponses(ctx *gin.Context) {
	session := r.SessionBuilder.Build(ctx)
	emplID := session.GetEmplID()
	candID := session.GetCandID()

	var responses []models.ResponseWithTitle
	var err error
	if candID != uuid.Nil && emplID == uuid.Nil {
		responses, err = r.UsecaseResponse.GetAllCandidateResponses(candID, nil)
		if err != nil {
			common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
			//ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else if candID == uuid.Nil && emplID != uuid.Nil {
		responses, err = r.UsecaseResponse.GetAllEmployerResponses(emplID, nil)
		if err != nil {
			common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
			//ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	} else {
		err := errors.New("this user cannot have responses")
		common.WriteErrResponse(ctx, http.StatusMethodNotAllowed, err.Error())
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListResponseWithTitle(responses), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResponseHandler) handlerGetAllResumeWithoutResponse(ctx *gin.Context) {
	candID, vacancyID, err := r.handlerGetAllEntityWithoutResponse(ctx, common.CandID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	resumes, err := r.UsecaseResponse.GetAllResumeWithoutResponse(candID, vacancyID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListBriefResumeInfo(resumes), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ResponseHandler) handlerGetAllVacancyWithoutResponse(ctx *gin.Context) {
	emplID, resumeID, err := r.handlerGetAllEntityWithoutResponse(ctx, common.EmplID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	vacancies, err := r.UsecaseResponse.GetAllVacancyWithoutResponse(emplID, resumeID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListVacancy(vacancies), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
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
	var userID uuid.UUID
	if userType == common.CandID {
		userID = session.GetCandID()
	} else if userType == common.EmplID {
		userID = session.GetEmplID()
	}
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return userID, entityID, nil
}

func getNewResponses(r *ResponseHandler, unId uuid.UUID, userType string, respIds []uuid.UUID) ([]models.ResponseWithTitle, int, error) {
	var (
		responses []models.ResponseWithTitle
		err       error
	)
	if userType == common.Candidate {
		responses, err = r.UsecaseResponse.GetAllCandidateResponses(unId, respIds)
		if err != nil {
			return nil, http.StatusInternalServerError, errors.New(common.DataBaseErr)
		}
	} else if userType == common.Employer {
		responses, err = r.UsecaseResponse.GetAllEmployerResponses(unId, respIds)
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
		notifications response2.RespNotifications
		err           error
		status        int
		daysFromNow   int
		req           response2.ReqNotify
	)

	session := r.SessionBuilder.Build(ctx)
	var unId uuid.UUID
	var userType string

	candID := session.GetCandID()
	emplID := session.GetEmplID()
	if candID != uuid.Nil {
		unId, userType = candID, common.CandID
	} else {
		unId, userType = emplID, common.EmplID
	}

	if err != nil {
		common.WriteErrResponse(ctx, http.StatusMethodNotAllowed, err.Error())
		return
	}
	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body,  &req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	if !req.OnlyVacCnt && req.ListEnd <= req.ListStart {
		common.WriteErrResponse(ctx, http.StatusBadRequest, "invalid 'vac_list_start' and 'vac_list_limit' params")
		return
	}
	if req.VacInLastNDays != nil {
		if daysFromNow = *req.VacInLastNDays; daysFromNow == 0 {
			daysFromNow = common.Week
		}
	}
	if req.OnlyRespCnt {
		notifications.UnreadRespCnt, err = r.UsecaseResponse.GetResponsesCnt(unId, userType)
	} else {
		notifications.UnreadResp, status, err = getNewResponses(r, unId, userType, req.NewRespNotifications)
		notifications.UnreadRespCnt = uint(len(notifications.UnreadResp))
	}
	unId = session.GetUserID()
	if req.OnlyVacCnt && req.VacInLastNDays != nil {
		if daysFromNow == 0 {
			daysFromNow = common.Week
		}
		notifications.RecommendedVacCnt, err = r.UsecaseResponse.GetRecommendedVacCnt(unId, daysFromNow)
	} else if daysFromNow > 0 {
		notifications.RecommendedVac, err = r.UsecaseResponse.GetRecommendedVacancies(unId, req.ListStart, req.ListEnd, daysFromNow)
		notifications.RecommendedVacCnt = uint(len(notifications.RecommendedVac))
	}
	if err != nil && err.Error() != common.NoRecommendation {
		ctx.JSON(status, models.RespError{Err: common.DataBaseErr})
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(&notifications, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
