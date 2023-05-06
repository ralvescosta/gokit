package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/tracing"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type (
	Publisher interface {
		SimplePublish(ctx context.Context, target string, msg any) error
		Publish(ctx context.Context, exchange, key string, msg any) error
	}

	publisher struct {
		logger  logging.Logger
		configs *configs.Configs
		channel AMQPChannel
	}
)

const (
	JsonContentType = "application/json"
)

func NewPublisher(logger logging.Logger, configs *configs.Configs, channel AMQPChannel) *publisher {
	return &publisher{logger, configs, channel}
}

func (p *publisher) SimplePublish(ctx context.Context, target string, msg any) error {
	return p.publish(ctx, target, "", msg)
}

func (p *publisher) Publish(ctx context.Context, exchange, key string, msg any) error {
	return p.publish(ctx, exchange, key, msg)
}

func (p *publisher) publish(ctx context.Context, exchange, key string, msg any) error {
	byt, err := json.Marshal(msg)
	if err != nil {
		p.logger.Error(LogMessage("publisher marshal"), zap.Error(err))
		return err
	}

	headers := amqp.Table{}
	tracing.AMQPPropagator.Inject(ctx, tracing.AMQPHeader(headers))

	return p.channel.Publish(exchange, key, false, false, amqp.Publishing{
		Headers:     headers,
		Type:        fmt.Sprintf("%T", msg),
		ContentType: JsonContentType,
		MessageId:   uuid.NewString(),
		UserId:      p.configs.RabbitMQConfigs.User,
		AppId:       p.configs.AppConfigs.AppName,
		Body:        byt,
	})
}
