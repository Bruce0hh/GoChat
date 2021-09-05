package websocket

import (
	"GoChat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
)

type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	DataQueue chan []byte
	Lock      sync.RWMutex
}

var (
	WsClient  interface{}
	H         = initHub()
	clientMap = make(map[string]*Client, 0)
	lock      sync.RWMutex
)

func initHub() *Hub {
	WsClient = CreateHubFactory()
	if wh, ok := WsClient.(*Hub); ok {
		go wh.Run()
	}
	return WsClient.(*Hub)
}

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

	client := &Client{
		Hub:       H,
		Conn:      connect,
		DataQueue: make(chan []byte, 100),
	}
	lock.RLock()
	H.Login <- client
	clientMap[userId] = client
	lock.RUnlock()

	go sendProcess(client)
	go receiveProcess(client)
	sendMessage(userId, []byte("hello, darling "+userId))
}

func sendMessage(userId string, message []byte) {
	lock.RLock()
	client, ok := clientMap[userId]
	lock.RUnlock()
	if ok {
		client.DataQueue <- message
	}
}

func sendProcess(client *Client) {
	for {
		select {
		case data := <-client.DataQueue:
			if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				zap.Error(err)
				return
			}
		}
	}
}

func receiveProcess(client *Client) {
	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			zap.Error(err)
			return
		}
		sendMessage("2", data)
	}
}

func Calculate(ctx *gin.Context) {
	utils.Success(ctx, nil, strconv.Itoa(len(H.Clients)))
}
