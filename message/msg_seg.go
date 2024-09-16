package message

import (
	"reflect"

	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/qq"
	"github.com/nekoite/go-napcat/utils"
	"github.com/tidwall/gjson"
)

// SegmentType 消息片段类型
type SegmentType string

type MusicType string

const (
	SegmentTypeText      SegmentType = "text"
	SegmentTypeFace      SegmentType = "face"
	SegmentTypeImage     SegmentType = "image"
	SegmentTypeRecord    SegmentType = "record"
	SegmentTypeVideo     SegmentType = "video" // MessageTypeVideo 短视频
	SegmentTypeAt        SegmentType = "at"    // MessageTypeAt @某人
	SegmentTypeRps       SegmentType = "rps"   // MessageTypeRps 猜拳魔法表情
	SegmentTypeDice      SegmentType = "dice"  // MessageTypeDice 掷骰子魔法表情
	SegmentTypeShake     SegmentType = "shake" // MessageTypeShake 窗口抖动
	SegmentTypePoke      SegmentType = "poke"
	SegmentTypeAnonymous SegmentType = "anonymous"
	SegmentTypeShare     SegmentType = "share"
	SegmentTypeContact   SegmentType = "contact"
	SegmentTypeLocation  SegmentType = "location"
	SegmentTypeMusic     SegmentType = "music"
	SegmentTypeReply     SegmentType = "reply"
	SegmentTypeForward   SegmentType = "forward"
	SegmentTypeNode      SegmentType = "node"
	SegmentTypeXml       SegmentType = "xml"
	SegmentTypeJson      SegmentType = "json"

	MusicTypeQQ     MusicType = "qq"
	MusicType163    MusicType = "163"
	MusicTypeXm     MusicType = "xm"
	MusicTypeCustom MusicType = "custom"

	ImageTypeFlash = "flash"
)

var (
	atAll = NewAt("all")
)

// Segment 消息片段
type Segment struct {
	Type SegmentType `json:"type"`
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
	UserId   qq.UserId `json:"user_id,string"`
	Nickname string    `json:"nickname"`
	Content  any       `json:"content"`
}

type XmlData struct {
	Data string `json:"data"`
}

type JsonData struct {
	Data string `json:"data"`
}

type UnknownData map[string]any

func (m Segment) AsChain() *Chain {
	return NewChain(m)
}

// GetTextData 获取文本消息数据，如果类型不匹配返回 nil
func (m Segment) GetTextData() *TextData {
	return GetMsgData[TextData](&m)
}

// GetImageData 获取图片消息数据，如果类型不匹配返回 nil
func (m Segment) GetImageData() *ImageData {
	return GetMsgData[ImageData](&m)
}

// GetFaceData 获取表情消息数据，如果类型不匹配返回 nil
func (m Segment) GetFaceData() *FaceData {
	return GetMsgData[FaceData](&m)
}

// GetRecordData 获取语音消息数据，如果类型不匹配返回 nil
func (m Segment) GetRecordData() *RecordData {
	return GetMsgData[RecordData](&m)
}

// GetVideoData 获取视频消息数据，如果类型不匹配返回 nil
func (m Segment) GetVideoData() *VideoData {
	return GetMsgData[VideoData](&m)
}

// GetAtData 获取 @ 消息数据，如果类型不匹配返回 nil
func (m Segment) GetAtData() *AtData {
	return GetMsgData[AtData](&m)
}

// GetRpsData 获取猜拳消息数据，如果类型不匹配返回 nil
func (m Segment) GetRpsData() *RpsData {
	return GetMsgData[RpsData](&m)
}

// GetDiceData 获取掷骰子消息数据，如果类型不匹配返回 nil
func (m Segment) GetDiceData() *DiceData {
	return GetMsgData[DiceData](&m)
}

// GetShakeData 获取窗口抖动消息数据，如果类型不匹配返回 nil
func (m Segment) GetShakeData() *ShakeData {
	return GetMsgData[ShakeData](&m)
}

// GetPokeData 获取戳一戳消息数据，如果类型不匹配返回 nil
func (m Segment) GetPokeData() *PokeData {
	return GetMsgData[PokeData](&m)
}

// GetShareData 获取分享消息数据，如果类型不匹配返回 nil
func (m Segment) GetShareData() *ShareData {
	return GetMsgData[ShareData](&m)
}

// GetContactData 获取联系人消息数据，如果类型不匹配返回 nil
func (m Segment) GetContactData() *ContactData {
	return GetMsgData[ContactData](&m)
}

// GetLocationData 获取位置消息数据，如果类型不匹配返回 nil
func (m Segment) GetLocationData() *LocationData {
	return GetMsgData[LocationData](&m)
}

// GetMusicData 获取音乐消息数据，如果类型不匹配返回 nil
func (m Segment) GetMusicData() *MusicData {
	return GetMsgData[MusicData](&m)
}

