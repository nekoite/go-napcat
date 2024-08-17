package api

import (
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/errors"
	"github.com/nekoite/go-napcat/message"
	"github.com/nekoite/go-napcat/utils"
	"github.com/nekoite/go-napcat/ws"
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
	conn   *ws.Client
	sendId atomic.Int64
	reqMap sync.Map
}

type Req struct {
	id     int64
	resp   chan map[string]any
	action Action
}

type apiReq struct {
	Action Action `json:"action"`
	Params any    `json:"params"`
	Echo   string `json:"echo"`
}

type SendMsgReqParams[T message.SendableMessage] struct {
	MessageType string `json:"message_type"`
	UserId      int64  `json:"user_id,omitempty"`
	GroupId     int64  `json:"group_id,omitempty"`
	Message     T      `json:"message"`
	AutoEscape  bool   `json:"auto_escape,omitempty"`
}

func NewSender(conn *ws.Client) *Sender {
	return &Sender{
		sendId: atomic.Int64{},
		conn:   conn,
	}
}

func (s *Sender) NewReq(action Action) *Req {
	return &Req{
		id:     s.sendId.Add(1),
		resp:   make(chan map[string]any, 1),
		action: action,
	}
}

func (s *Sender) HandleApiResp(resp []byte) error {
	r := make(map[string]any)
	if err := json.Unmarshal(resp, &r); err != nil {
		return err
	}
	id, err := strconv.Atoi(r["echo"].(string))
	if err != nil {
		return err
	}
	if req, ok := s.reqMap.Load(int64(id)); ok {
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
	resp := <-req.resp
	s.reqMap.Delete(req.id)
	return ParseResp(req.action, resp)
}

func (s *Sender) SendPrivateMsgString(userId int64, message string, autoEscape bool) (*Resp[RespDataMessageId], error) {
	return returnAsType[RespDataMessageId](s.SendRaw(ActionSendPrivateMsg, map[string]any{
		"user_id":     userId,
		"message":     message,
		"auto_escape": autoEscape,
	}))
}

func (s *Sender) SendPrivateMsg(userId int64, message message.Chain, autoEscape bool) (*Resp[RespDataMessageId], error) {
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

func (s *Sender) SendGroupMsg(groupId int64, message message.Chain, autoEscape bool) (*Resp[RespDataMessageId], error) {
	return returnAsType[RespDataMessageId](s.SendRaw(ActionSendGroupMsg, map[string]any{
		"group_id": groupId,
		"message":  message,
	}))
}

func (s *Sender) SendMsg(msg any, autoEscape bool) (*Resp[RespDataMessageId], error) {
	switch msg := msg.(type) {
	case SendMsgReqParams[string]:
		return returnAsType[RespDataMessageId](s.SendRaw(ActionSendMsg, msg))
	case SendMsgReqParams[message.Chain]:
		return returnAsType[RespDataMessageId](s.SendRaw(ActionSendMsg, msg))
	}
	return nil, errors.ErrInvalidMessage
}

func (s *Sender) DeleteMsg(messageId int64) (*Resp[utils.Void], error) {
	return returnAsType[utils.Void](s.SendRaw(ActionDeleteMsg, map[string]any{
		"message_id": messageId,
	}))
}

func returnAsType[T any](r IResp, err error) (*Resp[T], error) {
	return GetRespAs[T](r), err
}
