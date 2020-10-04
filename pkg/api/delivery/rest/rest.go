package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api/usecase"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/models"
	"github.com/google/uuid"
	"net/http"
)

type rest struct {
	Usecase usecase.Usecase
}

func NewRest(router *gin.RouterGroup, usecase usecase.Usecase) *rest {
	rest := &rest{Usecase: usecase}
	rest.routes(router)
	return rest
}

func (r *rest) routes(router *gin.RouterGroup) {
	router.GET("/nothing", r.handlerGetNothing)
	router.GET("/users/:user_id", r.handlerGetUserByID)
	router.POST("/users/add", r.handlerCreateUser)
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
		Name    string `json:"name" binding:"required"`
		Surname string `json:"surname" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := models.User{
		Name:    reqPayload.Name,
		Surname: reqPayload.Surname,
	}

	user, err := r.Usecase.CreateUser(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		User models.User `json:"user"`
	}

	c.JSON(http.StatusOK, Resp{User: user})
}
