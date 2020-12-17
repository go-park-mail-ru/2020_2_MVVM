package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/chat"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
	"time"
)

type ChatHandler struct {
	UseCaseChat    chat.IUseCaseChat
	SessionBuilder common.SessionBuilder
}

func NewRest(router *gin.RouterGroup,
	usecaseChat chat.IUseCaseChat,
	sessionBuilder common.SessionBuilder,
	AuthRequired gin.HandlerFunc) *ChatHandler {
	rest := &ChatHandler{
		UseCaseChat:    usecaseChat,
		SessionBuilder: sessionBuilder,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (r *ChatHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.Use(AuthRequired)
	{
		router.GET("/by/id/:chat_id", r.GetChatByID)
		router.GET("/list", r.ListChats)
		router.POST("/send", r.CreateMessage)
	}
}

func (r *ChatHandler) GetChatByID(ctx *gin.Context) {
	var params struct {
		Offset *uint      `json:"offset"`
		Limit  *uint      `json:"limit"`
		From   *time.Time `json:"from"`
		To     *time.Time `json:"to"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}

	var request struct {
		ChatID string `uri:"chat_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}

	chatID, err := uuid.Parse(request.ChatID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}

	var result models.ChatHistory
	session := r.SessionBuilder.Build(ctx)
	if session.GetCandID() != uuid.Nil {
		result, err = r.UseCaseChat.GetChatHistory(chatID, common.Candidate, params.From, params.To, params.Offset, params.Limit)
	} else {
		result, err = r.UseCaseChat.GetChatHistory(chatID, common.Employer, params.From, params.To, params.Offset, params.Limit)
	}

	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(result, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ChatHandler) CreateMessage(ctx *gin.Context) {
	mes := new(models.Message)

	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body, mes); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}

	session := r.SessionBuilder.Build(ctx)
	senderID := session.GetUserID()
	candID := session.GetCandID()
	emplID := session.GetEmplID()

	if candID != uuid.Nil && emplID == uuid.Nil {
		mes.Sender = common.Candidate
	} else if candID == uuid.Nil && emplID != uuid.Nil {
		mes.Sender = common.Employer
	} else {
		common.WriteErrResponse(ctx, http.StatusMethodNotAllowed, common.AuthRequiredErr)
		return
	}

	newMes, err := r.UseCaseChat.CreateMessage(*mes, senderID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(newMes, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ChatHandler) ListChats(ctx *gin.Context) {
	session := r.SessionBuilder.Build(ctx)
	var userType string
	userID := session.GetUserID()
	candID := session.GetCandID()
	emplID := session.GetEmplID()
	if candID != uuid.Nil && emplID == uuid.Nil {
		userType = common.Candidate
	} else if candID == uuid.Nil && emplID != uuid.Nil {
		userType = common.Employer
	} else {
		common.WriteErrResponse(ctx, http.StatusMethodNotAllowed, common.AuthRequiredErr)
		return
	}
	listChats, err := r.UseCaseChat.ListChats(userID, userType)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListChatSummary(listChats), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}

}
