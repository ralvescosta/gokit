// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// MQTTConfigs defines configuration parameters for MQTT broker connections.
// It provides settings for connection, authentication, and TLS/SSL security.
type MQTTConfigs struct {
	// Host specifies the MQTT broker hostname or IP address
	Host string
	// Port defines the network port on which the MQTT broker is listening
	Port int
	// User specifies the username for MQTT broker authentication
	User string
	// Password contains the authentication credential for the MQTT user
	Password string
	// Protocol specifies the MQTT protocol version to use (e.g., "mqtt", "mqtts")
	Protocol string

	// RootCaPath is the file path to the root CA certificate for TLS verification
	RootCaPath string
	// CertPath is the file path to the client certificate for mutual TLS authentication
	CertPath string
	// PrivateKeyPath is the file path to the private key associated with the client certificate
	PrivateKeyPath string
}
