package mqtt

import (
	"fmt"

	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	MQTTClient interface {
		Connect() error
		Client() myQTT.Client
	}

	mqttClient struct {
		logger logging.Logger
		cfgs   *configs.Configs
		client myQTT.Client
	}
)

func NewMQTTClient(cfgs *configs.Configs) MQTTClient {
	return &mqttClient{
		cfgs:   cfgs,
		logger: cfgs.Logger,
	}
}

func (c *mqttClient) Connect() error {
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

func (c *mqttClient) Client() myQTT.Client {
	return c.client
}

func (c *mqttClient) onConnectionEvent(_ myQTT.Client) {
	c.logger.Debug(LogMessage("received on connect event from mqtt broker"))
}

func (c *mqttClient) onDisconnectEvent(_ myQTT.Client, err error) {
	c.logger.Error(LogMessage("received disconnect event from mqtt broker"), zap.Error(err))
}

func (c *mqttClient) onReconnectionEvent(_ myQTT.Client, _ *myQTT.ClientOptions) {
	c.logger.Debug(LogMessage("received reconnection event - trying to reconnect"))
}
