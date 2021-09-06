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

// Message 消息结构体
type Message struct {
	Id        int64     //消息Id
	Sender    string    //消息发送者
	Receiver  string    //消息接收者
	Category  uint8     //消息类别：私聊|群聊
	TimeStamp time.Time // 消息时间戳
	Content   string    // 消息内容
	Image     string    // 消息缩略图
	Url       string    // 媒体 url
	Memo      string    // 消息备注
	Amount    string    // 数字相关
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

	//todo:客户端标识
	userId := ctx.Query("userId")
	// http 升级 WebSocket 协议
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
		DataQueue: make(chan []byte, 20480),
		flag:      userId,
	}
	//client 注册到 hub
	H.Login <- client

	go writePump(client)
	go readPump(client)
	//todo:不显示消息
	sendMessage(userId, []byte("hello, darling! "+userId+",你已上线"))
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

		//todo: 暂时 demo
		if client.flag == "admin" {
			broadcast(data)
		} else if client.flag == "2" {
			sendMessage("1", data)
		} else {
			sendMessage("2", data)
		}

	}
}

func broadcast(message []byte) {
	for i := range H.Clients {
		if i != "admin" {
			sendMessage(H.Clients[i].flag, message)
		}
	}
}

func Calculate(ctx *gin.Context) {
	for i := range H.Clients {
		fmt.Println(i)
	}
	utils.Success(ctx, nil, strconv.Itoa(len(H.Clients)))
}
