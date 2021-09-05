package websocket

import (
	"GoChat/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
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
		flag:      userId,
	}
	//client 注册到 hub
	H.Login <- client

	go writePump(client)
	go readPump(client)
	//todo:不显示消息
	sendMessage(userId, []byte("hello, darling! "+userId+" ,你已上线"))
}

func serverWs(client *Client) {

}

func sendMessage(userId string, message []byte) {
	lock.RLock()
	client, ok := H.Clients[userId]
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
func readPump(client *Client) {
	defer func() {
		client.Hub.Logout <- client
		client.Conn.Close()
	}()
	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			zap.Error(err)
			return
		}
		if client.flag == "2" {
			sendMessage("1", data)
		} else {
			sendMessage("2", data)
		}
	}
}

func Calculate(ctx *gin.Context) {
	for i := range H.Clients {
		fmt.Println(i)
	}
	utils.Success(ctx, nil, strconv.Itoa(len(H.Clients)))
}
