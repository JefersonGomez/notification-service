package controller

import (
	"net/http"
	"notification-service/worker"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var HubGlobal *worker.Hub

func WsHandler(c *gin.Context) {

	usuarioID := c.GetUint("id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		return
	}

	defer HubGlobal.Eliminar(usuarioID)

	HubGlobal.Registrar(usuarioID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

}
