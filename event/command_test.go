package event

import (
	"sync"
	"testing"

	"github.com/alecthomas/kong"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type testCommand struct {
	name             string
	mode             CmdNameMode
	splitBySpaceOnly bool
	onCommand        func(parseResult *ParseResult)
}

func (c *testCommand) GetName() (string, CmdNameMode) {
	return c.name, c.mode
}

func (c *testCommand) GetNew() any {
	return nil
}

func (c *testCommand) GetOptions() []kong.Option {
	return nil
}

func (c *testCommand) SplitBySpaceOnly() bool {
	return c.splitBySpaceOnly
}

func (c *testCommand) OnCommand(parseResult *ParseResult) {
	c.onCommand(parseResult)
}

func TestGetArgs1(t *testing.T) {
	assert := assert.New(t)
	msg := "arg1 arg2 arg3"
	args := getArgs(msg, true)
	assert.Equal([]string{"arg1", "arg2", "arg3"}, args)
}

func TestGetArgs2(t *testing.T) {
	assert := assert.New(t)
	msg := "[CQ:at,qq=123456]arg1 arg2 arg3"
	args := getArgs(msg, false)
	assert.Equal([]string{"[CQ:at,qq=123456]", "arg1", "arg2", "arg3"}, args)
}

func TestGetArgs3(t *testing.T) {
	assert := assert.New(t)
	msg := "[CQ:at,qq=123456]arg1 arg2 arg3"
	args := getArgs(msg, true)
	assert.Equal([]string{"[CQ:at,qq=123456]arg1", "arg2", "arg3"}, args)
}

func TestGetArgs4(t *testing.T) {
	assert := assert.New(t)
	msg := "[CQ:at,qq=123 45&amp;6]arg1 arg2&amp; arg3"
	args := getArgs(msg, false)
	assert.Equal([]string{"[CQ:at,qq=123 45&amp;6]", "arg1", "arg2&amp;", "arg3"}, args)
}

func TestGetArgs5(t *testing.T) {
	assert := assert.New(t)
	msg := `"abc [CQ:at,qq=123 45&amp;6]arg1" arg2`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc [CQ:at,qq=123 45&amp;6]arg1", "arg2"}, args)
}

func TestGetArgs6(t *testing.T) {
	assert := assert.New(t)
	msg := `"abc\"de\\f" arg2`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc\"de\\f", "arg2"}, args)
}

func TestGetArgs7(t *testing.T) {
	assert := assert.New(t)
	msg := `"abc[CQ:x,qq=abc"def\g\"hi]f" [CQ:x,qq=abc"def\g\"hi] arg2`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc[CQ:x,qq=abc\"def\\g\\\"hi]f", "[CQ:x,qq=abc\"def\\g\\\"hi]", "arg2"}, args)
}

func TestGetArgs8(t *testing.T) {
	assert := assert.New(t)
	msg := `abc[CQ:x,qq=abc"def\g\"hi][CQ:x,qq=abc"def\g\"hi] arg2`
	args := getArgs(msg, true)
	assert.Equal([]string{"abc[CQ:x,qq=abc\"def\\g\\\"hi][CQ:x,qq=abc\"def\\g\\\"hi]", "arg2"}, args)
}

func TestGetArgs9(t *testing.T) {
	assert := assert.New(t)
	msg := `abc[CQ:x,qq=abc"def\g\"hi][CQ:x,qq=abc"def\g\"hi] arg2`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc", "[CQ:x,qq=abc\"def\\g\\\"hi]", "[CQ:x,qq=abc\"def\\g\\\"hi]", "arg2"}, args)
}

func TestGetArgs10(t *testing.T) {
	assert := assert.New(t)
	msg := `abc[CQ:x,qq=abc"def\g\"\\hi] arg2 "arg3\a "`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc", "[CQ:x,qq=abc\"def\\g\\\"\\\\hi]", "arg2", "arg3\\a "}, args)
}

func TestGetPrefix1(t *testing.T) {
	assert := assert.New(t)
	msg := "prefix arg1 arg2 arg3"
	prefix := getPrefix(msg)
	assert.Equal("prefix", prefix)
}

