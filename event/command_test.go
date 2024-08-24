package event

import (
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

func (c *testCommand) Preprocess(remaining string) {}

func (c *testCommand) OnCommand(parseResult *ParseResult) {
	c.onCommand(parseResult)
}

func TestSplitMsg1(t *testing.T) {
	assert := assert.New(t)
	msg := "arg1 arg2 arg3"
	args := getArgs(msg, true)
	assert.Equal([]string{"arg1", "arg2", "arg3"}, args)
}

func TestSplitMsg2(t *testing.T) {
	assert := assert.New(t)
	msg := "[CQ:at,qq=123456]arg1 arg2 arg3"
	args := getArgs(msg, false)
	assert.Equal([]string{"[CQ:at,qq=123456]", "arg1", "arg2", "arg3"}, args)
}

func TestSplitMsg3(t *testing.T) {
	assert := assert.New(t)
	msg := "[CQ:at,qq=123456]arg1 arg2 arg3"
	args := getArgs(msg, true)
	assert.Equal([]string{"[CQ:at,qq=123456]arg1", "arg2", "arg3"}, args)
}

func TestSplitMsg4(t *testing.T) {
	assert := assert.New(t)
	msg := "[CQ:at,qq=123 45&amp;6]arg1 arg2&amp; arg3"
	args := getArgs(msg, false)
	assert.Equal([]string{"[CQ:at,qq=123 45&amp;6]", "arg1", "arg2&amp;", "arg3"}, args)
}

func TestSplitMsg5(t *testing.T) {
	assert := assert.New(t)
	msg := `"abc [CQ:at,qq=123 45&amp;6]arg1" arg2`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc [CQ:at,qq=123 45&amp;6]arg1", "arg2"}, args)
}

func TestSplitMsg6(t *testing.T) {
	assert := assert.New(t)
	msg := `"abc\"de\\f" arg2`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc\"de\\f", "arg2"}, args)
}

func TestSplitMsg7(t *testing.T) {
	assert := assert.New(t)
	msg := `"abc[CQ:x,qq=abc"def\g\"hi]f" [CQ:x,qq=abc"def\g\"hi] arg2`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc[CQ:x,qq=abc\"def\\g\\\"hi]f", "[CQ:x,qq=abc\"def\\g\\\"hi]", "arg2"}, args)
}

func TestSplitMsg8(t *testing.T) {
	assert := assert.New(t)
	msg := `abc[CQ:x,qq=abc"def\g\"hi][CQ:x,qq=abc"def\g\"hi] arg2`
	args := getArgs(msg, true)
	assert.Equal([]string{"abc[CQ:x,qq=abc\"def\\g\\\"hi][CQ:x,qq=abc\"def\\g\\\"hi]", "arg2"}, args)
}

func TestSplitMsg9(t *testing.T) {
	assert := assert.New(t)
	msg := `abc[CQ:x,qq=abc"def\g\"hi][CQ:x,qq=abc"def\g\"hi] arg2`
	args := getArgs(msg, false)
	assert.Equal([]string{"abc", "[CQ:x,qq=abc\"def\\g\\\"hi]", "[CQ:x,qq=abc\"def\\g\\\"hi]", "arg2"}, args)
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
