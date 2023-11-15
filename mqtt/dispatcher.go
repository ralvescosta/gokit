package mqtt

import (
	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	MQTTDispatcher interface{}

	subscription struct {
		qos     QoS
		topic   string
		handler Handler
	}

	Handler = func(topic string, qos QoS, payload []byte) error

	mqttDispatcher struct {
		logger      logging.Logger
		client      myQTT.Client
		subscribers []*subscription
	}
)

func NewDispatcher(logger logging.Logger, client myQTT.Client) MQTTDispatcher {
	return &mqttDispatcher{
		logger:      logger,
		client:      client,
		subscribers: []*subscription{},
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

func (d *mqttDispatcher) ConsumeBlocking() {
	for _, s := range d.subscribers {
		d.client.Subscribe(s.topic, byte(s.qos), d.defaultMessageHandler(s.handler))
	}
}

func (d *mqttDispatcher) defaultMessageHandler(handler Handler) myQTT.MessageHandler {
	return func(client myQTT.Client, msg myQTT.Message) {
		d.logger.Debug(LogMessage("received message from topic: ", msg.Topic()))
		msg.Ack()

		err := handler(msg.Topic(), QoSFromBytes(msg.Qos()), msg.Payload())
		if err != nil {
			d.logger.Error(LogMessage("failure to execute the topic handler"), zap.Error(err))
		}
	}
}
