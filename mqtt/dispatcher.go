package mqtt

import (
	"os"

	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	MQTTDispatcher interface {
		Register(topic string, qos QoS, handler Handler) error
		ConsumeBlocking(ch chan os.Signal)
	}

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

func NewMQTTDispatcher(logger logging.Logger, client myQTT.Client) MQTTDispatcher {
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

func (d *mqttDispatcher) defaultMessageHandler(handler Handler) myQTT.MessageHandler {
	return func(_ myQTT.Client, msg myQTT.Message) {
		d.logger.Debug(LogMessage("received message from topic: ", msg.Topic()))
		msg.Ack()

		err := handler(msg.Topic(), QoSFromBytes(msg.Qos()), msg.Payload())
		if err != nil {
			d.logger.Error(LogMessage("failure to execute the topic handler"), zap.Error(err))
		}

		d.logger.Debug(LogMessage("message processed successfully"))
	}
}
