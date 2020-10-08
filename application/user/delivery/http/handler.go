package http

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/common"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const IMAGE_DIR = ""

type UserHandler struct {
	UserUseCase user.IUseCaseUser
}

func NewRest(router *gin.RouterGroup, useCase user.IUseCaseUser, authMiddleware *jwt.GinJWTMiddleware) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router, authMiddleware)
	return rest
}

func (U *UserHandler) routes(router *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	router.GET("/by/id/:user_id", U.handlerGetUserByID)
	router.POST("/add", U.handlerCreateUser)
	router.PUT("/update/:user_id", U.handlerUpdateUser)
	router.Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/me", U.handlerGetCurrentUser)
		router.PUT("/update", U.handlerUpdateUser)
	}
}

func (U *UserHandler) handlerGetCurrentUser(ctx *gin.Context) {
	// move to constants
	identityKey := "myid"
	jwtuser, _ := ctx.Get(identityKey)
	userID := jwtuser.(*models.JWTUserData).ID

	userById, err := U.UserUseCase.GetUserByID(userID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		User models.User `json:"user"`
	}

	ctx.JSON(http.StatusOK, Resp{User: userById})
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
	type Resp struct {
		User models.User `json:"user"`
	}

	ctx.JSON(http.StatusOK, Resp{User: user})
}

func (U *UserHandler) handlerCreateUser(ctx *gin.Context) {
	var req struct {
		NickName string `form:"nickname" json:"nickname" binding:"required"`
		Name     string `form:"name" json:"name" binding:"required"`
		Surname  string `form:"surname" json:"surname" binding:"required"`
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
		//Avatar   multipart.FileHeader `form:"img" json:"img" binding:"required"`
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
		//AvatarPath:   uuid.New().String(),
	})

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		User models.User `json:"user"`
	}

	ctx.JSON(http.StatusOK, Resp{User: userNew})
}

func (U *UserHandler) handlerUpdateUser(ctx *gin.Context) {
	var req struct {
		NickName    string `json:"nickname"`
		Name        string `json:"name"`
		Surname     string `json:"surname"`
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	identityKey := "myid"
	jwtuser, _ := ctx.Get(identityKey)
	userID := jwtuser.(*models.JWTUserData).ID
	user, err := U.UserUseCase.UpdateUser(userID, req.NewPassword, req.OldPassword, req.NickName, req.Name, req.Surname, req.Email)
	if err != nil {
		if err == common.ErrInvalidUpdatePassword {
			ctx.AbortWithError(http.StatusForbidden, err)
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
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
