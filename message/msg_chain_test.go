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
	chain = NewChain(msg1.Segment())
	assert.NotNil(chain)
	assert.Len(chain.Messages, 1)
	assert.Equal(msg1.Segment(), chain.Messages[0])

	// Test creating a new chain with multiple messages
	msg2 := NewFace(1)
	msg3 := NewAt("123456")
	chain = NewChain(msg1.Segment(), msg2.Segment(), msg3.Segment())
	assert.NotNil(chain)
	assert.Len(chain.Messages, 3)
	assert.Equal(msg1.Segment(), chain.Messages[0])
	assert.Equal(msg2.Segment(), chain.Messages[1])
	assert.Equal(msg3.Segment(), chain.Messages[2])
}

func TestPrependSegment(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	msg1 := NewText("Hello")
	chain.PrependSegment(msg1.Segment())
	assert.Len(chain.Messages, 1)
	assert.Equal(msg1.Segment(), chain.Messages[0])

	msg2 := NewFace(1)
	chain.PrependSegment(msg2.Segment())
	assert.Len(chain.Messages, 2)
	assert.Equal(msg2.Segment(), chain.Messages[0])
	assert.Equal(msg1.Segment(), chain.Messages[1])
}

func TestAddSegment(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	msg1 := NewText("Hello")
	chain.AddSegment(msg1.Segment())
	assert.Len(chain.Messages, 1)
	assert.Equal(msg1.Segment(), chain.Messages[0])

	msg2 := NewFace(1)
	chain.AddSegment(msg2.Segment())
	assert.Len(chain.Messages, 2)
	assert.Equal(msg1.Segment(), chain.Messages[0])
	assert.Equal(msg2.Segment(), chain.Messages[1])
}

func TestAddSegments(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	msg1 := NewText("Hello")
	msg2 := NewFace(1)
	chain.AddSegments(msg1.Segment(), msg2.Segment())
	assert.Len(chain.Messages, 2)
	assert.Equal(msg1.Segment(), chain.Messages[0])
	assert.Equal(msg2.Segment(), chain.Messages[1])
}

func TestAddText(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	chain.AddText("Hello")
	assert.Len(chain.Messages, 1)
	assert.Equal(NewText("Hello").Segment(), chain.Messages[0])
}

func TestAddAt(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	chain.AddAt("123456")
	assert.Len(chain.Messages, 1)
	assert.Equal(NewAt("123456").Segment(), chain.Messages[0])
}

func TestAppendChain(t *testing.T) {
	assert := assert.New(t)

	chain1 := NewChain(NewText("Hello").Segment())
	chain2 := NewChain(NewFace(1).Segment(), NewAt("123456").Segment())

	chain1.AppendChain(*chain2)
	assert.Len(chain1.Messages, 3)
	assert.Equal(NewText("Hello").Segment(), chain1.Messages[0])
	assert.Equal(NewFace(1).Segment(), chain1.Messages[1])
	assert.Equal(NewAt("123456").Segment(), chain1.Messages[2])
}

func TestSetSendAsAnonymous(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment())
	chain.SetSendAsAnonymous(true)
	assert.Len(chain.Messages, 2)
	assert.IsType(&AnonymousData{}, chain.Messages[0].Data)
	assert.True(chain.Messages[0].Data.(*AnonymousData).Ignore)
}

func TestSetReplyTo(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment())
	chain.SetReplyTo(123456)
	assert.Len(chain.Messages, 2)
	assert.IsType(&ReplyData{}, chain.Messages[0].Data)
	assert.EqualValues(123456, chain.Messages[0].Data.(*ReplyData).Id)
}

func TestPrependChain(t *testing.T) {
	assert := assert.New(t)

	chain1 := NewChain(NewText("Hello").Segment())
	chain2 := NewChain(NewFace(1).Segment(), NewAt("123456").Segment())

	chain1.PrependChain(*chain2)
	assert.Len(chain1.Messages, 3)
	assert.Equal(NewFace(1).Segment(), chain1.Messages[0])
	assert.Equal(NewAt("123456").Segment(), chain1.Messages[1])
	assert.Equal(NewText("Hello").Segment(), chain1.Messages[2])
}

func TestAddAtUser(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	chain.AddAtUser(qq.UserId(123456))
	assert.Len(chain.Messages, 1)
	assert.Equal(NewAt("123456").Segment(), chain.Messages[0])
}

func TestAddAtAll(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain()
	chain.AddAtAll()
	assert.Len(chain.Messages, 1)
	assert.Equal(NewAtAll().Segment(), chain.Messages[0])
}
func TestLen(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment(), NewFace(1).Segment())
	assert.Equal(2, chain.Len())
}

func TestAt(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment(), NewFace(1).Segment())
	assert.Equal(chain.Messages[0], chain.At(0))
	assert.Equal(chain.Messages[1], chain.At(1))
}

func TestAtRef(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment(), NewFace(1).Segment())
	assert.Equal(&chain.Messages[0], chain.AtRef(0))
	assert.Equal(&chain.Messages[1], chain.AtRef(1))
}

func TestFirstOfType(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment(), NewFace(1).Segment(), NewText("World").Segment())
	assert.Equal(NewText("Hello").Segment(), chain.FirstOfType(SegmentTypeText))
	assert.Equal(NewFace(1).Segment(), chain.FirstOfType(SegmentTypeFace))
	assert.Empty(chain.FirstOfType(SegmentTypeAt))
	assert.True(chain.FirstOfType(SegmentTypeAt).IsInvalid())
}

func TestFirstOfTypeRef(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment(), NewFace(1).Segment(), NewText("World").Segment())
	assert.Equal(&chain.Messages[0], chain.FirstOfTypeRef(SegmentTypeText))
	assert.Equal(&chain.Messages[1], chain.FirstOfTypeRef(SegmentTypeFace))
	assert.Nil(chain.FirstOfTypeRef(SegmentTypeAt))
}

func TestFirstText(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment(), NewFace(1).Segment(), NewText("World").Segment())
	assert.Equal(chain.Messages[0].Data, chain.FirstText())

	chain = NewChain(NewFace(1).Segment(), NewFace(2).Segment())
	assert.Nil(chain.FirstText())
}

func TestFirstImage(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment(), NewFace(1).Segment(), NewText("World").Segment())
	assert.Nil(chain.FirstImage())

	chain = NewChain(NewText("Hello").Segment(), NewFace(1).Segment(), NewImage("image1.png").Segment(), NewImage("image2.png").Segment())
	assert.EqualValues(chain.Messages[2].Data, chain.FirstImage())
}

func TestClear(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment())
	chain.Clear()
	assert.Empty(chain.Messages)
}

func TestGetSegmentsWithType(t *testing.T) {
	assert := assert.New(t)

	chain := NewChain(NewText("Hello").Segment(), NewFace(1).Segment(), NewText("World").Segment())
	textSegments := chain.GetSegmentsWithType(SegmentTypeText)
	assert.Len(textSegments, 2)
	assert.Equal(NewText("Hello").Segment(), textSegments[0])
	assert.Equal(NewText("World").Segment(), textSegments[1])

	faceSegments := chain.GetSegmentsWithType(SegmentTypeFace)
	assert.Len(faceSegments, 1)
	assert.Equal(NewFace(1).Segment(), faceSegments[0])
}
