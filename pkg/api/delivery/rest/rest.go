package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/pkg/api/usecase"
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
