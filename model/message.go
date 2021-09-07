package model

import (
	"time"
)

// Message 消息结构体
type Message struct {
	Id        int64     `json:"id"`         // 消息Id
	Sender    string    `json:"sender"`     // 消息发送者
	Receiver  string    `json:"receiver"`   // 消息接收者
	Category  uint8     `json:"category"`   // 消息类别：私聊|群聊
	TimeStamp time.Time `json:"time_stamp"` // 消息时间戳
	Content   string    `json:"content"`    // 消息内容
	Image     string    `json:"image"`      // 消息缩略图
	Url       string    `json:"url"`        // 媒体 url
	Memo      string    `json:"memo"`       // 消息备注
	Amount    string    `json:"amount"`     // 数字相关
}
