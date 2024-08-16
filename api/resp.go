package api

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/nekoite/go-napcat/errors"
	"github.com/nekoite/go-napcat/message"
	"github.com/nekoite/go-napcat/utils"
)

type MessageType string

const (
	MessageTypePrivate MessageType = "private"
	MessageTypeGroup   MessageType = "group"
)

type IResp interface {
	GetStatus() string
	GetRetCode() int
	GetEcho() string
	GetData() any
}

type Resp[T any] struct {
	Status  string `json:"status"`
	RetCode int    `json:"retcode"`
	Echo    string `json:"echo"`
	Data    T      `json:"data"`
}

type RespDataMessageId struct {
	MessageId int64 `json:"message_id" mapstructure:"message_id"`
}

type RespDataMessage struct {
	RespDataMessageId
	Time        int64           `json:"time" mapstructure:"time"`
	MessageType MessageType     `json:"message_type" mapstructure:"message_type"`
	Message     message.Chain   `json:"message" mapstructure:"message"`
	RealId      int64           `json:"real_id" mapstructure:"real_id"`
	Sender      message.ISender `json:"sender" mapstructure:"sender"`
}

type ServerStatus struct {
	Online bool `json:"online" mapstructure:"online"`
	Good   bool `json:"good" mapstructure:"good"`
}

func (r *Resp[T]) GetStatus() string {
	return r.Status
}

func (r *Resp[T]) GetRetCode() int {
	return r.RetCode
}

func (r *Resp[T]) GetEcho() string {
	return r.Echo
}

func (r *Resp[T]) GetData() any {
	return r.Data
}

func ParseResp(action Action, data map[string]any) (IResp, error) {
	// status := data["status"].(string)
	retcode := int(data["retcode"].(float64))
	// echo := data["echo"].(string)
	if retcode != 0 {
		return nil, fmt.Errorf("%w code %d", errors.ErrApiResp, retcode)
	}
	var resp any
	switch action {
	case ActionSendPrivateMsg:
		fallthrough
	case ActionSendGroupMsg:
		resp = &Resp[RespDataMessageId]{}
	case ActionDeleteMsg:
		resp = &Resp[utils.Void]{}
	default:
		return nil, fmt.Errorf("%w : %s", errors.ErrUnknownAction, action)
	}
	if err := mapstructure.Decode(data, &resp); err != nil {
		return nil, err
	}
	return resp.(IResp), nil
}
