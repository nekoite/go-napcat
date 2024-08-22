package event

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/nekoite/go-napcat/message"
	"go.uber.org/zap"
)

type ParseResult struct {
	Ctx        *kong.Context
	Event      IMessageEvent
	ParsedArgs any
	Error      error
	ExitCode   int
	StdOut     string
	StdErr     string
}

type ICommand interface {
	GetName() string
	// GetNew 用于获取一个新的命令行参数定义结构体。
	// 它将被传入回调函数。具体请参考 kong 的文档。
	GetNew() any
	GetOptions() []kong.Option
	SplitBySpaceOnly() bool
	Preprocess(remaining string)
	OnCommand(parseResult *ParseResult)
}

type CommandCenter struct {
	logger *zap.Logger

	Commands map[string]ICommand
}

func NewParseResult() *ParseResult {
	return &ParseResult{}
}

func NewCommandCenter(logger *zap.Logger) *CommandCenter {
	return &CommandCenter{
		logger:   logger.Named("command"),
		Commands: make(map[string]ICommand),
	}
}

func (c *CommandCenter) RegisterCommand(command ICommand) {
	c.Commands[command.GetName()] = command
}

func (c *CommandCenter) onMessageRecv(event IMessageEvent) {
	if len(c.Commands) == 0 {
		return
	}
	rawMsg := event.GetRawMessage()
	cmd, prefix := c.isCommand(rawMsg)
	if cmd == nil {
		return
	}
	parseResult := NewParseResult()
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	options := []kong.Option{
		kong.Exit(func(i int) { parseResult.ExitCode = i }),
		kong.Writers(&stdout, &stderr),
	}
	k, err := kong.New(
		cmd.GetNew(),
		append(options, cmd.GetOptions()...)...,
	)
	if err != nil {
		c.logger.Error("failed to create kong", zap.Error(err))
		return
	}
	ctx, err := k.Parse(getArgs(rawMsg[len(prefix):], cmd.SplitBySpaceOnly()))
	if err != nil {
		parseResult.Error = err
	}
	parseResult.Ctx = ctx
	parseResult.Event = event
	parseResult.StdOut = stdout.String()
	parseResult.StdErr = stderr.String()
	cmd.OnCommand(parseResult)
}

func (c *CommandCenter) isCommand(raw string) (ICommand, string) {
	pref := getPrefix(raw)
	if len(pref) == 0 {
		return nil, ""
	}
	cmd, ok := c.Commands[message.UnescapeCQString(pref)]
	if !ok {
		return nil, ""
	}
	return cmd, pref
}

func getPrefix(raw string) string {
	pref, _, _ := strings.Cut(raw, "[CQ:")
	pref, _, _ = strings.Cut(pref, " ")
	if len(pref) == 0 {
		return ""
	}
	return pref
}

func getArgs(s string, splitBySpaceOnly bool) []string {
	res := make([]string, 0, 3)
	sb := strings.Builder{}
	i := 0
	appendSbToRes := func() {
		if sb.Len() > 0 {
			res = append(res, sb.String())
			sb.Reset()
		}
	}
	inCQ := false
	for i < len(s) {
		c := s[i]
	SW:
		switch c {
		case ' ':
			appendSbToRes()
		case '"':
			if inCQ {
				sb.WriteByte('"')
				break SW
			}
			i++
			inCQ2 := false
		LP2:
			for i < len(s) {
			SW2:
				switch s[i] {
				case '\\':
					if inCQ2 {
						sb.WriteByte('\\')
						break SW2
					}
					i++
					if i < len(s) {
						switch s[i] {
						case '"':
							sb.WriteByte('"')
						case '\\':
							sb.WriteByte('\\')
						default:
							sb.WriteByte('\\')
							sb.WriteByte(s[i])
						}
					} else {
						sb.WriteByte('\\')
						break LP2
					}
				case '"':
					if inCQ2 {
						sb.WriteByte('"')
					} else {
						break LP2
					}
				case '[':
					inCQ2 = true
					sb.WriteByte(s[i])
				case ']':
					inCQ2 = false
					sb.WriteByte(s[i])
				default:
					sb.WriteByte(s[i])
				}
				i++
			}
		case '[':
			inCQ = true
			if splitBySpaceOnly {
				sb.WriteByte(c)
			} else {
				appendSbToRes()
				for i < len(s) {
					sb.WriteByte(s[i])
					if s[i] == ']' {
						if !splitBySpaceOnly {
							appendSbToRes()
						}
						break
					}
					i++
				}
			}
		case ']':
			inCQ = false
			fallthrough
		default:
			sb.WriteByte(c)
		}
		i++
	}
	appendSbToRes()
	return res
}