func TestGetPrefix2(t *testing.T) {
	assert := assert.New(t)
	msg := "prefix[CQ:at,qq=123456]arg1 arg2 arg3"
	prefix := getPrefix(msg)
	assert.Equal("prefix", prefix)
}

func TestGetPrefix3(t *testing.T) {
	assert := assert.New(t)
	msg := "pre&amp;fix[CQ:at,qq=123456]arg1 arg2 arg3"
	prefix := getPrefix(msg)
	assert.Equal("pre&amp;fix", prefix)
}

func TestGetCommandModeNormal(t *testing.T) {
	assert := assert.New(t)
	c := NewCommandCenter(zap.NewNop())
	testCmd := &testCommand{
		name:             "prefix",
		mode:             CmdNameModeNormal,
		splitBySpaceOnly: true,
	}
	testCmd2 := &testCommand{
		name:             "pre&fix",
		mode:             CmdNameModeNormal,
		splitBySpaceOnly: true,
	}
	c.Commands["prefix"] = testCmd
	c.Commands["pre&fix"] = testCmd2
	actual, pref := c.getCommand("prefix arg1 arg2 arg3")
	assert.NotNil(actual)
	assert.Equal(testCmd, actual)
	assert.Equal("prefix", pref)

	actual, pref = c.getCommand("prefix[CQ:at,qq=123456]arg1 arg2 arg3")
	assert.NotNil(actual)
	assert.Equal(testCmd, actual)
	assert.Equal("prefix", pref)

	actual, pref = c.getCommand("pre&amp;fix[CQ:at,qq=123456]arg1 arg2 arg3")
	assert.NotNil(actual)
	assert.Equal(testCmd2, actual)
	assert.Equal("pre&amp;fix", pref)

	actual, pref = c.getCommand("prefixd arg1 arg2 arg3")
	assert.Nil(actual)
	assert.Equal("", pref)

	actual, pref = c.getCommand("[CQ:at,qq=123456]arg1 arg2 arg3")
	assert.Nil(actual)
	assert.Equal("", pref)
}

func TestGetCommandModePrefix(t *testing.T) {
	assert := assert.New(t)
	c := NewCommandCenter(zap.NewNop())
	testCmd := &testCommand{
		name:             "prefix",
		mode:             CmdNameModePrefix,
		splitBySpaceOnly: true,
	}
	testCmd2 := &testCommand{
		name:             "pre[fix",
		mode:             CmdNameModePrefix,
		splitBySpaceOnly: true,
	}
	c.PrefixCommands = append(c.PrefixCommands, testCmd, testCmd2)
	actual, pref := c.getCommand("prefix arg1 arg2 arg3")
	assert.NotNil(actual)
	assert.Equal(testCmd, actual)
	assert.Equal("prefix", pref)

	actual, pref = c.getCommand("prefix[CQ:at,qq=123456]arg1 arg2 arg3")
	assert.NotNil(actual)
	assert.Equal(testCmd, actual)
	assert.Equal("prefix", pref)

	actual, pref = c.getCommand("pre&#91;fix[CQ:at,qq=123456]arg1 arg2 arg3")
	assert.NotNil(actual)
	assert.Equal(testCmd2, actual)
	assert.Equal("pre&#91;fix", pref)

	actual, pref = c.getCommand("prefixd arg1 arg2 arg3")
	assert.NotNil(actual)
	assert.Equal(testCmd, actual)
	assert.Equal("prefix", pref)

	actual, pref = c.getCommand("pre&#91;fixd arg1 arg2 arg3")
	assert.NotNil(actual)
	assert.Equal(testCmd2, actual)
	assert.Equal("pre&#91;fix", pref)

	actual, pref = c.getCommand("[CQ:at,qq=123456]arg1 arg2 arg3")
	assert.Nil(actual)
	assert.Equal("", pref)
}

