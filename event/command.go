package event

import (
	"sort"
	"strings"
	"sync/atomic"

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
	logger         *zap.Logger
	globalPrefixes atomic.Pointer[[]string]

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
		c.storeGlobalPrefixes(nil)
		return
	}
	c.storeGlobalPrefixes([]string{prefix})
}

// SetGlobalCommandPrefixes 设置一组全局命令前缀，命中其中任意一个即可触发命令解析。
// 传入空切片表示不要求任何前缀。
// 该方法会按前缀长度倒序排序，确保发生重叠时最长前缀优先匹配（例如同时配置
// "!" 和 "!!" 时，"!!cmd" 会以 "!!" 命中），调用方传入顺序无影响。
// 方法可在任意 goroutine 中安全调用，与并发的命令分发不会发生数据竞争。
func (c *CommandCenter) SetGlobalCommandPrefixes(prefixes []string) {
	filtered := make([]string, 0, len(prefixes))
	for _, p := range prefixes {
		if len(p) == 0 {
			continue
		}
		filtered = append(filtered, p)
	}
	c.storeGlobalPrefixes(filtered)
}

// storeGlobalPrefixes 以 copy-on-write 的方式发布前缀列表，发布后切片不会再被修改，
// 读取侧（getCommand）通过 atomic.Pointer 拿到的就是不可变快照，避免数据竞争。
func (c *CommandCenter) storeGlobalPrefixes(prefixes []string) {
	copied := append([]string(nil), prefixes...)
	sort.SliceStable(copied, func(i, j int) bool {
		return len(copied[i]) > len(copied[j])
	})
	c.globalPrefixes.Store(&copied)
}

// loadGlobalPrefixes 返回当前发布的前缀快照；未配置时返回 nil。
func (c *CommandCenter) loadGlobalPrefixes() []string {
	p := c.globalPrefixes.Load()
	if p == nil {
		return nil
	}
	return *p
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
	if sp, ok := cmd.(ICommandStopPropagation); ok && sp.StopPropagation() {
		event.PreventDefault()
	}
}

func (c *CommandCenter) getCommand(raw string) (ICommand, string) {
	pref := getPrefix(raw)
	if len(pref) == 0 {
		return nil, ""
	}
	prefixes := c.loadGlobalPrefixes()
	if len(prefixes) == 0 {
		if cmd, matched := c.matchCommand(pref); cmd != nil {
			return cmd, matched
		}
		return nil, ""
	}
	for _, gp := range prefixes {
		if !strings.HasPrefix(pref, gp) {
			continue
		}
		stripped := pref[len(gp):]
		if cmd, matched := c.matchCommand(stripped); cmd != nil {
			return cmd, gp + matched
		}
	}
	return nil, ""
}

// matchCommand 在已经剥离全局前缀的字符串上尝试匹配 prefix 命令和 normal 命令。
// 返回匹配到的命令以及命中的命令名片段（不含全局前缀）。
func (c *CommandCenter) matchCommand(pref string) (ICommand, string) {
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
