package napcat

import (
	gonapcat "github.com/nekoite/go-napcat"
	"github.com/nekoite/go-napcat/api"
	"github.com/nekoite/go-napcat/qq"
	"github.com/nekoite/go-napcat/utils"
)

const (
	ActionSetQQAvatar            api.Action = "set_qq_avatar"
	ActionGetGroupSystemMsg      api.Action = "get_group_system_msg"
	ActionGetFile                api.Action = "get_file"
	ActionDownloadFile           api.Action = "download_file"
	ActionForwardFriendSingleMsg api.Action = "forward_friend_single_msg"
	ActionForwardGroupSingleMsg  api.Action = "forward_group_single_msg"
	ActionSetMsgEmojiLike        api.Action = "set_msg_emoji_like"
	ActionMarkPrivateMsgAsRead   api.Action = "mark_private_msg_as_read"
	ActionMarkGroupMsgAsRead     api.Action = "mark_group_msg_as_read"
	ActionGetRobotUinRange       api.Action = "get_robot_uin_range"
	ActionSetOnlineStatus        api.Action = "set_online_status"
	ActionGetFriendsWithCategory api.Action = "get_friends_with_category"
	ActionGetGroupFileCount      api.Action = "get_group_file_count"
	ActionGetGroupFileList       api.Action = "get_group_file_list"
	ActionSetGroupFileFolder     api.Action = "set_group_file_folder"
	ActionDelGroupFile           api.Action = "del_group_file"
	ActionDelGroupFileFolder     api.Action = "del_group_file_folder"
	ActionTranslateEn2Zh         api.Action = "translate_en2zh"
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

type RespDataGetFile struct {
	File     string `json:"file"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Base64   string `json:"base64"`
}

type RespDataGroupFileCount struct {
	Count int `json:"count"`
}

type RespDataGroupFileList struct {
	FileList []map[string]any `json:"FileList"`
}

type AnyResult = map[string]any

var (
	Extension = api.NewExtension("napcat").WithActions(map[api.Action]api.GetNewResultFunc{
		ActionSetQQAvatar:            func() any { return new(api.Resp[utils.Void]) },
		ActionGetGroupSystemMsg:      func() any { return new(api.Resp[RespDataGetGroupSystemMsg]) },
		ActionGetFile:                func() any { return new(api.Resp[RespDataGetFile]) },
		ActionForwardFriendSingleMsg: func() any { return new(api.Resp[utils.Void]) },
		ActionForwardGroupSingleMsg:  func() any { return new(api.Resp[utils.Void]) },
		ActionSetMsgEmojiLike:        func() any { return new(api.Resp[utils.Void]) },
		ActionMarkPrivateMsgAsRead:   func() any { return new(api.Resp[utils.Void]) },
		ActionMarkGroupMsgAsRead:     func() any { return new(api.Resp[utils.Void]) },
		ActionGetRobotUinRange:       func() any { return new(api.Resp[AnyResult]) },
		ActionSetOnlineStatus:        func() any { return new(api.Resp[utils.Void]) },
		ActionGetFriendsWithCategory: func() any { return new(api.Resp[AnyResult]) },
		ActionGetGroupFileCount:      func() any { return new(api.Resp[RespDataGroupFileCount]) },
		ActionGetGroupFileList:       func() any { return new(api.Resp[RespDataGroupFileList]) },
		ActionSetGroupFileFolder:     func() any { return new(api.Resp[AnyResult]) },
		ActionDelGroupFile:           func() any { return new(api.Resp[AnyResult]) },
		ActionDelGroupFileFolder:     func() any { return new(api.Resp[AnyResult]) },
		ActionTranslateEn2Zh:         func() any { return new(api.Resp[[]string]) },
	})
)

func returnAsType[T any](r api.IResp, err error) (*api.Resp[T], error) {
	return api.GetRespAs[T](r), err
}

func SetQQAvatar(bot *gonapcat.Bot, file string) (*api.Resp[utils.Void], error) {
	return returnAsType[utils.Void](bot.SendRaw(ActionSetQQAvatar, map[string]any{
		"file": file,
	}))
}

func GetGroupSystemMsg(bot *gonapcat.Bot, group int64) (*api.Resp[RespDataGetGroupSystemMsg], error) {
	return returnAsType[RespDataGetGroupSystemMsg](bot.SendRaw(ActionGetGroupSystemMsg, nil))
}

func GetFile(bot *gonapcat.Bot, fileId string) (*api.Resp[RespDataGetFile], error) {
	return returnAsType[RespDataGetFile](bot.SendRaw(ActionGetFile, map[string]any{
		"file_id": fileId,
	}))
}

func ForwardFriendSingleMsg(bot *gonapcat.Bot, userId int64, messageId qq.MessageId) (*api.Resp[utils.Void], error) {
	return returnAsType[utils.Void](bot.SendRaw(ActionForwardFriendSingleMsg, map[string]any{
		"user_id":    userId,
		"message_id": messageId,
	}))
}

func ForwardGroupSingleMsg(bot *gonapcat.Bot, groupId int64, messageId qq.MessageId) (*api.Resp[utils.Void], error) {
	return returnAsType[utils.Void](bot.SendRaw(ActionForwardGroupSingleMsg, map[string]any{
		"group_id":   groupId,
		"message_id": messageId,
	}))
}

// SetMsgEmojiLike 设置消息点赞。
//
// emojiId 列表：https://bot.q.qq.com/wiki/develop/api-v2/openapi/emoji/model.html#EmojiType
func SetMsgEmojiLike(bot *gonapcat.Bot, messageId qq.MessageId, emojiId int) (*api.Resp[utils.Void], error) {
	return returnAsType[utils.Void](bot.SendRaw(ActionSetMsgEmojiLike, map[string]any{
		"message_id": messageId,
		"emoji_id":   emojiId,
	}))
}

func MarkPrivateMsgAsRead(bot *gonapcat.Bot, userId int64) (*api.Resp[utils.Void], error) {
	return returnAsType[utils.Void](bot.SendRaw(ActionMarkPrivateMsgAsRead, map[string]any{
		"user_id": userId,
	}))
}

func MarkGroupMsgAsRead(bot *gonapcat.Bot, groupId int64) (*api.Resp[utils.Void], error) {
	return returnAsType[utils.Void](bot.SendRaw(ActionMarkGroupMsgAsRead, map[string]any{
		"group_id": groupId,
	}))
}

func GetRobotUinRange(bot *gonapcat.Bot) (*api.Resp[AnyResult], error) {
	return returnAsType[map[string]any](bot.SendRaw(ActionGetRobotUinRange, nil))
}

// SetOnlineStatus 设置在线状态。
//
// 在线状态列表：https://napneko.github.io/zh-CN/develop/status_list
func SetOnlineStatus(bot *gonapcat.Bot, status, extStatus, batteryStatus int) (*api.Resp[utils.Void], error) {
	return returnAsType[utils.Void](bot.SendRaw(ActionSetOnlineStatus, map[string]any{
		"status":        status,
		"extStatus":     extStatus,
		"batteryStatus": batteryStatus,
	}))
}

func GetFriendsWithCategory(bot *gonapcat.Bot) (*api.Resp[AnyResult], error) {
	return returnAsType[map[string]any](bot.SendRaw(ActionGetFriendsWithCategory, nil))
}

func GetGroupFileCount(bot *gonapcat.Bot, groupId int64) (*api.Resp[RespDataGroupFileCount], error) {
	return returnAsType[RespDataGroupFileCount](bot.SendRaw(ActionGetGroupFileCount, map[string]any{
		"group_id": groupId,
	}))
}

func GetGroupFileList(bot *gonapcat.Bot, groupId int64, startIndex, fileCount int) (*api.Resp[RespDataGroupFileList], error) {
	return returnAsType[RespDataGroupFileList](bot.SendRaw(ActionGetGroupFileList, map[string]any{
		"group_id":    groupId,
		"start_index": startIndex,
		"file_count":  fileCount,
	}))
}

func SetGroupFileFolder(bot *gonapcat.Bot, groupId int64, folderName string) (*api.Resp[AnyResult], error) {
	return returnAsType[AnyResult](bot.SendRaw(ActionSetGroupFileFolder, map[string]any{
		"group_id":    groupId,
		"folder_name": folderName,
	}))
}

func DelGroupFile(bot *gonapcat.Bot, groupId int64, fileId string) (*api.Resp[AnyResult], error) {
	return returnAsType[AnyResult](bot.SendRaw(ActionDelGroupFile, map[string]any{
		"group_id": groupId,
		"file_id":  fileId,
	}))
}

func DelGroupFileFolder(bot *gonapcat.Bot, groupId int64, folderId string) (*api.Resp[AnyResult], error) {
	return returnAsType[AnyResult](bot.SendRaw(ActionDelGroupFileFolder, map[string]any{
		"group_id":  groupId,
		"folder_id": folderId,
	}))
}

func TranslateEn2Zh(bot *gonapcat.Bot, words []string) (*api.Resp[[]string], error) {
	return returnAsType[[]string](bot.SendRaw(ActionTranslateEn2Zh, map[string]any{
		"words": words,
	}))
}
