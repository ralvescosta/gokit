// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type (
	// Dispatcher defines an interface for managing RabbitMQ message consumption.
	// It provides methods to register message handlers and consume messages in a blocking manner.
	Dispatcher interface {
		// Register associates a queue with a message type and a handler function.
		// It ensures that messages from the specified queue are processed by the handler.
		// Returns an error if the registration parameters are invalid or if the queue definition is not found.
		Register(queue string, typE any, handler ConsumerHandler) error

		// ConsumeBlocking starts consuming messages and dispatches them to the registered handlers.
		// This method blocks execution until the process is terminated by a signal.
		ConsumeBlocking()
	}

	// dispatcher is the concrete implementation of the Dispatcher interface.
	// It manages the registration and execution of message handlers for RabbitMQ queues.
	dispatcher struct {
		logger              logging.Logger
		channel             AMQPChannel
		queueDefinitions    map[string]*QueueDefinition
		consumersDefinition map[string]*ConsumerDefinition
		tracer              trace.Tracer
		signalCh            chan os.Signal
	}

	// ConsumerHandler is a function type that defines message handler callbacks.
	// It receives a context (for tracing), the unmarshaled message, and metadata about the delivery.
	// Returns an error if the message processing fails.
	ConsumerHandler = func(ctx context.Context, msg any, metadata any) error

	// ConsumerDefinition represents the configuration for a consumer.
	// It holds information about the queue, message type, and handler function.
	ConsumerDefinition struct {
		queue           string
		msgType         string
		reflect         *reflect.Value
		queueDefinition *QueueDefinition
		handler         ConsumerHandler
	}

	// deliveryMetadata contains metadata extracted from an AMQP delivery.
	// This includes message ID, retry count, message type, and headers.
	deliveryMetadata struct {
		MessageId string
		XCount    int64
		Type      string
		Headers   map[string]interface{}
	}
)

// NewDispatcher creates a new dispatcher instance with the provided configuration.
// It initializes signal handling and sets up the necessary components for message consumption.
func NewDispatcher(cfgs *configs.Configs, channel AMQPChannel, queueDefinitions map[string]*QueueDefinition) *dispatcher {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	return &dispatcher{
		logger:              cfgs.Logger,
		channel:             channel,
		queueDefinitions:    queueDefinitions,
		consumersDefinition: map[string]*ConsumerDefinition{},
		tracer:              otel.Tracer("rmq-dispatcher"),
		signalCh:            signalCh,
	}
}

// Register associates a queue with a message type and a handler function.
// It validates the parameters and ensures that the queue definition exists.
// Returns an error if the registration parameters are invalid or if the queue definition is not found.
func (d *dispatcher) Register(queue string, msg any, handler ConsumerHandler) error {
	if msg == nil || queue == "" {
		return InvalidDispatchParamsError
	}

	def, ok := d.queueDefinitions[queue]
	if !ok {
		return QueueDefinitionNotFoundError
	}

	ref := reflect.New(reflect.TypeOf(msg))
	msgType := fmt.Sprintf("%T", msg)

	d.consumersDefinition[msgType] = &ConsumerDefinition{
		queue:           queue,
		msgType:         msgType,
		reflect:         &ref,
		queueDefinition: def,
		handler:         handler,
	}

	return nil
}

// ConsumeBlocking starts consuming messages from all registered queues.
// It creates a goroutine for each consumer and blocks until a termination signal is received.
func (d *dispatcher) ConsumeBlocking() {
	for _, cd := range d.consumersDefinition {
		go d.consume(cd.queue, cd.msgType)
	}

	<-d.signalCh
	d.logger.Debug(LogMessage("signal received, closing dispatcher"))
}

