package utils_test

import (
	"testing"

	"github.com/nekoite/go-napcat/utils"
	"github.com/stretchr/testify/assert"
)

func TestDerefAny(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, utils.DerefAny(1))
	x := 1
	assert.Equal(1, utils.DerefAny(&x))
	var c struct{ a int } = struct{ a int }{a: 1}
	assert.True(c == utils.DerefAny(c))
	assert.True(c == utils.DerefAny(&c))
	var d any
	d = c
	assert.True(c == utils.DerefAny(d))
	d = &c
	assert.True(c == utils.DerefAny(d))
}
