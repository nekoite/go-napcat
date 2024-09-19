package gonapcat

import (
	errors2 "errors"
	"fmt"
	"time"

	"github.com/nekoite/go-napcat/api"
	"github.com/nekoite/go-napcat/config"
	"github.com/nekoite/go-napcat/errors"
	"github.com/nekoite/go-napcat/event"
	"github.com/nekoite/go-napcat/message"
	"github.com/nekoite/go-napcat/qq"
	"github.com/nekoite/go-napcat/utils"
	"github.com/nekoite/go-napcat/ws"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type botInfoStore struct {
	botUser qq.BasicUser
}

type BotLogger struct {
	logger *zap.Logger
}

type Bot struct {
	botInfoStore
	id         qq.UserId
	cfg        *config.BotConfig
	conn       *ws.Client
	dispatcher *event.Dispatcher
	api        *api.Sender

	logger *BotLogger
}

func NewBot(cfg *config.BotConfig) (*Bot, error) {
	logger := zap.L().Named(fmt.Sprint(cfg.Id))
	bot := &Bot{
		id:         qq.UserId(cfg.Id),
		cfg:        cfg,
		dispatcher: event.NewDispatcher(logger, cfg.UseGoroutine),

		logger: &BotLogger{logger: logger},
	}
	var err error
	bot.conn, err = ws.NewConn(logger, cfg, bot.onRecvWsMsg)
	if err != nil {
		return nil, err
	}
	bot.api = api.NewSender(logger, bot.conn, cfg.ApiTimeout)
	return bot, nil
}

func (b *Bot) Logger() *BotLogger {
	return b.logger
}

func (b *BotLogger) Log(level zapcore.Level, msg string, fields ...zap.Field) {
	b.logger.Log(level, msg, fields...)
}

func (b *BotLogger) Info(msg string, fields ...zap.Field) {
	b.logger.Info(msg, fields...)
}

func (b *BotLogger) Debug(msg string, fields ...zap.Field) {
	b.logger.Debug(msg, fields...)
}

func (b *BotLogger) Warn(msg string, fields ...zap.Field) {
	b.logger.Warn(msg, fields...)
}

func (b *BotLogger) Error(msg string, fields ...zap.Field) {
	b.logger.Error(msg, fields...)
}

func (b *BotLogger) Fatal(msg string, fields ...zap.Field) {
	b.logger.Fatal(msg, fields...)
}

func (b *BotLogger) Panic(msg string, fields ...zap.Field) {
	b.logger.Panic(msg, fields...)
}

func (b *BotLogger) SyncLogger() {
	b.logger.Sync()
}

func (b *Bot) RegisterHandler(h event.Handler) {
	b.dispatcher.RegisterHandlerAllTypes(h)
}

func (b *Bot) RegisterHandlerGroupMessage(h event.Handler) {
	b.dispatcher.RegisterHandlerGroupMessage(h)
}

func (b *Bot) RegisterHandlerPrivateMessage(h event.Handler) {
	b.dispatcher.RegisterHandlerPrivateMessage(h)
}

func (b *Bot) RegisterHandlerNotice(h event.Handler) {
	b.dispatcher.RegisterHandlerNotice(h)
}

func (b *Bot) RegisterHandlerMeta(h event.Handler) {
	b.dispatcher.RegisterHandlerMeta(h)
}

func (b *Bot) RegisterHandlerRequest(h event.Handler) {
	b.dispatcher.RegisterHandlerRequest(h)
}

func (b *Bot) RegisterCommand(c event.ICommand) {
	b.dispatcher.RegisterCommand(c)
}

func (b *Bot) SetGlobalCommandPrefix(prefix string) {
	b.dispatcher.SetGlobalCommandPrefix(prefix)
}

func (b *Bot) Start() error {
	b.conn.Start()
	err := b.initializeBotInfo()
	if err != nil {
		b.logger.Error("error initializing bot info, shutting down", zap.Error(err))
		b.Close()
		return err
	}
	return nil
}

func (b *Bot) Close() {
	b.conn.Close()
	b.logger.SyncLogger()
}

func (b *Bot) Api() *api.Sender {
	return b.api
}

func (b *Bot) SendRawString(msg string) {
	b.conn.Send([]byte(msg))
}

func (b *Bot) SendRaw(action api.Action, params map[string]any) (api.IResp, error) {
	return b.api.SendRaw(action, params)
}

func (b *Bot) SendPrivateMsgString(userId qq.UserId, message string, autoEscape bool) (qq.MessageId, error) {
	return extractRespMessageId(b.api.SendPrivateMsgString(userId, message, autoEscape))
}

func (b *Bot) SendPrivateMsg(userId qq.UserId, message *message.Chain) (qq.MessageId, error) {
	return extractRespMessageId(b.api.SendPrivateMsg(userId, message))
}

func (b *Bot) SendGroupMsgString(groupId qq.GroupId, message string, autoEscape bool) (qq.MessageId, error) {
	return extractRespMessageId(b.api.SendGroupMsgString(groupId, message, autoEscape))
}

func (b *Bot) SendGroupMsg(groupId qq.GroupId, message *message.Chain) (qq.MessageId, error) {
	return extractRespMessageId(b.api.SendGroupMsg(groupId, message))
}

func (b *Bot) SendMsg(msg api.SendMsgReqParams, autoEscape bool) (qq.MessageId, error) {
	return extractRespMessageId(b.api.SendMsg(msg, autoEscape))
}

func (b *Bot) DeleteMsg(messageId qq.MessageId) error {
	_, err := b.api.DeleteMsg(messageId)
	return err
}

func (b *Bot) GetMsg(messageId qq.MessageId) (*api.RespDataMessage, error) {
	resp, err := b.api.GetMsg(messageId)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (b *Bot) GetForwardMsg(id string) (*message.Chain, error) {
	resp, err := b.api.GetForwardMsg(id)
	if err != nil {
		return nil, err
	}
	return &resp.Data.Message, nil
}

func (b *Bot) SendLike(userId qq.UserId, times int) error {
	_, err := b.api.SendLike(userId, times)
	return err
}

func (b *Bot) SetGroupKick(groupId qq.GroupId, userId qq.UserId, rejectAddRequest bool) error {
	_, err := b.api.SetGroupKick(groupId, userId, rejectAddRequest)
	return err
}

func (b *Bot) SetGroupBan(groupId qq.GroupId, userId qq.UserId, duration int) error {
	_, err := b.api.SetGroupBan(groupId, userId, duration)
	return err
}

func (b *Bot) SetGroupAnonymousBan(groupId qq.GroupId, anonymous *qq.AnonymousData, duration int) error {
	_, err := b.api.SetGroupAnonymousBan(groupId, anonymous, duration)
	return err
}

func (b *Bot) SetGroupWholeBan(groupId qq.GroupId, enable bool) error {
	_, err := b.api.SetGroupWholeBan(groupId, enable)
	return err
}

func (b *Bot) SetGroupAdmin(groupId qq.GroupId, userId qq.UserId, enable bool) error {
	_, err := b.api.SetGroupAdmin(groupId, userId, enable)
	return err
}

func (b *Bot) SetGroupAnonymous(groupId qq.GroupId, enable bool) error {
	_, err := b.api.SetGroupAnonymous(groupId, enable)
	return err
}

func (b *Bot) SetGroupCard(groupId qq.GroupId, userId qq.UserId, card string) error {
	_, err := b.api.SetGroupCard(groupId, userId, card)
	return err
}

func (b *Bot) SetGroupName(groupId qq.GroupId, name string) error {
	_, err := b.api.SetGroupName(groupId, name)
	return err
}

func (b *Bot) LeaveGroup(groupId qq.GroupId, isDismiss bool) error {
	_, err := b.api.LeaveGroup(groupId, isDismiss)
	return err
}

func (b *Bot) SetGroupSpecialTitle(groupId qq.GroupId, userId qq.UserId, title string, duration int) error {
	_, err := b.api.SetGroupSpecialTitle(groupId, userId, title, duration)
	return err
}

func (b *Bot) SetFriendAddRequest(flag string, approve bool, remark string) error {
	_, err := b.api.SetFriendAddRequest(flag, approve, remark)
	return err
}

func (b *Bot) SetGroupAddRequest(flag string, subType string, approve bool, reason string) error {
	_, err := b.api.SetGroupAddRequest(flag, subType, approve, reason)
	return err
}

func (b *Bot) GetLoginInfo() (*qq.BasicUser, error) {
	resp, err := b.api.GetLoginInfo()
	if err != nil {
		return nil, err
	}
	return (*qq.BasicUser)(&resp.Data), nil
}

func (b *Bot) GetStrangerInfo(userId qq.UserId, noCache bool) (*qq.User, error) {
	resp, err := b.api.GetStrangerInfo(userId, noCache)
	if err != nil {
		return nil, err
	}
	return (*qq.User)(&resp.Data), nil
}

func (b *Bot) GetFriendList() ([]qq.BasicFriend, error) {
	resp, err := b.api.GetFriendList()
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (b *Bot) GetGroupInfo(groupId qq.GroupId, noCache bool) (*qq.Group, error) {
	resp, err := b.api.GetGroupInfo(groupId, noCache)
	if err != nil {
		return nil, err
	}
	return (*qq.Group)(&resp.Data), nil
}

func (b *Bot) GetGroupList() ([]qq.Group, error) {
	resp, err := b.api.GetGroupList()
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func (b *Bot) GetGroupMemberInfo(groupId qq.GroupId, userId qq.UserId, noCache bool) (*qq.GroupUser, error) {
	resp, err := b.api.GetGroupMemberInfo(groupId, userId, noCache)
	if err != nil {
		return nil, err
	}
	return (*qq.GroupUser)(&resp.Data), nil
}

func (b *Bot) GetGroupMemberList(groupId qq.GroupId) (*api.Resp[api.RespDataGroupMemberList], error) {
	return b.api.GetGroupMemberList(groupId)
}

func (b *Bot) GetGroupHonorInfo(groupId qq.GroupId) (*api.Resp[api.RespDataGroupHonorInfo], error) {
	return b.api.GetGroupHonorInfo(groupId)
}

func (b *Bot) Id() qq.UserId {
	return b.botUser.UserId
}

func (b *Bot) Nickname() string {
	return b.botUser.Nickname
}

func (b *Bot) initializeBotInfo() error {
	botUser, err := b.GetLoginInfo()
	if err != nil {
		return err
	}
	if botUser.UserId != b.id {
		b.logger.Warn("bot user id mismatch", zap.Int64("expected", int64(b.id)), zap.Int64("actual", int64(botUser.UserId)))
		b.id = botUser.UserId
	}
	b.botUser = *botUser
	return nil
}

func (b *Bot) onRecvWsMsg(msg []byte) {
	if utils.IsRawMessageApiResp(msg) {
		err := utils.TimedFunc(func() error {
			return b.api.HandleApiResp(msg)
		}, func(t time.Duration) {
			b.logger.Debug("handle api resp duration", zap.Duration("ms", t))
		})
		if err != nil {
			b.logger.Error("handle api resp", zap.Error(err))
		}
		return
	}
	e, err := event.ParseEvent(msg, b.api)
	if err != nil {
		if !errors2.Is(err, errors.ErrGoNapcat) {
			b.logger.Error("parse event", zap.Error(err))
			return
		}
		b.logger.Warn("parse event", zap.Error(err))
	}
	b.logger.Debug("received event", zap.Any("event", e))
	utils.TimedAction(func() {
		b.dispatcher.Dispatch(e)
	}, func(t time.Duration) {
		if b.cfg.UseGoroutine {
			return
		}
		if t > 5*time.Second {
			b.logger.Warn("long event execution duration", zap.Duration("ms", t))
		} else {
			b.logger.Debug("event execution duration", zap.Duration("ms", t))
		}
	})
}

func extractRespMessageId(r *api.Resp[api.RespDataMessageId], err error) (qq.MessageId, error) {
	if err != nil {
		return 0, err
	}
	return r.Data.MessageId, nil
}
