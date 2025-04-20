// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

type AWSConfigs struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

type DynamoDBConfigs struct {
	Endpoint string
	Region   string
	Table    string
}

type AWSSecretManagerConfigs struct {
	Region string
}
