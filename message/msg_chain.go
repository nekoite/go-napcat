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
	Messages []Segment
}

func NewChain(msg ...Segment) *Chain {
	return &Chain{Messages: msg}
}

func (mc *Chain) PrependSegment(msg Segment) {
	mc.Messages = append([]Segment{msg}, mc.Messages...)
}

func (mc *Chain) PrependChain(chain Chain) {
	mc.Messages = append(chain.Messages, mc.Messages...)
}

func (mc *Chain) AddSegment(msg Segment) {
	mc.Messages = append(mc.Messages, msg)
}

func (mc *Chain) AddSegments(msg ...Segment) {
	mc.Messages = append(mc.Messages, msg...)
}

func (mc *Chain) AddText(text string) {
	mc.AddSegments(NewText(text).Segment())
}

func (mc *Chain) AddAt(userId string) {
	mc.AddSegments(NewAt(userId).Segment())
}

func (mc *Chain) AddAtUser(userId qq.UserId) {
	mc.AddSegments(NewAt(userId.String()).Segment())
}

func (mc *Chain) AddAtAll() {
	mc.AddSegments(NewAtAll().Segment())
}

func (mc *Chain) AppendChain(chain Chain) {
	mc.AddSegments(chain.Messages...)
}

// SetSendAsAnonymous 设置是否匿名发送消息。当 ignore 为 true 时，将在无法匿名发送消息时继续发送消息。
func (mc *Chain) SetSendAsAnonymous(ignore bool) {
	mc.PrependSegment(NewAnonymous(ignore).Segment())
}

// SetReplyTo 设置回复的消息 ID。
func (mc *Chain) SetReplyTo(msgId qq.MessageId) {
	mc.PrependSegment(NewReply(msgId).Segment())
}

// At 返回位置在 idx 的消息的拷贝。
func (mc *Chain) At(idx int) Segment {
	return mc.Messages[idx]
}

// AtRef 返回位置在 idx 的消息的引用。
func (mc *Chain) AtRef(idx int) *Segment {
	return &mc.Messages[idx]
}

func (mc *Chain) Len() int {
	return len(mc.Messages)
}

// FirstOfType 返回第一个类型为 msgType 的消息的拷贝。如果没有找到，返回空消息。
func (mc *Chain) FirstOfType(msgType SegmentType) Segment {
	for _, msg := range mc.Messages {
		if msg.Type == msgType {
			return msg
		}
	}
	return Segment{}
}

// FirstOfTypeRef 返回第一个类型为 msgType 的消息的引用。如果没有找到，返回 nil。
func (mc *Chain) FirstOfTypeRef(msgType SegmentType) *Segment {
	for i := 0; i < len(mc.Messages); i++ {
		if mc.Messages[i].Type == msgType {
			return &mc.Messages[i]
		}
	}
	return nil
}

func (mc *Chain) FirstText() *TextData {
	m := mc.FirstOfTypeRef(SegmentTypeText)
	if m == nil {
		return nil
	}
	return m.Data.(*TextData)
}

func (mc *Chain) FirstImage() *ImageData {
	m := mc.FirstOfTypeRef(SegmentTypeImage)
	if m == nil {
		return nil
	}
	return m.Data.(*ImageData)
}

func (mc *Chain) Clear() {
	mc.Messages = nil
}

func (mc *Chain) GetSegmentsWithType(msgType SegmentType) []Segment {
	return utils.Filter(mc.Messages, func(msg Segment) bool {
		return msg.Type == msgType
	})
}

func (mc *Chain) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	mc.Messages = make([]Segment, len(raw))
	for i, r := range raw {
		var msg Segment
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
