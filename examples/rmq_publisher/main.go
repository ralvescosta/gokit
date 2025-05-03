// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package main

import (
	"context"
	"log"
	"time"

	configsBuilder "github.com/ralvescosta/gokit/configs_builder"
	"github.com/ralvescosta/gokit/rabbitmq"
)

func main() {
	cfgs, err := configsBuilder.
		NewConfigsBuilder().
		RabbitMQ().
		Build()
	if err != nil {
		cfgs.Logger.Fatal(err.Error())
	}

	_, channel, err := rabbitmq.NewConnection(cfgs)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ channel: %v", err)
	}

	publisher := rabbitmq.NewPublisher(cfgs, channel)
	for {
		message := "Hello World!"
		to := "exchange_name"
		key := "routing_key"

		err := publisher.Publish(context.Background(), &to, nil, &key, []byte(message))
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		} else {
			log.Printf("Published message: %s", message)
		}

		time.Sleep(1 * time.Second)
	}
}
