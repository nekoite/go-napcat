package api

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/errors"
	"github.com/nekoite/go-napcat/message"
	"github.com/nekoite/go-napcat/qq"
	"github.com/nekoite/go-napcat/utils"
	"github.com/nekoite/go-napcat/ws"
	"go.uber.org/zap"
)

type Action string

const (
	ActionSendPrivateMsg      Action = "send_private_msg"
	ActionSendGroupMsg        Action = "send_group_msg"
	ActionSendMsg             Action = "send_msg"
	ActionDeleteMsg           Action = "delete_msg"
	ActionGetMsg              Action = "get_msg"
	ActionGetForwardMsg       Action = "get_forward_msg"
	ActionSendLike            Action = "send_like"
	ActionSetGroupKick        Action = "set_group_kick"
	ActionSetGroupBan         Action = "set_group_ban"
	ActionSetGroupWholeBan    Action = "set_group_whole_ban"
	ActionSetGroupAdmin       Action = "set_group_admin"
	ActionSetGroupCard        Action = "set_group_card"
	ActionSetGroupName        Action = "set_group_name"
	ActionSetGroupLeave       Action = "set_group_leave"
	ActionSetFriendAddRequest Action = "set_friend_add_request"
	ActionSetGroupAddRequest  Action = "set_group_add_request"
	ActionGetLoginInfo        Action = "get_login_info"
	ActionGetStrangerInfo     Action = "get_stranger_info"
	ActionGetFriendList       Action = "get_friend_list"
	ActionGetGroupList        Action = "get_group_list"
	ActionGetGroupInfo        Action = "get_group_info"
	ActionGetGroupMemberInfo  Action = "get_group_member_info"
	ActionGetGroupMemberList  Action = "get_group_member_list"
	ActionGetGroupHonorInfo   Action = "get_group_honor_info"
	ActionGetCookies          Action = "get_cookies"
	ActionGetRecord           Action = "get_record"
	ActionGetImage            Action = "get_image"
	ActionCanSendImage        Action = "can_send_image"
	ActionCanSendRecord       Action = "can_send_record"
	ActionGetStatus           Action = "get_status"
	ActionGetVersionInfo      Action = "get_version_info"
	ActionSetRestart          Action = "set_restart"
	ActionCleanCache          Action = "clean_cache"

	ActionHandleQuickOperation Action = ".handle_quick_operation"

	// Deprecated:
	ActionSetGroupAnonymous Action = "set_group_anonymous"
	// Deprecated:
	ActionSetGroupAnonymousBan Action = "set_group_anonymous_ban"
	// Deprecated:
	ActionSetGroupSpecialTitle Action = "set_group_special_title"
	// Deprecated:
	ActionGetCsrfToken Action = "get_csrf_token"
	// Deprecated:
	ActionGetCredentials Action = "get_credentials"
)

type Sender struct {
	logger  *zap.Logger
	conn    *ws.Client
	timeout int
	sendId  atomic.Int64
	reqMap  sync.Map
}

type Req struct {
	id     int64
	resp   chan apiResp
	action Action
}

type apiReq struct {
	Action Action `json:"action"`
	Params any    `json:"params"`
	Echo   string `json:"echo"`
}

type apiResp struct {
	Status  string `json:"status"`
	Echo    string `json:"echo"`
	RetCode int    `json:"retcode"`
	Raw     []byte `json:"-"`
}

type SendMsgReqParams[T message.SendableMessage] struct {
	MessageType string `json:"message_type"`
	UserId      int64  `json:"user_id,omitempty"`
	GroupId     int64  `json:"group_id,omitempty"`
	Message     T      `json:"message"`
	AutoEscape  bool   `json:"auto_escape,omitempty"`
}

func NewSender(logger *zap.Logger, conn *ws.Client, timeout int) *Sender {
	return &Sender{
		logger:  logger.Named("api"),
		sendId:  atomic.Int64{},
		conn:    conn,
		timeout: timeout,
	}
}

func (s *Sender) NewReq(action Action) *Req {
	return &Req{
		id:     s.sendId.Add(1),
		resp:   make(chan apiResp, 1),
		action: action,
	}
}

func (s *Sender) HandleApiResp(resp []byte) error {
	var r apiResp
	if err := json.Unmarshal(resp, &r); err != nil {
		return err
	}
	r.Raw = resp
	id, err := strconv.Atoi(r.Echo)
	if err != nil {
		return err
	}
	if req, ok := s.reqMap.LoadAndDelete(int64(id)); ok {
		req.(*Req).resp <- r
		return nil
	}
	return errors.ErrUnknownResponse
}

func (s *Sender) SendRaw(action Action, params any) (IResp, error) {
	req := s.NewReq(action)
	s.reqMap.Store(req.id, req)
	apiReq := &apiReq{
		Action: req.action,
		Params: params,
		Echo:   strconv.FormatInt(req.id, 10),
	}
	raw, err := json.Marshal(apiReq)
	if err != nil {
		return nil, err
	}
	s.conn.Send(raw)
	select {
	case resp := <-req.resp:
		return parseResp(req.action, resp)
	case <-time.After(time.Duration(s.timeout) * time.Second):
		s.logger.Error("timeout", zap.String("action", string(action)), zap.Any("params", params), zap.Int("echo", int(req.id)))
		return nil, errors.ErrTimeout
	}
}

