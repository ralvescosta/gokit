package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ralvescosta/gokit/errors"
	"github.com/ralvescosta/gokit/otel/trace"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/codes"
)

func (m *RabbitMQMessaging) RegisterDispatcher(queue string, handler ConsumerHandler, t any) error {
	if t == nil || queue == "" {
		return errors.ErrorAMQPRegisterDispatcher
	}

	var conf *Topology

	for _, v := range m.topologies {
		if v.Queue.Name == queue {
			conf = v
			break
		}
	}

	dispatch := &dispatcher{
		Queue:         queue,
		Topology:      conf,
		Handler:       handler,
		MsgType:       fmt.Sprintf("%T", t),
		ReflectedType: reflect.New(reflect.TypeOf(t).Elem()),
	}

	m.dispatchers = append(m.dispatchers, dispatch)

	return nil
}

func (m *RabbitMQMessaging) startConsumer(d *Dispatcher, shotdown chan error) {
	delivery, err := m.channel.Consume(d.Topology.Queue.Name, d.Topology.Binding.RoutingKey, false, false, false, false, nil)
	if err != nil {
		shotdown <- err
	}

	for received := range delivery {
		metadata, err := m.validateAndExtractMetadataFromDeliver(&received, d)
		if err != nil {
			received.Nack(true, false)
			continue
		}

		if metadata == nil {
			received.Nack(true, true)
			continue
		}

		m.logger.Info(LogMsgWithType("message received ", d.MsgType, received.MessageId))

		ptr := d.ReflectedType.Interface()
		err = json.Unmarshal(received.Body, ptr)
		if err != nil {
			m.logger.Error(LogMsgWithMessageId("unmarshal error", received.MessageId))
			received.Nack(true, false)
			continue
		}

		ctx, span, err := trace.SpanFromAMQPTraceparent(m.tracer, metadata.Traceparent, metadata.Type, d.Topology.Exchange.Name, d.Topology.Queue.Name)
		if err != nil {
			m.logger.Error("could not create a span")
			continue
		}

		if d.Topology.Queue.Retryable != nil && metadata.XCount > d.Topology.Queue.Retryable.NumberOfRetry {
			m.logger.Warn("message reprocessed to many times, sending to dead letter")
			received.Nack(true, false)
			span.SetStatus(codes.Error, "dead letter")
			span.End()
			continue
		}

		err = d.Handler(ptr, metadata)
		if err != nil {
			if d.Topology.Queue.Retryable == nil || err != errors.ErrorAMQPRetryable {
				received.Nack(true, false)
				span.SetStatus(codes.Error, "process failure without a retry")
				span.End()
				continue
			}

			m.logger.Warn(LogMessage("send message to process latter"))

			m.publishToDelayed(ctx, metadata, d.Topology, &received)

			received.Ack(true)
			span.End()
			continue
		}

		m.logger.Info(LogMsgWithMessageId("message processed properly", received.MessageId))
		received.Ack(true)
		span.SetStatus(codes.Ok, "success")
		span.End()
	}
}

func (m *RabbitMQMessaging) validateAndExtractMetadataFromDeliver(delivery *amqp.Delivery, d *Dispatcher) (*DeliveryMetadata, error) {
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

	if typ != d.MsgType {
		return nil, nil
	}

	return &DeliveryMetadata{
		MessageId:   msgID,
		Type:        typ,
		XCount:      xCount,
		Traceparent: traceparent.(string),
		Headers:     delivery.Headers,
	}, nil
}

func (m *RabbitMQMessaging) publishToDelayed(ctx context.Context, metadata *DeliveryMetadata, t *Topology, received *amqp.Delivery) error {

	return m.ch.Publish(t.delayed.ExchangeName, t.delayed.RoutingKey, false, false, amqp.Publishing{
		Headers: amqp.Table{
			AMQPHeaderNumberOfRetry: metadata.XCount + 1,
			AMQPHeaderTraceparent:   metadata.Traceparent,
			AMQPHeaderDelay:         t.Queue.Retryable.DelayBetween.Milliseconds(),
		},
		Type:        received.Type,
		ContentType: received.ContentType,
		MessageId:   received.MessageId,
		UserId:      received.UserId,
		AppId:       received.AppId,
		Body:        received.Body,
	})
}
