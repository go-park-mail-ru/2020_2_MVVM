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
		router.GET("/by/id/:chat_id", r.HandlerGetChatByID)
		router.GET("/list", r.HandlerListChats)
		router.GET("/messenger/:chat_id", r.PollingMessages)
		router.POST("/send", r.CreateMessage)
	}
}

func (r *ChatHandler) PollingMessages(ctx *gin.Context) {
	chat := r.GetUnreadMessages(ctx)
	listChats := r.ListChats(ctx)

	if chat != nil && listChats != nil {
		result := models.Messager{
			ListChatSummary: *listChats,
			ChatHistory:     *chat,
		}
		if _, _, err := easyjson.MarshalToHTTPResponseWriter(result, ctx.Writer); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}


func (r *ChatHandler) GetUnreadMessages(ctx *gin.Context) *models.ChatHistory {
	var request struct {
		ChatID string `uri:"chat_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return nil
	}

	chatID, err := uuid.Parse(request.ChatID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return nil
	}

	var result models.ChatHistory
	session := r.SessionBuilder.Build(ctx)
	if session.GetCandID() != uuid.Nil {
		result, err = r.UseCaseChat.GetUnreadMessages(chatID, common.Candidate)
	} else {
		result, err = r.UseCaseChat.GetUnreadMessages(chatID, common.Employer)
	}

	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return nil
	}

	return &result
}


func (r *ChatHandler) HandlerGetChatByID(ctx *gin.Context) {
	result := r.GetChatByID(ctx)
	if result != nil {
		if _, _, err := easyjson.MarshalToHTTPResponseWriter(result, ctx.Writer); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}

func (r *ChatHandler) HandlerListChats(ctx *gin.Context) {
	result := r.ListChats(ctx)
	if result != nil {
		if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListChatSummary(*result), ctx.Writer); err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}

func (r *ChatHandler) GetChatByID(ctx *gin.Context) *models.ChatHistory {
	var params struct {
		Offset *uint      `json:"offset"`
		Limit  *uint      `json:"limit"`
		From   *time.Time `json:"from"`
		To     *time.Time `json:"to"`
	}

	if err := ctx.ShouldBindJSON(&params); err != nil && err.Error() != "EOF"{
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return nil
	}

	var request struct {
		ChatID string `uri:"chat_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return nil
	}

	chatID, err := uuid.Parse(request.ChatID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return nil
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
		return nil
	}

	return &result
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
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ChatHandler) ListChats(ctx *gin.Context) *[]models.ChatSummary {
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
		return nil
	}
	listChats, err := r.UseCaseChat.ListChats(userID, userType)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return nil
	}

	return &listChats

}

