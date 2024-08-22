package message

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nekoite/go-napcat/errors"
	"github.com/nekoite/go-napcat/utils"
)

var (
	cqEscapeReplacer = strings.NewReplacer(
		"[", "&#91;",
		"]", "&#93;",
		"&", "&amp;",
		",", "&#44;",
	)
	cqUnescapeReplacer = strings.NewReplacer(
		"&#91;", "[",
		"&#93;", "]",
		"&amp;", "&",
		"&#44;", ",",
	)
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
	if m.Type == MessageTypeText {
		return EscapeCQString(m.GetTextData().Text)
	}
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
	case UnknownData:
		for k, v := range m.Data.(UnknownData) {
			sb.WriteString(fmt.Sprintf(",%s=", k))
			sb.WriteString(EscapeCQString(fmt.Sprintf("%v", v)))
		}
		return
	}
	utils.WalkStructWithTag(&m.Data, func(v reflect.Value, tagPath []reflect.StructTag) error {
		tag := tagPath[len(tagPath)-1]
		if tag.Get("json") == "-" {
			return nil
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
		return nil
	})
}

func (mc *Chain) String() string {
	sb := strings.Builder{}
	for _, msg := range mc.Messages {
		sb.WriteString(msg.String())
	}
	return sb.String()
}

func ParseCQString(s string) (*Chain, error) {
	chain := NewChain()
	idx := 0
	for {
		start := strings.Index(s[idx:], "[CQ:")
		if start == -1 {
			break
		}
		start += idx
		if start-idx > 0 {
			chain.Messages = append(chain.Messages, NewText(UnescapeCQString(s[idx:start])).Message())
		}
		end := strings.Index(s[start:], "]")
		if end == -1 {
			return chain, errors.ErrInvalidCQString
		}
		end += start
		msg, err := parseSingleCQComponent(s[start : end+1])
		if err != nil {
			return chain, err
		}
		chain.Messages = append(chain.Messages, msg)
		idx = end + 1
	}
	if idx < len(s) {
		chain.Messages = append(chain.Messages, NewText(UnescapeCQString(s[idx:])).Message())
	}
	return chain, nil
}

func parseSingleCQComponent(s string) (Message, error) {
	if !strings.HasPrefix(s, "[CQ:") || !strings.HasSuffix(s, "]") {
		return Message{}, errors.ErrInvalidCQString
	}
	s = s[4 : len(s)-1]
	parts := strings.Split(s, ",")
	if len(parts) == 0 {
		return Message{}, errors.ErrInvalidCQString
	}
	ty := MessageType(parts[0])
	parts = parts[1:]
	partsMap := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return Message{}, errors.ErrInvalidCQString
		}
		partsMap[kv[0]] = kv[1]
	}
	return parseMessagePart(ty, partsMap)
}

func parseMessagePart(ty MessageType, parts map[string]string) (Message, error) {
	m := Message{Type: ty}
	if len(parts) == 0 {
		return m, nil
	}
	switch ty {
	case MessageTypeText:
		m.Data = new(TextData)
	case MessageTypeFace:
		m.Data = new(FaceData)
	case MessageTypeImage:
		m.Data = new(ImageData)
	case MessageTypeRecord:
		m.Data = new(RecordData)
	case MessageTypeVideo:
		m.Data = new(VideoData)
	case MessageTypeAt:
		m.Data = new(AtData)
	case MessageTypeRps:
		m.Data = new(RpsData)
	case MessageTypeDice:
		m.Data = new(DiceData)
	case MessageTypeShake:
		m.Data = new(ShakeData)
	case MessageTypePoke:
		m.Data = new(PokeData)
	case MessageTypeAnonymous:
		m.Data = new(AnonymousData)
	case MessageTypeShare:
		m.Data = new(ShareData)
	case MessageTypeContact:
		m.Data = new(ContactData)
	case MessageTypeLocation:
		m.Data = new(LocationData)
	case MessageTypeMusic:
		if parts["type"] == "custom" {
			m.Data = new(CustomMusicData)
		} else {
			m.Data = new(MusicData)
		}
	case MessageTypeReply:
		m.Data = new(ReplyData)
	case MessageTypeForward:
		m.Data = new(ForwardData)
	case MessageTypeNode:
		if parts["id"] != "" {
			m.Data = new(IdNodeData)
		} else {
			d := new(CustomNodeData)
			d.Content = NewChain()
			m.Data = d
		}
	case MessageTypeXml:
		m.Data = XmlData{}
	case MessageTypeJson:
		m.Data = JsonData{}
	default:
		m.Data = parts
		return m, nil
	}
	err := utils.WalkStructWithTag(m.Data, func(v reflect.Value, tagPath []reflect.StructTag) error {
		tag := tagPath[len(tagPath)-1]
		if tag.Get("json") == "-" {
			return nil
		}
		tags := strings.Split(tag.Get("json"), ",")
		v = reflect.Indirect(v)
		if v.CanSet() {
			if val, ok := parts[tags[0]]; ok {
				val = UnescapeCQString(val)
				switch v.Kind() {
				case reflect.String:
					v.SetString(val)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					// todo: throw errors
					i := utils.MustAtoi(val)
					v.SetInt(i)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					i := utils.MustAtoui(val)
					v.SetUint(i)
				case reflect.Float32, reflect.Float64:
					f := utils.MustAtof(val)
					v.SetFloat(f)
				case reflect.Bool:
					b := utils.MustAtob(val)
					v.SetBool(b)
				case reflect.Struct:
					switch v.Interface().(type) {
					case Chain:
						chain, err := ParseCQString(val)
						if err != nil {
							return err
						}
						v.Set(reflect.ValueOf(chain).Elem())
					default:
						return errors.ErrInvalidCQString
					}
				default:
					return errors.ErrInvalidCQString
				}
			}
		}
		return nil
	})
	m.Data = utils.DerefAny(m.Data)
	return m, err
}

func EscapeCQString(s string) string {
	return cqEscapeReplacer.Replace(s)
}

func UnescapeCQString(s string) string {
	return cqUnescapeReplacer.Replace(s)
}