func TestGetCommandWithGlobalPrefixes(t *testing.T) {
	assert := assert.New(t)
	c := NewCommandCenter(zap.NewNop())
	topCmd := &testCommand{
		name:             "top",
		mode:             CmdNameModeNormal,
		splitBySpaceOnly: true,
	}
	helpCmd := &testCommand{
		name:             "help",
		mode:             CmdNameModePrefix,
		splitBySpaceOnly: true,
	}
	c.Commands["top"] = topCmd
	c.PrefixCommands = append(c.PrefixCommands, helpCmd)
	c.SetGlobalCommandPrefixes([]string{".", "。"})

	actual, pref := c.getCommand(".top 次 10")
	assert.Equal(topCmd, actual)
	assert.Equal(".top", pref)

	actual, pref = c.getCommand("。top 次 10")
	assert.Equal(topCmd, actual)
	assert.Equal("。top", pref)

	actual, pref = c.getCommand("。helpme")
	assert.Equal(helpCmd, actual)
	assert.Equal("。help", pref)

	actual, pref = c.getCommand("top 次 10")
	assert.Nil(actual)
	assert.Equal("", pref)

	actual, pref = c.getCommand(".unknown")
	assert.Nil(actual)
	assert.Equal("", pref)
}

func TestSetGlobalCommandPrefixCompat(t *testing.T) {
	assert := assert.New(t)
	c := NewCommandCenter(zap.NewNop())
	topCmd := &testCommand{
		name:             "top",
		mode:             CmdNameModeNormal,
		splitBySpaceOnly: true,
	}
	c.Commands["top"] = topCmd

	c.SetGlobalCommandPrefix(".")
	actual, pref := c.getCommand(".top")
	assert.Equal(topCmd, actual)
	assert.Equal(".top", pref)

	actual, pref = c.getCommand("。top")
	assert.Nil(actual)
	assert.Equal("", pref)

	c.SetGlobalCommandPrefix("")
	actual, pref = c.getCommand("top")
	assert.Equal(topCmd, actual)
	assert.Equal("top", pref)
}

// TestGetCommandPrefixLongestWins 验证重叠前缀场景下，最长前缀优先匹配，与传入顺序无关。
func TestGetCommandPrefixLongestWins(t *testing.T) {
	assert := assert.New(t)
	c := NewCommandCenter(zap.NewNop())
	fooCmd := &testCommand{name: "foo", mode: CmdNameModeNormal, splitBySpaceOnly: true}
	bfooCmd := &testCommand{name: "Bfoo", mode: CmdNameModeNormal, splitBySpaceOnly: true}
	c.Commands["foo"] = fooCmd
	c.Commands["Bfoo"] = bfooCmd

	// 故意以“短前缀在前”的顺序传入，验证排序后“AB”优先于“A”被匹配：
	// - 未排序时："A" 先命中 → stripped="Bfoo" → 返回 bfooCmd。
	// - 排序后："AB" 先命中 → stripped="foo" → 返回 fooCmd。
	c.SetGlobalCommandPrefixes([]string{"A", "AB"})
	actual, pref := c.getCommand("ABfoo")
	assert.Equal(fooCmd, actual)
	assert.Equal("ABfoo", pref)
}

// TestGlobalPrefixesConcurrentAccess 主要依赖 -race 标记运行，验证读写不会发生数据竞争。
func TestGlobalPrefixesConcurrentAccess(t *testing.T) {
	c := NewCommandCenter(zap.NewNop())
	cmd := &testCommand{name: "top", mode: CmdNameModeNormal, splitBySpaceOnly: true}
	c.Commands["top"] = cmd

	var wg sync.WaitGroup
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		toggle := false
		for {
			select {
			case <-stop:
				return
			default:
			}
			if toggle {
				c.SetGlobalCommandPrefixes([]string{".", "。"})
			} else {
				c.SetGlobalCommandPrefix("!")
			}
			toggle = !toggle
		}
	}()

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				_, _ = c.getCommand(".top")
				_, _ = c.getCommand("!top")
				_, _ = c.getCommand("。top")
			}
		}()
	}

	close(stop)
	wg.Wait()
}
