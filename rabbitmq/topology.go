package rabbitmq

func NewTopology() Topology {
	return &topologyImpl{}
}

func (t *topologyImpl) Exchange(opts *ExchangeOpts) Topology {
	t.exchanges = append(t.exchanges, opts)

	return t
}

func (t *topologyImpl) FanoutExchanges(exchanges ...string) Topology {
	for _, name := range exchanges {
		t.exchanges = append(t.exchanges, &ExchangeOpts{name, FANOUT_EXCHANGE})
	}

	return t
}

func (s *topologyImpl) Queue(opts *QueueOpts) Topology {
	s.queues = append(s.queues, opts)

	return s
}

func (s *topologyImpl) GetQueueOpts(queue string) *QueueOpts {
	for _, opts := range s.queues {
		if opts.name == queue {
			return opts
		}
	}

	return nil
}
