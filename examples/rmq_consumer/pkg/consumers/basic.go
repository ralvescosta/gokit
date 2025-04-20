// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package consumers

import (
	"context"

	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/rabbitmq"
	"go.uber.org/zap"
)

type (
	BasicMessage struct{}

	BasicConsumer interface {
		Install(dispatcher rabbitmq.Dispatcher)
	}

	basicConsumer struct {
		logger    logging.Logger
		queueName string
	}
)

func NewBasicMessage(logger logging.Logger, queue string) *basicConsumer {
	return &basicConsumer{logger, queue}
}

func (b *basicConsumer) Install(dispatcher rabbitmq.Dispatcher) {
	b.logger.Debug("Installing BasicConsumer...")
	dispatcher.Register(b.queueName, BasicMessage{}, b.basicConsumer)
}

func (b *basicConsumer) basicConsumer(ctx context.Context, msg any, metadata any) error {
	basic := msg.(BasicMessage)

	b.logger.Info("Basic Consumer", zap.String("msg", basic.String()))
	return nil
}

func (BasicMessage) String() string {
	return "BasicMessage"
}
