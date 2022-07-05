package rabbitmq

import (
	"errors"
	"testing"

	// "github.com/ralvescostati/pkgs/messaging/rabbitmq/mock"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logging"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/suite"
)

type RabbitMQMessagingSuiteTest struct {
	suite.Suite

	amqpConn    *MockAMQPConnection
	amqpConnErr error
	amqpChannel *MockAMQPChannel
	messaging   *RabbitMQMessaging
}

func TestRabbitMQMessagingSuiteTest(t *testing.T) {
	suite.Run(t, new(RabbitMQMessagingSuiteTest))
}

func (s *RabbitMQMessagingSuiteTest) SetupTest() {
	s.amqpConn = NewMockAMQPConnection()
	s.amqpConnErr = nil
	s.amqpChannel = NewMockAMQPChannel()

	dial = func(cfg *env.Configs) (AMQPConnection, error) {
		return s.amqpConn, s.amqpConnErr
	}

	s.messaging = &RabbitMQMessaging{
		logger: logging.NewMockLogger(),
		conn:   s.amqpConn,
		ch:     s.amqpChannel,
		config: &env.Configs{},
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
