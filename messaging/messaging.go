package messaging

import "context"

type (
	Handler = func(msg any, opts map[string]any) error

	IMessageBroker[Params any] interface {
		Publisher(ctx context.Context, params *Params, msg any, opts map[string]any) error
		Subscriber(ctx context.Context, params *Params) error
		AddDispatcher(event string, handler Handler, msgType any) error
	}
)
