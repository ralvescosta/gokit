// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

type KafkaConfigs struct {
	Host             string
	Port             int
	SecurityProtocol string
	SASLMechanisms   string
	User             string
	Password         string
}
