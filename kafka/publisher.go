// Package kafka provides an implementation of the messaging.Publisher interface for Kafka.
// It allows publishing messages to Kafka topics using the segmentio/kafka-go library.
package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/messaging"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// kafkaPublisher is the concrete implementation of the messaging.Publisher interface.
// It uses a Kafka writer to send messages to Kafka topics.
//
// Fields:
// - logger: A structured logger for logging events and errors.
// - writer: A Kafka writer instance for sending messages.
type kafkaPublisher struct {
	logger logging.Logger
	writer *kafka.Writer
}

// NewPublisher creates a new instance of kafkaPublisher.
//
// Parameters:
// - configs: Configuration settings including Kafka host and logger.
//
// Returns:
// - A new instance of kafkaPublisher that implements the messaging.Publisher interface.
func NewPublisher(configs *configs.Configs) messaging.Publisher {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(configs.KafkaConfigs.Host),
		Balancer: &kafka.LeastBytes{},
	}

	return &kafkaPublisher{
		logger: configs.Logger,
		writer: writer,
	}
}

// Publish sends a message to the specified Kafka topic.
//
// Parameters:
// - ctx: The context for managing deadlines, cancellations, and other request-scoped values.
// - to: The destination topic where the message should be sent (optional).
// - from: The source or origin of the message (optional).
// - key: A routing key or identifier for the message (optional).
// - msg: The message payload to be sent.
// - options: Additional dynamic parameters for the message (optional).
//
// Returns:
// - An error if the message could not be sent.
func (p *kafkaPublisher) Publish(ctx context.Context, to, from, key *string, msg any, options ...*messaging.Option) error {
	if to == nil || *to == "" {
		return fmt.Errorf("destination topic cannot be empty")
	}

	topic := *to
	messageKey := ""
	if key != nil {
		messageKey = *key
	}

	message := kafka.Message{
		Topic: topic,
		Key:   []byte(messageKey),
		Value: []byte(fmt.Sprintf("%v", msg)),
	}

	p.logger.Info("Publishing message", zap.String("topic", topic), zap.String("key", messageKey))
	return p.writer.WriteMessages(ctx, message)
}

// PublishDeadline sends a message to the specified Kafka topic with a deadline.
//
// This method ensures that the message is sent within the context's deadline.
//
// Parameters:
// - ctx: The context for managing deadlines, cancellations, and other request-scoped values.
// - to: The destination topic where the message should be sent (optional).
// - from: The source or origin of the message (optional).
// - key: A routing key or identifier for the message (optional).
// - msg: The message payload to be sent.
// - options: Additional dynamic parameters for the message (optional).
//
// Returns:
// - An error if the message could not be sent within the deadline.
func (p *kafkaPublisher) PublishDeadline(ctx context.Context, to, from, key *string, msg any, options ...*messaging.Option) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10) // Example timeout
	defer cancel()

	return p.Publish(ctx, to, from, key, msg, options...)
}
