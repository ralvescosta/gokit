// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package mqtt

import (
	"fmt"

	myQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

// MQTTClient defines the interface for an MQTT client.
type MQTTClient interface {
	// Connect establishes a connection to the MQTT broker.
	Connect() error
	// Client returns the underlying MQTT client instance.
	Client() myQTT.Client
}

// mqttClient is the concrete implementation of the MQTTClient interface.
type mqttClient struct {
	logger logging.Logger
	cfgs   *configs.Configs
	client myQTT.Client
}

// NewMQTTClient creates a new instance of mqttClient.
func NewMQTTClient(cfgs *configs.Configs) MQTTClient {
	return &mqttClient{
		cfgs:   cfgs,
		logger: cfgs.Logger,
	}
}

// Connect establishes a connection to the MQTT broker.
func (c *mqttClient) Connect() error {
	c.logger.Debug(LogMessage("connecting to the MQTT broker..."))

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
	if !token.Wait() || token.Error() != nil {
		c.logger.Error(LogMessage("connection failure"), zap.Error(token.Error()))
		return ConnectionFailureError
	}

	c.client = client
	c.logger.Debug(LogMessage("MQTT broker connected successfully"))
	return nil
}

// Client returns the underlying MQTT client instance.
func (c *mqttClient) Client() myQTT.Client {
	return c.client
}

// onConnectionEvent handles the MQTT broker connection event.
func (c *mqttClient) onConnectionEvent(_ myQTT.Client) {
	c.logger.Debug(LogMessage("connected to the MQTT broker"))
}

// onDisconnectEvent handles the MQTT broker disconnection event.
func (c *mqttClient) onDisconnectEvent(_ myQTT.Client, err error) {
	c.logger.Error(LogMessage("disconnected from the MQTT broker"), zap.Error(err))
}

// onReconnectionEvent handles the MQTT broker reconnection event.
func (c *mqttClient) onReconnectionEvent(_ myQTT.Client, _ *myQTT.ClientOptions) {
	c.logger.Debug(LogMessage("attempting to reconnect to the MQTT broker"))
}
