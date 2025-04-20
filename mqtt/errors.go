// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package mqtt

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

func NewError(msg string) error {
	return &Error{msg}
}

var (
	ConnectionFailureError = NewError("connection failure")
	EmptyTopicError        = NewError("subscribe top cannot be empty string")
	NillHandlerError       = NewError("subscribe handler cannot be nil")
	NillPayloadError       = NewError("publish payload cannot be nil")
	InvalidQoSError        = NewError("qos must be one of: byte(0), byte(1) or byte(2)")
)
