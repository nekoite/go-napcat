package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextMessageToString(t *testing.T) {
	assert := assert.New(t)
	msg := NewText("Hello, w=orld!&?")
	assert.Equal("[CQ:text,text=Hello&#44; w=orld!&amp;?]", msg.Message().String())
}

func TestFaceMessageToString(t *testing.T) {
	assert := assert.New(t)
	msg := NewFace(1)
	assert.Equal("[CQ:face,id=1]", msg.Message().String())
}
