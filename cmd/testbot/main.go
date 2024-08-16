package main

import (
	"os"
	"os/signal"
	"time"

	gonapcat "github.com/nekoite/go-napcat"
	"github.com/nekoite/go-napcat/config"
	"github.com/nekoite/go-napcat/event"
	"go.uber.org/zap"
)

func main() {
	gonapcat.Init(config.DefaultLogConfig().WithStderr().WithLevel("debug"))
	bot := gonapcat.NewBot(config.DefaultBotConfig(1341400490, "114514"))
	bot.RegisterHandlerPrivateMessage(func(event event.IEvent) {
		bot.Logger.Info("Received private message", zap.Any("event", event))
		resp, err := bot.SendPrivateMsgString(714026292, "你好", false)
		if err != nil {
			bot.Logger.Error("Failed to send private message", zap.Error(err))
			return
		}
		bot.Logger.Info("Sent private message", zap.Any("resp", resp))
		<-time.After(1 * time.Second)
		resp2, err := bot.DeleteMsg(resp.Data.MessageId)
		if err != nil {
			bot.Logger.Error("Failed to delete message", zap.Error(err))
			return
		}
		bot.Logger.Info("Deleted message", zap.Any("resp", resp2))
	})
	bot.Start()
	defer gonapcat.Finalize()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	bot.Close()
	<-time.After(1 * time.Second)
}
