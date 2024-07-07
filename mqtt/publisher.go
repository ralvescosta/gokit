package mqtt

import (
	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	Publisher interface {
		Pub(topic string, qos QoS, payload any) error
		PubRetained(topic string, qos QoS, payload any) error
	}

	mqttPublisher struct {
		logger logging.Logger
		client myQTT.Client
	}
)

func NewPublisher(logger logging.Logger, client myQTT.Client) Publisher {
	return &mqttPublisher{logger, client}
}

func (p *mqttPublisher) Pub(topic string, qos QoS, payload any) error {
	if err := p.validate(topic, qos, payload); err != nil {
		return err
	}

	p.client.Publish(topic, byte(qos), false, payload)
	p.logger.Debug(LogMessage("msg published on topic: ", topic))

	return nil
}

func (p *mqttPublisher) PubRetained(topic string, qos QoS, payload any) error {
	if err := p.validate(topic, qos, payload); err != nil {
		return err
	}

	p.client.Publish(topic, byte(qos), true, payload)
	p.logger.Debug(LogMessage("msg published on topic: ", topic))

	return nil
}

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
