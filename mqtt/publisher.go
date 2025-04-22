// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package mqtt

import (
	"context"
	"time"

	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

// Publisher defines the interface for publishing messages to MQTT topics.
type Publisher interface {
	// PubCtx publishes a message to the specified topic with the given QoS.
	// This method handles context timeout to ensure the operation does not wait longer than 1 second.
	// Returns an error if the topic, QoS, or payload is invalid.
	PubCtx(ctx context.Context, topic string, qos QoS, payload any) error

	// PubRetainedCtx publishes a retained message to the specified topic with the given QoS.
	// This method handles context timeout to ensure the operation does not wait longer than 1 second.
	// Returns an error if the topic, QoS, or payload is invalid.
	PubRetainedCtx(ctx context.Context, topic string, qos QoS, payload any) error

	// Pub publishes a message to the specified topic with the given QoS.
	// This method does not handle context deadlines or errors from client.Publish.
	Pub(topic string, qos QoS, payload any) error

	// PubRetained publishes a retained message to the specified topic with the given QoS.
	// This method does not handle context deadlines or errors from client.Publish.
	PubRetained(topic string, qos QoS, payload any) error
}

// mqttPublisher is the concrete implementation of the Publisher interface.
type mqttPublisher struct {
	logger logging.Logger
	client myQTT.Client
}

// NewPublisher creates a new instance of mqttPublisher.
func NewPublisher(logger logging.Logger, client myQTT.Client) Publisher {
	return &mqttPublisher{logger, client}
}

// PubCtx publishes a message to the specified topic with the given QoS.
// This method handles context timeout to ensure the operation does not wait longer than 1 second.
func (p *mqttPublisher) PubCtx(ctx context.Context, topic string, qos QoS, payload any) error {
	if err := p.validate(topic, qos, payload); err != nil {
		return err
	}

	// Create a context with a 1-second deadline.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	token := p.client.Publish(topic, byte(qos), false, payload)
	select {
	case <-ctx.Done():
		p.logger.Error(LogMessage("publish operation timed out"), zap.String("topic", topic))
		return ctx.Err()
	case <-token.Done():
		if token.Error() != nil {
			p.logger.Error(LogMessage("failed to publish message"), zap.Error(token.Error()))
			return token.Error()
		}
	}

	p.logger.Debug(LogMessage("msg published on topic: ", topic))
	return nil
}

// PubRetainedCtx publishes a retained message to the specified topic with the given QoS.
// This method handles context timeout to ensure the operation does not wait longer than 1 second.
func (p *mqttPublisher) PubRetainedCtx(ctx context.Context, topic string, qos QoS, payload any) error {
	if err := p.validate(topic, qos, payload); err != nil {
		return err
	}

	// Create a context with a 1-second deadline.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	token := p.client.Publish(topic, byte(qos), true, payload)
	select {
	case <-ctx.Done():
		p.logger.Error(LogMessage("publish operation timed out"), zap.String("topic", topic))
		return ctx.Err()
	case <-token.Done():
		if token.Error() != nil {
			p.logger.Error(LogMessage("failed to publish retained message"), zap.Error(token.Error()))
			return token.Error()
		}
	}

	p.logger.Debug(LogMessage("msg published on topic: ", topic))
	return nil
}

// Pub publishes a message to the specified topic with the given QoS.
// This method does not handle context deadlines or errors from client.Publish.
func (p *mqttPublisher) Pub(topic string, qos QoS, payload any) error {
	if err := p.validate(topic, qos, payload); err != nil {
		return err
	}

	p.client.Publish(topic, byte(qos), false, payload)
	p.logger.Debug(LogMessage("msg published on topic: ", topic))
	return nil
}

// PubRetained publishes a retained message to the specified topic with the given QoS.
// This method does not handle context deadlines or errors from client.Publish.
func (p *mqttPublisher) PubRetained(topic string, qos QoS, payload any) error {
	if err := p.validate(topic, qos, payload); err != nil {
		return err
	}

	p.client.Publish(topic, byte(qos), true, payload)
	p.logger.Debug(LogMessage("msg published on topic: ", topic))
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
