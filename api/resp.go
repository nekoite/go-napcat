package api

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/errors"
	"github.com/nekoite/go-napcat/message"
	"github.com/nekoite/go-napcat/qq"
	"github.com/nekoite/go-napcat/utils"
	"github.com/tidwall/gjson"
)

type MessageType string

const (
	MessageTypePrivate MessageType = "private"
	MessageTypeGroup   MessageType = "group"
)

type IResp interface {
	GetStatus() string
	GetRetCode() int
	GetEcho() string
	GetData() any
}

type Resp[T any] struct {
	Status  string `json:"status"`
	RetCode int    `json:"retcode"`
	Echo    string `json:"echo"`
	Data    T      `json:"data"`
}

type RespDataMessageId struct {
	MessageId qq.MessageId `json:"message_id"`
}

type RespDataMessageOnly struct {
	Message message.Chain `json:"message"`
}

type RespDataMessage struct {
	RespDataMessageId
	RespDataMessageOnly
	Time        int64       `json:"time"`
	MessageType MessageType `json:"message_type"`
	RealId      int64       `json:"real_id"`
	Sender      qq.IUser    `json:"-"`
}

type RespDataLoginInfo qq.BasicUser

type RespDataStrangerInfo qq.User

type RespDataFriendList []qq.BasicFriend

type RespDataGroupInfo qq.Group

type RespDataGroupList []qq.Group

type RespDataGroupMemberInfo qq.GroupUser

type RespDataGroupMemberList []qq.GroupUser

type GroupHonorListElement struct {
	qq.BasicUserWithAvatar
	Description string `json:"description"`
}

type RespDataGroupHonorInfo struct {
	GroupId          int64 `json:"group_id"`
	CurrentTalkative struct {
		qq.BasicUserWithAvatar
		DayCount int32 `json:"day_count"`
	} `json:"current_talkative"`
	TalkativeList    []GroupHonorListElement `json:"talkative_list"`
	PerformerList    []GroupHonorListElement `json:"performer_list"`
	LegendList       []GroupHonorListElement `json:"legend_list"`
	StrongNewbieList []GroupHonorListElement `json:"strong_newbie_list"`
	EmotionList      []GroupHonorListElement `json:"emotion_list"`
}

type RespDataCookies struct {
	Cookies string `json:"cookies"`
}

type RespDataCsrfToken struct {
	CsrfToken int32 `json:"csrf_token"`
}

type RespDataCredentials struct {
	RespDataCookies
	RespDataCsrfToken
}

type RespDataFile struct {
	File string `json:"file"`
}

type RespDataYesOrNo struct {
	Yes bool `json:"yes"`
}

type ServerStatus struct {
	Online bool `json:"online"`
	Good   bool `json:"good"`
}

type RespDataVersionInfo struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	ProtocolVersion string `json:"protocol_version"`
}

func (r *Resp[T]) GetStatus() string {
	return r.Status
}

func (r *Resp[T]) GetRetCode() int {
	return r.RetCode
}

func (r *Resp[T]) GetEcho() string {
	return r.Echo
}

func (r *Resp[T]) GetData() any {
	return r.Data
}

func GetDataAs[T any](r IResp) T {
	return r.GetData().(T)
}

func GetRespAs[T any](r IResp) *Resp[T] {
	if r == nil {
		return nil
	}
	return r.(*Resp[T])
}

func parseResp(action Action, data apiResp) (IResp, error) {
	if data.RetCode != 0 {
		return nil, fmt.Errorf("%w code %d", errors.ErrApiResp, data.RetCode)
	}
	var resp any
	switch action {
	case ActionSendPrivateMsg:
		fallthrough
	case ActionSendGroupMsg:
		resp = &Resp[RespDataMessageId]{}
	case ActionDeleteMsg:
		resp = &Resp[utils.Void]{}
	case ActionGetMsg:
		resp = &Resp[RespDataMessage]{}
	case ActionGetForwardMsg:
		resp = &Resp[RespDataMessageOnly]{}
	case ActionSendLike:
		fallthrough
	case ActionSetGroupKick:
		fallthrough
	case ActionSetGroupBan:
		fallthrough
	case ActionSetGroupAnonymousBan:
		fallthrough
	case ActionSetGroupWholeBan:
		fallthrough
	case ActionSetGroupAdmin:
		fallthrough
	case ActionSetGroupAnonymous:
		fallthrough
	case ActionSetGroupCard:
		fallthrough
	case ActionSetGroupLeave:
		fallthrough
	case ActionSetGroupName:
		fallthrough
	case ActionSetGroupSpecialTitle:
		fallthrough
	case ActionSetFriendAddRequest:
		fallthrough
	case ActionSetGroupAddRequest:
		resp = &Resp[utils.Void]{}
	case ActionGetLoginInfo:
		resp = &Resp[RespDataLoginInfo]{}
	case ActionGetStrangerInfo:
		resp = &Resp[RespDataStrangerInfo]{}
	case ActionGetFriendList:
		resp = &Resp[RespDataFriendList]{}
	case ActionGetGroupInfo:
		resp = &Resp[RespDataGroupInfo]{}
	case ActionGetGroupList:
		resp = &Resp[RespDataGroupList]{}
	case ActionGetGroupMemberInfo:
		resp = &Resp[RespDataGroupMemberInfo]{}
	case ActionGetGroupMemberList:
		resp = &Resp[RespDataGroupMemberList]{}
	case ActionGetGroupHonorInfo:
		resp = &Resp[RespDataGroupHonorInfo]{}
	case ActionGetCookies:
		resp = &Resp[RespDataCookies]{}
	case ActionGetCsrfToken:
		resp = &Resp[RespDataCsrfToken]{}
	case ActionGetCredentials:
		resp = &Resp[RespDataCredentials]{}
	case ActionGetRecord:
		fallthrough
	case ActionGetImage:
		resp = &Resp[RespDataFile]{}
	case ActionCanSendImage:
		fallthrough
	case ActionCanSendRecord:
		resp = &Resp[RespDataYesOrNo]{}
	case ActionGetStatus:
		resp = &Resp[ServerStatus]{}
	case ActionGetVersionInfo:
		resp = &Resp[RespDataVersionInfo]{}
	case ActionSetRestart:
		fallthrough
	case ActionCleanCache:
		resp = &Resp[utils.Void]{}
	default:
		if act, ok := extActions[action]; ok {
			resp = act.GetNewResultFunc()
		} else {
			resp = &Resp[any]{}
		}
	}
	if err := json.Unmarshal(data.Raw, resp); err != nil {
		return nil, err
	}
	return resp.(IResp), nil
}

func (r *RespDataMessage) UnmarshalJSON(data []byte) error {
	fields := gjson.GetManyBytes(data, "message_id", "message", "time", "message_type", "real_id", "sender")
	messageType := MessageType(fields[3].String())
	var s qq.IUser
	if messageType == MessageTypePrivate {
		s = &qq.Friend{}
	} else if messageType == MessageTypeGroup {
		s = &qq.GroupUser{}
	} else {
		s = &qq.BasicUser{}
	}
	if err := json.Unmarshal([]byte(fields[5].Raw), s); err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(fields[1].Raw), &r.Message); err != nil {
		return err
	}
	r.MessageId = qq.MessageId(fields[0].Int())
	r.Time = fields[2].Int()
	r.MessageType = messageType
	r.RealId = fields[4].Int()
	r.Sender = utils.DerefAny(s).(qq.IUser)
	return nil
}
