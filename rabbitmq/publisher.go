// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/messaging"
	"github.com/ralvescosta/gokit/tracing"
	"go.uber.org/zap"
)

type (
	// publisher is the concrete implementation of the Publisher interface.
	// It handles the details of marshaling messages, setting headers, and publishing to RabbitMQ.
	publisher struct {
		logger  logging.Logger
		configs *configs.Configs
		channel AMQPChannel
	}
)

// JsonContentType is the MIME type used for JSON message content.
const (
	JsonContentType = "application/json"
)

// NewPublisher creates a new publisher instance with the provided configuration and AMQP channel.
func NewPublisher(configs *configs.Configs, channel AMQPChannel) messaging.Publisher {
	return &publisher{configs.Logger, configs, channel}
}

// SimplePublish publishes a message directly to a target queue.
// The exchange is left empty, which means the default exchange is used.
func (p *publisher) SimplePublish(ctx context.Context, target string, msg any) error {
	return p.publish(ctx, target, "", msg)
}

// Refactored Publish method to align with messaging.Publisher interface
func (p *publisher) Publish(ctx context.Context, to, from, key *string, msg any, options ...*messaging.Option) error {
	if to == nil || *to == "" {
		return fmt.Errorf("exchange cannot be empty")
	}

	exchange := *to
	routingKey := ""
	if key != nil {
		routingKey = *key
	}

	return p.publish(ctx, exchange, routingKey, msg)
}

// Refactored PublishDeadline method to align with messaging.Publisher interface
func (p *publisher) PublishDeadline(ctx context.Context, to, from, key *string, msg any, options ...*messaging.Option) error {
	if to == nil || *to == "" {
		return fmt.Errorf("exchange cannot be empty")
	}

	exchange := *to
	routingKey := ""
	if key != nil {
		routingKey = *key
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return p.publish(ctx, exchange, routingKey, msg)
}

// publish is the internal method that handles the details of publishing a message.
// It marshals the message to JSON, sets headers for tracing, and publishes to RabbitMQ.
func (p *publisher) publish(ctx context.Context, exchange, key string, msg any) error {
	byt, err := json.Marshal(msg)
	if err != nil {
		p.logger.Error(LogMessage("publisher marshal"), zap.Error(err))
		return err
	}

	headers := amqp.Table{}
	tracing.AMQPPropagator.Inject(ctx, tracing.AMQPHeader(headers))

	return p.channel.Publish(exchange, key, false, false, amqp.Publishing{
		Headers:     headers,
		Type:        fmt.Sprintf("%T", msg),
		ContentType: JsonContentType,
		MessageId:   uuid.NewString(),
		UserId:      p.configs.RabbitMQConfigs.User,
		AppId:       p.configs.AppConfigs.AppName,
		Body:        byt,
	})
}
