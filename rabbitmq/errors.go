package rabbitmq

import "errors"

var (
	Error                             = errors.New("error")
	ErrorAMQPBadTraceparent           = errors.New("bad traceparent")
	ErrorAMQPConnection               = errors.New("messaging failure to connect to rabbitmq")
	ErrorAMQPChannel                  = errors.New("messaging error to stablish amqp channel")
	ErrorAMQPRegisterDispatcher       = errors.New("messaging unformatted dispatcher params")
	ErrorAMQPRetryable                = errors.New("messaging failure to process send to retry latter")
	ErrorAMQPReceivedMessageValidator = errors.New("messaging unformatted received message")
	ErrorAMQPQueueDeclaration         = errors.New("to use dql feature the bind exchanges must be declared first")
)
