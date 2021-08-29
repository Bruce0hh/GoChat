package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func chat(ctx *gin.Context) {
	ws := websocket.Upgrader{}
	ws.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	upgrade, _ := ws.Upgrade(ctx.Writer, ctx.Request, nil)
	defer func(upgrade *websocket.Conn) {
		err := upgrade.Close()
		if err != nil {

		}
	}(upgrade)
}
