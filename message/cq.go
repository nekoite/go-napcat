package message

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nekoite/go-napcat/utils"
)

var cqEscapeReplacer = strings.NewReplacer(
	"[", "&#91;",
	"]", "&#93;",
	"&", "&amp;",
	",", "&#44;",
)

type CQStringer fmt.Stringer

// CQMessage QQ 消息
type CQMessage interface {
	// GetType 获取消息类型
	GetType() MessageType
	// String CQ 码的消息字符串
	CQStringer
}

func (m Message) String() string {
	sb := strings.Builder{}
	sb.WriteString("[CQ:")
	sb.WriteString(string(m.Type))
	m.buildDataSegmentString(&sb)
	sb.WriteString("]")
	return sb.String()
}

func (m Message) GetType() MessageType {
	return m.Type
}

func (m Message) buildDataSegmentString(sb *strings.Builder) {
	switch m.Data.(type) {
	case map[string]any:
		for k, v := range m.Data.(map[string]any) {
			sb.WriteString(fmt.Sprintf(",%s=", k))
			sb.WriteString(EscapeCQString(fmt.Sprintf("%v", v)))
		}
		return
	}
	utils.WalkStructWithTag(&m.Data, func(v reflect.Value, tagPath []reflect.StructTag) {
		tag := tagPath[len(tagPath)-1]
		if tag.Get("json") == "-" {
			return
		}
		tags := strings.Split(tag.Get("json"), ",")
		remaining := utils.NewSetFrom(tags[1:]...)
		sb.WriteString(fmt.Sprintf(",%s=", tags[0]))
		if remaining.Contains("omitempty") {
			if !v.IsZero() {
				sb.WriteString(EscapeCQString(fmt.Sprintf("%v", v.Interface())))
			}
		} else {
			sb.WriteString(EscapeCQString(fmt.Sprintf("%v", v.Interface())))
		}
	})
}

func EscapeCQString(s string) string {
	return cqEscapeReplacer.Replace(s)
}
