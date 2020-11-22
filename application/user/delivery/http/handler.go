package http

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/google/uuid"
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
	router.POST("/login", u.LoginHandler)
	router.POST("/", u.CreateUserHandler)
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
		ctx.JSON(http.StatusForbidden, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err := common.ReqValidation(&reqUser); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}
	user, err := u.UserUseCase.Login(reqUser)
	if err != nil {
		if errMsg := err.Error(); errMsg == common.AuthErr {
			ctx.JSON(http.StatusConflict, common.RespError{Err: errMsg})
		} else {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		}
		return
	}
	session := sessions.Default(ctx)
	if user.UserType == "candidate" {
		cand, err := u.UserUseCase.GetCandidateByID(user.ID.String())
		if err != nil {
			fmt.Println(err)
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
			return
		}
		session.Set("cand_id", cand.ID.String())
		session.Set("empl_id", nil)

	} else if user.UserType == "employer" {
		empl, err := u.UserUseCase.GetEmployerByID(user.ID.String())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
			return
		}
		session.Set("empl_id", empl.ID.String())
		session.Set("cand_id", nil)
	} else {
		ctx.JSON(http.StatusMethodNotAllowed, common.RespError{Err: common.AuthErr})
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
		Password      string `json:"password" binding:"required" valid:"stringlength(5|25)~длина пароля должна быть от 5 до 25 символов."`
		Name          string `json:"name" binding:"required" valid:"utfletter~имя должно содержать только буквы,stringlength(3|25)~длина имени должна быть от 3 до 25 символов."`
		Surname       string `json:"surname" binding:"required" valid:"utfletter~фамилия должна содержать только буквы,stringlength(3|25)~длина фамилии должна быть от 3 до 25 символов."`
		Email         string `json:"email" binding:"required" valid:"email"`
		Phone         string `json:"phone" valid:"numeric~номер телефона должен состоять только из цифр.,stringlength(7|18)~номер телефона от 7 до 18 цифр"`
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
		if errMsg := err.Error(); errMsg == common.UserExistErr {
			ctx.JSON(http.StatusConflict, common.RespError{Err: errMsg})
		} else {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		}
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: userNew})
}

// TODO: fix update on frontend require password for update and deletion
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
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: common.EmptyFieldErr})
		return
	}
	if err := common.ReqValidation(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.RespError{Err: err.Error()})
		return
	}

	session := sessions.Default(ctx).Get("user_id")
	if session == nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.SessionErr})
		return
	}
	userID, errSession := uuid.Parse(session.(string))
	if errSession != nil {
		ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.SessionErr})
		return
	}
	userUpdate, err := u.UserUseCase.UpdateUser(models.User{ID: userID, Name: req.Name, Surname: req.Surname,
		Phone: &req.Phone, Email: req.Email, SocialNetwork: &req.SocialNetwork})
	if err != nil {
		if errMsg := err.Error(); errMsg == common.WrongPasswd {
			ctx.JSON(http.StatusConflict, common.RespError{Err: errMsg})
		} else {
			ctx.JSON(http.StatusInternalServerError, common.RespError{Err: common.DataBaseErr})
		}
		return
	}
	ctx.JSON(http.StatusOK, userUpdate)
}
