package rabbitmq

func News() Topology {
	return &topology{}
}

func (t *topology) Exchange(opts *ExchangeOpts) Topology {
	t.exchanges = append(t.exchanges, opts)

	return t
}

func (t *topology) FanoutExchanges(exchanges ...string) Topology {
	for _, name := range exchanges {
		t.exchanges = append(t.exchanges, &ExchangeOpts{name, FANOUT_EXCHANGE})
	}

	return t
}

func (s *topology) Queue(opts *QueueOpts) Topology {
	s.queues = append(s.queues, opts)

	return s
}

func (s *topology) GetQueueOpts(queue string) *QueueOpts {
	for _, opts := range s.queues {
		if opts.name == queue {
			return opts
		}
	}

	return nil
}
