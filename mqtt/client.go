package mqtt

import (
	"fmt"

	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	RabbitMQClient interface {
		Connect() error
		Client() myQTT.Client
	}

	rabbitMQClient struct {
		logger logging.Logger
		cfgs   *configs.Configs
		client myQTT.Client
	}
)

func NewMQTTClient(cfgs *configs.Configs, logger logging.Logger) RabbitMQClient {
	return &rabbitMQClient{
		cfgs:   cfgs,
		logger: logger,
	}
}

func (c *rabbitMQClient) Connect() error {
	c.logger.Debug(LogMessage("connecting to the mqtt broker..."))

	clientOpts := myQTT.NewClientOptions()

	clientOpts.AddBroker(fmt.Sprintf("%s://%s:%v", "tcp", c.cfgs.MQTTConfigs.Host, c.cfgs.MQTTConfigs.Port))
	clientOpts.SetUsername(c.cfgs.MQTTConfigs.User)
	clientOpts.SetPassword(c.cfgs.MQTTConfigs.Password)
	clientOpts.SetClientID(c.cfgs.AppConfigs.AppName)
	clientOpts.Order = false
	clientOpts.OnConnect = c.onConnectionEvent
	clientOpts.OnConnectionLost = c.onDisconnectEvent
	clientOpts.OnReconnecting = c.onReconnectionEvent

	client := myQTT.NewClient(clientOpts)

	token := client.Connect()
	if !token.Wait() {
		c.logger.Error(LogMessage("connection failure"))
		return ConnectionFailureError
	}

	c.client = client

	c.logger.Debug(LogMessage("mqtt broker was connected"))
	return nil
}

func (c *rabbitMQClient) Client() myQTT.Client {
	return c.client
}

func (c *rabbitMQClient) onConnectionEvent(clint myQTT.Client) {
	c.logger.Debug(LogMessage("received on connect event from mqtt broker"))
}

func (c *rabbitMQClient) onDisconnectEvent(clint myQTT.Client, err error) {
	c.logger.Error(LogMessage("received disconnect event from mqtt broker"), zap.Error(err))
}

func (c *rabbitMQClient) onReconnectionEvent(clint myQTT.Client, co *myQTT.ClientOptions) {
	c.logger.Debug(LogMessage("received reconnection event - trying to reconnect"))
}