// GetReplyData 获取回复消息数据，如果类型不匹配返回 nil
func (m Segment) GetReplyData() *ReplyData {
	return GetMsgData[ReplyData](&m)
}

// GetForwardData 获取转发消息数据，如果类型不匹配返回 nil
func (m Segment) GetForwardData() *ForwardData {
	return GetMsgData[ForwardData](&m)
}

// GetCustomNodeData 获取自定义节点消息数据，如果类型不匹配返回 nil
func (m Segment) GetCustomNodeData() any {
	return GetMsgData[CustomNodeData](&m)
}

// GetXmlData 获取 XML 消息数据，如果类型不匹配返回 nil
func (m Segment) GetXmlData() *XmlData {
	return GetMsgData[XmlData](&m)
}

// GetJsonData 获取 JSON 消息数据，如果类型不匹配返回 nil
func (m Segment) GetJsonData() *JsonData {
	return GetMsgData[JsonData](&m)
}

func (m Segment) IsInvalid() bool {
	return m.Type == ""
}

func (m Segment) GetDataPtr() any {
	return reflect.ValueOf(m.Data).Addr().Interface()
}

func (m *Segment) UnmarshalJSON(data []byte) error {
	var d any
	fields := gjson.ParseBytes(data)
	m.Type = SegmentType(fields.Get("type").String())
	switch m.Type {
	case SegmentTypeText:
		d = new(TextData)
	case SegmentTypeFace:
		d = new(FaceData)
	case SegmentTypeImage:
		d = new(ImageData)
	case SegmentTypeRecord:
		d = new(RecordData)
	case SegmentTypeVideo:
		d = new(VideoData)
	case SegmentTypeAt:
		d = new(AtData)
	case SegmentTypeRps:
		d = new(RpsData)
	case SegmentTypeDice:
		d = new(DiceData)
	case SegmentTypeShake:
		d = new(ShakeData)
	case SegmentTypePoke:
		d = new(PokeData)
	case SegmentTypeAnonymous:
		d = new(AnonymousData)
	case SegmentTypeShare:
		d = new(ShareData)
	case SegmentTypeContact:
		d = new(ContactData)
	case SegmentTypeLocation:
		d = new(LocationData)
	case SegmentTypeMusic:
		musicType := MusicType(fields.Get("data").Get("type").String())
		switch musicType {
		case MusicTypeCustom:
			d = new(CustomMusicData)
		default:
			d = new(MusicData)
		}
	case SegmentTypeReply:
		d = new(ReplyData)
	case SegmentTypeForward:
		d = new(ForwardData)
	case SegmentTypeNode:
		hasId := fields.Get("data").Get("id").Exists()
		if hasId {
			d = new(IdNodeData)
		} else {
			// special handling for custom node
			nd := new(CustomNodeData)
			nd.Content = NewChain()
			d = nd
		}
	case SegmentTypeXml:
		d = new(XmlData)
	case SegmentTypeJson:
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

// GetMsgDataUnsafe 获取类型为 T 的消息数据，如果类型不匹配会引发 panic
func GetMsgDataUnsafe[T any](msg *Segment) *T {
	return msg.Data.(*T)
}

// GetMsgData 安全获取类型为 T 的消息数据，返回 nil 表示类型不匹配
func GetMsgData[T any](msg *Segment) *T {
	data, ok := msg.Data.(*T)
	if !ok {
		return nil
	}
	return data
}

func NewText(text string) *TextData {
	return &TextData{Text: text}
}

func NewTextSegment(text string) Segment {
	return NewText(text).Segment()
}

func NewFace(id int64) *FaceData {
	return &FaceData{Id: id}
}

func NewFaceSegment(id int64) Segment {
	return NewFace(id).Segment()
}

func NewImage(file string) *ImageData {
	return &ImageData{BasicFileData: BasicFileData{File: file}}
}

func NewImageSegment(file string) Segment {
	return NewImage(file).Segment()
}

func NewRecord(file string) *RecordData {
	return &RecordData{BasicFileData: BasicFileData{File: file}}
}

func NewRecordSegment(file string) Segment {
	return NewRecord(file).Segment()
}

func NewVideo(file string) *VideoData {
	return &VideoData{File: file}
}

func NewVideoSegment(file string) Segment {
	return NewVideo(file).Segment()
}

// NewAt 创建 @ 消息，qq 为被 @ 的用户 QQ 号或 all 表示 @ 所有人
func NewAt(qq string) *AtData {
	return &AtData{QQ: qq}
}

func NewAtMessage(qq string) Segment {
	return NewAt(qq).Segment()
}

// NewAtAll 创建 @ 所有人 消息
func NewAtAll() *AtData {
	return atAll
}

func NewAtUser(id qq.UserId) *AtData {
	return NewAt(id.String())
}

func NewRps() *RpsData {
	return &RpsData{}
}

func NewDice() *DiceData {
	return &DiceData{}
}

func NewShake() *ShakeData {
	return &ShakeData{}
}

func NewRpsMessage() Segment {
	return Segment{Type: SegmentTypeRps}
}

func NewDiceMessage() Segment {
	return Segment{Type: SegmentTypeDice}
}

func NewShakeMessage() Segment {
	return Segment{Type: SegmentTypeShake}
}

func NewAnonymous(ignore bool) *AnonymousData {
	return &AnonymousData{Ignore: ignore}
}

func NewShare(title, url string) *ShareData {
	return &ShareData{Title: title, Url: url}
}

func NewContact(t string, id int64) *ContactData {
	return &ContactData{Type: t, BasicIdData: BasicIdData{Id: id}}
}

func NewLocation(lat, lon float64) *LocationData {
	return &LocationData{Lat: lat, Lon: lon}
}

func NewMusic(t MusicType, id int64) *MusicData {
	return &MusicData{BasicMusicData: BasicMusicData{Type: t}, BasicIdData: BasicIdData{Id: id}}
}

func NewCustomMusic(title, url, audio string) *CustomMusicData {
	return &CustomMusicData{BasicMusicData: BasicMusicData{Type: MusicTypeCustom}, Title: title, Url: url, Audio: audio}
}

func NewReply(id qq.MessageId) *ReplyData {
	return &ReplyData{Id: int64(id)}
}

func NewNode(id int64) *IdNodeData {
	return &IdNodeData{Id: id}
}

func NewCustomNode[T SendableMessage](userId qq.UserId, nickname string, content T) *CustomNodeData {
	return &CustomNodeData{UserId: userId, Nickname: nickname, Content: content}
}

func NewXml(data string) *XmlData {
	return &XmlData{Data: data}
}

func NewJson(data string) *JsonData {
	return &JsonData{Data: data}
}

func (d *TextData) Segment() Segment {
	return Segment{
		Type: SegmentTypeText,
		Data: d,
	}
}

func (d *FaceData) Segment() Segment {
	return Segment{
		Type: SegmentTypeFace,
		Data: d,
	}
}

func (d *ImageData) Segment() Segment {
	return Segment{
		Type: SegmentTypeImage,
		Data: d,
	}
}

func (d *ImageData) SetFlash() {
	d.Type = ImageTypeFlash
}

func (d *ImageData) IsFlash() bool {
	return d.Type == ImageTypeFlash
}

func (d *RecordData) Segment() Segment {
	return Segment{
		Type: SegmentTypeRecord,
		Data: d,
	}
}

func (d *VideoData) Segment() Segment {
	return Segment{
		Type: SegmentTypeVideo,
		Data: d,
	}
}

func (d *AtData) Segment() Segment {
	return Segment{
		Type: SegmentTypeAt,
		Data: d,
	}
}

func (d *RpsData) Segment() Segment {
	return Segment{
		Type: SegmentTypeRps,
		Data: d,
	}
}

func (d *DiceData) Segment() Segment {
	return Segment{
		Type: SegmentTypeDice,
		Data: d,
	}
}

func (d *ShakeData) Segment() Segment {
	return Segment{
		Type: SegmentTypeShake,
		Data: d,
	}
}

func (d *PokeData) Segment() Segment {
	return Segment{
		Type: SegmentTypePoke,
		Data: d,
	}
}

func (d *AnonymousData) Segment() Segment {
	return Segment{
		Type: SegmentTypeAnonymous,
		Data: d,
	}
}

func (d *ShareData) Segment() Segment {
	return Segment{
		Type: SegmentTypeShare,
		Data: d,
	}
}

func (d *ContactData) Segment() Segment {
	return Segment{
		Type: SegmentTypeContact,
		Data: d,
	}
}

func (d *LocationData) Segment() Segment {
	return Segment{
		Type: SegmentTypeLocation,
		Data: d,
	}
}

func (d *MusicData) Segment() Segment {
	return Segment{
		Type: SegmentTypeMusic,
		Data: d,
	}
}

func (d *CustomMusicData) Segment() Segment {
	return Segment{
		Type: SegmentTypeMusic,
		Data: d,
	}
}

func (d *ReplyData) Segment() Segment {
	return Segment{
		Type: SegmentTypeReply,
		Data: d,
	}
}

func (d *ForwardData) Segment() Segment {
	return Segment{
		Type: SegmentTypeForward,
		Data: d,
	}
}

func (d *IdNodeData) Segment() Segment {
	return Segment{
		Type: SegmentTypeNode,
		Data: d,
	}
}

func (d *CustomNodeData) Segment() Segment {
	return Segment{
		Type: SegmentTypeNode,
		Data: d,
	}
}

func (d *XmlData) Segment() Segment {
	return Segment{
		Type: SegmentTypeXml,
		Data: d,
	}
}

func (d *JsonData) Segment() Segment {
	return Segment{
		Type: SegmentTypeJson,
		Data: d,
	}
}
