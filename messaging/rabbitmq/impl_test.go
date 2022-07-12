package rabbitmq

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ralvescosta/toolkit/env"
	"github.com/ralvescosta/toolkit/logging"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RabbitMQMessagingSuiteTest struct {
	suite.Suite

	amqpConn    *MockAMQPConnection
	amqpConnErr error
	amqpChannel *MockAMQPChannel
	cfg         *env.Configs
	messaging   *RabbitMQMessaging
}

func TestRabbitMQMessagingSuiteTest(t *testing.T) {
	suite.Run(t, new(RabbitMQMessagingSuiteTest))
}

func (s *RabbitMQMessagingSuiteTest) SetupTest() {
	s.amqpConn = NewMockAMQPConnection()
	s.amqpConnErr = nil
	s.amqpChannel = NewMockAMQPChannel()
	s.cfg = &env.Configs{}

	dial = func(cfg *env.Configs) (AMQPConnection, error) {
		return s.amqpConn, s.amqpConnErr
	}

	s.messaging = &RabbitMQMessaging{
		logger: logging.NewMockLogger(),
		conn:   s.amqpConn,
		ch:     s.amqpChannel,
		config: s.cfg,
	}
}

func (s *RabbitMQMessagingSuiteTest) TestNew() {
	s.amqpConn.
		On("Channel").
		Return(&amqp.Channel{}, nil)

	msg := New(&env.Configs{}, logging.NewMockLogger())
	conn, err := msg.Build()

	s.NotNil(conn)
	s.NoError(err)
}

func (s *RabbitMQMessagingSuiteTest) TestNewConnErr() {
	s.amqpConnErr = errors.New("some err")

	msg := New(&env.Configs{}, logging.NewMockLogger())
	conn, err := msg.Build()

	s.Nil(conn)
	s.Error(err)
}

func (s *RabbitMQMessagingSuiteTest) TestNewChannelErr() {
	s.amqpConn.
		On("Channel").
		Return(&amqp.Channel{}, errors.New("some error"))

	msg := New(&env.Configs{}, logging.NewMockLogger())
	conn, err := msg.Build()

	s.Nil(conn)
	s.Error(err)
}

func (s *RabbitMQMessagingSuiteTest) TestDeclare() {
	s.messaging.Declare(&Topology{
		Exchange:   &ExchangeOpts{},
		Queue:      &QueueOpts{},
		Binding:    &BindingOpts{},
		isBindable: true,
	})

	s.NotNil(s.messaging.topologies)
	s.Len(s.messaging.topologies, 1)
}

func (s *RabbitMQMessagingSuiteTest) TestDeclareErr() {
	s.messaging.Err = errors.New("some error")

	s.messaging.Declare(&Topology{})

	s.Nil(s.messaging.topologies)
}

func (s *RabbitMQMessagingSuiteTest) TestBind() {
	tp := &Topology{
		Exchange: &ExchangeOpts{
			Name: "",
			Type: DIRECT_EXCHANGE,
		},
		Queue: &QueueOpts{
			WithDeadLatter: true,
			Retryable:      &Retry{},
		},
	}

	s.messaging.bind(tp)

	s.NotNil(tp.deadLetter)
	s.NotNil(tp.delayed)
}

func (s *RabbitMQMessagingSuiteTest) TestApplyBinds() {
	tp := &Topology{
		Exchange: &ExchangeOpts{
			Name: "",
			Type: DIRECT_EXCHANGE,
		},
		Queue: &QueueOpts{
			WithDeadLatter: true,
			Retryable:      &Retry{},
		},
	}

	s.messaging.Declare(tp).ApplyBinds()

	s.NotNil(tp.deadLetter)
	s.NotNil(tp.delayed)
}

func (s *RabbitMQMessagingSuiteTest) TestApplyBindsErr() {
	s.messaging.Err = errors.New("some error")
	tp := &Topology{}

	s.messaging.Declare(tp).ApplyBinds()

	s.Nil(tp.deadLetter)
	s.Nil(tp.delayed)
}

