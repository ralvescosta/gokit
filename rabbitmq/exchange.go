package rabbitmq

func Fanout(name string) *ExchangeOpts {
	return &ExchangeOpts{name, FANOUT_EXCHANGE}
}

func Direct(name string) *ExchangeOpts {
	return &ExchangeOpts{name, DIRECT_EXCHANGE}
}
