// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// KafkaConfigs defines configuration parameters for Apache Kafka connections.
// It provides settings for connection, authentication, and security protocols.
type KafkaConfigs struct {
	// Host specifies the Kafka broker hostname or IP address
	Host string
	// Port defines the network port on which the Kafka broker is listening
	Port int
	// SecurityProtocol defines the protocol used to communicate with brokers
	// (e.g., "PLAINTEXT", "SSL", "SASL_PLAINTEXT", "SASL_SSL")
	SecurityProtocol string
	// SASLMechanisms specifies the SASL mechanism to use for authentication
	// (e.g., "PLAIN", "SCRAM-SHA-256", "SCRAM-SHA-512")
	SASLMechanisms string
	// User specifies the username for Kafka broker authentication when using SASL
	User string
	// Password contains the authentication credential for the Kafka user
	Password string
}
