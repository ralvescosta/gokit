package rabbitmq

import (
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
	return &dispatcher{
		logger:    logger,
		messaging: messaging,
		topology:  topology,
		tracer:    otel.Tracer("dispatcher"),
	}
}

func (d *dispatcher) RegisterDispatcher(queue string, msg any, handler ConsumerHandler) error {
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

func (d *dispatcher) ConsumeBlocking(ch chan os.Signal) {
	for i, q := range d.queues {
		go d.consume(q, d.msgsTypes[i], d.reflectedTypes[i], d.handlers[i])
	}

	<-ch
}

func (d *dispatcher) consume(queue, msgType string, reflected *reflect.Value, handler ConsumerHandler) {
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

		ctx, span := tracing.NewConsumerSpan(d.tracer, metadata.Headers, metadata.Type)

		ptr := reflected.Interface()
		err = json.Unmarshal(received.Body, ptr)
		if err != nil {
			span.RecordError(err)
			d.logger.Error(Message("unmarshal error"), zap.String("messageId", metadata.MessageId), tracing.Format(ctx))
			received.Nack(true, false)
			span.End()
			continue
		}

		if queueOpts.retry != nil && metadata.XCount > queueOpts.retry.NumberOfRetry {
			d.logger.Warn("message reprocessed to many times, sending to dead letter", tracing.Format(ctx))
			span.RecordError(err)
			received.Nack(true, false)
			span.End()
			continue
		}

		err = handler(ctx, ptr, metadata)
		if err != nil {
			if queueOpts.retry == nil || err != errors.ErrorAMQPRetryable {
				span.RecordError(err)
				received.Nack(true, false)
				span.End()
				continue
			}

			d.logger.Warn(Message("send message to process latter"), tracing.Format(ctx))

			received.Ack(true)
			span.End()
			continue
		}

		d.logger.Info(Message("message processed properly"), zap.String("messageId", metadata.MessageId), tracing.Format(ctx))
		received.Ack(true)
		span.SetStatus(codes.Ok, "success")
		span.End()
	}
}

func (m *dispatcher) extractMetadataFromDeliver(delivery *amqp.Delivery) (*DeliveryMetadata, error) {
	typ := delivery.Type
	if typ == "" {
		m.logger.Error(Message("unformatted amqp delivery - missing type parameter - send message to DLQ"), zap.String("messageId", delivery.MessageId))
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	xCount, ok := delivery.Headers[AMQPHeaderNumberOfRetry].(int64)
	if !ok {
		m.logger.Error(Message("unformatted amqp delivery - missing x-count header - send message to DLQ"), zap.String("messageId", delivery.MessageId))
	}

	return &DeliveryMetadata{
		MessageId: delivery.MessageId,
		Type:      typ,
		XCount:    xCount,
		Headers:   delivery.Headers,
	}, nil
}

// func (m *RabbitMQMessaging) publishToDelayed(ctx context.Context, metadata *DeliveryMetadata, t *Topology, received *amqp.Delivery) error {

// 	return m.ch.Publish(t.delayed.ExchangeName, t.delayed.RoutingKey, false, false, amqp.Publishing{
// 		Headers: amqp.Table{
// 			AMQPHeaderNumberOfRetry: metadata.XCount + 1,
// 			AMQPHeaderTraceparent:   metadata.Traceparent,
// 			AMQPHeaderDelay:         t.Queue.Retryable.DelayBetween.Milliseconds(),
// 		},
// 		Type:        received.Type,
// 		ContentType: received.ContentType,
// 		MessageId:   received.MessageId,
// 		UserId:      received.UserId,
// 		AppId:       received.AppId,
// 		Body:        received.Body,
// 	})
// }
