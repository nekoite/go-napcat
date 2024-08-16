package gonapcat

import (
	"fmt"

	"github.com/nekoite/go-napcat/api"
	"github.com/nekoite/go-napcat/config"
	"github.com/nekoite/go-napcat/event"
	"github.com/nekoite/go-napcat/message"
	"github.com/nekoite/go-napcat/utils"
	"github.com/nekoite/go-napcat/ws"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BotLogger struct {
	logger *zap.Logger
}

type Bot struct {
	conn       *ws.Client
	dispatcher *event.Dispatcher
	sender     *api.Sender

	Logger *BotLogger
}

func NewBot(cfg *config.BotConfig) *Bot {
	logger := zap.L().Named(fmt.Sprint(cfg.Id))
	bot := &Bot{
		dispatcher: event.NewDispatcher(cfg.UseGoroutine),

		Logger: &BotLogger{logger: logger},
	}
	bot.conn = ws.NewConn(logger, cfg, bot.onRecvWsMsg)
	bot.sender = api.NewSender(bot.conn)
	return bot
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

func (b *BotLogger) Sync() {
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

func (b *Bot) Start() {
	b.conn.Start()
}

func (b *Bot) Close() {
	b.conn.Close()
	b.Logger.Sync()
}

func (b *Bot) SendRawString(msg string) {
	b.conn.Send([]byte(msg))
}

func (b *Bot) SendRaw(action api.Action, params map[string]any) (api.IResp, error) {
	return b.sender.SendRaw(action, params)
}

func (b *Bot) SendPrivateMsgString(userId int64, message string, autoEscape bool) (*api.Resp[api.RespDataMessageId], error) {
	return b.sender.SendPrivateMsgString(userId, message, autoEscape)
}

func (b *Bot) SendPrivateMsg(userId int64, message message.Chain, autoEscape bool) (*api.Resp[api.RespDataMessageId], error) {
	return b.sender.SendPrivateMsg(userId, message, autoEscape)
}

func (b *Bot) SendGroupMsgString(groupId int64, message string, autoEscape bool) (*api.Resp[api.RespDataMessageId], error) {
	return b.sender.SendGroupMsgString(groupId, message, autoEscape)
}

func (b *Bot) SendGroupMsg(groupId int64, message message.Chain, autoEscape bool) (*api.Resp[api.RespDataMessageId], error) {
	return b.sender.SendGroupMsg(groupId, message, autoEscape)
}

func (b *Bot) SendMsg(msg any, autoEscape bool) (*api.Resp[api.RespDataMessageId], error) {
	return b.sender.SendMsg(msg, autoEscape)
}

func (b *Bot) DeleteMsg(messageId int64) (*api.Resp[utils.Void], error) {
	return b.sender.DeleteMsg(messageId)
}

func (b *Bot) onRecvWsMsg(msg []byte) {
	if utils.IsRawMessageApiResp(msg) {
		err := b.sender.HandleApiResp(msg)
		if err != nil {
			b.Logger.Error("handle api resp", zap.Error(err))
		}
		return
	}
	e, err := event.ParseEvent(msg)
	if err != nil {
		b.Logger.Error("parse event", zap.Error(err))
		return
	}
	b.Logger.Debug("received event", zap.Any("event", e))
	b.dispatcher.Dispatch(e)
}
