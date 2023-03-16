package rabbitmq

type (
	AMQPConnection interface {
		Channel() (AMQPChannel, error)
	}
)
