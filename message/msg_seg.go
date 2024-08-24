package message

import (
	"strconv"

	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/qq"
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
	Magic int `json:"magic,omitempty"`
}

type VideoData BasicFileData

type AtData struct {
	QQ string `json:"qq"`
}

type RpsData utils.Void

type DiceData utils.Void

type ShakeData utils.Void

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

type ReplyData BasicIdData

type ForwardData BasicIdData

type IdNodeData BasicIdData

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

type UnknownData map[string]any

func (m Message) AsChain() *Chain {
	return NewChain(m)
}

func (m *Message) GetTextData() TextData {
	return GetMsgData[TextData](m)
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
		d = new(RpsData)
	case MessageTypeDice:
		d = new(DiceData)
	case MessageTypeShake:
		d = new(ShakeData)
	case MessageTypePoke:
		d = new(PokeData)
	case MessageTypeAnonymous:
		d = new(AnonymousData)
	case MessageTypeShare:
		d = new(ShareData)
	case MessageTypeContact:
		d = new(ContactData)
	case MessageTypeLocation:
		d = new(LocationData)
	case MessageTypeMusic:
		musicType := MusicType(fields.Get("data").Get("type").String())
		switch musicType {
		case MusicTypeCustom:
			d = new(CustomMusicData)
		default:
			d = new(MusicData)
		}
	case MessageTypeReply:
		d = new(ReplyData)
	case MessageTypeForward:
		d = new(ForwardData)
	case MessageTypeNode:
		hasId := fields.Get("data").Get("id").Exists()
		if hasId {
			d = new(IdNodeData)
		} else {
			// special handling for custom node
			nd := new(CustomNodeData)
			nd.Content = NewChain()
			d = nd
		}
	case MessageTypeXml:
		d = new(XmlData)
	case MessageTypeJson:
		d = new(JsonData)
	default:
		d = make(UnknownData)
		if err := json.Unmarshal([]byte(fields.Get("data").Raw), &d); err != nil {
			return err
		}
		m.Data = d
		return nil
	}
	if err := json.Unmarshal([]byte(fields.Get("data").Raw), d); err != nil {
		return err
	}
	m.Data = utils.DerefAny(d)
	return nil
}

// func (m *Message) unmarshalForCustomNode(fields *gjson.Result) error {
// 	var d CustomNodeData
// 	data := fields.Get("data")
// 	d.UserId = data.Get("user_id").Int()
// 	d.Nickname = data.Get("nickname").String()
// 	if data.Get("content").IsArray() {
// 		var content Chain
// 		if err := json.Unmarshal([]byte(data.Get("content").Raw), &content); err != nil {
// 			return err
// 		}
// 		d.Content = content
// 	}
// 	m.Data = d
// 	return nil
// }

func GetMsgData[T any](msg *Message) T {
	return msg.Data.(T)
}

func NewText(text string) TextData {
	return TextData{Text: text}
}

func NewFace(id int64) FaceData {
	return FaceData{Id: id}
}

func NewImage(file string) ImageData {
	return ImageData{BasicFileData: BasicFileData{File: file}}
}

func NewRecord(file string) RecordData {
	return RecordData{BasicFileData: BasicFileData{File: file}}
}

func NewVideo(file string) VideoData {
	return VideoData{File: file}
}

// NewAt 创建 @ 消息，qq 为被 @ 的用户 QQ 号或 all 表示 @ 所有人
func NewAt(qq string) AtData {
	return AtData{QQ: qq}
}

func NewAtAll() AtData {
	return AtData{QQ: "all"}
}

func NewAtUser(id int64) AtData {
	return NewAt(strconv.Itoa(int(id)))
}

func NewRps() RpsData {
	return RpsData{}
}

func NewDice() DiceData {
	return DiceData{}
}

func NewShake() ShakeData {
	return ShakeData{}
}

func NewRpsMessage() Message {
	return Message{Type: MessageTypeRps}
}

func NewDiceMessage() Message {
	return Message{Type: MessageTypeDice}
}

func NewShakeMessage() Message {
	return Message{Type: MessageTypeShake}
}

