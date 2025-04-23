// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package messaging

import "context"

// ConsumerHandler defines the function signature for handling consumed messages.
// It processes a message along with its metadata in a given context.
//
// Parameters:
// - ctx: The context for managing deadlines, cancellations, and other request-scoped values.
// - msg: The message payload to be processed.
// - metadata: Additional metadata associated with the message.
//
// Returns:
// - An error if the message processing fails.
type ConsumerHandler = func(ctx context.Context, msg any, metadata any) error

// Dispatcher defines an interface for registering message handlers and consuming
// messages in a blocking manner. It abstracts the logic for dispatching messages
// to appropriate handlers based on their type or source.
type Dispatcher interface {
	// Register associates a message type and source with a specific ConsumerHandler.
	//
	// Parameters:
	// - from: The source or topic from which the message originates.
	// - msgType: The type of the message to be handled.
	// - handler: The ConsumerHandler function to process the message.
	//
	// Returns:
	// - An error if the registration fails.
	Register(from string, msgType any, handler ConsumerHandler) error

	// ConsumeBlocking starts consuming messages and dispatches them to the
	// appropriate registered handlers. This method blocks the execution and
	// continues running until the process is terminated.
	ConsumeBlocking()
}
