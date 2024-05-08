package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
		tracer              trace.Tracer
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
		tracer:              otel.Tracer("dispatcher"),
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

	ref := reflect.New(reflect.TypeOf(msg))
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
			LogMessage("failure to declare consumer"),
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

		d.logger.Debug(
			LogMessage("received message: ", metadata.Type),
			zap.String("messageId", metadata.MessageId),
		)

		def, ok := d.consumersDefinition[msgType]

		if !ok {
			d.logger.Warn(
				LogMessage("could not find any consumer for this msg type"),
				zap.String("type", metadata.Type),
				zap.String("messageId", metadata.MessageId),
			)
			if err := received.Ack(false); err != nil {
				d.logger.Error(
					LogMessage("failed to ack msg"),
					zap.String("messageId", received.MessageId),
				)
			}
			continue
		}

		ctx, span := tracing.NewConsumerSpan(d.tracer, received.Headers, received.Type)

		ptr := def.reflect.Interface()
		if err = json.Unmarshal(received.Body, ptr); err != nil {
			span.RecordError(err)
			d.logger.Error(
				LogMessage("unmarshal error"),
				zap.String("messageId", received.MessageId),
				tracing.Format(ctx),
			)
			received.Nack(true, false)
			span.End()
			continue
		}

		if def.queueDefinition.withRetry && metadata.XCount > def.queueDefinition.retires {
			d.logger.Warn(
				LogMessage("message reprocessed to many times, sending to dead letter"),
				tracing.Format(ctx),
			)
			received.Ack(false)

			if err = d.publishToDlq(def, &received); err != nil {
				span.RecordError(err)
				d.logger.Error(
					LogMessage("failure to publish to dlq"),
					zap.String("messageId", received.MessageId),
					tracing.Format(ctx),
				)
			}

			span.End()
			continue
		}

		if err = def.handler(ctx, ptr, metadata); err != nil {
			d.logger.Error(
				LogMessage("error to process message"),
				zap.Error(err),
				tracing.Format(ctx),
			)

			if def.queueDefinition.withDLQ || err != RetryableError {
				span.RecordError(err)
				received.Ack(false)

				if err = d.publishToDlq(def, &received); err != nil {
					span.RecordError(err)
					d.logger.Error(
						LogMessage("failure to publish to dlq"),
						zap.String("messageId", received.MessageId),
						tracing.Format(ctx),
					)
				}

				span.End()
				continue
			}

			d.logger.Warn(
				LogMessage("send message to process latter"),
				tracing.Format(ctx),
			)

			received.Nack(false, false)
			span.End()
			continue
		}

		d.logger.Debug(LogMessage("message processed properly"), zap.String("messageId", received.MessageId), tracing.Format(ctx))
		received.Ack(true)
		span.SetStatus(codes.Ok, "success")
		span.End()
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

func (m *dispatcher) publishToDlq(definition *ConsumerDefinition, received *amqp.Delivery) error {
	return m.channel.Publish("", definition.queueDefinition.dqlName, false, false, amqp.Publishing{
		Headers:     received.Headers,
		Type:        received.Type,
		ContentType: received.ContentType,
		MessageId:   received.MessageId,
		UserId:      received.UserId,
		AppId:       received.AppId,
		Body:        received.Body,
	})
}
