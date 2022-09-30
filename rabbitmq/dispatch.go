package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/ralvescosta/gokit/errors"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/tracing"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func NewDispatcher(logger logging.Logger, messaging Messaging, topology Topology) Dispatcher {
	return &dispatcherImpl{
		logger:    logger,
		messaging: messaging,
		topology:  topology,
		tracer:    otel.Tracer("dispatcher"),
	}
}

func (d *dispatcherImpl) RegisterDispatcher(queue string, msg any, handler ConsumerHandler) error {
	if msg == nil || queue == "" {
		return errors.ErrorAMQPRegisterDispatcher
	}

	elem := reflect.TypeOf(msg)
	ref := reflect.New(elem)

	d.queues = append(d.queues, queue)
	d.msgsTypes = append(d.msgsTypes, fmt.Sprintf("%T", msg))
	d.handlers = append(d.handlers, handler)
	d.reflectedTypes = append(d.reflectedTypes, &ref)

	return nil
}

func (d *dispatcherImpl) ConsumeBlocking(ch chan os.Signal) {
	for i, q := range d.queues {
		go d.consume(q, d.msgsTypes[i], d.reflectedTypes[i], d.handlers[i])
	}

	<-ch
}

func (d *dispatcherImpl) consume(queue, msgType string, reflected *reflect.Value, handler ConsumerHandler) {
	delivery, err := d.messaging.Channel().Consume(queue, msgType, false, false, false, false, nil)
	if err != nil {
		return
	}

	queueOpts := d.topology.GetQueueOpts(queue)

	for received := range delivery {
		metadata, err := d.extractMetadataFromDeliver(&received)
		if err != nil {
			received.Ack(false)
			continue
		}

		d.logger.Info(MessageType("message received: ", msgType, received.MessageId))

		valid := false
		for _, typ := range d.msgsTypes {
			if typ == metadata.Type {
				valid = true
				break
			}
		}

		if !valid {
			received.Ack(false)
			continue
		}

		ctx, span := tracing.NewConsumerSpan(d.tracer, received.Headers, received.Type)

		ptr := reflected.Interface()
		err = json.Unmarshal(received.Body, ptr)
		if err != nil {
			span.RecordError(err)
			d.logger.Error(Message("unmarshal error"), zap.String("messageId", received.MessageId), tracing.Format(ctx))
			received.Nack(true, false)
			span.End()
			continue
		}

		if queueOpts.retry != nil && metadata.XCount > queueOpts.retry.NumberOfRetry {
			d.logger.Warn("message reprocessed to many times, sending to dead letter", tracing.Format(ctx))
			received.Ack(false)
			d.publishToDlq(ctx, queueOpts, &received)
			span.End()
			continue
		}

		err = handler(ctx, ptr, metadata)
		if err != nil {
			if queueOpts.retry == nil || err != errors.ErrorAMQPRetryable {
				span.RecordError(err)
				received.Ack(false)
				d.publishToDlq(ctx, queueOpts, &received)
				span.End()
				continue
			}

			d.logger.Warn(Message("send message to process latter"), tracing.Format(ctx))

			received.Nack(false, false)
			span.End()
			continue
		}

		d.logger.Info(Message("message processed properly"), zap.String("messageId", received.MessageId), tracing.Format(ctx))
		received.Ack(true)
		span.SetStatus(codes.Ok, "success")
		span.End()
	}
}

func (m *dispatcherImpl) extractMetadataFromDeliver(delivery *amqp.Delivery) (*DeliveryMetadata, error) {
	typ := delivery.Type
	if typ == "" {
		m.logger.Error(Message("unformatted amqp delivery - missing type parameter - send message to DLQ"), zap.String("messageId", delivery.MessageId))
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	var xCount int64 = 0
	if xDeath, ok := delivery.Headers["x-death"]; ok {
		v, _ := xDeath.([]interface{})
		table, _ := v[0].(amqp.Table)
		count, _ := table["count"].(int64)
		xCount = count
	}

	return &DeliveryMetadata{
		MessageId: delivery.MessageId,
		Type:      typ,
		XCount:    xCount,
		Headers:   delivery.Headers,
	}, nil
}

func (m *dispatcherImpl) publishToDlq(ctx context.Context, queueOpts *QueueOpts, received *amqp.Delivery) error {
	return m.messaging.Channel().Publish("", queueOpts.DqlName(), false, false, amqp.Publishing{
		Headers:     received.Headers,
		Type:        received.Type,
		ContentType: received.ContentType,
		MessageId:   received.MessageId,
		UserId:      received.UserId,
		AppId:       received.AppId,
		Body:        received.Body,
	})
}
