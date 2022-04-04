package websocket

import (
	"GoChat/config"
	"GoChat/model"
	"GoChat/service"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

// Client Websocket 客户端结构
type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	DataQueue chan []byte
	Lock      sync.RWMutex
	flag      string
}

// 初始化变量
var (
	WsClient interface{}
	H        = initHub()
	lock     sync.RWMutex // 读写锁
)

// 全局初始化 Hub
func initHub() *Hub {
	WsClient = CreateHubFactory()
	if wh, ok := WsClient.(*Hub); ok {
		go wh.Run()
	}
	return WsClient.(*Hub)
}

func Chat(ctx *gin.Context) {

	userId := ctx.Query("sender")
	// http 升级 WebSocket 协议
	upgrade := websocket.Upgrader{
		//todo: token 判断
		CheckOrigin: func(r *http.Request) bool {
			return service.CheckUserId(userId, ctx)
		},
	}
	connect, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		zap.Error(err)
		return
	}

	client := &Client{
		Hub:       H,
		Conn:      connect,
		DataQueue: make(chan []byte, 20480),
		flag:      userId,
	}
	//client 注册到 hub
	H.Login <- client

	go writePump(client)
	go readPump(client, ctx)

}

func sendMessage(receiverId string, message []byte) {
	lock.RLock()
	client, ok := H.Clients[receiverId]
	lock.RUnlock()
	if ok {
		client.DataQueue <- message
	}
}

// 服务端--send-->client 消息接收方
func writePump(client *Client) {
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

// 服务端<--receive--client 消息发送方
func readPump(client *Client, ctx *gin.Context) {
	defer func() {
		client.Hub.Logout <- client
		err := client.Conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			zap.Error(err)
			return
		}
		//todo: 从前端传来的 data 应该为一个 Message 结构体的 Json
		// 解析 Message
		dispatcher(client, data, ctx)

	}
}

// 广播功能
func broadcast(message []byte) {
	for i := range H.Clients {
		if i != "admin" {
			sendMessage(H.Clients[i].flag, message)
		}
	}
}

// 解析消息并存入 MongoDB
func dispatcher(client *Client, data []byte, ctx *gin.Context) {
	message := &model.Message{
		Content:   string(data),
		TimeStamp: time.Now(),
		Receiver:  ctx.Query("receiverId"),
		Sender:    client.flag,
	}
	//todo: 暂时 demo
	if client.flag == "admin" {
		message.Category = 8
		broadcast(data)
	} else {
		message.Category = 0
		receiver := ctx.Query("receiverId")
		if receiver != "" {
			sendMessage(receiver, data)
		}
	}
	_, err := config.MongoDB.Collection("message").InsertOne(context.Background(), message)
	if err != nil {
		return
	}
}
