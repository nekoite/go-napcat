package message

import (
	"testing"

	"github.com/nekoite/go-napcat/qq"
	"github.com/stretchr/testify/assert"
)

func TestNewChain(t *testing.T) {
	assert := assert.New(t)

	// Test creating a new chain with no messages
	chain := NewChain()
	assert.NotNil(chain)
	assert.Empty(chain.Messages)

	// Test creating a new chain with one message
	msg1 := NewText("Hello")
	chain = NewChain(msg1.Message())
	assert.NotNil(chain)
	assert.Len(chain.Messages, 1)
	assert.Equal(msg1.Message(), chain.Messages[0])

	// Test creating a new chain with multiple messages
	msg2 := NewFace(1)
	msg3 := NewAt("123456")
	chain = NewChain(msg1.Message(), msg2.Message(), msg3.Message())
	assert.NotNil(chain)
	assert.Len(chain.Messages, 3)
	assert.Equal(msg1.Message(), chain.Messages[0])
	assert.Equal(msg2.Message(), chain.Messages[1])
	assert.Equal(msg3.Message(), chain.Messages[2])
}

func TestPrependMessage(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	msg1 := NewText("Hello")
	chain.PrependMessage(msg1.Message())
	assert.Len(chain.Messages, 1)
	assert.Equal(msg1.Message(), chain.Messages[0])

	msg2 := NewFace(1)
	chain.PrependMessage(msg2.Message())
	assert.Len(chain.Messages, 2)
	assert.Equal(msg2.Message(), chain.Messages[0])
	assert.Equal(msg1.Message(), chain.Messages[1])
}

func TestAddMessage(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	msg1 := NewText("Hello")
	chain.AddMessage(msg1.Message())
	assert.Len(chain.Messages, 1)
	assert.Equal(msg1.Message(), chain.Messages[0])

	msg2 := NewFace(1)
	chain.AddMessage(msg2.Message())
	assert.Len(chain.Messages, 2)
	assert.Equal(msg1.Message(), chain.Messages[0])
	assert.Equal(msg2.Message(), chain.Messages[1])
}

func TestAddText(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	chain.AddText("Hello")
	assert.Len(chain.Messages, 1)
	assert.Equal(NewText("Hello").Message(), chain.Messages[0])
}

func TestAddAt(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	chain.AddAt("123456")
	assert.Len(chain.Messages, 1)
	assert.Equal(NewAt("123456").Message(), chain.Messages[0])
}

func TestAppendChain(t *testing.T) {
	assert := assert.New(t)

	chain1 := NewChain(NewText("Hello").Message())
	chain2 := NewChain(NewFace(1).Message(), NewAt("123456").Message())

	chain1.AppendChain(*chain2)
	assert.Len(chain1.Messages, 3)
	assert.Equal(NewText("Hello").Message(), chain1.Messages[0])
	assert.Equal(NewFace(1).Message(), chain1.Messages[1])
	assert.Equal(NewAt("123456").Message(), chain1.Messages[2])
}

func TestSetSendAsAnonymous(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message())
	chain.SetSendAsAnonymous(true)
	assert.Len(chain.Messages, 2)
	assert.IsType(&AnonymousData{}, chain.Messages[0].Data)
	assert.True(chain.Messages[0].Data.(*AnonymousData).Ignore)
}

func TestSetReplyTo(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message())
	chain.SetReplyTo(123456)
	assert.Len(chain.Messages, 2)
	assert.IsType(ReplyData{}, chain.Messages[0].Data)
	assert.EqualValues(123456, chain.Messages[0].Data.(ReplyData).Id)
}

func TestPrependChain(t *testing.T) {
	assert := assert.New(t)

	chain1 := NewChain(NewText("Hello").Message())
	chain2 := NewChain(NewFace(1).Message(), NewAt("123456").Message())

	chain1.PrependChain(*chain2)
	assert.Len(chain1.Messages, 3)
	assert.Equal(NewFace(1).Message(), chain1.Messages[0])
	assert.Equal(NewAt("123456").Message(), chain1.Messages[1])
	assert.Equal(NewText("Hello").Message(), chain1.Messages[2])
}

func TestAddAtUser(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	chain.AddAtUser(qq.UserId(123456))
	assert.Len(chain.Messages, 1)
	assert.Equal(NewAt("123456").Message(), chain.Messages[0])
}

func TestAddAtAll(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	chain.AddAtAll()
	assert.Len(chain.Messages, 1)
	assert.Equal(NewAtAll().Message(), chain.Messages[0])
}
func TestLen(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message(), NewFace(1).Message())
	assert.Equal(2, chain.Len())
}

func TestAt(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message(), NewFace(1).Message())
	assert.Equal(chain.Messages[0], chain.At(0))
	assert.Equal(chain.Messages[1], chain.At(1))
}

func TestAtRef(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message(), NewFace(1).Message())
	assert.Equal(&chain.Messages[0], chain.AtRef(0))
	assert.Equal(&chain.Messages[1], chain.AtRef(1))
}

func TestFirstOfType(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message(), NewFace(1).Message(), NewText("World").Message())
	assert.Equal(NewText("Hello").Message(), chain.FirstOfType(MessageTypeText))
	assert.Equal(NewFace(1).Message(), chain.FirstOfType(MessageTypeFace))
	assert.Empty(chain.FirstOfType(MessageTypeAt))
	assert.True(chain.FirstOfType(MessageTypeAt).IsInvalid())
}

func TestFirstOfTypeRef(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message(), NewFace(1).Message(), NewText("World").Message())
	assert.Equal(&chain.Messages[0], chain.FirstOfTypeRef(MessageTypeText))
	assert.Equal(&chain.Messages[1], chain.FirstOfTypeRef(MessageTypeFace))
	assert.Nil(chain.FirstOfTypeRef(MessageTypeAt))
}

func TestFirstText(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message(), NewFace(1).Message(), NewText("World").Message())
	assert.Equal(chain.Messages[0].Data, chain.FirstText())

	chain = NewChain(NewFace(1).Message(), NewFace(2).Message())
	assert.Nil(chain.FirstText())
}

func TestFirstImage(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Message(), NewFace(1).Message(), NewText("World").Message())
	assert.Nil(chain.FirstImage())

	chain = NewChain(NewText("Hello").Message(), NewFace(1).Message(), NewImage("image1.png").Message(), NewImage("image2.png").Message())
	assert.EqualValues(chain.Messages[2].Data, chain.FirstImage())
}
