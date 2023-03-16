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
	NullableChannelError                      = NewRabbitMQError("channel cant be null")
	NotFoundQueueDefinitionError              = NewRabbitMQError("not found queue definition")
	InvalidDispatchParamsError                = NewRabbitMQError("register dispatch with invalid parameters")
	QueueDefinitionNotFoundError              = NewRabbitMQError("any queue definition was founded to the given queue")
	ReceivedMessageWithUnformattedHeaderError = NewRabbitMQError("received message with unformatted headers")
)
