package rabbitmq

import (
	"fmt"
	"time"
)

type QueueDefinition struct {
	name      string
	durable   bool
	delete    bool
	exclusive bool
	withTTL   bool
	ttl       time.Duration
	withDLQ   bool
	dqlName   string
	withRetry bool
	retryTTL  time.Duration
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

func (q *QueueDefinition) WithTTL(ttl time.Duration) *QueueDefinition {
	q.withTTL = true
	q.ttl = ttl
	return q
}

func (q *QueueDefinition) WithDQL() *QueueDefinition {
	q.withDLQ = true
	q.dqlName = fmt.Sprintf("%s-dlq", q.name)
	return q
}

func (q *QueueDefinition) WithRetry(ttl time.Duration, retries int32) *QueueDefinition {
	q.withRetry = true
	q.retryTTL = ttl
	q.retires = retries
	return q
}

func (q *QueueDefinition) DLQName() string {
	return fmt.Sprintf("%s-dlq", q.name)
}

func (q *QueueDefinition) RetryName() string {
	return fmt.Sprintf("%s-retry", q.name)
}
