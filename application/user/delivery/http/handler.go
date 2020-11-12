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
	UserUseCase user.UseCase
}

type Resp struct {
	User *models.User `json:"user"`
}

func NewRest(router *gin.RouterGroup, useCase user.UseCase, AuthRequired gin.HandlerFunc) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
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
	session := sessions.Default(ctx)
	userID := session.Get("user_id")

	userById, err := u.UserUseCase.GetUserByID(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetUserByID(req.UserID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetCandByID(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	user, err := u.UserUseCase.GetEmplByID(req.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, Resp{User: user})
}

func (u *UserHandler) LoginHandler(ctx *gin.Context) {
	var reqUser models.UserLogin
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		if errMsg := err.Error(); errMsg == "missing Email or Password" {
			ctx.JSON(http.StatusConflict, common.RespError{Err: errMsg})
		} else {
			ctx.JSON(http.StatusForbidden, common.RespError{Err: common.EmptyFieldErr})
		}
		return
	}
	if err := common.ReqValidation(&reqUser); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}

	user, err := u.UserUseCase.Login(reqUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	session := sessions.Default(ctx)
	if user.UserType == "candidate" {
		cand, err := u.UserUseCase.GetCandidateByID(user.ID.String())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		session.Set("cand_id", cand.ID.String())
		session.Set("empl_id", nil)

	} else if user.UserType == "employer" {
		empl, err := u.UserUseCase.GetEmployerByID(user.ID.String())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
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
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.SessionErr})
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: user})

}

func (u *UserHandler) LogoutHandler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	err := session.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.SessionErr})
		return
	}
	ctx.Status(http.StatusOK)
}

func (u *UserHandler) CreateUserHandler(ctx *gin.Context) {
	var req struct {
		UserType      string `json:"user_type" binding:"required"`
		Password      string `json:"password" binding:"required" valid:"utfletternum~пароль содержит неразрешенные символы,stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
		Name          string `json:"name" binding:"required" valid:"utfletter~имя должно содержать только буквы,stringlength(1|25)~длина имени должна быть от 1 до 25 символов."`
		Surname       string `json:"surname" binding:"required" valid:"utfletter~фамилия должна содержать только буквы,stringlength(1|25)~длина фамилии должна быть от 1 до 25 символов."`
		Email         string `json:"email" binding:"required" valid:"email"`
		Phone         string `json:"phone" valid:"numeric~номер телефона должен состоять только из цифр.,stringlength(1|18)~номер телефона от 1 до 18 цифр"`
		SocialNetwork string `json:"social_network"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err := common.ReqValidation(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
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
		if errMsg := err.Error(); errMsg == "user already exists" {
			ctx.JSON(http.StatusConflict, common.RespError{Err: errMsg})
		} else {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		}
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: userNew})
}

func (u *UserHandler) UpdateUserHandler(ctx *gin.Context) {
	var req struct {
		Name          string `json:"name" valid:"optional,utfletter~имя должно содержать только буквы,stringlength(3|25)~длина имени должна быть от 3 до 25 символов."`
		Surname       string `json:"surname" valid:"optional,utfletter~фамилия должна содержать только буквы,stringlength(5|25)~длина фамилии должна быть от 3 до 25 символов."`
		Email         string `json:"email" valid:"optional, email"`
		NewPassword   string `json:"new_password" valid:"optional, utfletternum~пароль содержит неразрешенные символы,stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
		OldPassword   string `json:"old_password" valid:"optional, utfletternum~пароль содержит неразрешенные символы,stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
		Phone         string `json:"phone" valid:"optional, numeric~номер телефона должен состоять только из цифр.,stringlength(4|18)~номер телефона от 4 до 18 цифр"`
		SocialNetwork string `json:"social_network" valid:"-"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}

	session := sessions.Default(ctx)
	userID := session.Get("user_id")
	userUpdate, err := u.UserUseCase.UpdateUser(userID.(string), req.NewPassword, req.OldPassword, req.Name,
		req.Surname, req.Email, req.Phone, req.SocialNetwork)
	if err != nil {
		if err == common.ErrInvalidUpdatePassword {
			ctx.JSON(http.StatusForbidden, err)
			return
		}
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		return
	}
	ctx.JSON(http.StatusOK, userUpdate)
}
