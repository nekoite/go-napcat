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

type BotLogger struct {
	logger *zap.Logger
}

type Bot struct {
	cfg        *config.BotConfig
	conn       *ws.Client
	dispatcher *event.Dispatcher
	api        *api.Sender

	Logger *BotLogger
}

func NewBot(cfg *config.BotConfig) *Bot {
	logger := zap.L().Named(fmt.Sprint(cfg.Id))
	bot := &Bot{
		cfg:        cfg,
		dispatcher: event.NewDispatcher(logger, cfg.UseGoroutine),

		Logger: &BotLogger{logger: logger},
	}
	bot.conn = ws.NewConn(logger, cfg, bot.onRecvWsMsg)
	bot.api = api.NewSender(logger, bot.conn, cfg.ApiTimeout)
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

func (b *Bot) RegisterCommand(c event.ICommand) {
	b.dispatcher.RegisterCommand(c)
}

func (b *Bot) Start() {
	b.conn.Start()
}

func (b *Bot) Close() {
	b.conn.Close()
	b.Logger.Sync()
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

func (b *Bot) SendPrivateMsgString(userId int64, message string, autoEscape bool) (*api.Resp[api.RespDataMessageId], error) {
	return b.api.SendPrivateMsgString(userId, message, autoEscape)
}

func (b *Bot) SendPrivateMsg(userId int64, message *message.Chain) (*api.Resp[api.RespDataMessageId], error) {
	return b.api.SendPrivateMsg(userId, message)
}

func (b *Bot) SendGroupMsgString(groupId int64, message string, autoEscape bool) (*api.Resp[api.RespDataMessageId], error) {
	return b.api.SendGroupMsgString(groupId, message, autoEscape)
}

func (b *Bot) SendGroupMsg(groupId int64, message *message.Chain) (*api.Resp[api.RespDataMessageId], error) {
	return b.api.SendGroupMsg(groupId, message)
}

func (b *Bot) SendMsg(msg any, autoEscape bool) (*api.Resp[api.RespDataMessageId], error) {
	return b.api.SendMsg(msg, autoEscape)
}

func (b *Bot) DeleteMsg(messageId qq.MessageId) (*api.Resp[utils.Void], error) {
	return b.api.DeleteMsg(messageId)
}

func (b *Bot) GetMsg(messageId qq.MessageId) (*api.Resp[api.RespDataMessage], error) {
	return b.api.GetMsg(messageId)
}

func (b *Bot) GetForwardMsg(id string) (*api.Resp[api.RespDataMessageOnly], error) {
	return b.api.GetForwardMsg(id)
}

func (b *Bot) SendLike(userId int64, times int) (*api.Resp[utils.Void], error) {
	return b.api.SendLike(userId, times)
}

func (b *Bot) onRecvWsMsg(msg []byte) {
	if utils.IsRawMessageApiResp(msg) {
		err := utils.TimedFunc(func() error {
			return b.api.HandleApiResp(msg)
		}, func(t time.Duration) {
			b.Logger.Debug("handle api resp duration", zap.Duration("ms", t))
		})
		if err != nil {
			b.Logger.Error("handle api resp", zap.Error(err))
		}
		return
	}
	e, err := event.ParseEvent(msg, b.api)
	if err != nil {
		if !errors2.Is(err, errors.ErrGoNapcat) {
			b.Logger.Error("parse event", zap.Error(err))
			return
		}
		b.Logger.Warn("parse event", zap.Error(err))
	}
	b.Logger.Debug("received event", zap.Any("event", e))
	utils.TimedAction(func() {
		b.dispatcher.Dispatch(e)
	}, func(t time.Duration) {
		if b.cfg.UseGoroutine {
			return
		}
		if t > 5*time.Second {
			b.Logger.Warn("long event execution duration", zap.Duration("ms", t))
		} else {
			b.Logger.Debug("event execution duration", zap.Duration("ms", t))
		}
	})
}
