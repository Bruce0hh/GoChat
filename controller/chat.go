package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

func Chat(ctx *gin.Context) {

	upgrade := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	if _, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		zap.Error(err)
	}

}
