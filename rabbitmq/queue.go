// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package rabbitmq

import (
	"fmt"
	"time"
)

// QueueDefinition represents the configuration of a RabbitMQ queue.
// It encapsulates properties such as name, durability, auto-delete behavior,
// exclusivity, TTL, DLQ (Dead Letter Queue), and retry mechanisms.
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
	retires   int64
}

// NewQueue creates a new queue definition with the given name.
// By default, queues are durable, not auto-deleted, and not exclusive.
func NewQueue(name string) *QueueDefinition {
	return &QueueDefinition{name: name, durable: true, delete: false, exclusive: false}
}

// Durable sets the durability flag for the queue.
// Durable queues survive broker restarts.
func (q *QueueDefinition) Durable(d bool) *QueueDefinition {
	q.durable = d
	return q
}

// Delete sets the auto-delete flag for the queue.
// Auto-deleted queues are removed when no longer in use.
func (q *QueueDefinition) Delete(d bool) *QueueDefinition {
	q.delete = d
	return q
}

// Exclusive sets the exclusive flag for the queue.
// Exclusive queues can only be used by the connection that created them
// and are deleted when that connection closes.
func (q *QueueDefinition) Exclusive(e bool) *QueueDefinition {
	q.exclusive = e
	return q
}

// WithTTL sets a Time-To-Live (TTL) for messages in the queue.
// Messages that remain in the queue longer than the TTL will be automatically removed.
func (q *QueueDefinition) WithTTL(ttl time.Duration) *QueueDefinition {
	q.withTTL = true
	q.ttl = ttl
	return q
}

// WithDQL enables a Dead Letter Queue (DLQ) for this queue.
// Messages that are rejected, expired, or exceed max length will be routed to the DLQ.
// The DLQ name is automatically generated as "<queue-name>-dlq".
func (q *QueueDefinition) WithDQL() *QueueDefinition {
	q.withDLQ = true
	q.dqlName = fmt.Sprintf("%s-dlq", q.name)
	return q
}

// WithRetry enables a retry mechanism for this queue.
// Failed messages will be moved to a retry queue for the specified TTL duration
// and then requeued up to the specified number of retries.
func (q *QueueDefinition) WithRetry(ttl time.Duration, retries int64) *QueueDefinition {
	q.withRetry = true
	q.retryTTL = ttl
	q.retires = retries
	return q
}

// DLQName returns the name of the Dead Letter Queue associated with this queue.
// The DLQ name follows the pattern "<queue-name>-dlq".
func (q *QueueDefinition) DLQName() string {
	return fmt.Sprintf("%s-dlq", q.name)
}

// RetryName returns the name of the Retry Queue associated with this queue.
// The Retry Queue name follows the pattern "<queue-name>-retry".
func (q *QueueDefinition) RetryName() string {
	return fmt.Sprintf("%s-retry", q.name)
}
