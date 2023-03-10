package rabbitmq

import "fmt"

type QueueDefinition struct {
	name      string
	durable   bool
	delete    bool
	exclusive bool
	withTTL   bool
	ttl       int32
	withDLQ   bool
	dqlName   string
	withRetry bool
	retryTTL  int32
	retires   int32
}

func NewQueue(name string) *QueueDefinition {
	return &QueueDefinition{name: name, durable: true, delete: false, exclusive: false}
}

func (q *QueueDefinition) Durable(d bool) *QueueDefinition {
	q.durable = d
	return q
}

func (q *QueueDefinition) Delete(d bool) *QueueDefinition {
	q.delete = d
	return q
}

func (q *QueueDefinition) Exclusive(e bool) *QueueDefinition {
	q.exclusive = e
	return q
}

func (q *QueueDefinition) WithTTL(ttl int32) *QueueDefinition {
	q.withTTL = true
	q.ttl = ttl
	return q
}

func (q *QueueDefinition) WithDQL() *QueueDefinition {
	q.withDLQ = true
	q.dqlName = fmt.Sprintf("%s-dlq", q.name)
	return q
}

func (q *QueueDefinition) WithRetry(ttl, retries int32) *QueueDefinition {
	q.withRetry = true
	q.retryTTL = ttl
	q.retires = retries
	return q
}
