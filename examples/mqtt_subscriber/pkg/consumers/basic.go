// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package consumers

import (
	"context"
	"encoding/json"

	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/mqtt"
	"go.uber.org/zap"
)

type (
	BasicMessage struct{}

	BasicConsumer interface {
		Install(dispatcher mqtt.Dispatcher)
	}

	basicConsumer struct {
		logger logging.Logger
		topic  string
	}
)

func NewBasicMessage(logger logging.Logger, topic string) *basicConsumer {
	return &basicConsumer{logger, topic}
}

func (b *basicConsumer) Install(dispatcher mqtt.Dispatcher) {
	b.logger.Debug("Installing BasicConsumer...")
	dispatcher.Register(b.topic, mqtt.AtMostOnce, b.basicConsumer)
}

func (b *basicConsumer) basicConsumer(ctx context.Context, topic string, qos mqtt.QoS, payload []byte) error {
	basic := BasicMessage{}
	if err := json.Unmarshal(payload, &basic); err != nil {
		b.logger.Error("failed to unmarshal message", zap.Error(err))
		return err
	}

	b.logger.Info("Basic Consumer", zap.String("msg", basic.String()))
	return nil
}

func (BasicMessage) String() string {
	return "BasicMessage"
}
