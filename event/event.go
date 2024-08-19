package event

import (
	"encoding/json"

	"github.com/nekoite/go-napcat/api"
	"github.com/nekoite/go-napcat/errors"
	"github.com/nekoite/go-napcat/message"
	"github.com/nekoite/go-napcat/qq"
	"github.com/tidwall/gjson"
)

type EventType string

type MetaEventType string

type MetaEventSubtype string

type MessageEventType api.MessageType

type MessageEventSubtype string

type NoticeEventType string

type NoticeEventSubtype string

type RequestEventType string

type GroupRequestSubtype string

type HonorType string

const (
	EventTypeMessage EventType = "message"
	EventTypeNotice  EventType = "notice"
	EventTypeRequest EventType = "request"
	EventTypeMeta    EventType = "meta_event"

	MetaEventTypeLifecycle MetaEventType = "lifecycle"
	MetaEventTypeHeartbeat MetaEventType = "heartbeat"

	MetaEventSubtypeConnect MetaEventSubtype = "connect"
	MetaEventSubtypeDisable MetaEventSubtype = "disable"
	MetaEventSubtypeEnable  MetaEventSubtype = "enable"
	MetaEventSubtypeNone    MetaEventSubtype = ""

	MessageEventTypePrivate MessageEventType = MessageEventType(api.MessageTypePrivate)
	MessageEventTypeGroup   MessageEventType = MessageEventType(api.MessageTypeGroup)

	MessageEventSubtypeFriend    MessageEventSubtype = "friend"
	MessageEventSubtypeGroup     MessageEventSubtype = "group"
	MessageEventSubtypeOther     MessageEventSubtype = "other"
	MessageEventSubtypeNormal    MessageEventSubtype = "normal"
	MessageEventSubtypeAnonymous MessageEventSubtype = "anonymous"
	MessageEventSubtypeNotice    MessageEventSubtype = "notice"

	NoticeEventTypeGroupUpload   NoticeEventType = "group_upload"
	NoticeEventTypeGroupAdmin    NoticeEventType = "group_admin"
	NoticeEventTypeGroupDecrease NoticeEventType = "group_decrease"
	NoticeEventTypeGroupIncrease NoticeEventType = "group_increase"
	NoticeEventTypeGroupBan      NoticeEventType = "group_ban"
	NoticeEventTypeGroupRecall   NoticeEventType = "group_recall"
	NoticeEventTypeFriendAdd     NoticeEventType = "friend_add"
	NoticeEventTypeFriendRecall  NoticeEventType = "friend_recall"
	NoticeEventTypeNotify        NoticeEventType = "notify"

	NoticeEventSubtypeApprove   NoticeEventSubtype = "approve"
	NoticeEventSubtypeInvite    NoticeEventSubtype = "invite"
	NoticeEventSubtypeKick      NoticeEventSubtype = "kick"
	NoticeEventSubtypeLeave     NoticeEventSubtype = "leave"
	NoticeEventSubtypeKickMe    NoticeEventSubtype = "kick_me"
	NoticeEventSubtypeSet       NoticeEventSubtype = "set"
	NoticeEventSubtypeUnset     NoticeEventSubtype = "unset"
	NoticeEventSubtypePoke      NoticeEventSubtype = "poke"
	NoticeEventSubtypeLuckyKing NoticeEventSubtype = "lucky_king"
	NoticeEventSubtypeHonor     NoticeEventSubtype = "honor"
	NoticeEventSubtypeBan       NoticeEventSubtype = "ban"
	NoticeEventSubtypeLiftBan   NoticeEventSubtype = "lift_ban"

	RequestEventTypeFriend RequestEventType = "friend"
	RequestEventTypeGroup  RequestEventType = "group"

	GroupRequestSubtypeAdd    GroupRequestSubtype = "add"
	GroupRequestSubtypeInvite GroupRequestSubtype = "invite"

	HonorTypeTalkative HonorType = "talkative"
	HonorTypePerformer HonorType = "performer"
	HonorTypeEmotion   HonorType = "emotion"
)

