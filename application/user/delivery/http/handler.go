package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/models"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/user"
	"github.com/google/uuid"
	"net/http"
	"os"
)

type UserHandler struct {
	UserUseCase user.IUseCaseUser
}

func NewRest(router *gin.RouterGroup, useCase user.IUseCaseUser) *UserHandler {
	rest := &UserHandler{UserUseCase: useCase}
	rest.routes(router)
	return rest
}

func (U *UserHandler) routes(router *gin.RouterGroup) {
	router.GET("/users/:user_id", U.handlerGetUserByID)
	router.POST("/users/add", U.handlerCreateUser)
	router.PUT("/users/update/:user_id", U.handlerUpdateUser)
}

func (U *UserHandler) handlerGetUserByID(c *gin.Context) {
	var req struct {
		UserID uuid.UUID `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := U.UserUseCase.GetUserByID(req.UserID.String())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		User models.User `json:"user"`
	}

	c.JSON(http.StatusOK, Resp{User: user})
}

func (U *UserHandler) handlerCreateUser(c *gin.Context) {
	var reqPayload struct {
		Name    string `json:"name" binding:"required"`
		Surname string `json:"surname" binding:"required"`
		Avatar  os.File ``
	}
	if err := c.ShouldBindJSON(&reqPayload); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user := models.User{
		Name:    reqPayload.Name,
		Surname: reqPayload.Surname,
	}

	user, err := U.UserUseCase.CreateUser(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	type Resp struct {
		User models.User `json:"user"`
	}

	c.JSON(http.StatusOK, Resp{User: user})
}

func (U *UserHandler) handlerUpdateUser(ctx *gin.Context) {

}
