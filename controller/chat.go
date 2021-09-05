package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Hub struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	ClientMap map[string]interface{}
}

var (
	rwlock    sync.RWMutex
	clientMap = make(map[string]*Hub, 0)
)

func Chat(ctx *gin.Context) {

	userId := ctx.Query("userId")
	upgrade := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	connect, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		zap.Error(err)
		return
	}

	hub := &Hub{
		Conn:      connect,
		DataQueue: make(chan []byte, 100),
		ClientMap: make(map[string]interface{}),
	}
	rwlock.RLock()
	clientMap[userId] = hub
	rwlock.RUnlock()

	go sendProcess(hub)
	go receiveProcess(hub)
	sendMessage(userId, []byte("hello, darling "+userId))

}

func sendMessage(userId string, message []byte) {
	rwlock.RLock()
	hub, ok := clientMap[userId]
	rwlock.RUnlock()
	if ok {
		hub.DataQueue <- message
	}
}

func sendProcess(hub *Hub) {
	for {
		select {
		case data := <-hub.DataQueue:
			if err := hub.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				zap.Error(err)
				return
			}
		}
	}
}

func receiveProcess(hub *Hub) {
	for {
		_, data, err := hub.Conn.ReadMessage()
		if err != nil {
			zap.Error(err)
			return
		}
		sendMessage("2", data)
		fmt.Println("========================", string(data))
	}

}