func (s *Sender) SendPrivateMsgString(userId int64, message string, autoEscape bool) (*Resp[RespDataMessageId], error) {
	return returnAsType[RespDataMessageId](s.SendRaw(ActionSendPrivateMsg, map[string]any{
		"user_id":     userId,
		"message":     message,
		"auto_escape": autoEscape,
	}))
}

func (s *Sender) SendPrivateMsg(userId int64, message *message.Chain) (*Resp[RespDataMessageId], error) {
	return returnAsType[RespDataMessageId](s.SendRaw(ActionSendPrivateMsg, map[string]any{
		"user_id": userId,
		"message": message,
	}))
}

func (s *Sender) SendGroupMsgString(groupId int64, message string, autoEscape bool) (*Resp[RespDataMessageId], error) {
	return returnAsType[RespDataMessageId](s.SendRaw(ActionSendGroupMsg, map[string]any{
		"group_id":    groupId,
		"message":     message,
		"auto_escape": autoEscape,
	}))
}

func (s *Sender) SendGroupMsg(groupId int64, message *message.Chain) (*Resp[RespDataMessageId], error) {
	return returnAsType[RespDataMessageId](s.SendRaw(ActionSendGroupMsg, map[string]any{
		"group_id": groupId,
		"message":  message,
	}))
}

func (s *Sender) SendMsg(msg any, autoEscape bool) (*Resp[RespDataMessageId], error) {
	switch msg := msg.(type) {
	case SendMsgReqParams[string]:
		return returnAsType[RespDataMessageId](s.SendRaw(ActionSendMsg, msg))
	case SendMsgReqParams[*message.Chain]:
		return returnAsType[RespDataMessageId](s.SendRaw(ActionSendMsg, msg))
	}
	return nil, errors.ErrInvalidMessage
}

func (s *Sender) DeleteMsg(messageId qq.MessageId) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionDeleteMsg, map[string]any{
		"message_id": messageId,
	}))
}

func (s *Sender) GetMsg(messageId qq.MessageId) (*Resp[RespDataMessage], error) {
	return returnAsType[RespDataMessage](s.SendRaw(ActionGetMsg, map[string]any{
		"message_id": messageId,
	}))
}

func (s *Sender) GetForwardMsg(id string) (*Resp[RespDataMessageOnly], error) {
	return returnAsType[RespDataMessageOnly](s.SendRaw(ActionGetForwardMsg, map[string]any{
		"id": id,
	}))
}

func (s *Sender) SendLike(userId int64, times int) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSendLike, map[string]any{
		"user_id": userId,
		"times":   times,
	}))
}

// SetGroupKick 将用户踢出群组。groupId 为群组 ID，userId 为要踢的用户 QQ，rejectAddRequest 为是否拒绝此人的加群请求。
func (s *Sender) SetGroupKick(groupId int64, userId int64, rejectAddRequest bool) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupKick, map[string]any{
		"group_id":           groupId,
		"user_id":            userId,
		"reject_add_request": rejectAddRequest,
	}))
}

// SetGroupBan 禁言用户。groupId 为群组 ID，userId 为要禁言的用户 QQ，duration 为禁言时长（秒），0 为解除禁言。
func (s *Sender) SetGroupBan(groupId int64, userId int64, duration int) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupBan, map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"duration": duration,
	}))
}

func (s *Sender) SetGroupAnonymousBan(groupId int64, anonymous qq.AnonymousData, duration int) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupAnonymousBan, map[string]any{
		"group_id":  groupId,
		"anonymous": anonymous,
		"duration":  duration,
	}))
}

func (s *Sender) SetGroupWholeBan(groupId int64, enable bool) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupWholeBan, map[string]any{
		"group_id": groupId,
		"enable":   enable,
	}))
}

func (s *Sender) SetGroupAdmin(groupId int64, userId int64, enable bool) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupAdmin, map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"enable":   enable,
	}))
}

func (s *Sender) SetGroupAnonymous(groupId int64, enable bool) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupAnonymous, map[string]any{
		"group_id": groupId,
		"enable":   enable,
	}))
}

func (s *Sender) SetGroupCard(groupId int64, userId int64, card string) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupCard, map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"card":     card,
	}))
}

func (s *Sender) SetGroupName(groupId int64, name string) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupName, map[string]any{
		"group_id": groupId,
		"name":     name,
	}))
}

func (s *Sender) SetGroupLeave(groupId int64, isDismiss bool) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupLeave, map[string]any{
		"group_id":   groupId,
		"is_dismiss": isDismiss,
	}))
}

func (s *Sender) SetGroupSpecialTitle(groupId int64, userId int64, title string, duration int) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupSpecialTitle, map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"title":    title,
		"duration": duration,
	}))
}

