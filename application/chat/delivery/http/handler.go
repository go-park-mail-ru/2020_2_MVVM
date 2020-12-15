package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/chat"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
)

type ChatHandler struct {
	UsecaseChat     chat.IUseCaseChat
	SessionBuilder  common.SessionBuilder
}

func NewRest(router *gin.RouterGroup,
	usecaseChat chat.IUseCaseChat,
	sessionBuilder common.SessionBuilder,
	AuthRequired gin.HandlerFunc) *ChatHandler {
	rest := &ChatHandler{
		UsecaseChat:     usecaseChat,
		SessionBuilder:  sessionBuilder,
	}
	rest.routes(router, AuthRequired)
	return rest
}

func (r *ChatHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.Use(AuthRequired)
	{
		router.GET("/by/id/:chat_id", r.GetChatByID)
		//router.GET("/list", r.ListChats)
		router.POST("/send", r.CreateMessage)
	}
}

func (r *ChatHandler) GetChatByID(ctx *gin.Context) {
	var params struct {
		Start uint `form:"start" default:"0"`
		Limit uint `form:"limit" default:"20"`
	}

	if err := ctx.ShouldBindQuery(&params); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}

	var request struct {
		ChatID string `uri:"chat_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&request); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	chatID, _ := uuid.Parse(request.ChatID)

	result, err := r.UsecaseChat.GetByID(chatID, params.Start, params.Limit)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(models.ListMessage(result), ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r *ChatHandler) CreateMessage(ctx *gin.Context)  {
	mes := new(models.Message)

	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body, mes); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	session := r.SessionBuilder.Build(ctx)
	var sender string
	senderID := session.GetUserID()
	candID := session.GetCandID()
	emplID := session.GetEmplID()
	if candID != uuid.Nil && emplID == uuid.Nil {
		sender = "candidate"
	} else if candID == uuid.Nil && emplID != uuid.Nil {
		sender = "employer"
	} else {
		common.WriteErrResponse(ctx, http.StatusMethodNotAllowed, common.AuthRequiredErr)
		//ctx.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}
	mes.Sender = sender
	newMes, err := r.UsecaseChat.CreateMessage(*mes, senderID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, _, err := easyjson.MarshalToHTTPResponseWriter(newMes, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
//
//func (r *ChatHandler) ListChats(ctx *gin.Context) {
//	session := r.SessionBuilder.Build(ctx)
//	var sender string
//	senderID := session.GetUserID()
//	candID := session.GetCandID()
//	emplID := session.GetEmplID()
//	if candID != uuid.Nil && emplID == uuid.Nil {
//		sender = "candidate"
//	} else if candID == uuid.Nil && emplID != uuid.Nil {
//		sender = "employer"
//	} else {
//		common.WriteErrResponse(ctx, http.StatusMethodNotAllowed, common.AuthRequiredErr)
//		//ctx.AbortWithError(http.StatusMethodNotAllowed, err)
//		return
//	}
//}


