// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package mqtt

import (
	"context"
	"time"

	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/messaging"
	"go.uber.org/zap"
)

// mqttPublisher is the concrete implementation of the messaging.Publisher interface.
type mqttPublisher struct {
	logger logging.Logger
	client myQTT.Client
}

// NewPublisher creates a new instance of mqttPublisher.
func NewPublisher(configs *configs.Configs, client myQTT.Client) messaging.Publisher {
	return &mqttPublisher{logger: configs.Logger, client: client}
}

// Refactored Publish method to align with messaging.Publisher interface
func (p *mqttPublisher) Publish(ctx context.Context, to, from, key *string, msg any, options ...*messaging.Option) error {
	if to == nil || *to == "" {
		return EmptyTopicError
	}

	topic := *to
	qos := p.qosFromOptions(options...)
	retain := p.retainFromOptions(options...)

	if err := p.validate(topic, qos, msg); err != nil {
		p.logger.Error(LogMessage("validation error"), zap.String("topic", topic), zap.Error(err))
	}

	token := p.client.Publish(topic, byte(qos), retain, msg)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Refactored PublishDeadline method to align with messaging.Publisher interface
func (p *mqttPublisher) PublishDeadline(ctx context.Context, to, from, key *string, msg any, options ...*messaging.Option) error {
	if to == nil || *to == "" {
		return EmptyTopicError
	}

	topic := *to
	qos := p.qosFromOptions(options...)
	retain := p.retainFromOptions(options...)

	if err := p.validate(topic, qos, msg); err != nil {
		p.logger.Error(LogMessage("validation error"), zap.String("topic", topic), zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	token := p.client.Publish(topic, byte(qos), retain, msg)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		if token.Error() != nil {
			return token.Error()
		}
	}

	return nil
}

// validate checks the validity of the topic, QoS, and payload for publishing.
func (p *mqttPublisher) validate(topic string, qos QoS, payload any) error {
	if !ValidateQoS(qos) {
		p.logger.Error(LogMessage("invalid qos"), zap.Int("qos", int(qos)))
		return InvalidQoSError
	}

	if topic == "" {
		p.logger.Error(LogMessage("invalid topic"), zap.String("topic", topic))
		return EmptyTopicError
	}

	if payload == nil {
		p.logger.Error(LogMessage("empty payload"))
		return NillPayloadError
	}

	return nil
}

func (p *mqttPublisher) qosFromOptions(options ...*messaging.Option) QoS {
	for _, option := range options {
		if option.Key == "qos" {
			switch option.Value {
			case "0":
				return QoS(0)
			case "1":
				return QoS(1)
			case "2":
				return QoS(2)
			default:
				return QoS(0)
			}
		}
	}

	// Default QoS if not specified
	return QoS(0)
}

func (p *mqttPublisher) retainFromOptions(options ...*messaging.Option) bool {
	for _, option := range options {
		if option.Key == "retain" {
			return option.Value == "true"
		}
	}

	// Default retain value if not specified
	return false
}
