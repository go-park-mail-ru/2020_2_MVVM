package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/authmicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/models/models"
	user2 "github.com/go-park-mail-ru/2020_2_MVVM.git/models/user"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
)

type UserHandler struct {
	UserUseCase    user.UseCase
	authClient     authmicro.AuthClient
	cookieConfig   common.AuthCookieConfig
	SessionBuilder common.SessionBuilder
}

func NewRest(router *gin.RouterGroup,
	useCase user.UseCase,
	authClient authmicro.AuthClient,
	authCookieConfig common.AuthCookieConfig,
	sessionBuilder common.SessionBuilder,
	AuthRequired gin.HandlerFunc) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase, cookieConfig: authCookieConfig,
		SessionBuilder: sessionBuilder, authClient: authClient}
	rest.routes(router, AuthRequired)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:user_id", u.GetUserByIdHandler)
	router.GET("cand/by/id/:cand_id", u.GetCandByIdHandler)
	router.GET("empl/by/id/:empl_id", u.GetEmplByIdHandler)
	router.POST("/", u.CreateUserHandler)
	router.POST("/login", u.LoginHandler)
	router.Use(AuthRequired)
	{
		router.POST("/logout", u.LogoutHandler)
		router.GET("/me", u.GetCurrentUserHandler)
		router.PUT("/", u.UpdateUserHandler)
	}
}

func (u *UserHandler) GetCurrentUserHandler(ctx *gin.Context) {
	session := u.SessionBuilder.Build(ctx)
	userID := session.GetUserID()

	userById, err := u.UserUseCase.GetUserByID(userID.String())
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp := models.RespUser{User: userById}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (u *UserHandler) GetUserByIdHandler(ctx *gin.Context) {
	var req struct {
		UserID string `uri:"user_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetUserByID(req.UserID)

	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp := models.RespUser{User: user}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (u *UserHandler) GetCandByIdHandler(ctx *gin.Context) {
	var req struct {
		UserID string `uri:"cand_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetCandByID(req.UserID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp := models.RespUser{User: user}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (u *UserHandler) GetEmplByIdHandler(ctx *gin.Context) {
	var req struct {
		UserID string `uri:"empl_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		//ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetEmplByID(req.UserID)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		//ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp := models.RespUser{User: user}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (u *UserHandler) login(ctx *gin.Context, reqUser models.UserLogin) error {
	session, err := u.authClient.Login(reqUser.Email, reqUser.Password)
	if err != nil {
		return err
	}

	// Save session id to the cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     u.cookieConfig.Key,
		Value:    url.QueryEscape(session.GetSessionID()),
		MaxAge:   u.cookieConfig.MaxAge,
		Path:     u.cookieConfig.Path,
		Domain:   u.cookieConfig.Domain,
		SameSite: u.cookieConfig.SameSite,
		Secure:   u.cookieConfig.Secure,
		HttpOnly: u.cookieConfig.HttpOnly,
	})
	return nil
}

func (u *UserHandler) LoginHandler(ctx *gin.Context) {
	var reqUser models.UserLogin

	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body,  &reqUser); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	if err := common.ReqValidation(&reqUser); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	err := u.login(ctx, reqUser)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(nil, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (u *UserHandler) LogoutHandler(ctx *gin.Context) {
	session := u.SessionBuilder.Build(ctx)
	// clear cookie
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     u.cookieConfig.Key,
		Value:    "",
		MaxAge:   -1,
		Path:     u.cookieConfig.Path,
		Domain:   u.cookieConfig.Domain,
		SameSite: u.cookieConfig.SameSite,
		Secure:   u.cookieConfig.Secure,
		HttpOnly: u.cookieConfig.HttpOnly,
	})
	err := u.authClient.Logout(session.GetSessionID())
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.SessionErr)
		return
	}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(nil, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (u *UserHandler) CreateUserHandler(ctx *gin.Context) {
	var req user2.Register
	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body,  &req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	if err := common.ReqValidation(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}

	var uuidComp *uuid.UUID = nil
	if req.Company != "" {
		uuidTmp, err := uuid.Parse(req.Company)
		uuidComp = &uuidTmp
		if err != nil {
			common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
			return
		}
	}

	userNew, err := u.UserUseCase.CreateUser(models.User{
		UserType:      req.UserType,
		Name:          req.Name,
		Surname:       req.Surname,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		Phone:         &req.Phone,
		SocialNetwork: &req.SocialNetwork,
	}, uuidComp)
	if err != nil {
		if errMsg := err.Error(); errMsg == common.UserExistErr {
			common.WriteErrResponse(ctx, http.StatusConflict, errMsg)
		} else {
			common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		}
		return
	}

	reqUser := models.UserLogin{
		Email:    userNew.Email,
		Password: req.Password,
	}

	err = u.login(ctx, reqUser)
	if err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := models.RespUser{User: userNew}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (u *UserHandler) UpdateUserHandler(ctx *gin.Context) {
	var req user2.Update
	if err := common.UnmarshalFromReaderWithNilCheck(ctx.Request.Body,  &req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, common.EmptyFieldErr)
		return
	}
	if err := common.ReqValidation(&req); err != nil {
		common.WriteErrResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	session := u.SessionBuilder.Build(ctx)
	userIDFromSession := session.GetUserID()
	userID, errSession := uuid.Parse(userIDFromSession.String())
	if errSession != nil {
		common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		return
	}
	userUpdate, err := u.UserUseCase.UpdateUser(models.User{ID: userID, Name: req.Name, Surname: req.Surname,
		Phone: &req.Phone, Email: req.Email, SocialNetwork: &req.SocialNetwork})
	if err != nil {
		if errMsg := err.Error(); errMsg == common.WrongPasswd {
			common.WriteErrResponse(ctx, http.StatusConflict, errMsg)
		} else {
			common.WriteErrResponse(ctx, http.StatusInternalServerError, common.DataBaseErr)
		}
		return
	}
	resp := models.RespUser{User: userUpdate}
	if _, _, err := easyjson.MarshalToHTTPResponseWriter(resp, ctx.Writer); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
}
