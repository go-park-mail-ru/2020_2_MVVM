package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2020_2_MVVM.git/application/chat/api/usecase"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

// Client represents the websocket client at the server
type Client struct {
	// The actual websocket connection.
	conn *websocket.Conn
}

func newClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

type rest struct {
	usecase usecase.Usecase
}

func NewRest(router *gin.RouterGroup, usecase usecase.Usecase) *rest {
	rest := &rest{usecase: usecase}
	rest.routes(router)
	return rest
}

func (r *rest) routes(router *gin.RouterGroup) {
	router.GET("/ws", r.handlerWS)
}

func (r *rest) handlerWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error is ", err)
		//c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	client := newClient(conn)

	fmt.Println("New Client joined the hub!")
	fmt.Println(client)
	c.Status(200)
}
