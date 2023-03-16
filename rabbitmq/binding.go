package rabbitmq

type (
	ExchangeBindingDefinition struct {
		source      string
		destination string
		routingKey  string
		args        map[string]interface{}
	}

	QueueBindingDefinition struct {
		routingKey string
		queue      string
		exchange   string
		args       map[string]interface{}
	}
)

func NewExchangeBiding() *ExchangeBindingDefinition {
	return &ExchangeBindingDefinition{}
}

func NewQueueBinding() *QueueBindingDefinition {
	return &QueueBindingDefinition{}
}

func (b *QueueBindingDefinition) RoutingKey(key string) *QueueBindingDefinition {
	b.routingKey = key
	return b
}

func (b *QueueBindingDefinition) Queue(name string) *QueueBindingDefinition {
	b.queue = name
	return b
}

func (b *QueueBindingDefinition) Exchange(name string) *QueueBindingDefinition {
	b.exchange = name
	return b
}
