package rabbitmq

func NewFanout(name string) *ExchangeOpts {
	return &ExchangeOpts{name, FANOUT_EXCHANGE}
}

func NewDirect(name string) *ExchangeOpts {
	return &ExchangeOpts{name, DIRECT_EXCHANGE}
}
