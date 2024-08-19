package message

import (
	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/utils"
)

type SendableMessage interface {
	string | Chain
}

type Chain struct {
	Messages []Message
}

func NewMessageChain() *Chain {
	return &Chain{}
}

func (mc *Chain) AddMessages(msg ...Message) {
	mc.Messages = append(mc.Messages, msg...)
}

func (mc *Chain) AddText(text string) {
	mc.AddMessages(NewText(text).Message())
}

func (mc *Chain) AppendChain(chain Chain) {
	mc.AddMessages(chain.Messages...)
}

// SetSendAsAnonymous 设置是否匿名发送消息。当 ignore 为 true 时，将在无法匿名发送消息时继续发送消息。
func (mc *Chain) SetSendAsAnonymous(ignore bool) {
	mc.AddMessages(NewAnonymous(ignore).Message())
}

func (mc *Chain) Clear() {
	clear(mc.Messages)
}

func (mc *Chain) GetMessagesWithType(msgType MessageType) []Message {
	return utils.Filter(mc.Messages, func(msg Message) bool {
		return msg.Type == msgType
	})
}

func (mc *Chain) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	mc.Messages = make([]Message, len(raw))
	for i, r := range raw {
		var msg Message
		if err := json.Unmarshal(r, &msg); err != nil {
			return err
		}
		mc.Messages[i] = msg
	}
	return nil
}
