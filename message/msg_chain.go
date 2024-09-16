package message

import (
	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/qq"
	"github.com/nekoite/go-napcat/utils"
)

type SendableMessage interface {
	string | *Chain
}

type Chain struct {
	Messages []Message
}

func NewChain(msg ...Message) *Chain {
	return &Chain{Messages: msg}
}

func (mc *Chain) PrependMessage(msg Message) {
	mc.Messages = append([]Message{msg}, mc.Messages...)
}

func (mc *Chain) PrependChain(chain Chain) {
	mc.Messages = append(chain.Messages, mc.Messages...)
}

func (mc *Chain) AddMessage(msg Message) {
	mc.Messages = append(mc.Messages, msg)
}

func (mc *Chain) AddMessages(msg ...Message) {
	mc.Messages = append(mc.Messages, msg...)
}

func (mc *Chain) AddText(text string) {
	mc.AddMessages(NewText(text).Message())
}

func (mc *Chain) AddAt(userId string) {
	mc.AddMessages(NewAt(userId).Message())
}

func (mc *Chain) AddAtUser(userId qq.UserId) {
	mc.AddMessages(NewAt(userId.String()).Message())
}

func (mc *Chain) AddAtAll() {
	mc.AddMessages(NewAtAll().Message())
}

func (mc *Chain) AppendChain(chain Chain) {
	mc.AddMessages(chain.Messages...)
}

// SetSendAsAnonymous 设置是否匿名发送消息。当 ignore 为 true 时，将在无法匿名发送消息时继续发送消息。
func (mc *Chain) SetSendAsAnonymous(ignore bool) {
	mc.PrependMessage(NewAnonymous(ignore).Message())
}

func (mc *Chain) SetReplyTo(msgId qq.MessageId) {
	mc.PrependMessage(NewReply(msgId).Message())
}

// At 返回位置在 idx 的消息的拷贝。
func (mc *Chain) At(idx int) Message {
	return mc.Messages[idx]
}

// AtRef 返回位置在 idx 的消息的引用。
func (mc *Chain) AtRef(idx int) *Message {
	return &mc.Messages[idx]
}

func (mc *Chain) Len() int {
	return len(mc.Messages)
}

// FirstOfType 返回第一个类型为 msgType 的消息的拷贝。如果没有找到，返回空消息。
func (mc *Chain) FirstOfType(msgType MessageType) Message {
	for _, msg := range mc.Messages {
		if msg.Type == msgType {
			return msg
		}
	}
	return Message{}
}

// FirstOfTypeRef 返回第一个类型为 msgType 的消息的引用。如果没有找到，返回 nil。
func (mc *Chain) FirstOfTypeRef(msgType MessageType) *Message {
	for i := 0; i < len(mc.Messages); i++ {
		if mc.Messages[i].Type == msgType {
			return &mc.Messages[i]
		}
	}
	return nil
}

func (mc *Chain) FirstText() *TextData {
	m := mc.FirstOfTypeRef(MessageTypeText)
	if m == nil {
		return nil
	}
	return m.Data.(*TextData)
}

func (mc *Chain) FirstImage() *ImageData {
	m := mc.FirstOfTypeRef(MessageTypeImage)
	if m == nil {
		return nil
	}
	return m.Data.(*ImageData)
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

func (mc *Chain) MarshalJSON() ([]byte, error) {
	return json.Marshal(mc.Messages)
}
