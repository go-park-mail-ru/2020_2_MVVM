package http

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
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

func NewRest(router *gin.RouterGroup, useCase user.IUseCaseUser, authMiddleware *jwt.GinJWTMiddleware) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router, authMiddleware)
	return rest
}

func (U *UserHandler) routes(router *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	router.GET("/users/by/id/:user_id", U.handlerGetUserByID)
	router.POST("/users/add", U.handlerCreateUser)
	router.PUT("/users/update/:user_id", U.handlerUpdateUser)
	router.Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/users/me", U.handlerGetCurrentUser)
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

	ctx.JSON(http.StatusOK, Resp{User: userById})
}

func (U *UserHandler) handlerGetUserByID(ctx *gin.Context) {
	var req struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userById, err := U.UserUseCase.GetUserByID(req.UserID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: userById})
}

func (U *UserHandler) handlerCreateUser(ctx *gin.Context) {
	var req struct {
		NickName string               `form:"nickname" json:"nickname" binding:"required"`
		Name     string               `form:"name" json:"name" binding:"required"`
		Surname  string               `form:"surname" json:"surname" binding:"required"`
		Email    string               `form:"email" json:"email" binding:"required"`
		Password string               `form:"password" json:"password" binding:"required"`
		Avatar   multipart.FileHeader `form:"img" json:"img"`
	}

	if err := ctx.ShouldBind(&req); err != nil {
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
			ctx.JSON(http.StatusOK, RespError{Err: errMsg})
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	img, err := req.Avatar.Open()
	if err := addOrUpdateUserImage(userNew.ID.String(), img); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, Resp{User: userNew})
}

func (U *UserHandler) handlerUpdateUser(ctx *gin.Context) {
	var req struct {
		UserId   string               `form:"user_id" json:"user_id" binding:"required"`
		NickName string               `form:"nickname" json:"nickname"`
		Name     string               `form:"name" json:"name"`
		Surname  string               `form:"surname" json:"surname"`
		Email    string               `form:"email" json:"email"`
		Password string               `form:"password" json:"password"`
		Avatar   multipart.FileHeader `form:"img" json:"img"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	uid, _ := uuid.Parse(req.UserId)
	userNew, err := U.UserUseCase.UpdateUser(models.User{
		ID:           uid,
		Nickname:     req.NickName,
		Name:         req.Name,
		Surname:      req.Surname,
		Email:        req.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if errMsg := err.Error(); errMsg == "user already exists" {
			ctx.JSON(http.StatusOK, RespError{Err: errMsg})
		} else {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}
	img, err := req.Avatar.Open()
	if err := addOrUpdateUserImage(userNew.ID.String(), img); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, Resp{User: userNew})
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