// consume starts consuming messages from a specific queue.
// It handles message unmarshaling, error handling, retries, and dead-letter queuing.
func (d *dispatcher) consume(queue, msgType string) {
	delivery, err := d.channel.Consume(queue, msgType, false, false, false, false, nil)
	if err != nil {
		d.logger.Error(
			LogMessage("failure to declare consumer"),
			zap.String("queue", queue),
			zap.Error(err),
		)
		return
	}

	for received := range delivery {
		metadata, err := d.extractMetadata(&received)
		if err != nil {
			_ = received.Ack(false)
			continue
		}

		d.logger.Debug(
			LogMessage("received message: ", metadata.Type),
			zap.String("messageId", metadata.MessageId),
		)

		def, ok := d.consumersDefinition[msgType]

		if !ok {
			d.logger.Warn(
				LogMessage("could not find any consumer for this msg type"),
				zap.String("type", metadata.Type),
				zap.String("messageId", metadata.MessageId),
			)
			if err := received.Ack(false); err != nil {
				d.logger.Error(
					LogMessage("failed to ack msg"),
					zap.String("messageId", received.MessageId),
				)
			}
			continue
		}

		ctx, span := tracing.NewConsumerSpan(d.tracer, received.Headers, received.Type)

		ptr := def.reflect.Interface()
		if err = json.Unmarshal(received.Body, ptr); err != nil {
			span.RecordError(err)
			d.logger.Error(
				LogMessage("unmarshal error"),
				zap.String("messageId", received.MessageId),
				tracing.Format(ctx),
			)
			_ = received.Nack(true, false)
			span.End()
			continue
		}

		if def.queueDefinition.withRetry && metadata.XCount > def.queueDefinition.retires {
			d.logger.Warn(
				LogMessage("message reprocessed to many times, sending to dead letter"),
				tracing.Format(ctx),
			)
			_ = received.Ack(false)

			if err = d.publishToDlq(def, &received); err != nil {
				span.RecordError(err)
				d.logger.Error(
					LogMessage("failure to publish to dlq"),
					zap.String("messageId", received.MessageId),
					tracing.Format(ctx),
				)
			}

			span.End()
			continue
		}

		if err = def.handler(ctx, ptr, metadata); err != nil {
			d.logger.Error(
				LogMessage("error to process message"),
				zap.Error(err),
				tracing.Format(ctx),
			)

			if def.queueDefinition.withDLQ || err != RetryableError {
				span.RecordError(err)
				_ = received.Ack(false)

				if err = d.publishToDlq(def, &received); err != nil {
					span.RecordError(err)
					d.logger.Error(
						LogMessage("failure to publish to dlq"),
						zap.String("messageId", received.MessageId),
						tracing.Format(ctx),
					)
				}

				span.End()
				continue
			}

			d.logger.Warn(
				LogMessage("send message to process latter"),
				tracing.Format(ctx),
			)

			_ = received.Nack(false, false)
			span.End()
			continue
		}

		d.logger.Debug(LogMessage("message processed properly"), zap.String("messageId", received.MessageId), tracing.Format(ctx))
		_ = received.Ack(true)
		span.SetStatus(codes.Ok, "success")
		span.End()
	}
}

// extractMetadata extracts relevant metadata from an AMQP delivery.
// This includes the message ID, type, and retry count.
// Returns an error if the message has unformatted headers.
func (d *dispatcher) extractMetadata(delivery *amqp.Delivery) (*deliveryMetadata, error) {
	typ := delivery.Type
	if typ == "" {
		d.logger.Error(
			LogMessage("unformatted amqp delivery - missing type parameter"),
			zap.String("messageId", delivery.MessageId),
		)
		return nil, ReceivedMessageWithUnformattedHeaderError
	}

	var xCount int64
	if xDeath, ok := delivery.Headers["x-death"]; ok {
		v, _ := xDeath.([]interface{})
		table, _ := v[0].(amqp.Table)
		count, _ := table["count"].(int64)
		xCount = count
	}

	return &deliveryMetadata{
		MessageId: delivery.MessageId,
		Type:      typ,
		XCount:    xCount,
		Headers:   delivery.Headers,
	}, nil
}

// publishToDlq publishes a message to the dead-letter queue.
// It preserves the original message properties and headers.
func (m *dispatcher) publishToDlq(definition *ConsumerDefinition, received *amqp.Delivery) error {
	return m.channel.Publish("", definition.queueDefinition.dqlName, false, false, amqp.Publishing{
		Headers:     received.Headers,
		Type:        received.Type,
		ContentType: received.ContentType,
		MessageId:   received.MessageId,
		UserId:      received.UserId,
		AppId:       received.AppId,
		Body:        received.Body,
	})
}