func (s *Sender) SetFriendAddRequest(flag string, approve bool, remark string) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetFriendAddRequest, map[string]any{
		"flag":    flag,
		"approve": approve,
		"remark":  remark,
	}))
}

func (s *Sender) SetGroupAddRequest(flag string, subType string, approve bool, reason string) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetGroupAddRequest, map[string]any{
		"flag":     flag,
		"sub_type": subType,
		"approve":  approve,
		"reason":   reason,
	}))
}

func (s *Sender) GetLoginInfo() (*Resp[RespDataLoginInfo], error) {
	return returnAsType[RespDataLoginInfo](s.SendRaw(ActionGetLoginInfo, nil))
}

func (s *Sender) GetStrangerInfo(userId int64, noCache bool) (*Resp[RespDataStrangerInfo], error) {
	return returnAsType[RespDataStrangerInfo](s.SendRaw(ActionGetStrangerInfo, map[string]any{
		"user_id":  userId,
		"no_cache": noCache,
	}))
}

func (s *Sender) GetFriendList() (*Resp[RespDataFriendList], error) {
	return returnAsType[RespDataFriendList](s.SendRaw(ActionGetFriendList, nil))
}

func (s *Sender) GetGroupList() (*Resp[RespDataGroupList], error) {
	return returnAsType[RespDataGroupList](s.SendRaw(ActionGetGroupList, nil))
}

func (s *Sender) GetGroupInfo(groupId int64, noCache bool) (*Resp[RespDataGroupInfo], error) {
	return returnAsType[RespDataGroupInfo](s.SendRaw(ActionGetGroupInfo, map[string]any{
		"group_id": groupId,
		"no_cache": noCache,
	}))
}

func (s *Sender) GetGroupMemberInfo(groupId int64, userId int64, noCache bool) (*Resp[RespDataGroupMemberInfo], error) {
	return returnAsType[RespDataGroupMemberInfo](s.SendRaw(ActionGetGroupMemberInfo, map[string]any{
		"group_id": groupId,
		"user_id":  userId,
		"no_cache": noCache,
	}))
}

func (s *Sender) GetGroupMemberList(groupId int64) (*Resp[RespDataGroupMemberList], error) {
	return returnAsType[RespDataGroupMemberList](s.SendRaw(ActionGetGroupMemberList, map[string]any{
		"group_id": groupId,
	}))
}

func (s *Sender) GetGroupHonorInfo(groupId int64) (*Resp[RespDataGroupHonorInfo], error) {
	return returnAsType[RespDataGroupHonorInfo](s.SendRaw(ActionGetGroupHonorInfo, map[string]any{
		"group_id": groupId,
		"type":     "all",
	}))
}

func (s *Sender) GetCookies(domain string) (*Resp[RespDataCookies], error) {
	return returnAsType[RespDataCookies](s.SendRaw(ActionGetCookies, map[string]any{
		"domain": domain,
	}))
}

func (s *Sender) GetCsrfToken() (*Resp[RespDataCsrfToken], error) {
	return returnAsType[RespDataCsrfToken](s.SendRaw(ActionGetCsrfToken, nil))
}

func (s *Sender) GetCredentials(domain string) (*Resp[RespDataCredentials], error) {
	return returnAsType[RespDataCredentials](s.SendRaw(ActionGetCredentials, map[string]any{
		"domain": domain,
	}))
}

func (s *Sender) GetRecord(file string, outFormat string) (*Resp[RespDataFile], error) {
	return returnAsType[RespDataFile](s.SendRaw(ActionGetRecord, map[string]any{
		"file":       file,
		"out_format": outFormat,
	}))
}

func (s *Sender) GetImage(file string) (*Resp[RespDataFile], error) {
	return returnAsType[RespDataFile](s.SendRaw(ActionGetImage, map[string]any{
		"file": file,
	}))
}

func (s *Sender) CanSendImage() (*Resp[RespDataYesOrNo], error) {
	return returnAsType[RespDataYesOrNo](s.SendRaw(ActionCanSendImage, nil))
}

func (s *Sender) CanSendRecord() (*Resp[RespDataYesOrNo], error) {
	return returnAsType[RespDataYesOrNo](s.SendRaw(ActionCanSendRecord, nil))
}

func (s *Sender) GetStatus() (*Resp[ServerStatus], error) {
	return returnAsType[ServerStatus](s.SendRaw(ActionGetStatus, nil))
}

func (s *Sender) GetVersionInfo() (*Resp[RespDataVersionInfo], error) {
	return returnAsType[RespDataVersionInfo](s.SendRaw(ActionGetVersionInfo, nil))
}

func (s *Sender) SetRestart() (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionSetRestart, nil))
}

func (s *Sender) CleanCache() (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionCleanCache, nil))
}

func (s *Sender) QuickOp(context any, operation any) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionHandleQuickOperation, map[string]any{
		"context":   context,
		"operation": operation,
	}))
}

func returnAsType[T any](r IResp, err error) (*Resp[T], error) {
	return GetRespAs[T](r), err
}