func (s *RabbitMQMessagingSuiteTest) TestBuild() {
	tp := &Topology{
		Exchange: &ExchangeOpts{
			Name:     "exchange",
			Type:     DIRECT_EXCHANGE,
			Bindings: []string{"to-bind"},
		},
		Queue: &QueueOpts{
			Name:           "queue",
			WithDeadLatter: true,
			Retryable: &Retry{
				NumberOfRetry: 3,
				DelayBetween:  10,
			},
		},
	}

	msg := s.messaging.Declare(tp).ApplyBinds()

	s.amqpChannel.
		On("ExchangeDeclare", tp.Exchange.Name, string(tp.Exchange.Type), true, false, false, false, amqp.Table(nil)).
		Return(nil).
		Once()
	s.amqpChannel.
		On("ExchangeDeclare", tp.delayed.ExchangeName, string(DELAY_EXCHANGE), true, false, false, false, amqp.Table{
			"x-delayed-type": "direct",
		}).
		Return(nil).
		Once()

	s.amqpChannel.
		On("ExchangeBind", tp.Exchange.Bindings[0], s.messaging.newRoutingKey(tp.Exchange.Name, tp.Exchange.Bindings[0]), tp.Exchange.Name, false, amqp.Table(nil)).
		Return(nil)

	s.amqpChannel.
		On("QueueDeclare", tp.deadLetter.QueueName, true, false, false, false, amqp.Table(nil)).
		Return(amqp.Queue{}, nil).
		Once()
	s.amqpChannel.
		On("QueueDeclare", tp.Queue.Name, true, false, false, false, amqp.Table{
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": tp.deadLetter.QueueName,
		}).
		Return(amqp.Queue{}, nil).
		Once()

	s.amqpChannel.
		On("QueueBind", tp.Queue.Name, tp.Binding.RoutingKey, tp.Exchange.Name, false, amqp.Table(nil)).
		Return(nil).
		Once()
	s.amqpChannel.
		On("QueueBind", tp.delayed.QueueName, tp.Binding.delayedRoutingKey, tp.delayed.ExchangeName, false, amqp.Table(nil)).
		Return(nil).
		Once()

	msg.Build()

	s.amqpChannel.AssertExpectations(s.T())
}

func (s *RabbitMQMessagingSuiteTest) TestBuildErr() {
	s.messaging.Err = errors.New("some error")
	tp := &Topology{}

	_, err := s.messaging.Declare(tp).ApplyBinds().Build()

	s.Error(err)
}

func (s *RabbitMQMessagingSuiteTest) TestBuildDeclareExchangeErr() {
	tp := &Topology{
		Exchange: &ExchangeOpts{
			Name:     "exchange",
			Type:     DIRECT_EXCHANGE,
			Bindings: []string{"to-bind"},
		},
		Queue: &QueueOpts{
			Name:           "queue",
			WithDeadLatter: true,
			Retryable: &Retry{
				NumberOfRetry: 3,
				DelayBetween:  10,
			},
		},
	}

	msg := s.messaging.Declare(tp).ApplyBinds()

	s.amqpChannel.
		On("ExchangeDeclare", tp.Exchange.Name, string(tp.Exchange.Type), true, false, false, false, amqp.Table(nil)).
		Return(errors.New("some error")).
		Once()

	_, err := msg.Build()

	s.Error(err)
	s.amqpChannel.AssertExpectations(s.T())

	s.amqpChannel.
		On("ExchangeDeclare", tp.Exchange.Name, string(tp.Exchange.Type), true, false, false, false, amqp.Table(nil)).
		Return(nil).
		Once()
	s.amqpChannel.
		On("ExchangeDeclare", tp.delayed.ExchangeName, string(DELAY_EXCHANGE), true, false, false, false, amqp.Table{
			"x-delayed-type": "direct",
		}).
		Return(errors.New("some error")).
		Once()

	_, err = msg.Build()

	s.Error(err)
	s.amqpChannel.AssertExpectations(s.T())
}

func (s *RabbitMQMessagingSuiteTest) TestBuildBindExchangeErr() {
	tp := &Topology{
		Exchange: &ExchangeOpts{
			Name:     "exchange",
			Type:     DIRECT_EXCHANGE,
			Bindings: []string{"to-bind"},
		},
		Queue: &QueueOpts{
			Name:           "queue",
			WithDeadLatter: true,
			Retryable: &Retry{
				NumberOfRetry: 3,
				DelayBetween:  10,
			},
		},
	}

	msg := s.messaging.Declare(tp).ApplyBinds()

	s.amqpChannel.
		On("ExchangeDeclare", tp.Exchange.Name, string(tp.Exchange.Type), true, false, false, false, amqp.Table(nil)).
		Return(nil).
		Once()
	s.amqpChannel.
		On("ExchangeDeclare", tp.delayed.ExchangeName, string(DELAY_EXCHANGE), true, false, false, false, amqp.Table{
			"x-delayed-type": "direct",
		}).
		Return(nil).
		Once()

	s.amqpChannel.
		On("ExchangeBind", tp.Exchange.Bindings[0], s.messaging.newRoutingKey(tp.Exchange.Name, tp.Exchange.Bindings[0]), tp.Exchange.Name, false, amqp.Table(nil)).
		Return(errors.New("some error"))

	_, err := msg.Build()

	s.amqpChannel.AssertExpectations(s.T())
	s.Error(err)
}

