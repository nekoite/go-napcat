package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
