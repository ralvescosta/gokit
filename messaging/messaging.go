package messaging

import "context"

type (
	IMessageBroker[Params any] interface {
		Publisher(ctx context.Context, params *Params, msg any, opts map[string]any) error
		Subscriber(ctx context.Context, params *Params) error
		RegisterHandler(handler func(msg any, opts map[string]any) error, msgType any)
	}
)

func A[T any](a T) {}
