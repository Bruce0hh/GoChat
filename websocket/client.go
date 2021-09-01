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
	WebsocketHub interface{}
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
		zap.Error(err.(error))
		return nil, false
	} else {
		if hub, ok := WebsocketHub.(*Hub); ok {
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

	//设置最大消息读取长度
	c.Conn.SetReadLimit(65535)
	c.Hub.Login <- c
	c.Status = 1
	return c, true

}
