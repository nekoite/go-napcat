package utils_test

import (
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/nekoite/go-napcat/utils"
	"github.com/stretchr/testify/assert"
)

func TestTimedFunc(t *testing.T) {
	assert := assert.New(t)
	patch := gomonkey.ApplyFunc(time.Since, func(t time.Time) time.Duration {
		return time.Second
	})
	defer patch.Reset()
	timed := utils.TimedFunc(func() int {
		return 1
	}, func(d time.Duration) {
		assert.True(d == time.Second)
	})
	assert.Equal(1, timed)
}

func TestTimedAction(t *testing.T) {
	assert := assert.New(t)
	patch := gomonkey.ApplyFunc(time.Since, func(t time.Time) time.Duration {
		return time.Second
	})
	defer patch.Reset()
	x := 0
	utils.TimedAction(func() {
		x = 1
	}, func(d time.Duration) {
		assert.True(d == time.Second)
	})
	assert.Equal(1, x)
}