func (s *RabbitMQMessagingSuiteTest) TestPublisher() {
	exchange := "exchange"
	routingKey := "key"
	msg := make(map[string]interface{})

	s.amqpChannel.
		On("Publish", exchange, routingKey, false, false, mock.AnythingOfType("amqp.Publishing")).
		Return(nil).
		Once()

	err := s.messaging.Publisher(exchange, routingKey, msg, nil)

	s.NoError(err)
	s.amqpChannel.AssertExpectations(s.T())
}

// func (s *RabbitMQMessagingSuiteTest) TestPublisherErr() {
// 	exchange := "exchange"
// 	routingKey := "key"
// 	// msg := make(map[string]interface{})

// 	// s.amqpChannel.
// 	// 	On("Publish", exchange, routingKey, false, false, mock.AnythingOfType("amqp.Publishing")).
// 	// 	Return(nil).
// 	// 	Once()

// 	err := s.messaging.Publisher(exchange, routingKey, errors.New(""), nil)

// 	s.NoError(err)
// 	s.amqpChannel.AssertNotCalled(s.T(), "Publish")
// }

func (s *RabbitMQMessagingSuiteTest) TestRegisterDispatcher() {
	queue := "queue"
	handler := func(msg any, metadata *DeliveryMetadata) error {
		return nil
	}
	s.messaging.topologies = []*Topology{{
		Queue: &QueueOpts{
			Name: queue,
		},
	}}
	msg := make(map[string]interface{})

	err := s.messaging.RegisterDispatcher(queue, handler, msg)

	s.NoError(err)
	s.Len(s.messaging.dispatchers, 1)
}

func (s *RabbitMQMessagingSuiteTest) TestRegisterDispatcherErr() {
	queue := "queue"
	handler := func(msg any, metadata *DeliveryMetadata) error {
		return nil
	}

	msg := make(map[string]interface{})

	err := s.messaging.RegisterDispatcher("", handler, msg)

	s.Error(err)

	err = s.messaging.RegisterDispatcher(queue, handler, nil)

	s.Error(err)
}

func (s *RabbitMQMessagingSuiteTest) TestConsumer() {
	queue := "queue"
	key := "key"
	typ := "type"
	s.messaging.dispatchers = []*Dispatcher{{
		Queue: queue,
		Topology: &Topology{
			Queue: &QueueOpts{
				Name: queue,
			},
			Binding: &BindingOpts{
				RoutingKey: key,
			},
		},
		MsgType: typ,
	}}

	s.amqpChannel.
		On("Consume", queue, key, false, false, false, false, amqp.Table(nil)).
		Return(make(<-chan amqp.Delivery), errors.New("some error"))

	err := s.messaging.Consume()

	s.Error(err)
}

func (s *RabbitMQMessagingSuiteTest) TestConsumerErr() {
	s.messaging.Err = errors.New("some error")

	err := s.messaging.Consume()

	s.Error(err)
}

type MsgBody struct {
	Name string
}

func (s *RabbitMQMessagingSuiteTest) TestStartConsumer() {
	shotdown := make(chan error)
	d, rootChan, fakeDelivery := s.senary(nil)

	var deliveryChan <-chan amqp.Delivery = rootChan

	s.amqpChannel.
		On("Consume", d.Queue, d.Topology.Binding.RoutingKey, false, false, false, false, amqp.Table(nil)).
		Return(deliveryChan, nil)

	go s.messaging.startConsumer(d, shotdown)
	rootChan <- fakeDelivery
	rootChan = nil

	time.Sleep(time.Second * 1)
	s.amqpChannel.AssertExpectations(s.T())
}

func (s *RabbitMQMessagingSuiteTest) TestStartConsumerRetry() {
	d, rootChan, fakeDelivery := s.senary(ErrorRetryable)

	var deliveryChan <-chan amqp.Delivery = rootChan

	s.amqpChannel.
		On("Consume", d.Queue, d.Topology.Binding.RoutingKey, false, false, false, false, amqp.Table(nil)).
		Return(deliveryChan, nil)

	s.amqpChannel.
		On("Publish", d.Topology.Exchange.Name, d.Topology.Binding.RoutingKey, false, false, mock.AnythingOfType("amqp.Publishing")).
		Return(nil)

	shotdown := make(chan error)
	go s.messaging.startConsumer(d, shotdown)

	rootChan <- fakeDelivery

	time.Sleep(time.Second * 1)
	s.amqpChannel.AssertExpectations(s.T())
}

