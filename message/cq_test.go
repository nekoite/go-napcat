package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextSegmentToString(t *testing.T) {
	assert := assert.New(t)
	msg := NewText("Hello, w=orld!&?")
	assert.Equal("Hello&#44; w=orld!&amp;?", msg.Segment().String())
}

func TestFaceSegmentToString(t *testing.T) {
	assert := assert.New(t)
	msg := NewFace(1)
	assert.Equal("[CQ:face,id=1]", msg.Segment().String())
}

func TestUnknownSegmentToString(t *testing.T) {
	assert := assert.New(t)
	msg := Segment{Type: SegmentType("unknown"), Data: UnknownData{"key": "value"}}
	assert.Equal("[CQ:unknown,key=value]", msg.String())
}

func TestChainToString(t *testing.T) {
	assert := assert.New(t)
	chain := Chain{
		Messages: []Segment{
			{Type: SegmentTypeText, Data: &TextData{Text: "Hello, w=orld!&?"}},
			{Type: SegmentTypeFace, Data: &FaceData{Id: 1}},
			{Type: SegmentTypeAt, Data: &AtData{QQ: "123456"}},
			{Type: SegmentTypeText, Data: &TextData{Text: "another&text"}},
		},
	}
	assert.Equal("Hello&#44; w=orld!&amp;?[CQ:face,id=1][CQ:at,qq=123456]another&amp;text", chain.String())
}

func TestTextCQToChain(t *testing.T) {
	assert := assert.New(t)
	cq := "Hello&#44; w=orld!&amp;?"
	expected := &Chain{
		Messages: []Segment{
			{Type: SegmentTypeText, Data: &TextData{Text: "Hello, w=orld!&?"}},
		},
	}
	actual, err := ParseCQString(cq)
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestFaceCQToChain(t *testing.T) {
	assert := assert.New(t)
	cq := "[CQ:face,id=1]"
	expected := &Chain{
		Messages: []Segment{
			{Type: SegmentTypeFace, Data: &FaceData{Id: 1}},
		},
	}
	actual, err := ParseCQString(cq)
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestCustomNodeCQToChain(t *testing.T) {
	assert := assert.New(t)
	cq := "[CQ:node,user_id=10001000,nickname=某人,content=&#91;CQ:face&#44;id=123&#93;哈喽&amp;amp;~]"
	expected := &Chain{
		Messages: []Segment{
			{Type: SegmentTypeNode, Data: &CustomNodeData{
				UserId:   10001000,
				Nickname: "某人",
				Content: &Chain{
					Messages: []Segment{
						{Type: SegmentTypeFace, Data: &FaceData{Id: 123}},
						{Type: SegmentTypeText, Data: &TextData{Text: "哈喽&~"}},
					},
				},
			}},
		},
	}
	actual, err := ParseCQString(cq)
	assert.Nil(err)
	assert.Equal(expected, actual)
}

func TestChainCQToChain(t *testing.T) {
	assert := assert.New(t)
	cq := "Hello&#44; w=orld!&amp;?[CQ:face,id=1][CQ:at,qq=123456]another&amp;text"
	expected := &Chain{
		Messages: []Segment{
			{Type: SegmentTypeText, Data: &TextData{Text: "Hello, w=orld!&?"}},
			{Type: SegmentTypeFace, Data: &FaceData{Id: 1}},
			{Type: SegmentTypeAt, Data: &AtData{QQ: "123456"}},
			{Type: SegmentTypeText, Data: &TextData{Text: "another&text"}},
		},
	}
	actual, err := ParseCQString(cq)
	assert.Nil(err)
	assert.Equal(expected, actual)
}
