package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/microservices/auth/authmicro"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/google/uuid"
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

type Resp struct {
	User *models.User `json:"user"`
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
		ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: userById})
}

func (u *UserHandler) GetUserByIdHandler(ctx *gin.Context) {
	var req struct {
		UserID string `uri:"user_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetUserByID(req.UserID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: user})
}

func (u *UserHandler) GetCandByIdHandler(ctx *gin.Context) {
	var req struct {
		UserID string `uri:"cand_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetCandByID(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, Resp{User: user})
}

func (u *UserHandler) GetEmplByIdHandler(ctx *gin.Context) {
	var req struct {
		UserID string `uri:"empl_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetEmplByID(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, Resp{User: user})
}

func (u *UserHandler) login(ctx *gin.Context, reqUser models.UserLogin) {
	session, err := u.authClient.Login(reqUser.Email, reqUser.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: err.Error()})
		return
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
}

func (u *UserHandler) LoginHandler(ctx *gin.Context) {
	var reqUser models.UserLogin

	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err := common.ReqValidation(&reqUser); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: err.Error()})
		return
	}
	u.login(ctx, reqUser)
	ctx.JSON(http.StatusOK, nil)
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
		ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.SessionErr})
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func (u *UserHandler) CreateUserHandler(ctx *gin.Context) {
	var req struct {
		UserType      string `json:"user_type" binding:"required"`
		Password      string `json:"password" binding:"required" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
		Name          string `json:"name" binding:"required" valid:"utfletter~имя должно содержать только буквы,stringlength(3|25)~длина имени должна быть от 3 до 25 символов."`
		Surname       string `json:"surname" binding:"required" valid:"utfletter~фамилия должна содержать только буквы,stringlength(3|25)~длина фамилии должна быть от 3 до 25 символов."`
		Email         string `json:"email" binding:"required" valid:"email"`
		Phone         string `json:"phone" valid:"numeric~номер телефона должен состоять только из цифр.,stringlength(7|18)~номер телефона от 7 до 18 цифр"`
		SocialNetwork string `json:"social_network"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err := common.ReqValidation(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: err.Error()})
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.DataBaseErr})
		return
	}
	userNew, err := u.UserUseCase.CreateUser(models.User{
		UserType:      req.UserType,
		Name:          req.Name,
		Surname:       req.Surname,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		Phone:         &req.Phone,
		SocialNetwork: &req.SocialNetwork,
	})
	if err != nil {
		if errMsg := err.Error(); errMsg == common.UserExistErr {
			ctx.JSON(http.StatusConflict, models.RespError{Err: errMsg})
		} else {
			ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.DataBaseErr})
		}
		return
	}

	reqUser := models.UserLogin{
		Email:    userNew.Email,
		Password: req.Password,
	}
	u.login(ctx, reqUser)
	ctx.JSON(http.StatusOK, Resp{User: userNew})
}

func (u *UserHandler) UpdateUserHandler(ctx *gin.Context) {
	var req struct {
		Name          string `json:"name" valid:"utfletter~имя должно содержать только буквы,stringlength(3|25)~длина имени должна быть от 3 до 25 символов."`
		Surname       string `json:"surname" valid:"utfletter~фамилия должна содержать только буквы,stringlength(3|25)~длина фамилии должна быть от 3 до 25 символов."`
		Email         string `json:"email" valid:"email"`
		NewPassword   string `json:"new_password" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
		OldPassword   string `json:"old_password" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
		Phone         string `json:"phone" valid:"numeric~номер телефона должен состоять только из цифр.,stringlength(4|18)~номер телефона от 4 до 18 цифр"`
		SocialNetwork string `json:"social_network"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err := common.ReqValidation(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.RespError{Err: err.Error()})
		return
	}

	session := u.SessionBuilder.Build(ctx)
	userIDFromSession := session.GetUserID()
	userID, errSession := uuid.Parse(userIDFromSession.String())
	if errSession != nil {
		ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.SessionErr})
		return
	}
	userUpdate, err := u.UserUseCase.UpdateUser(models.User{ID: userID, Name: req.Name, Surname: req.Surname,
		Phone: &req.Phone, Email: req.Email, SocialNetwork: &req.SocialNetwork})
	if err != nil {
		if errMsg := err.Error(); errMsg == common.WrongPasswd {
			ctx.JSON(http.StatusConflict, models.RespError{Err: errMsg})
		} else {
			ctx.JSON(http.StatusInternalServerError, models.RespError{Err: common.DataBaseErr})
		}
		return
	}
	ctx.JSON(http.StatusOK, Resp{User: userUpdate})
}
