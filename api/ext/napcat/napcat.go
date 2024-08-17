package napcat

import "github.com/nekoite/go-napcat/api"

const (
	ActionSetQQAvatar       api.Action = "set_qq_avatar"
	ActionGetGroupSystemMsg api.Action = "get_group_system_msg"
	ActionGetFile           api.Action = "get_file"
	ActionDownloadFile      api.Action = "download_file"
)

var (
	Extension = api.NewExtension("napcat").WithActions(map[api.Action]api.GetNewResultFunc{
		ActionSetQQAvatar: func() any { return new(api.Resp[map[string]any]) },
	})
)
