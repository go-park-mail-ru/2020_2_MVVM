package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserHandler struct {
	UserUseCase user.IUseCaseUser
}

type Resp struct {
	User *models.User `json:"user"`
}

func NewRest(router *gin.RouterGroup, useCase user.IUseCaseUser, AuthRequired gin.HandlerFunc) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router, AuthRequired)
	return rest
}

func (u *UserHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:user_id", u.handlerGetUserByID)
	router.POST("/", u.handlerCreateUser)
	router.POST("/login", u.handlerLogin)
	router.Use(AuthRequired)
	{
		router.POST("/logout", u.handlerLogout)
		router.GET("/me", u.handlerGetCurrentUser)
		router.PUT("/", u.handlerUpdateUser)
	}
}

func (u *UserHandler) handlerGetCurrentUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")

	userById, err := u.UserUseCase.GetUserByID(userID.(string))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: userById})
}

func (u *UserHandler) handlerGetUserByID(ctx *gin.Context) {
	var req struct {
		UserID string `uri:"user_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetUserByID(req.UserID)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: user})
}

func (u *UserHandler) handlerLogin(ctx *gin.Context) {
	var reqUser models.UserLogin
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		if errMsg := err.Error(); errMsg == "missing Nickname, Password, or Email" {
			ctx.JSON(http.StatusConflict, common.RespError{Err: errMsg})
		} else {
			ctx.AbortWithError(http.StatusForbidden, err)
		}
		return
	}

	user, err := u.UserUseCase.Login(reqUser)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session := sessions.Default(ctx)
	if user.UserType == "candidate" {
		cand, err := u.UserUseCase.GetCandidateByID(user.ID.String())
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		session.Set("cand_id", cand.ID.String())
		session.Set("empl_id", nil)

	} else if user.UserType == "employer" {
		empl, err := u.UserUseCase.GetEmployerByID(user.ID.String())
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		session.Set("empl_id", empl.ID.String())
		session.Set("cand_id", nil)
	} else {
		errMsg := "cannot login, undefined user type"
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: errMsg})
	}

	session.Set("user_id", user.ID.String())
	err = session.Save()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: user})

}

func (u *UserHandler) handlerLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	err := session.Save()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (u *UserHandler) handlerCreateUser(ctx *gin.Context) {
	var req struct {
		UserType      string `json:"user_type" binding:"required"`
		NickName      string `json:"nickname" binding:"required"`
		Password      string `json:"password" binding:"required"`
		Name          string `json:"name" binding:"required"`
		Surname       string `json:"surname" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Phone         string `json:"phone"`
		SocialNetwork string `json:"social_network"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	userNew, err := u.UserUseCase.CreateUser(models.User{
		UserType:      req.UserType,
		Nickname:      req.NickName,
		Name:          req.Name,
		Surname:       req.Surname,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		Phone:         &req.Phone,
		SocialNetwork: &req.SocialNetwork,
	})
	if err != nil {
		if errMsg := err.Error(); errMsg == "user already exists" {
			ctx.JSON(http.StatusConflict, common.RespError{Err: errMsg})
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: userNew})
}

func (u *UserHandler) handlerUpdateUser(ctx *gin.Context) {
	var req struct {
		NickName      string `json:"nickname"`
		Name          string `json:"name"`
		Surname       string `json:"surname"`
		Email         string `json:"email"`
		NewPassword   string `json:"new_password"`
		OldPassword   string `json:"old_password"`
		Phone         string `json:"phone"`
		SocialNetwork string `json:"social_network"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	userUpdate, err := u.UserUseCase.UpdateUser(userID.(string), req.NewPassword, req.OldPassword, req.NickName, req.Name,
		req.Surname, req.Email, req.Phone, req.SocialNetwork)
	if err != nil {
		if err == common.ErrInvalidUpdatePassword {
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}
	}
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, userUpdate)
}
