package event

import "github.com/nekoite/go-napcat/errors"

func GetAs[T any](e any) *T {
	t, ok := e.(*T)
	if !ok {
		return nil
	}
	return t
}

func GetAsUnsafe[T any](e any) *T {
	return e.(*T)
}

func GetAsOrError[T any](e any) (*T, error) {
	t, ok := e.(*T)
	if !ok {
		return nil, errors.ErrTypeAssertion
	}
	return t, nil
}

func (e *PrivateMessageEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return e
}

func (e *PrivateMessageEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *PrivateMessageEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *PrivateMessageEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *PrivateMessageEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *PrivateMessageEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *PrivateMessageEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *PrivateMessageEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *PrivateMessageEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *PrivateMessageEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *PrivateMessageEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *PrivateMessageEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *PrivateMessageEvent) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *GroupMessageEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *GroupMessageEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return e
}

func (e *GroupMessageEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *GroupMessageEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *GroupMessageEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *GroupMessageEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *GroupMessageEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *GroupMessageEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *GroupMessageEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *GroupMessageEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *GroupMessageEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *GroupMessageEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *GroupMessageEvent) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *MessageEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *MessageEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *MessageEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *MessageEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *MessageEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *MessageEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *MessageEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *MessageEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *MessageEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *MessageEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *MessageEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *MessageEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *MessageEvent) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *NoticeEventGroupUpload) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEventGroupUpload) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEventGroupUpload) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return e
}

func (e *NoticeEventGroupUpload) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *NoticeEventGroupUpload) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *NoticeEventGroupUpload) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *NoticeEventGroupUpload) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *NoticeEventGroupUpload) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *NoticeEventGroupUpload) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *NoticeEventGroupUpload) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *NoticeEventGroupUpload) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEventGroupUpload) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEventGroupUpload) AsMetaEvent() *MetaEvent {
	return nil
}
func (e *NoticeEventGroupOperation) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEventGroupOperation) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEventGroupOperation) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *NoticeEventGroupOperation) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return e
}

func (e *NoticeEventGroupOperation) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *NoticeEventGroupOperation) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *NoticeEventGroupOperation) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *NoticeEventGroupOperation) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *NoticeEventGroupOperation) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *NoticeEventGroupOperation) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *NoticeEventGroupOperation) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEventGroupOperation) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEventGroupOperation) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *NoticeEventGroupBan) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEventGroupBan) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEventGroupBan) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *NoticeEventGroupBan) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *NoticeEventGroupBan) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return e
}

func (e *NoticeEventGroupBan) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *NoticeEventGroupBan) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *NoticeEventGroupBan) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *NoticeEventGroupBan) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *NoticeEventGroupBan) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *NoticeEventGroupBan) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEventGroupBan) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEventGroupBan) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *NoticeEventGroupRecall) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEventGroupRecall) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEventGroupRecall) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *NoticeEventGroupRecall) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *NoticeEventGroupRecall) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *NoticeEventGroupRecall) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return e
}

func (e *NoticeEventGroupRecall) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *NoticeEventGroupRecall) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *NoticeEventGroupRecall) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *NoticeEventGroupRecall) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *NoticeEventGroupRecall) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEventGroupRecall) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEventGroupRecall) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *NoticeEventFriendRecall) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEventFriendRecall) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEventFriendRecall) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *NoticeEventFriendRecall) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *NoticeEventFriendRecall) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *NoticeEventFriendRecall) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *NoticeEventFriendRecall) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return e
}

func (e *NoticeEventFriendRecall) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *NoticeEventFriendRecall) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *NoticeEventFriendRecall) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *NoticeEventFriendRecall) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEventFriendRecall) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEventFriendRecall) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *GroupNoticeEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *GroupNoticeEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *GroupNoticeEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *GroupNoticeEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *GroupNoticeEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *GroupNoticeEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *GroupNoticeEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *GroupNoticeEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *GroupNoticeEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *GroupNoticeEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *GroupNoticeEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *GroupNoticeEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *GroupNoticeEvent) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *NoticeEventFriendAdd) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEventFriendAdd) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEventFriendAdd) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *NoticeEventFriendAdd) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return e
}

