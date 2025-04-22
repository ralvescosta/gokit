// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package mqtt

// Error represents a custom error type for the MQTT package.
type Error struct {
	msg string
}

// Error returns the error message as a string.
func (e *Error) Error() string {
	return e.msg
}

// NewError creates a new instance of the custom Error type with the provided message.
func NewError(msg string) error {
	return &Error{msg}
}

var (
	// ConnectionFailureError indicates a failure to connect to the MQTT broker.
	ConnectionFailureError = NewError("connection failure")
	// EmptyTopicError indicates that the topic for a subscription cannot be an empty string.
	EmptyTopicError = NewError("subscribe topic cannot be an empty string")
	// NillHandlerError indicates that the handler for a subscription cannot be nil.
	NillHandlerError = NewError("subscribe handler cannot be nil")
	// NillPayloadError indicates that the payload for a publish operation cannot be nil.
	NillPayloadError = NewError("publish payload cannot be nil")
	// InvalidQoSError indicates that the provided QoS value is invalid.
	InvalidQoSError = NewError("qos must be one of: byte(0), byte(1), or byte(2)")
)
