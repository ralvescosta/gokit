// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package mqtt

import (
	"context"
	"os"

	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type (
	// Dispatcher is an interface for managing MQTT subscriptions and consuming messages.
	Dispatcher interface {
		// Register adds a new subscription to the dispatcher with the specified topic, QoS, and handler.
		// Returns an error if the topic is empty, the handler is nil, or the QoS is invalid.
		Register(topic string, qos QoS, handler Handler) error

		// ConsumeBlocking starts consuming messages for all registered subscriptions.
		// Blocks until a signal is received on the provided channel, at which point it unsubscribes from all topics.
		ConsumeBlocking(ch chan os.Signal)
	}

	subscription struct {
		qos     QoS
		topic   string
		handler Handler
	}

	// Updated Handler type to include context.Context as the first argument.
	Handler = func(ctx context.Context, topic string, qos QoS, payload []byte) error

	// mqttDispatcher is the concrete implementation of the Dispatcher interface.
	mqttDispatcher struct {
		logger      logging.Logger
		client      myQTT.Client
		subscribers []*subscription
		tracer      trace.Tracer
	}
)

// NewDispatcher initializes a new mqttDispatcher with the provided logger and MQTT client.
func NewDispatcher(logger logging.Logger, client myQTT.Client) Dispatcher {
	tracer := otel.Tracer("gokit/mqtt")
	return &mqttDispatcher{
		logger:      logger,
		client:      client,
		subscribers: []*subscription{},
		tracer:      tracer,
	}
}

func (d *mqttDispatcher) Register(topic string, qos QoS, handler Handler) error {
	if topic == "" {
		return EmptyTopicError
	}

	if handler == nil {
		return NillHandlerError
	}

	if !ValidateQoS(qos) {
		return InvalidQoSError
	}

	d.subscribers = append(d.subscribers, &subscription{qos, topic, handler})

	return nil
}

func (d *mqttDispatcher) ConsumeBlocking(ch chan os.Signal) {
	for _, s := range d.subscribers {
		d.logger.Debug(LogMessage("subscribing to topic: ", s.topic))
		d.client.Subscribe(s.topic, 1, d.defaultMessageHandler(s.handler))
	}

	<-ch

	d.logger.Warn(LogMessage("received stop signal, unsubscribing..."))

	for _, s := range d.subscribers {
		d.logger.Warn(LogMessage("unsubscribing to topic: ", s.topic))
		d.client.Unsubscribe(s.topic)
	}

	d.logger.Debug(LogMessage("stopping consumer..."))
}

// defaultMessageHandler wraps a Handler with additional functionality, such as tracing.
func (d *mqttDispatcher) defaultMessageHandler(handler Handler) myQTT.MessageHandler {
	return func(_ myQTT.Client, msg myQTT.Message) {
		d.logger.Debug(LogMessage("received message from topic: ", msg.Topic()))
		msg.Ack()

		// Create a new context with an OpenTelemetry span using the dispatcher tracer.
		ctx, span := d.tracer.Start(context.Background(), msg.Topic())
		defer span.End()

		err := handler(ctx, msg.Topic(), QoSFromBytes(msg.Qos()), msg.Payload())
		if err != nil {
			d.logger.Error(LogMessage("failure to execute the topic handler"), zap.Error(err))
		}

		d.logger.Debug(LogMessage("message processed successfully"))
	}
}
