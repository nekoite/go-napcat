package message

import "fmt"

// CQMessage QQ 消息
type CQMessage interface {
	// GetType 获取消息类型
	GetType() string
	// String CQ 码的消息字符串
	fmt.Stringer
}
