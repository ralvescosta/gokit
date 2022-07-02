package main

import (
	"errors"
	"fmt"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logger"
	"github.com/ralvescostati/pkgs/messaging/rabbitmq"
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

	log, _ := logger.NewDefaultLogger(cfg)

	exch := &rabbitmq.DeclareExchangeParams{
		ExchangeName: "example-exchange",
		ExchangeType: rabbitmq.DIRECT_EXCHANGE,
	}

	qe := &rabbitmq.DeclareQueueParams{
		QueueName:      "example-queue",
		WithDeadLatter: true,
		Retryable: &rabbitmq.Retry{
			NumberOfRetry: 3,
			DelayBetween:  400,
		},
	}

	bind := &rabbitmq.BindQueueParams{
		QueueName:    "example-queue",
		ExchangeName: "example-exchange",
	}

	messaging, err := rabbitmq.
		New(cfg, log).
		DeclareExchange(exch).
		DeclareQueue(qe).
		BindQueue(bind).
		Build()

	if err != nil {
		log.Error(err.Error())
	}

	err = messaging.RegisterDispatcher(qe.QueueName, handler, &ExampleMessage{})

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
	return errors.New("")
}
