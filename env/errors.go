package env

import "fmt"

type ConfigsError struct {
	msg string
}

func NewConfigsError(msg string) error {
	return &ConfigsError{
		msg: fmt.Sprintf("configs builder error - %s", msg),
	}
}

func (e *ConfigsError) Error() string {
	return e.msg
}

func NewErrRequiredConfig(env string) error {
	return NewConfigsError(fmt.Sprintf("%s is required", env))
}

var (
	ErrUnknownEnv = NewConfigsError("unknown env")
)
