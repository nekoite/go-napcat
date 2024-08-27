package main

import (
	"os"
	"os/signal"
	"time"

	gonapcat "github.com/nekoite/go-napcat"
	"github.com/nekoite/go-napcat/config"
	"github.com/nekoite/go-napcat/event"
	"github.com/nekoite/go-napcat/message"
	"go.uber.org/zap"
)

func main() {
	gonapcat.Init(config.DefaultLogConfig().WithStderr().WithLevel("debug"))
	bot := gonapcat.NewBot(config.DefaultBotConfig(1341400490, "114514"))
	bot.RegisterHandlerPrivateMessage(func(e event.IEvent) {
		bot.Logger.Info("Received private message", zap.Any("event", e.(*event.PrivateMessageEvent)))
		msgId, err := e.(*event.PrivateMessageEvent).Reply(message.NewText("你好").Message().AsChain(), true)
		if err != nil {
			bot.Logger.Error("Failed to send private message", zap.Error(err))
			return
		}
		bot.Logger.Info("Sent private message", zap.Any("id", msgId))
		resp2, err := bot.GetMsg(msgId)
		if err != nil {
			bot.Logger.Error("Failed to get message", zap.Error(err))
			return
		}
		bot.Logger.Info("Got message", zap.Any("resp", resp2))
		<-time.After(1 * time.Second)
		resp3, err := bot.DeleteMsg(msgId)
		if err != nil {
			bot.Logger.Error("Failed to delete message", zap.Error(err))
			return
		}
		bot.Logger.Info("Deleted message", zap.Any("resp", resp3))
	})
	bot.Start()
	defer gonapcat.Finalize()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	bot.Close()
	<-time.After(1 * time.Second)
}