func (s *RabbitMQMessagingSuiteTest) TestStartConsumerRetryExceeded() {
	d, rootChan, fakeDelivery := s.senary(ErrorRetryable)

	var deliveryChan <-chan amqp.Delivery = rootChan

	s.amqpChannel.
		On("Consume", d.Queue, d.Topology.Binding.RoutingKey, false, false, false, false, amqp.Table(nil)).
		Return(deliveryChan, nil)

	shotdown := make(chan error)
	go s.messaging.startConsumer(d, shotdown)

	fakeDelivery.Headers[AMQPHeaderNumberOfRetry] = int64(4)
	rootChan <- fakeDelivery

	time.Sleep(time.Second * 1)
	s.amqpChannel.AssertNotCalled(s.T(), "Publish")
}

func (s *RabbitMQMessagingSuiteTest) TestValidateAndExtractMetadataFromDeliver() {
	delivery := &amqp.Delivery{
		MessageId: "id",
		Type:      "type",
		Headers: amqp.Table{
			AMQPHeaderNumberOfRetry: int64(0),
			AMQPHeaderTraceID:       "id",
		},
	}
	dispatcher := &Dispatcher{
		MsgType: "type",
	}

	m, err := s.messaging.validateAndExtractMetadataFromDeliver(delivery, dispatcher)
	s.NotNil(m)
	s.NoError(err)

	delivery.MessageId = ""
	m, err = s.messaging.validateAndExtractMetadataFromDeliver(delivery, dispatcher)
	s.Nil(m)
	s.Error(err)

	delivery.MessageId = "id"
	delivery.Type = ""
	m, err = s.messaging.validateAndExtractMetadataFromDeliver(delivery, dispatcher)
	s.Nil(m)
	s.Error(err)

	delivery.Type = "type"
	delivery.Headers = amqp.Table{}
	m, err = s.messaging.validateAndExtractMetadataFromDeliver(delivery, dispatcher)
	s.Nil(m)
	s.Error(err)

	delivery.Headers = amqp.Table{
		AMQPHeaderNumberOfRetry: int64(0),
	}
	m, err = s.messaging.validateAndExtractMetadataFromDeliver(delivery, dispatcher)
	s.Nil(m)
	s.Error(err)

	delivery.Type = "t"
	delivery.Headers = amqp.Table{
		AMQPHeaderNumberOfRetry: int64(0),
		AMQPHeaderTraceID:       "id",
	}
	m, err = s.messaging.validateAndExtractMetadataFromDeliver(delivery, dispatcher)
	s.Nil(m)
	s.NoError(err)
}

func (s *RabbitMQMessagingSuiteTest) senary(handlerErr error) (*Dispatcher, chan amqp.Delivery, amqp.Delivery) {
	queue := "queue"
	exch := "exchange"
	key := "key"
	typ := "type"
	msg := &MsgBody{}
	msgByt, _ := json.Marshal(msg)

	dispatcher := &Dispatcher{
		Queue: queue,
		Topology: &Topology{
			Queue: &QueueOpts{
				Name: queue,
				Retryable: &Retry{
					NumberOfRetry: 3,
					DelayBetween:  300,
				},
			},
			Exchange: &ExchangeOpts{
				Name: exch,
			},
			Binding: &BindingOpts{
				RoutingKey: key,
			},
			delayed: &DelayedOpts{
				QueueName:    queue,
				ExchangeName: exch,
				RoutingKey:   key,
			},
		},
		Handler: func(msg any, metadata *DeliveryMetadata) error {
			return handlerErr
		},
		MsgType:       typ,
		ReflectedType: reflect.ValueOf(msg),
	}

	rootChn := make(chan amqp.Delivery)

	delivery := amqp.Delivery{
		MessageId: "id",
		Type:      typ,
		UserId:    "id",
		AppId:     "id",
		Body:      msgByt,
		Headers: amqp.Table{
			AMQPHeaderNumberOfRetry: int64(0),
			AMQPHeaderDelay:         "20",
			AMQPHeaderTraceID:       "id",
		},
	}

	return dispatcher, rootChn, delivery
}
