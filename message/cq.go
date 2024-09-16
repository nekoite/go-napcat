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
	GetType() SegmentType
	// String CQ 码的消息字符串
	CQStringer
}

func (m Segment) String() string {
	if m.Type == SegmentTypeText {
		return EscapeCQString(m.GetTextData().Text)
	}
	sb := strings.Builder{}
	sb.WriteString("[CQ:")
	sb.WriteString(string(m.Type))
	m.buildDataSegmentString(&sb)
	sb.WriteString("]")
	return sb.String()
}

func (m Segment) GetType() SegmentType {
	return m.Type
}

func (m Segment) buildDataSegmentString(sb *strings.Builder) {
	switch m.Data.(type) {
	case UnknownData:
		for k, v := range m.Data.(UnknownData) {
			sb.WriteString(fmt.Sprintf(",%s=", k))
			sb.WriteString(EscapeCQString(fmt.Sprintf("%v", v)))
		}
		return
	}
	utils.WalkStructLeafWithTag(m.Data, func(v reflect.Value, tagPath []reflect.StructTag) error {
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
			chain.Messages = append(chain.Messages, NewText(UnescapeCQString(s[idx:start])).Segment())
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
		chain.Messages = append(chain.Messages, NewText(UnescapeCQString(s[idx:])).Segment())
	}
	return chain, nil
}

func parseSingleCQComponent(s string) (Segment, error) {
	if !strings.HasPrefix(s, "[CQ:") || !strings.HasSuffix(s, "]") {
		return Segment{}, errors.ErrInvalidCQString
	}
	s = s[4 : len(s)-1]
	parts := strings.Split(s, ",")
	if len(parts) == 0 {
		return Segment{}, errors.ErrInvalidCQString
	}
	ty := SegmentType(parts[0])
	parts = parts[1:]
	partsMap := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return Segment{}, errors.ErrInvalidCQString
		}
		partsMap[kv[0]] = kv[1]
	}
	return parseMessagePart(ty, partsMap)
}

func parseMessagePart(ty SegmentType, parts map[string]string) (Segment, error) {
	m := Segment{Type: ty}
	if len(parts) == 0 {
		return m, nil
	}
	switch ty {
	case SegmentTypeText:
		m.Data = new(TextData)
	case SegmentTypeFace:
		m.Data = new(FaceData)
	case SegmentTypeImage:
		m.Data = new(ImageData)
	case SegmentTypeRecord:
		m.Data = new(RecordData)
	case SegmentTypeVideo:
		m.Data = new(VideoData)
	case SegmentTypeAt:
		m.Data = new(AtData)
	case SegmentTypeRps:
		m.Data = new(RpsData)
	case SegmentTypeDice:
		m.Data = new(DiceData)
	case SegmentTypeShake:
		m.Data = new(ShakeData)
	case SegmentTypePoke:
		m.Data = new(PokeData)
	case SegmentTypeAnonymous:
		m.Data = new(AnonymousData)
	case SegmentTypeShare:
		m.Data = new(ShareData)
	case SegmentTypeContact:
		m.Data = new(ContactData)
	case SegmentTypeLocation:
		m.Data = new(LocationData)
	case SegmentTypeMusic:
		if parts["type"] == "custom" {
			m.Data = new(CustomMusicData)
		} else {
			m.Data = new(MusicData)
		}
	case SegmentTypeReply:
		m.Data = new(ReplyData)
	case SegmentTypeForward:
		m.Data = new(ForwardData)
	case SegmentTypeNode:
		if parts["id"] != "" {
			m.Data = new(IdNodeData)
		} else {
			d := new(CustomNodeData)
			d.Content = NewChain()
			m.Data = d
		}
	case SegmentTypeXml:
		m.Data = XmlData{}
	case SegmentTypeJson:
		m.Data = JsonData{}
	default:
		m.Data = parts
		return m, nil
	}
	err := utils.WalkStructLeafWithTag(m.Data, func(v reflect.Value, tagPath []reflect.StructTag) error {
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
	return m, err
}

func EscapeCQString(s string) string {
	return cqEscapeReplacer.Replace(s)
}

func UnescapeCQString(s string) string {
	return cqUnescapeReplacer.Replace(s)
}
