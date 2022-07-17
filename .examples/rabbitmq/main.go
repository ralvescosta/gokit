package main

import (
	"fmt"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/messaging/rabbitmq"
)

type ExampleMessage struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	cfg := &env.Configs{
		GO_ENV:          env.DEVELOPMENT_ENV,
		LOG_LEVEL:       env.DEBUG_L,
		APP_NAME:        "examples",
		RABBIT_HOST:     "localhost",
		RABBIT_PORT:     "5672",
		RABBIT_USER:     "admin",
		RABBIT_PASSWORD: "password",
		RABBIT_VHOST:    "",
	}

	log, _ := logging.NewDefaultLogger(cfg)

	topology := &rabbitmq.Topology{
		Queue: &rabbitmq.QueueOpts{
			Name: "my-service.try",
			Retryable: &rabbitmq.Retry{
				NumberOfRetry: 3,
				DelayBetween:  time.Duration(30) * time.Second,
			},
			WithDeadLatter: true,
		},
		Exchange: &rabbitmq.ExchangeOpts{
			Name: "ex-my-service.try",
			Type: rabbitmq.DIRECT_EXCHANGE,
		},
	}

	messaging, err := rabbitmq.
		New(cfg, log).
		Declare(topology).
		ApplyBinds().
		Build()

	if err != nil {
		log.Error(err.Error())
	}

	err = messaging.RegisterDispatcher(topology.Queue.Name, handler, &ExampleMessage{})

	if err != nil {
		log.Error(err.Error())
	}

	err = messaging.Consume()

	if err != nil {
		log.Error(err.Error())
	}
}

func handler(msg any, metadata *rabbitmq.DeliveryMetadata) error {
	c := msg.(*ExampleMessage)
	fmt.Println("EXECUTED")
	fmt.Println(c)
	// return errors.New("")
	return nil
}
