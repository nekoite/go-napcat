package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsChain(t *testing.T) {
	assert := assert.New(t)

	message := Segment{Type: SegmentTypeText, Data: &TextData{Text: "hello"}}

	assert.Equal(&Chain{Messages: []Segment{message}}, message.AsChain())
}

func TestGetTextData(t *testing.T) {
	assert := assert.New(t)

	message := Segment{Type: SegmentTypeText, Data: &TextData{Text: "hello"}}
	textData := message.GetTextData()
	assert.Equal(&TextData{Text: "hello"}, textData)

	message = Segment{Type: SegmentTypeFace, Data: FaceData{Id: 1}}
	assert.Nil(message.GetTextData())
}

func TestGetImageData(t *testing.T) {
	assert := assert.New(t)

	data := &ImageData{BasicFileData: BasicFileData{File: "http://example.com/image.png"}}
	message := Segment{Type: SegmentTypeImage, Data: data}
	imageData := message.GetImageData()
	assert.Equal(data, imageData)

	message = Segment{Type: SegmentTypeFace, Data: FaceData{Id: 1}}
	assert.Nil(message.GetImageData())
}

func TestGetFaceData(t *testing.T) {
	assert := assert.New(t)

	data := &FaceData{Id: 1}
	message := Segment{Type: SegmentTypeFace, Data: data}
	faceData := message.GetFaceData()
	assert.Equal(data, faceData)

	message = Segment{Type: SegmentTypeText, Data: &TextData{Text: "hello"}}
	assert.Nil(message.GetFaceData())
}

func TestNewText(t *testing.T) {
	assert := assert.New(t)

	text := NewText("hello")
	assert.IsType(&TextData{}, text)
	assert.Equal("hello", text.Text)
}

func TestNewTextSegment(t *testing.T) {
	assert := assert.New(t)

	segment := NewText("hello").Segment()
	assert.Equal(SegmentTypeText, segment.Type)
	assert.Equal("hello", segment.Data.(*TextData).Text)
}

func TestNewFace(t *testing.T) {
	assert := assert.New(t)

	face := NewFace(1)
	assert.IsType(&FaceData{}, face)
	assert.EqualValues(1, face.Id)
}

func TestNewFaceSegment(t *testing.T) {
	assert := assert.New(t)

	segment := NewFace(1).Segment()
	assert.Equal(SegmentTypeFace, segment.Type)
	assert.EqualValues(1, segment.Data.(*FaceData).Id)
}

func TestNewAt(t *testing.T) {
	assert := assert.New(t)

	at := NewAt("123456")
	assert.IsType(&AtData{}, at)
	assert.Equal("123456", at.QQ)
}

func TestNewAtSegment(t *testing.T) {
	assert := assert.New(t)

	segment := NewAt("123456").Segment()
	assert.Equal(SegmentTypeAt, segment.Type)
	assert.Equal("123456", segment.Data.(*AtData).QQ)
}

func TestNewAtAll(t *testing.T) {
	assert := assert.New(t)

	atAll := NewAtAll()
	assert.IsType(&AtData{}, atAll)
	assert.Equal("all", atAll.QQ)
}

func TestNewAtAllSegment(t *testing.T) {
	assert := assert.New(t)

	segment := NewAtAll().Segment()
	assert.Equal(SegmentTypeAt, segment.Type)
	assert.Equal("all", segment.Data.(*AtData).QQ)
}

func TestNewAtUser(t *testing.T) {
	assert := assert.New(t)

	atUser := NewAtUser(123456)
	assert.IsType(&AtData{}, atUser)
	assert.Equal("123456", atUser.QQ)
}

func TestNewAtUserSegment(t *testing.T) {
	assert := assert.New(t)

	segment := NewAtUser(123456).Segment()
	assert.Equal(SegmentTypeAt, segment.Type)
	assert.Equal("123456", segment.Data.(*AtData).QQ)
}

func TestNewReply(t *testing.T) {
	assert := assert.New(t)

	reply := NewReply(123456)
	assert.IsType(&ReplyData{}, reply)
	assert.EqualValues(123456, reply.Id)
}

func TestNewReplySegment(t *testing.T) {
	assert := assert.New(t)

	segment := NewReply(123456).Segment()
	assert.Equal(SegmentTypeReply, segment.Type)
	assert.EqualValues(123456, segment.Data.(*ReplyData).Id)
}

func TestNewImage(t *testing.T) {
	assert := assert.New(t)

	image := NewImage("http://example.com/image.png")
	assert.IsType(&ImageData{}, image)
	assert.Equal("http://example.com/image.png", image.File)
}

func TestNewImageSegment(t *testing.T) {
	assert := assert.New(t)

	segment := NewImage("http://example.com/image.png").Segment()
	assert.Equal(SegmentTypeImage, segment.Type)
	assert.Equal("http://example.com/image.png", segment.Data.(*ImageData).File)
}

func TestNewVideo(t *testing.T) {
	assert := assert.New(t)

	video := NewVideo("http://example.com/video.mp4")
	assert.IsType(&VideoData{}, video)
	assert.Equal("http://example.com/video.mp4", video.File)
}
