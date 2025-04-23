// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package rabbitmq provides a comprehensive set of utilities for working with RabbitMQ in Go applications.
// It offers abstractions for connections, channels, exchanges, queues, bindings, publishers, and message consumers.
package rabbitmq

// LogMessage formats a consistent log message with a RabbitMQ package prefix.
// It concatenates all provided strings with the "[gokit::rabbitmq]" prefix.
// This helps with identifying RabbitMQ-related log entries throughout the application.
func LogMessage(msg ...string) string {
	f := "[gokit::rabbitmq] "

	for _, s := range msg {
		f += s
	}

	return f
}
