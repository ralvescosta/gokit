package rabbitmq

type (
	Topology struct {
		channel          AMQPChannel
		queues           []*QueueDefinition
		exchanges        []*ExchangeDefinition
		exchangeBindings []*ExchangeBindingDefinition
		queueBindings    []*QueueBindingDefinition
	}
)

func NewTopology() *Topology {
	return &Topology{}
}

func (t *Topology) Channel(c AMQPChannel) *Topology {
	t.channel = c
	return t
}

func (t *Topology) Queue(q *QueueDefinition) *Topology {
	t.queues = append(t.queues, q)
	return t
}

func (t *Topology) Queues(q []*QueueDefinition) *Topology {
	t.queues = append(t.queues, q...)
	return t
}

func (t *Topology) Exchange(e *ExchangeDefinition) *Topology {
	t.exchanges = append(t.exchanges, e)
	return t
}

func (t *Topology) Exchanges(e []*ExchangeDefinition) *Topology {
	t.exchanges = append(t.exchanges, e...)
	return t
}

func (t *Topology) ExchangeBinding(b *ExchangeBindingDefinition) *Topology {
	t.exchangeBindings = append(t.exchangeBindings, b)
	return t
}

func (t *Topology) QueueBinding(b *QueueBindingDefinition) *Topology {
	t.queueBindings = append(t.queueBindings, b)
	return t
}

func (t *Topology) Apply() error {
	if t.channel == nil {
		return NullableChannel
	}

	for _, exch := range t.exchanges {
		if err := t.channel.ExchangeDeclare(exch.name, exch.kind.String(), exch.durable, exch.delete, false, false, exch.params); err != nil {
			return err
		}
	}

	// for _, queue := range t.queues {}

	return nil
}
