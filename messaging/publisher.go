// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package messaging

import "context"

// Option represents a key-value pair for additional dynamic parameters in publishing messages.
type Option struct {
	Key   string
	Value string
}

// Publisher defines an interface for publishing messages to a messaging system.
// It provides methods for sending messages with optional metadata such as
// destination, source, and routing keys.
type Publisher interface {
	// Publish sends a message to the specified destination.
	//
	// Parameters:
	// - ctx: The context for managing deadlines, cancellations, and other request-scoped values.
	// - to: The destination or topic where the message should be sent (optional).
	// - from: The source or origin of the message (optional).
	// - key: A routing key or identifier for the message (optional).
	// - msg: The message payload to be sent.
	// - options: Additional dynamic parameters for the message (optional).
	//
	// Returns:
	// - An error if the message could not be sent.
	Publish(ctx context.Context, to, from, key *string, msg any, options ...*Option) error

	// PublishDeadline sends a message to the specified destination with a deadline.
	// This method ensures that the message is sent within the context's deadline.
	//
	// Parameters:
	// - ctx: The context for managing deadlines, cancellations, and other request-scoped values.
	// - to: The destination or topic where the message should be sent (optional).
	// - from: The source or origin of the message (optional).
	// - key: A routing key or identifier for the message (optional).
	// - msg: The message payload to be sent.
	// - options: Additional dynamic parameters for the message (optional).
	//
	// Returns:
	// - An error if the message could not be sent within the deadline.
	PublishDeadline(ctx context.Context, to, from, key *string, msg any, options ...*Option) error
}
