package event

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/nekoite/go-napcat/message"
	"go.uber.org/zap"
)

type CmdNameMode string

const (
	CmdNameModePrefix CmdNameMode = "prefix"
	CmdNameModeNormal CmdNameMode = "normal"
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
	GetName() (string, CmdNameMode)
	// GetNew 用于获取一个新的命令行参数定义结构体。
	// 它将被传入回调函数。具体请参考 kong 的文档。
	GetNew() any
	GetOptions() []kong.Option
	SplitBySpaceOnly() bool
	OnCommand(parseResult *ParseResult)
}

type ICommandWithPreprocess interface {
	Preprocess(remaining string) string
}

type ICommandStopPropagation interface {
	// StopPropagation 返回 true 时，事件将在指令处理完成后停止继续处理其它处理器。
	StopPropagation() bool
}

type CommandCenter struct {
	logger       *zap.Logger
	globalPrefix string

	Commands       map[string]ICommand
	PrefixCommands []ICommand
}

func NewParseResult() *ParseResult {
	return &ParseResult{}
}

func NewCommandCenter(logger *zap.Logger) *CommandCenter {
	return &CommandCenter{
		logger:         logger.Named("command"),
		Commands:       make(map[string]ICommand),
		PrefixCommands: make([]ICommand, 0),
	}
}

func (c *CommandCenter) RegisterCommand(command ICommand) {
	name, mode := command.GetName()
	switch mode {
	case CmdNameModePrefix:
		c.PrefixCommands = append(c.PrefixCommands, command)
	case CmdNameModeNormal:
		c.Commands[name] = command
	default:
		c.logger.Error("unknown command name mode", zap.String("mode", string(mode)))
	}
}

func (c *CommandCenter) SetGlobalCommandPrefix(prefix string) {
	if len(prefix) == 0 {
		return
	}
	c.globalPrefix = prefix
}

func (c *CommandCenter) onMessageRecv(event IMessageEvent) {
	if len(c.Commands) == 0 && len(c.PrefixCommands) == 0 {
		return
	}
	rawMsg := event.GetRawMessage()
	cmd, prefix := c.getCommand(rawMsg)
	if cmd == nil {
		return
	}
	cmdName, _ := cmd.GetName()
	parseResult := NewParseResult()
	stdout := strings.Builder{}
	stderr := strings.Builder{}
	options := []kong.Option{
		kong.Exit(func(i int) { parseResult.ExitCode = i }),
		kong.Writers(&stdout, &stderr),
		kong.Name(cmdName),
	}
	gram := cmd.GetNew()
	k, err := kong.New(
		gram,
		append(options, cmd.GetOptions()...)...,
	)
	if err != nil {
		c.logger.Error("failed to create kong", zap.Error(err))
		return
	}
	remaining := rawMsg[len(prefix):]
	if preprocessCmd, ok := cmd.(ICommandWithPreprocess); ok {
		remaining = preprocessCmd.Preprocess(remaining)
	}
	ctx, err := k.Parse(getArgs(remaining, cmd.SplitBySpaceOnly()))
	if err != nil {
		parseResult.Error = err
	}
	parseResult.ParsedArgs = gram
	parseResult.Ctx = ctx
	parseResult.Event = event
	parseResult.StdOut = stdout.String()
	parseResult.StdErr = stderr.String()
	cmd.OnCommand(parseResult)
	if cmd.(ICommandStopPropagation).StopPropagation() {
		event.PreventDefault()
	}
}

func (c *CommandCenter) getCommand(raw string) (ICommand, string) {
	pref := getPrefix(raw)
	if len(pref) == 0 {
		return nil, ""
	}
	if c.globalPrefix != "" && strings.HasPrefix(pref, c.globalPrefix) {
		pref = pref[len(c.globalPrefix):]
	}
	for _, cmd := range c.PrefixCommands {
		p, _ := cmd.GetName()
		escapedP := message.EscapeCQString(p)
		if strings.HasPrefix(pref, escapedP) {
			return cmd, pref[:len(escapedP)]
		}
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
			if !inCQ {
				appendSbToRes()
			} else {
				sb.WriteByte(c)
			}
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
			if !splitBySpaceOnly {
				appendSbToRes()
			}
			sb.WriteByte(c)
		case ']':
			inCQ = false
			sb.WriteByte(c)
			if !splitBySpaceOnly {
				appendSbToRes()
			}
		default:
			sb.WriteByte(c)
		}
		i++
	}
	appendSbToRes()
	return res
}
