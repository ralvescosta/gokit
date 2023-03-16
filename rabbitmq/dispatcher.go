package rabbitmq

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/ralvescosta/gokit/logging"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type (
	Dispatcher interface {
		Register(queue string, msg any, handler ConsumerHandler) error
		ConsumeBlocking(ch chan os.Signal)
	}

	dispatcher struct {
		logger              logging.Logger
		channel             AMQPChannel
		queueDefinitions    map[string]*QueueDefinition
		consumersDefinition map[string]*ConsumerDefinition
	}

	ConsumerHandler = func(ctx context.Context, msg any, metadata any) error

	ConsumerDefinition struct {
		queue           string
		msgType         string
		reflect         *reflect.Value
		queueDefinition *QueueDefinition
		handler         ConsumerHandler
	}

	deliveryMetadata struct {
		MessageId string
		XCount    int64
		Type      string
		Headers   map[string]interface{}
	}
)

func NewDispatcher(logger logging.Logger, channel AMQPChannel, queueDefinitions map[string]*QueueDefinition) *dispatcher {
	return &dispatcher{
		logger:              logger,
		channel:             channel,
		queueDefinitions:    queueDefinitions,
		consumersDefinition: map[string]*ConsumerDefinition{},
	}
}

func (d *dispatcher) Register(queue string, msg any, handler ConsumerHandler) error {
	if msg == nil || queue == "" {
		return InvalidDispatchParamsError
	}

	def, ok := d.queueDefinitions[queue]
	if !ok {
		return QueueDefinitionNotFoundError
	}

	elem := reflect.TypeOf(msg)
	ref := reflect.New(elem)
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

func (d *dispatcher) ConsumeBlocking(ch chan os.Signal) {
	for _, cd := range d.consumersDefinition {
		go d.consume(cd.queue, cd.msgType)
	}

	<-ch
}

func (d *dispatcher) consume(queue, msgType string) {
	delivery, err := d.channel.Consume(queue, msgType, false, false, false, false, nil)
	if err != nil {
		d.logger.Error(
			"failure to declare consumer",
			zap.String("queue", queue),
			zap.Error(err),
		)
		return
	}

	for received := range delivery {
		metadata, err := d.extractMetadata(&received)
		if err != nil {
			received.Ack(false)
			continue
		}

		d.logger.Debug(LogMessage("received message: ", metadata.Type, "messageId: ", metadata.MessageId))

	}
}

func (d *dispatcher) extractMetadata(delivery *amqp.Delivery) (*deliveryMetadata, error) {
	typ := delivery.Type
	if typ == "" {
		d.logger.Error(
			LogMessage("unformatted amqp delivery - missing type parameter"),
			zap.String("messageId", delivery.MessageId),
		)
		return nil, ReceivedMessageWithUnformattedHeaderError
	}

	var xCount int64 = 0
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
