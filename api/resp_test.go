package api

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/nekoite/go-napcat/message"
	"github.com/nekoite/go-napcat/qq"
	"github.com/nekoite/go-napcat/utils"
	"github.com/stretchr/testify/assert"
)

func constructSuccessRespJson[T any](t *testing.T, data T) []byte {
	b, err := json.Marshal(Resp[T]{
		Status:  "ok",
		RetCode: 0,
		Echo:    "123456",
		Data:    data,
	})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func makeSuccessApiResp(raw []byte) apiResp {
	return apiResp{
		Status:  "ok",
		RetCode: 0,
		Echo:    "123456",
		Raw:     raw,
	}
}

func TestParseRespActionSendPrivateMsg(t *testing.T) {
	assert := assert.New(t)
	resp := constructSuccessRespJson(t, map[string]any{"message_id": 123456})
	r, err := parseResp(ActionSendPrivateMsg, makeSuccessApiResp(resp))
	assert.Nil(err)
	assert.NotNil(r)
	if r0, ok := r.(*Resp[RespDataMessageId]); ok {
		assert.EqualValues(123456, r0.Data.MessageId)
	} else {
		assert.Failf("unexpected type", "%T", r)
	}
}

func TestParseRespActionSendGroupMsg(t *testing.T) {
	assert := assert.New(t)
	resp := constructSuccessRespJson(t, map[string]any{"message_id": 123456})
	r, err := parseResp(ActionSendGroupMsg, makeSuccessApiResp(resp))
	assert.Nil(err)
	assert.NotNil(r)
	if r0, ok := r.(*Resp[RespDataMessageId]); ok {
		assert.EqualValues(123456, r0.Data.MessageId)
	} else {
		assert.Failf("unexpected type", "%T", r)
	}
}

func TestParseRespActionDeleteMsg(t *testing.T) {
	assert := assert.New(t)
	resp := constructSuccessRespJson[any](t, nil)
	r, err := parseResp(ActionDeleteMsg, makeSuccessApiResp(resp))
	assert.Nil(err)
	assert.NotNil(r)
	if r0, ok := r.(*Resp[utils.Void]); ok {
		assert.Equal(utils.Void{}, r0.Data)
	} else {
		assert.Failf("unexpected type", "%T", r)
	}
}

func TestParseRespActionGetMsgPrivate(t *testing.T) {
	assert := assert.New(t)
	resp := constructSuccessRespJson(t, map[string]any{
		"message_id": 123456,
		"message": []map[string]any{
			{
				"type": "text",
				"data": map[string]any{
					"text": "text",
				},
			},
			{
				"type": "face",
				"data": map[string]any{
					"id": "123456",
				},
			},
		},
		"time": 123456,
		"sender": map[string]any{
			"user_id":  123456,
			"nickname": "nickname",
			"sex":      qq.SexMale,
			"age":      18,
		},
		"message_type": MessageTypePrivate,
		"real_id":      123456,
	})
	r, err := parseResp(ActionGetMsg, makeSuccessApiResp(resp))
	assert.Nil(err)
	assert.NotNil(r)
	if r0, ok := r.(*Resp[RespDataMessage]); ok {
		assert.EqualValues(123456, r0.Data.MessageId)
		assert.Equal(int64(123456), r0.Data.RealId)
		assert.Equal(int64(123456), r0.Data.Time)
		assert.Equal(MessageTypePrivate, r0.Data.MessageType)
		if s, ok := r0.Data.Sender.(qq.Friend); ok {
			assert.Equal(int64(123456), s.UserId)
			assert.Equal("nickname", s.Nickname)
			assert.Equal(qq.SexMale, s.Sex)
			assert.Equal(int32(18), s.Age)
		} else {
			assert.Failf("unexpected type", "%T", r0.Data.Sender)
		}
		assert.Len(r0.Data.Message.Messages, 2)
		if m0, ok := r0.Data.Message.Messages[0].Data.(message.TextData); ok {
			assert.Equal("text", m0.Text)
		} else {
			assert.Failf("unexpected type", "%T", r0.Data.Message.Messages[0].Data)
		}
		if m1, ok := r0.Data.Message.Messages[1].Data.(message.FaceData); ok {
			assert.Equal(int64(123456), m1.Id)
		} else {
			assert.Failf("unexpected type", "%T", r0.Data.Message.Messages[1].Data)
		}
	} else {
		assert.Failf("unexpected type", "%T", r)
	}
}

func TestParseRespActionGetMsgGroup(t *testing.T) {
	assert := assert.New(t)
	resp := constructSuccessRespJson(t, map[string]any{
		"message_id": 123456,
		"message": []map[string]any{
			{
				"type": "text",
				"data": map[string]any{
					"text": "text",
				},
			},
			{
				"type": "face",
				"data": map[string]any{
					"id": "123456",
				},
			},
		},
		"time": 123456,
		"sender": map[string]any{
			"user_id":  123456,
			"nickname": "nickname",
			"sex":      qq.SexMale,
			"age":      18,
		},
		"message_type": MessageTypeGroup,
		"real_id":      123456,
	})
	r, err := parseResp(ActionGetMsg, makeSuccessApiResp(resp))
	assert.Nil(err)
	assert.NotNil(r)
	if r0, ok := r.(*Resp[RespDataMessage]); ok {
		assert.EqualValues(123456, r0.Data.MessageId)
		assert.Equal(int64(123456), r0.Data.RealId)
		assert.Equal(int64(123456), r0.Data.Time)
		assert.Equal(MessageTypeGroup, r0.Data.MessageType)
		if s, ok := r0.Data.Sender.(qq.GroupUser); ok {
			assert.Equal(int64(123456), s.UserId)
			assert.Equal("nickname", s.Nickname)
			assert.Equal(qq.SexMale, s.Sex)
			assert.Equal(int32(18), s.Age)
		} else {
			assert.Failf("unexpected type", "%T", r0.Data.Sender)
		}
		assert.Len(r0.Data.Message.Messages, 2)
		if m0, ok := r0.Data.Message.Messages[0].Data.(message.TextData); ok {
			assert.Equal("text", m0.Text)
		} else {
			assert.Failf("unexpected type", "%T", r0.Data.Message.Messages[0].Data)
		}
		if m1, ok := r0.Data.Message.Messages[1].Data.(message.FaceData); ok {
			assert.Equal(int64(123456), m1.Id)
		} else {
			assert.Failf("unexpected type", "%T", r0.Data.Message.Messages[1].Data)
		}
	} else {
		assert.Failf("unexpected type", "%T", r)
	}
}

func TestParseRespInvalidJson(t *testing.T) {
	assert := assert.New(t)
	resp := []byte(`{"status":"ok","ret`)
	r, err := parseResp(ActionSendPrivateMsg, makeSuccessApiResp(resp))
	assert.NotNil(err)
	assert.Nil(r)
}
