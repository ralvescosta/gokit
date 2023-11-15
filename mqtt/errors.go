package mqtt

type MQTTError struct {
	msg string
}

func (e *MQTTError) Error() string {
	return e.msg
}

func NewMQTTError(msg string) error {
	return &MQTTError{msg}
}

var (
	ConnectionFailureError = NewMQTTError("connection failure")
	EmptyTopicError        = NewMQTTError("subscribe top cannot be empty string")
	NillHandlerError       = NewMQTTError("subscribe handler cannot be nil")
	NillPayloadError       = NewMQTTError("publish payload cannot be nil")
	InvalidQoSError        = NewMQTTError("qos must be one of: byte(0), byte(1) or byte(2)")
)
