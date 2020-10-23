package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const IMAGE_DIR = "static"

type UserHandler struct {
	UserUseCase user.IUseCaseUser
}

type Resp struct {
	User models.User `json:"user"`
}

type RespError struct {
	Err string `json:"error"`
}

func NewRest(router *gin.RouterGroup, useCase user.IUseCaseUser, AuthRequired gin.HandlerFunc) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router, AuthRequired)
	return rest
}

func (U *UserHandler) routes(router *gin.RouterGroup, AuthRequired gin.HandlerFunc) {
	router.GET("/by/id/:user_id", U.handlerGetUserByID)
	router.POST("/add", U.handlerCreateUser)
	router.POST("/login", U.handlerLogin)
	router.Use(AuthRequired)
	{
		router.POST("/logout", U.handlerLogout)
		router.GET("/me", U.handlerGetCurrentUser)
		router.PUT("/update", U.handlerUpdateUser)
	}
}



func (U *UserHandler) handlerGetCurrentUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userID := session.Get("user_id")

	//identityKey := "myid"
	//jwtuser, _ := ctx.Get(identityKey)
	//userID := jwtuser.(*models.JWTUserData).ID

	userById, err := U.UserUseCase.GetUserByID(userID.(string))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: *userById})
}

func (U *UserHandler) handlerGetUserByID(ctx *gin.Context) {
	var req struct {
		UserID string `uri:"user_id" binding:"required,uuid"`
	}

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := U.UserUseCase.GetUserByID(req.UserID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: *user})
}

func (U *UserHandler) handlerLogin(ctx *gin.Context) {
	var reqUser models.UserLogin
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
		//return "", errors.New("missing Username, Password, or Email") // make error constant
	}

	user, err := U.UserUseCase.Login(reqUser)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session := sessions.Default(ctx)
	session.Set("user_id", user.ID.String())
	err = session.Save()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: *user})

}

func (U *UserHandler) handlerLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	err := session.Save()
	if err != nil {
		ctx.AbortWithError(http.StatusForbidden, err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (U *UserHandler) handlerCreateUser(ctx *gin.Context) {
	var req struct {
		NickName string `form:"nickname" json:"nickname" binding:"required"`
		Name     string `form:"name" json:"name" binding:"required"`
		Surname  string `form:"surname" json:"surname" binding:"required"`
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
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
	userNew, err := U.UserUseCase.CreateUser(models.User{
		Nickname:     req.NickName,
		Name:         req.Name,
		Surname:      req.Surname,
		Email:        req.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errMsg := err.Error(); errMsg == "user already exists" {
			ctx.JSON(http.StatusConflict, RespError{Err: errMsg})
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: *userNew})
}

func (U *UserHandler) handlerUpdateUser(ctx *gin.Context) {
	var req struct {
		NickName      string   `form:"nickname" json:"nickname"`
		Name          string   `form:"name" json:"name"`
		Surname       string   `form:"surname" json:"surname"`
		Email         string   `form:"email" json:"email"`
		NewPassword   string   `form:"new_password" json:"new_password"`
		OldPassword   string   `form:"old_password" json:"old_password"`
		Phone         string   `form:"phone" json:"phone"`
		AreaSearch    string `form:"area_search" json:"area_search"`
		SocialNetwork []string `form:"social_network" json:"social_network"`
		Avatar        string   `form:"avatar" json:"avatar"`
		//Avatar   multipart.FileHeader `form:"img" json:"img"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	session := sessions.Default(ctx)
	userID := session.Get("user_id")

	//identityKey := "myid"
	//jwtuser, _ := ctx.Get(identityKey)
	//userID := jwtuser.(*models.JWTUserData).ID
	userUpdate, err := U.UserUseCase.UpdateUser(userID.(uuid.UUID), req.NewPassword, req.OldPassword, req.NickName, req.Name,
												req.Surname, req.Email, req.Phone, req.AreaSearch, req.SocialNetwork)
	if err != nil {
		if err == common.ErrInvalidUpdatePassword {
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}
	}
	if err != nil {
		if errMsg := err.Error(); errMsg == "user already exists" {
			ctx.JSON(http.StatusOK, RespError{Err: errMsg})
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	/*img, err := req.Avatar.Open()
	if err := addOrUpdateUserImage(userNew.ID.String(), img); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, RespResume{User: userNew})*/
	ctx.JSON(http.StatusOK, userUpdate)
}

func addOrUpdateUserImage(imgPath string, data io.Reader) error {
	path := filepath.Join(IMAGE_DIR, imgPath)

	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, data); err != nil {
		return err
	}
	return nil
}
