package mqtt

import (
	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/logging"
)

type (
	Publisher interface {
		Pub(topic string, qos byte, payload any)
		PubRetained(topic string, qos byte, payload any)
	}

	publisher struct {
		logger logging.Logger
		client myQTT.Client
	}
)

func NewPublisher(logger logging.Logger, client myQTT.Client) Publisher {
	return &publisher{logger, client}
}

func (p *publisher) Pub(topic string, qos byte, payload any) {
	p.client.Publish(topic, qos, false, payload)
	p.logger.Debug(LogMessage("msg published on topic: ", topic))
}

func (p *publisher) PubRetained(topic string, qos byte, payload any) {
	p.client.Publish(topic, qos, true, payload)
	p.logger.Debug(LogMessage("msg published on topic: ", topic))
}