func (e *NoticeEventFriendAdd) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *NoticeEventFriendAdd) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *NoticeEventFriendAdd) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEventFriendAdd) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEventFriendAdd) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *NoticeEventFriendAdd) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *NoticeEventFriendAdd) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *NoticeEventFriendAdd) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *NoticeEventFriendAdd) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *NoticeEventGroupNotify) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEventGroupNotify) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEventGroupNotify) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *NoticeEventGroupNotify) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *NoticeEventGroupNotify) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *NoticeEventGroupNotify) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *NoticeEventGroupNotify) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *NoticeEventGroupNotify) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *NoticeEventGroupNotify) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return e
}

func (e *NoticeEventGroupNotify) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *NoticeEventGroupNotify) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEventGroupNotify) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEventGroupNotify) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *NoticeEventGroupHonor) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEventGroupHonor) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEventGroupHonor) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *NoticeEventGroupHonor) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *NoticeEventGroupHonor) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *NoticeEventGroupHonor) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *NoticeEventGroupHonor) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *NoticeEventGroupHonor) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *NoticeEventGroupHonor) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *NoticeEventGroupHonor) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return e
}

func (e *NoticeEventGroupHonor) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEventGroupHonor) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEventGroupHonor) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *FriendRequestEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *FriendRequestEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *FriendRequestEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *FriendRequestEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *FriendRequestEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *FriendRequestEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *FriendRequestEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *FriendRequestEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *FriendRequestEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *FriendRequestEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *FriendRequestEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return e
}

func (e *FriendRequestEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *FriendRequestEvent) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *GroupRequestEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *GroupRequestEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *GroupRequestEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *GroupRequestEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *GroupRequestEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *GroupRequestEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *GroupRequestEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *GroupRequestEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *GroupRequestEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *GroupRequestEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *GroupRequestEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *GroupRequestEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return e
}

func (e *GroupRequestEvent) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *MetaEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *MetaEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *MetaEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *MetaEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *MetaEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *MetaEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *MetaEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *MetaEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *MetaEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *MetaEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *MetaEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *MetaEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *MetaEvent) AsMetaEvent() *MetaEvent {
	return e
}

func (e *NoticeEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *NoticeEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *NoticeEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *NoticeEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *NoticeEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *NoticeEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *NoticeEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *NoticeEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *NoticeEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *NoticeEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *NoticeEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *NoticeEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *NoticeEvent) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *RequestEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *RequestEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *RequestEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *RequestEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *RequestEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *RequestEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *RequestEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *RequestEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *RequestEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *RequestEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *RequestEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *RequestEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *RequestEvent) AsMetaEvent() *MetaEvent {
	return nil
}

func (e *BaseEvent) AsPrivateMessageEvent() *PrivateMessageEvent {
	return nil
}

func (e *BaseEvent) AsGroupMessageEvent() *GroupMessageEvent {
	return nil
}

func (e *BaseEvent) AsNoticeEventGroupUpload() *NoticeEventGroupUpload {
	return nil
}

func (e *BaseEvent) AsNoticeEventGroupOperation() *NoticeEventGroupOperation {
	return nil
}

func (e *BaseEvent) AsNoticeEventGroupBan() *NoticeEventGroupBan {
	return nil
}

func (e *BaseEvent) AsNoticeEventGroupRecall() *NoticeEventGroupRecall {
	return nil
}

func (e *BaseEvent) AsNoticeEventFriendRecall() *NoticeEventFriendRecall {
	return nil
}

func (e *BaseEvent) AsNoticeEventFriendAdd() *NoticeEventFriendAdd {
	return nil
}

func (e *BaseEvent) AsNoticeEventGroupNotify() *NoticeEventGroupNotify {
	return nil
}

func (e *BaseEvent) AsNoticeEventGroupHonor() *NoticeEventGroupHonor {
	return nil
}

func (e *BaseEvent) AsFriendRequestEvent() *FriendRequestEvent {
	return nil
}

func (e *BaseEvent) AsGroupRequestEvent() *GroupRequestEvent {
	return nil
}

func (e *BaseEvent) AsMetaEvent() *MetaEvent {
	return nil
}
