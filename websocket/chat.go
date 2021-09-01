package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WsClient struct {
	WsClient *Client
}

//OnConnection 开启 WebSocket 连接
func (ws *WsClient) OnConnection(ctx *gin.Context) (*WsClient, bool) {
	if client, ok := (&Client{}).OnConnection(ctx); ok {
		ws.WsClient = client
		go ws.WsClient.Heartbeat()
		return ws, true
	} else {
		return nil, false
	}
}

//OnMessage 处理业务消息
func (ws *WsClient) OnMessage(ctx *gin.Context) {
	go ws.WsClient.ReadPump(func(messageType int, receivedMessage []byte) {
		//messageType 消息类型	receivedMessage 服务器收到客户端发来的数据
		tempMessage := "服务器已经收到========>" + string(receivedMessage)
		if err := ws.WsClient.SendMessage(messageType, tempMessage); err != nil {
			zap.Error(err)
		}
	}, ws.OnError, ws.OnClose)
}

// OnError 错误处理
func (ws *WsClient) OnError(err error) {
	ws.WsClient.Status = 0
	zap.Error(err)
}

//OnClose 关闭客户端
func (ws *WsClient) OnClose() {
	ws.WsClient.Hub.Logout <- ws.WsClient
}

//BroadcastMessage 广播消息
func (ws *WsClient) BroadcastMessage(message string) {
	for onlineClients := range ws.WsClient.Hub.Clients {
		if err := onlineClients.SendMessage(websocket.TextMessage, message); err != nil {
			zap.Error(err)
		}
	}
}
