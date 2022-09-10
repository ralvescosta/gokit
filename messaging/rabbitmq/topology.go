package rabbitmq

func News() TopologyBuilder {
	return &Topology{}
}

func (t *Topology) Exchange(opts *ExchangeOpts) TopologyBuilder {
	t.exchanges = append(t.exchanges, opts)

	return t
}

func (t *Topology) FanoutExchanges(exchanges ...string) TopologyBuilder {
	for _, name := range exchanges {
		t.exchanges = append(t.exchanges, &ExchangeOpts{name, FANOUT_EXCHANGE})
	}

	return t
}

func (s *Topology) Queue(opts *QueueOpts) TopologyBuilder {
	s.queues = append(s.queues, opts)

	return s
}
