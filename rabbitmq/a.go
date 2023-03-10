package rabbitmq

func Ex() {
	NewTopology().
		Exchange(NewDirectExchange("batatinha")).
		Exchange(NewFanoutExchange("other batatinha")).
		Queue(NewQueue("b").WithDQL().WithRetry(1000, 10)).
		QueueBinding(NewQueueBinding().RoutingKey("key").Queue("b").Exchange("batatinha")).
		Apply()
}