type IEvent interface {
	GetTime() int64
	GetSelfId() int64
	GetEventType() EventType

	AsPrivateMessageEvent() *PrivateMessageEvent
	AsGroupMessageEvent() *GroupMessageEvent
	AsNoticeEventGroupUpload() *NoticeEventGroupUpload
	AsNoticeEventGroupOperation() *NoticeEventGroupOperation
	AsNoticeEventGroupBan() *NoticeEventGroupBan
	AsNoticeEventGroupRecall() *NoticeEventGroupRecall
	AsNoticeEventFriendRecall() *NoticeEventFriendRecall
	AsNoticeEventFriendAdd() *NoticeEventFriendAdd
	AsNoticeEventGroupNotify() *NoticeEventGroupNotify
	AsNoticeEventGroupHonor() *NoticeEventGroupHonor
	AsFriendRequestEvent() *FriendRequestEvent
	AsGroupRequestEvent() *GroupRequestEvent
	AsMetaEvent() *MetaEvent

	PreventDefault()
	GetError() error

	isDefaultPrevented() bool
	setApiSender(*api.Sender)
	setError(error)
}

type IMessageEvent interface {
	IEvent
	GetMessageEventType() MessageEventType
}

type BaseEvent struct {
	Time      int64     `json:"time"`
	SelfId    int64     `json:"self_id"`
	EventType EventType `json:"post_type"`

	isPrevented bool        `json:"-"`
	apiSender   *api.Sender `json:"-"`
	error       error       `json:"-"`
}

func (e *BaseEvent) GetTime() int64 {
	return e.Time
}

func (e *BaseEvent) GetSelfId() int64 {
	return e.SelfId
}

func (e *BaseEvent) GetEventType() EventType {
	return e.EventType
}

// PreventDefault 阻止事件继续传播。别问我为什么取这个名字。
// 当 BotConfig.UseGoroutine 为 true 时，这个函数无效。
func (e *BaseEvent) PreventDefault() {
	e.isPrevented = true
}

func (e *BaseEvent) GetError() error {
	return e.error
}

func (e *BaseEvent) isDefaultPrevented() bool {
	return e.isPrevented
}

func (e *BaseEvent) setApiSender(s *api.Sender) {
	e.apiSender = s
}

func (e *BaseEvent) setError(err error) {
	e.error = err
}

type AnonymousData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

type MessageEvent struct {
	BaseEvent
	MessageType MessageEventType    `json:"message_type"`
	SubType     MessageEventSubtype `json:"sub_type"`
	MessageId   int32               `json:"message_id"`
	UserId      int64               `json:"user_id"`
	Message     message.Chain       `json:"message"`
	RawMessage  string              `json:"raw_message"`
	Font        int32               `json:"font"`
}

type PrivateMessageEvent struct {
	MessageEvent
	Sender qq.User `json:"sender"`
}

type GroupMessageEvent struct {
	MessageEvent
	GroupId   int64         `json:"group_id"`
	Anonymous AnonymousData `json:"anonymous"`
	Sender    qq.GroupUser  `json:"sender"`
}

type NoticeEvent struct {
	BaseEvent
	NoticeType NoticeEventType    `json:"notice_type"`
	SubType    NoticeEventSubtype `json:"sub_type"`
	UserId     int64              `json:"user_id"`
}

type GroupNoticeEvent struct {
	NoticeEvent
	GroupId int64 `json:"group_id"`
}

type NoticeEventGroupUpload struct {
	GroupNoticeEvent
	File struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Size  int64  `json:"size"`
		Busid int64  `json:"busid"`
	} `json:"file"`
}

type NoticeEventGroupOperation struct {
	GroupNoticeEvent
	OperatorId int64 `json:"operator_id"`
}

type NoticeEventGroupBan struct {
	NoticeEventGroupOperation
	Duration int64 `json:"duration"`
}

type NoticeEventGroupRecall struct {
	GroupNoticeEvent
	MessageId int64 `json:"message_id"`
}