func NewAnonymous(ignore bool) AnonymousData {
	return AnonymousData{Ignore: ignore}
}

func NewShare(title, url string) ShareData {
	return ShareData{Title: title, Url: url}
}

func NewContact(t string, id int64) ContactData {
	return ContactData{Type: t, BasicIdData: BasicIdData{Id: id}}
}

func NewLocation(lat, lon float64) LocationData {
	return LocationData{Lat: lat, Lon: lon}
}

func NewMusic(t MusicType, id int64) MusicData {
	return MusicData{BasicMusicData: BasicMusicData{Type: t}, BasicIdData: BasicIdData{Id: id}}
}

func NewCustomMusic(title, url, audio string) CustomMusicData {
	return CustomMusicData{BasicMusicData: BasicMusicData{Type: MusicTypeCustom}, Title: title, Url: url, Audio: audio}
}

func NewReply(id qq.MessageId) ReplyData {
	return ReplyData{Id: int64(id)}
}

func NewNode(id int64) IdNodeData {
	return IdNodeData{Id: id}
}

func NewCustomNode[T SendableMessage](userId int64, nickname string, content T) CustomNodeData {
	return CustomNodeData{UserId: userId, Nickname: nickname, Content: content}
}

func NewXml(data string) XmlData {
	return XmlData{Data: data}
}

func NewJson(data string) JsonData {
	return JsonData{Data: data}
}

func (d TextData) Message() Message {
	return Message{
		Type: MessageTypeText,
		Data: d,
	}
}

func (d FaceData) Message() Message {
	return Message{
		Type: MessageTypeFace,
		Data: d,
	}
}

func (d ImageData) Message() Message {
	return Message{
		Type: MessageTypeImage,
		Data: d,
	}
}

func (d RecordData) Message() Message {
	return Message{
		Type: MessageTypeRecord,
		Data: d,
	}
}

func (d VideoData) Message() Message {
	return Message{
		Type: MessageTypeVideo,
		Data: d,
	}
}

func (d AtData) Message() Message {
	return Message{
		Type: MessageTypeAt,
		Data: d,
	}
}

func (d RpsData) Message() Message {
	return Message{
		Type: MessageTypeRps,
		Data: d,
	}
}

func (d DiceData) Message() Message {
	return Message{
		Type: MessageTypeDice,
		Data: d,
	}
}

func (d ShakeData) Message() Message {
	return Message{
		Type: MessageTypeShake,
		Data: d,
	}
}

func (d PokeData) Message() Message {
	return Message{
		Type: MessageTypePoke,
		Data: d,
	}
}

func (d AnonymousData) Message() Message {
	return Message{
		Type: MessageTypeAnonymous,
		Data: d,
	}
}

func (d ShareData) Message() Message {
	return Message{
		Type: MessageTypeShare,
		Data: d,
	}
}

func (d ContactData) Message() Message {
	return Message{
		Type: MessageTypeContact,
		Data: d,
	}
}

func (d LocationData) Message() Message {
	return Message{
		Type: MessageTypeLocation,
		Data: d,
	}
}

func (d MusicData) Message() Message {
	return Message{
		Type: MessageTypeMusic,
		Data: d,
	}
}

func (d CustomMusicData) Message() Message {
	return Message{
		Type: MessageTypeMusic,
		Data: d,
	}
}

func (d ReplyData) Message() Message {
	return Message{
		Type: MessageTypeReply,
		Data: d,
	}
}

func (d ForwardData) Message() Message {
	return Message{
		Type: MessageTypeForward,
		Data: d,
	}
}

func (d IdNodeData) Message() Message {
	return Message{
		Type: MessageTypeNode,
		Data: d,
	}
}

func (d CustomNodeData) Message() Message {
	return Message{
		Type: MessageTypeNode,
		Data: d,
	}
}

func (d XmlData) Message() Message {
	return Message{
		Type: MessageTypeXml,
		Data: d,
	}
}

func (d JsonData) Message() Message {
	return Message{
		Type: MessageTypeJson,
		Data: d,
	}
}
