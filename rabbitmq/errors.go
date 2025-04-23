// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

// RabbitMQError represents a custom error type for RabbitMQ-related operations.
// It encapsulates an error message describing the specific error condition.
type RabbitMQError struct {
	msg string
}

// Error implements the error interface and returns the error message.
func (e *RabbitMQError) Error() string {
	return e.msg
}

// NewRabbitMQError creates a new RabbitMQError instance with the provided message.
// Returns the error as an error interface.
func NewRabbitMQError(msg string) error {
	return &RabbitMQError{msg}
}

var (
	// rabbitMQDialError is a function that wraps a connection error into a RabbitMQError.
	rabbitMQDialError = func(err error) error { return NewRabbitMQError(err.Error()) }

	// getChannelError is a function that wraps a channel creation error into a RabbitMQError.
	getChannelError = func(err error) error { return NewRabbitMQError(err.Error()) }

	// NullableChannelError is returned when a channel operation is attempted on a nil channel.
	NullableChannelError = NewRabbitMQError("channel cant be null")

	// NotFoundQueueDefinitionError is returned when a queue definition cannot be found.
	NotFoundQueueDefinitionError = NewRabbitMQError("not found queue definition")

	// InvalidDispatchParamsError is returned when invalid parameters are provided to a dispatch operation.
	InvalidDispatchParamsError = NewRabbitMQError("register dispatch with invalid parameters")

	// QueueDefinitionNotFoundError is returned when no queue definition is found for a specified queue.
	QueueDefinitionNotFoundError = NewRabbitMQError("any queue definition was founded to the given queue")

	// ReceivedMessageWithUnformattedHeaderError is returned when a message has incorrectly formatted headers.
	ReceivedMessageWithUnformattedHeaderError = NewRabbitMQError("received message with unformatted headers")

	// RetryableError indicates that a message processing failed but can be retried later.
	RetryableError = NewRabbitMQError("error to process this message, retry latter")
)
