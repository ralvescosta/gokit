package rabbitmq

type RabbitMQError struct {
	msg string
}

func (e *RabbitMQError) Error() string {
	return e.msg
}

func NewRabbitMQError(msg string) error {
	return &RabbitMQError{msg}
}

var (
	NullableChannel = NewRabbitMQError("channel cant be null")
)
