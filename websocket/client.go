package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
)

var (
	wsHub interface{}
)

type Client struct {
	Hub                   *Hub //负责处理客户端上下线、在线管理
	Conn                  *websocket.Conn
	Send                  chan []byte   //存储自己的消息管道
	PingPeriod            time.Duration //心跳检测时间
	ReadDeadline          time.Duration //读取消息最大失败时间
	WriteDeadline         time.Duration //发送消息最大失败时间
	HeartbeatFailureTimes int           //心跳检测最大失败次数
	Status                uint8         // websocket状态码
	sync.RWMutex                        //读写锁
}

//OnConnection WebSocket连接
func (c *Client) OnConnection(ctx *gin.Context) (*Client, bool) {
	//处理错误，恢复现场
	defer func() {
		err := recover()
		if err != nil {
			if val, ok := err.(error); ok {
				zap.Error(val)
			}
		}
	}()

	//1.升级http至websocket
	upgrade := websocket.Upgrader{
		ReadBufferSize:  20480,
		WriteBufferSize: 20480,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//2.初始化一个WebSocket长连接客户端
	if connect, err := upgrade.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
		zap.Error(err)
		return nil, false
	} else {
		if hub, ok := wsHub.(*Hub); ok {
			c.Hub = hub
		}
		c.Conn = connect
		//读写缓存区分配字节
		c.Send = make(chan []byte, 20480)
		//心跳包频率
		c.PingPeriod = time.Second * 20
		//业务消息最小间隔时间
		c.ReadDeadline = time.Second * 100
		//消息单次写入超时时间
		c.WriteDeadline = time.Second * 35
	}

	if err := c.SendMessage(websocket.TextMessage, "WebSocket connect Success!"); err != nil {
		zap.Error(err)
	}

	//设置最大消息读取长度
	c.Conn.SetReadLimit(65535)
	c.Hub.Login <- c
	c.Status = 1
	return c, true

}

//SendMessage 发送消息
func (c *Client) SendMessage(messageType int, message string) (err error) {
	c.Lock()
	defer func() {
		c.Unlock()
	}()
	//发送消息，设置本次消息的最大允许时间（s）
	if err = c.Conn.SetWriteDeadline(time.Now().Add(c.WriteDeadline)); err != nil {
		zap.Error(err)
		return err
	}
	if err = c.Conn.WriteMessage(messageType, []byte(message)); err != nil {
		return err
	} else {
		return nil
	}
}

// ReadPump 实时接收消息
func (c *Client) ReadPump(callbackOnMessage func(messageType int, receivedMessage []byte),
	callbackOnError func(err error), callbackOnClose func()) {
	//回调 OnClose()
	defer func() {
		err := recover()
		if err != nil {
			zap.Error(err.(error))
		}
		callbackOnClose()
	}()

	for {
		if c.Status == 1 {
			message, receivedMessage, err := c.Conn.ReadMessage()
			if err == nil {
				callbackOnMessage(message, receivedMessage)
			} else {
				callbackOnError(err)
				break
			}
		} else {
			break
		}
	}
}

//Heartbeat 心跳检测
func (c *Client) Heartbeat() {
	//设置一个定时器，周期性地发送心跳包
	ticker := time.NewTicker(c.PingPeriod)
	defer func() {
		err := recover()
		if err != nil {
			zap.Error(err.(error))
		}
		ticker.Stop()
	}()
	//重置 ReadDeadline
	if c.ReadDeadline == 0 {
		_ = c.Conn.SetReadDeadline(time.Time{})
	} else {
		_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
	}
	// 接收到的消息——pong 服务器发送的消息——ping 实际两者是同一个数据
	c.Conn.SetPongHandler(func(receivedPong string) error {
		if c.ReadDeadline <= time.Nanosecond {
			_ = c.Conn.SetReadDeadline(time.Time{})
		} else {
			_ = c.Conn.SetReadDeadline(time.Now().Add(c.ReadDeadline))
		}
		return nil
	})

	for {
		select {
		case <-ticker.C:
			if c.Status == 1 {
				if err := c.SendMessage(websocket.PingMessage, "Server->Ping->Client"); err != nil {
					c.HeartbeatFailureTimes++
					//连续4次未检测到 ping, 即表示下线
					if c.HeartbeatFailureTimes > 4 {
						c.Status = 0
						zap.Error(err)
						return
					}
				} else {
					if c.HeartbeatFailureTimes > 0 {
						c.HeartbeatFailureTimes--
					}
				}
			} else {
				return
			}
		}
	}
}
