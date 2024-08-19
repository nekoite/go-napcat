package napcat

import (
	"github.com/nekoite/go-napcat/api"
	"github.com/nekoite/go-napcat/utils"
)

const (
	ActionSetQQAvatar       api.Action = "set_qq_avatar"
	ActionGetGroupSystemMsg api.Action = "get_group_system_msg"
	ActionGetFile           api.Action = "get_file"
	ActionDownloadFile      api.Action = "download_file"
)

type RespDataGetGroupSystemMsg struct {
	InvitedRequests []struct {
		GroupId     int64  `json:"group_id"`
		GroupName   string `json:"group_name"`
		RequestId   string `json:"request_id"`
		InvitorUin  int64  `json:"invitor_uin"`
		InvitorNick string `json:"invitor_nick"`
		Checked     bool   `json:"checked"`
		Actor       int64  `json:"actor"`
	} `json:"InvitedRequest"`
	JoinRequests []struct {
		GroupId       int64  `json:"group_id"`
		GroupName     string `json:"group_name"`
		RequestId     string `json:"request_id"`
		RequesterUin  int64  `json:"requester_uin"`
		RequesterNick string `json:"requester_nick"`
		Checked       bool   `json:"checked"`
		Actor         int64  `json:"actor"`
	} `json:"join_requests"`
}

var (
	Extension = api.NewExtension("napcat").WithActions(map[api.Action]api.GetNewResultFunc{
		ActionSetQQAvatar:       func() any { return new(api.Resp[utils.Void]) },
		ActionGetGroupSystemMsg: func() any { return new(RespDataGetGroupSystemMsg) },
	})
)
