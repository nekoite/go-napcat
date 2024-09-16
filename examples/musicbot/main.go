package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/alecthomas/kong"
	gonapcat "github.com/nekoite/go-napcat"
	"github.com/nekoite/go-napcat/config"
	"github.com/nekoite/go-napcat/event"
	"github.com/nekoite/go-napcat/extensions/napcat"
	"github.com/nekoite/go-napcat/message"
	"github.com/tidwall/gjson"
)

type MusicCommandArgs struct {
	Platform string `short:"p" enum:"163,qq" optional:"" default:"qq" help:"音乐平台，163 或 qq"`
	SongName string `arg:"" required:"" help:"歌曲名"`
}

type MusicCommand struct {
	bot *gonapcat.Bot
}

func (c *MusicCommand) GetName() (string, event.CmdNameMode) {
	return "music", event.CmdNameModeNormal
}

func (c *MusicCommand) GetNew() any {
	return &MusicCommandArgs{}
}

func (c *MusicCommand) GetOptions() []kong.Option {
	return nil
}

func (c *MusicCommand) SplitBySpaceOnly() bool {
	return false
}

func (c *MusicCommand) Preprocess(remaining string) {}

func (c *MusicCommand) OnCommand(parseResult *event.ParseResult) {
	args := parseResult.ParsedArgs.(*MusicCommandArgs)
	if parseResult.Error != nil {
		if parseResult.StdOut != "" {
			parseResult.Event.Reply(message.NewText(parseResult.StdOut).Message().AsChain(), true)
		} else {
			parseResult.Event.Reply(message.NewText(parseResult.Error.Error()).Message().AsChain(), true)
		}
		return
	}
	switch args.Platform {
	case "qq":
		id, err := getQQMusicId(args.SongName)
		if err != nil {
			parseResult.Event.Reply(message.NewText(err.Error()).Message().AsChain(), true)
			return
		}
		napcat.SetMsgEmojiLike(c.bot, parseResult.Event.GetMessageId(), 128166)
		parseResult.Event.Reply(message.NewMusic(message.MusicTypeQQ, id).Message().AsChain(), false)
	default:
		parseResult.Event.Reply(message.NewText("暂不支持").Message().AsChain(), true)
	}
}

func getQQMusicId(songName string) (int64, error) {
	resp, err := http.Get(fmt.Sprintf("https://c6.y.qq.com/splcloud/fcgi-bin/smartbox_new.fcg?_=1724470252605&cv=4747474&ct=24&format=json&inCharset=utf-8&outCharset=utf-8&notice=0&platform=yqq.json&needNewCode=1&uin=0&g_tk_new_20200303=1198146162&g_tk=1198146162&hostUin=0&is_xml=0&key=%s", url.QueryEscape(songName)))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	res := gjson.GetBytes(body, "data.song.itemlist.0.id").Int()
	if res == 0 {
		return 0, fmt.Errorf("未找到")
	}
	return res, nil
}

func main() {
	napcat.Extension.Register()
	gonapcat.Init(config.DefaultLogConfig().WithStderr().WithLevel("debug"))
	bot, err := gonapcat.NewBot(config.DefaultBotConfig(1341400490, "114514"))
	if err != nil {
		panic(err)
	}
	bot.RegisterCommand(&MusicCommand{bot: bot})
	bot.Start()
	defer gonapcat.Finalize()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	<-interrupt
	bot.Close()
	<-time.After(1 * time.Second)
}