type NoticeEventFriendRecall struct {
	NoticeEvent
	MessageId int64 `json:"message_id"`
}

type NoticeEventFriendAdd NoticeEvent

type NoticeEventGroupNotify struct {
	GroupNoticeEvent
	TargetId int64 `json:"target_id"`
}

type NoticeEventGroupHonor struct {
	GroupNoticeEvent
	HonorType HonorType `json:"honor_type"`
}

type RequestEvent struct {
	BaseEvent
	RequestType RequestEventType `json:"request_type"`
	UserId      int64            `json:"user_id"`
	Comment     string           `json:"comment"`
	Flag        string           `json:"flag"`
}

type FriendRequestEvent RequestEvent

type GroupRequestEvent struct {
	RequestEvent
	SubType GroupRequestSubtype `json:"sub_type"`
	GroupId int64               `json:"group_id"`
}

type MetaEvent struct {
	BaseEvent
	MetaEventType MetaEventType     `json:"meta_event_type"`
	SubType       MetaEventSubtype  `json:"sub_type"`
	Status        *api.ServerStatus `json:"status"`
	Interval      int64             `json:"interval"`
}

func (e *MessageEvent) GetMessageEventType() MessageEventType {
	return e.MessageType
}

func ParseEvent(data []byte, apiSender *api.Sender) (IEvent, error) {
	typeInfos := gjson.GetManyBytes(data, "post_type", "message_type", "notice_type", "request_type", "sub_type")
	var e IEvent
	var err error

	switch EventType(typeInfos[0].String()) {
	case EventTypeMessage:
		switch MessageEventType(typeInfos[1].String()) {
		case MessageEventTypePrivate:
			e = new(PrivateMessageEvent)
		case MessageEventTypeGroup:
			e = new(GroupMessageEvent)
		default:
			err = errors.ErrUnknownMessageEvent
			e = new(MessageEvent)
		}
	case EventTypeNotice:
		switch NoticeEventType(typeInfos[2].String()) {
		case NoticeEventTypeGroupUpload:
			e = new(NoticeEventGroupUpload)
		case NoticeEventTypeGroupIncrease:
			fallthrough
		case NoticeEventTypeGroupDecrease:
			e = new(NoticeEventGroupOperation)
		case NoticeEventTypeGroupBan:
			e = new(NoticeEventGroupBan)
		case NoticeEventTypeGroupRecall:
			e = new(NoticeEventGroupRecall)
		case NoticeEventTypeFriendRecall:
			e = new(NoticeEventFriendRecall)
		case NoticeEventTypeGroupAdmin:
			e = new(GroupNoticeEvent)
		case NoticeEventTypeNotify:
			switch NoticeEventSubtype(typeInfos[4].String()) {
			case NoticeEventSubtypeHonor:
				e = new(NoticeEventGroupHonor)
			case NoticeEventSubtypeLuckyKing:
				fallthrough
			case NoticeEventSubtypePoke:
				e = new(NoticeEventGroupNotify)
			default:
				err = errors.ErrUnknownNoticeEvent
				e = new(NoticeEvent)
			}
		case NoticeEventTypeFriendAdd:
			e = new(NoticeEventFriendAdd)
		default:
			err = errors.ErrUnknownNoticeEvent
			e = new(NoticeEvent)
		}
	case EventTypeRequest:
		if RequestEventType(typeInfos[3].String()) == RequestEventTypeGroup {
			e = new(GroupRequestEvent)
		} else if RequestEventType(typeInfos[3].String()) == RequestEventTypeFriend {
			e = new(FriendRequestEvent)
		} else {
			err = errors.ErrUnknownRequestEvent
			e = new(RequestEvent)
		}
	case EventTypeMeta:
		e = new(MetaEvent)
	default:
		err = errors.ErrUnknownEvent
		e = new(BaseEvent)
	}
	e.setError(err)
	e.setApiSender(apiSender)
	if err := json.Unmarshal(data, e); err != nil {
		return e, err
	}
	return e, err
}
