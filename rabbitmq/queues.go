package rabbitmq

import "time"

func NewQueue(name string) *QueueOpts {
	return &QueueOpts{
		name: name,
	}
}

func (q *QueueOpts) WithDql() *QueueOpts {
	q.withDeadLatter = true

	return q
}

func (q *QueueOpts) WithRetry(numberOfTry int64, delayBetween time.Duration) *QueueOpts {
	q.retry = &Retry{
		NumberOfRetry: numberOfTry,
		DelayBetween:  delayBetween,
	}

	return q
}

func (q *QueueOpts) WithTTL(ttl time.Duration) *QueueOpts {
	q.ttl = ttl

	return q
}

func (q *QueueOpts) Binding(exchange, key string) *QueueOpts {
	q.bindings = append(q.bindings, &BindingOpts{exchange, key})

	return q
}

func (q *QueueOpts) BindingFanout(exchange string) *QueueOpts {
	q.bindings = append(q.bindings, &BindingOpts{exchange, ""})

	return q
}
