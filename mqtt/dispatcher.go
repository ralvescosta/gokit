package mqtt

import (
	"errors"

	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/logging"
)

type (
	Dispatcher interface{}

	dispatcher struct {
		logger      logging.Logger
		client      myQTT.Client
		subscribers map[string]myQTT.MessageHandler
	}
)

func NewDispatcher(logger logging.Logger, client myQTT.Client) Dispatcher {
	return &dispatcher{
		logger:      logger,
		client:      client,
		subscribers: make(map[string]myQTT.MessageHandler),
	}
}

func (d *dispatcher) Register(topic string, qos byte, handler myQTT.MessageHandler) error {
	if topic == "" {
		return errors.New("")
	}

	if handler == nil {
		return errors.New("")
	}

	//TODO: we need to store the qos to
	d.subscribers[topic] = handler

	return nil
}

func (d *dispatcher) ConsumeBlocking() {}
