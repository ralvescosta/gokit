package rabbitmq

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ralvescosta/gokit/errors"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/tracing"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/codes"
)

func NewDispatcher(logger logging.ILogger, messaging Messaging, topology Topology) Dispatcher {
	return &dispatcher{
		logger:    logger,
		messaging: messaging,
		topology:  topology,
	}
}

func (d *dispatcher) RegisterDispatcher(queue string, msg any, handler ConsumerHandler) error {
	if msg == nil || queue == "" {
		return errors.ErrorAMQPRegisterDispatcher
	}

	ref := reflect.New(reflect.TypeOf(msg).Elem())

	d.queues = append(d.queues, queue)
	d.msgsTypes = append(d.msgsTypes, fmt.Sprintf("%T", msg))
	d.handlers = append(d.handlers, handler)
	d.reflectedTypes = append(d.reflectedTypes, &ref)

	return nil
}

func (d *dispatcher) ConsumeBlocking() {
	for i, q := range d.queues {
		go d.consume(q, d.msgsTypes[i], d.reflectedTypes[i], d.handlers[i])
	}
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

		d.logger.Info(LogMsgWithType("message received: ", msgType, received.MessageId))

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

		ptr := reflected.Interface()
		err = json.Unmarshal(received.Body, ptr)
		if err != nil {
			d.logger.Error(LogMsgWithMessageId("unmarshal error", received.MessageId))
			received.Nack(true, false)
			continue
		}

		_, span, err := tracing.SpanFromAMQPTraceparent(d.tracer, metadata.Traceparent, metadata.Type, received.Exchange, queue)
		if err != nil {
			d.logger.Error("could not create a span")
			continue
		}

		if queueOpts.retry != nil && metadata.XCount > queueOpts.retry.NumberOfRetry {
			d.logger.Warn("message reprocessed to many times, sending to dead letter")
			received.Nack(true, false)
			span.SetStatus(codes.Error, "dead letter")
			span.End()
			continue
		}

		err = handler(ptr, metadata)
		if err != nil {
			if queueOpts.retry == nil || err != errors.ErrorAMQPRetryable {
				received.Nack(true, false)
				span.SetStatus(codes.Error, "process failure without a retry")
				span.End()
				continue
			}

			d.logger.Warn(LogMessage("send message to process latter"))

			// d.publishToDelayed(ctx, metadata, d.Topology, &received)

			received.Ack(true)
			span.End()
			continue
		}

		d.logger.Info(LogMsgWithMessageId("message processed properly", received.MessageId))
		received.Ack(true)
		span.SetStatus(codes.Ok, "success")
		span.End()
	}
}

func (m *dispatcher) extractMetadataFromDeliver(delivery *amqp.Delivery) (*DeliveryMetadata, error) {
	msgID := delivery.MessageId
	if msgID == "" {
		m.logger.Error("unformatted amqp delivery - missing messageId parameter - send message to DLQ")
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	typ := delivery.Type
	if typ == "" {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing type parameter - send message to DLQ", delivery.MessageId))
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	xCount, ok := delivery.Headers[AMQPHeaderNumberOfRetry].(int64)
	if !ok {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing x-count header - send message to DLQ", delivery.MessageId))
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	traceparent, ok := delivery.Headers[AMQPHeaderTraceparent]
	if !ok {
		m.logger.Error(LogMsgWithMessageId("unformatted amqp delivery - missing x-trace-id header - send message to DLQ", delivery.MessageId))
		return nil, errors.ErrorAMQPReceivedMessageValidator
	}

	return &DeliveryMetadata{
		MessageId:   msgID,
		Type:        typ,
		XCount:      xCount,
		Traceparent: traceparent.(string),
		Headers:     delivery.Headers,
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
