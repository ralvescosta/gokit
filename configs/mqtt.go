// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

type MQTTConfigs struct {
	Host     string
	Port     int
	User     string
	Password string
	Protocol string

	RootCaPath     string
	CertPath       string
	PrivateKeyPath string
}
