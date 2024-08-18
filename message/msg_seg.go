package message

import (
	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/utils"
	"github.com/tidwall/gjson"
)

// MessageType 消息类型
type MessageType string

type MusicType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeFace   MessageType = "face"
	MessageTypeImage  MessageType = "image"
	MessageTypeRecord MessageType = "record"
	// MessageTypeVideo 短视频
	MessageTypeVideo MessageType = "video"
	// MessageTypeAt @某人
	MessageTypeAt MessageType = "at"
	// MessageTypeRps 猜拳魔法表情
	MessageTypeRps MessageType = "rps"
	// MessageTypeDice 掷骰子魔法表情
	MessageTypeDice MessageType = "dice"
	// MessageTypeShake 窗口抖动
	MessageTypeShake     MessageType = "shake"
	MessageTypePoke      MessageType = "poke"
	MessageTypeAnonymous MessageType = "anonymous"
	MessageTypeShare     MessageType = "share"
	MessageTypeContact   MessageType = "contact"
	MessageTypeLocation  MessageType = "location"
	MessageTypeMusic     MessageType = "music"
	MessageTypeReply     MessageType = "reply"
	MessageTypeForward   MessageType = "forward"
	MessageTypeNode      MessageType = "node"
	MessageTypeXml       MessageType = "xml"
	MessageTypeJson      MessageType = "json"

	MusicTypeQQ     MusicType = "qq"
	MusicType163    MusicType = "163"
	MusicTypeXm     MusicType = "xm"
	MusicTypeCustom MusicType = "custom"
)

type Message struct {
	Type MessageType `json:"type"`
	Data any         `json:"data"`
}

type TextData struct {
	Text string `json:"text"`
}

type BasicIdData struct {
	Id int64 `json:"id,string"`
}

type BasicFileData struct {
	File    string `json:"file"`
	URL     string `json:"url,omitempty"`
	Cache   *int   `json:"cache,omitempty"`
	Proxy   *int   `json:"proxy,omitempty"`
	Timeout *int   `json:"timeout,omitempty"`
}

type FileData struct {
	BasicFileData
	// Name 文件名【NapCat 扩展】
	Name string `json:"name,omitempty"`
}

type FaceData BasicIdData

type ImageData struct {
	BasicFileData
	// Summary 自定义显示的文件名【LLOneBot 扩展】
	Summary string `json:"summary,omitempty"`
	Type    string `json:"type,omitempty"`
}

type RecordData struct {
	BasicFileData
	Magic *int `json:"magic,omitempty"`
}

type VideoData BasicFileData

type AtData struct {
	QQ string `json:"qq"`
}

type PokeData struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type AnonymousData struct {
	Ignore bool `json:"ignore,omitempty"`
}

type ShareData struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Content string `json:"content,omitempty"`
	Image   string `json:"image,omitempty"`
}

type ContactData struct {
	Type string `json:"type"`
	BasicIdData
}

type LocationData struct {
	Lat     float64 `json:"lat,string"`
	Lon     float64 `json:"lon,string"`
	Title   string  `json:"title,omitempty"`
	Content string  `json:"content,omitempty"`
}

type BasicMusicData struct {
	Type MusicType `json:"type"`
}

type MusicData struct {
	BasicMusicData
	BasicIdData
}

type CustomMusicData struct {
	BasicMusicData
	Title string `json:"title"`
	Url   string `json:"url"`
	Audio string `json:"audio"`
}

type ReplyData struct {
	BasicIdData
}

type ForwardData struct {
	BasicIdData
}

type IdNodeData struct {
	BasicIdData
}

type CustomNodeData struct {
	UserId   int64  `json:"user_id,string"`
	Nickname string `json:"nickname"`
	Content  any    `json:"content"`
}

type XmlData struct {
	Data string `json:"data"`
}

type JsonData struct {
	Data string `json:"data"`
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var d any
	fields := gjson.ParseBytes(data)
	m.Type = MessageType(fields.Get("type").String())
	switch m.Type {
	case MessageTypeText:
		d = new(TextData)
	case MessageTypeFace:
		d = new(FaceData)
	case MessageTypeImage:
		d = new(ImageData)
	case MessageTypeRecord:
		d = new(RecordData)
	case MessageTypeVideo:
		d = new(VideoData)
	case MessageTypeAt:
		d = new(AtData)
	case MessageTypeRps:
		d = new(BasicIdData)
	default:
		d = make(map[string]any)
		if err := json.Unmarshal([]byte(fields.Get("data").Raw), &d); err != nil {
			return err
		}
		m.Data = d
		return nil
	}
	if err := json.Unmarshal([]byte(fields.Get("data").Raw), &d); err != nil {
		return err
	}
	m.Data = utils.DerefAny(d)
	return nil
}

func GetMsgData[T any](msg *Message) T {
	return msg.Data.(T)
}

func (m *Message) GetTextData() TextData {
	return GetMsgData[TextData](m)
}

func NewText(text string) Message {
	return Message{
		Type: MessageTypeText,
		Data: TextData{Text: text},
	}
}
