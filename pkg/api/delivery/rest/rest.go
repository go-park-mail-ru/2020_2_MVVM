package rest

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api/usecase"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type rest struct {
	Usecase usecase.Usecase
}

func NewRest(router *gin.RouterGroup, usecase usecase.Usecase, authMiddleware *jwt.GinJWTMiddleware) *rest {
	rest := &rest{Usecase: usecase}
	rest.routes(router, authMiddleware)
	return rest
}

func (r *rest) routes(router *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	router.GET("/nothing", r.handlerGetNothing)
	router.GET("/users/by/id/:user_id", r.handlerGetUserByID)
	router.POST("/users/add", r.handlerCreateUser)

	router.Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/users/me", r.handlerGetCurrentUser)
	}
}

func (r *rest) handlerGetNothing(c *gin.Context) {
	err := r.Usecase.DoNothing()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		Status string
	}

	c.JSON(http.StatusOK, Resp{Status: "ok"})
}

func (r *rest) handlerGetCurrentUser(c *gin.Context) {
	// move to constants
	identityKey := "myid"
	jwtuser, _ := c.Get(identityKey)
	userID := jwtuser.(*models.JWTUserData).ID

	user, err := r.Usecase.GetUserByID(userID.String())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		User models.User `json:"user"`
	}

	c.JSON(http.StatusOK, Resp{User: user})
}

func (r *rest) handlerGetUserByID(c *gin.Context) {
	var req struct {
		UserID string `uri:"user_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := r.Usecase.GetUserByID(userID.String())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		User models.User `json:"user"`
	}

	c.JSON(http.StatusOK, Resp{User: user})
}

func (r *rest) handlerCreateUser(c *gin.Context) {
	var reqPayload struct {
		Nickname string `form:"nickname" json:"nickname" binding:"required"`
		Name     string `form:"name" json:"name" binding:"required"`
		Surname  string `form:"surname" json:"surname" binding:"required"`
		Email    string `form:"email" json:"email" binding:"required"`
		Pasword  string `form:"password" json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&reqPayload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(reqPayload.Pasword), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	user := models.User{
		Nickname:     reqPayload.Nickname,
		Name:         reqPayload.Name,
		Surname:      reqPayload.Surname,
		Email:        reqPayload.Email,
		PasswordHash: passwordHash,
	}

	user, err = r.Usecase.CreateUser(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		User models.User `json:"user"`
	}

	c.JSON(http.StatusOK, Resp{User: user})
}
