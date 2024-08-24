package errors

import (
	"errors"
	"fmt"
)

var (
	ErrGoNapcat = errors.New("go-napcat error")

	ErrUnknownEvent        = fmt.Errorf("%w: unknown event", ErrGoNapcat)
	ErrUnknownMessageEvent = fmt.Errorf("%w: unknown message event", ErrGoNapcat)
	ErrUnknownNoticeEvent  = fmt.Errorf("%w: unknown notice event", ErrGoNapcat)
	ErrUnknownRequestEvent = fmt.Errorf("%w: unknown request event", ErrGoNapcat)
	ErrUnknownMetaEvent    = fmt.Errorf("%w: unknown meta event", ErrGoNapcat)

	ErrUnknownResponse = fmt.Errorf("%w: unknown response", ErrGoNapcat)
	ErrUnknownAction   = fmt.Errorf("%w: unknown action", ErrGoNapcat)
	ErrApiResp         = fmt.Errorf("%w: api response", ErrGoNapcat)

	ErrInvalidMessage  = fmt.Errorf("%w: invalid message", ErrGoNapcat)
	ErrInvalidCQString = fmt.Errorf("%w: invalid CQ string", ErrGoNapcat)

	ErrExtensionAlreadyRegistered = fmt.Errorf("%w: extension already registered", ErrGoNapcat)
	ErrActionAlreadyRegistered    = fmt.Errorf("%w: action already registered", ErrGoNapcat)

	ErrUnsupportedOperation = fmt.Errorf("%w: unsupported operation", ErrGoNapcat)
	ErrTimeout              = fmt.Errorf("%w: timeout", ErrGoNapcat)

	ErrTypeAssertion = errors.New("type assertion failed")
)
