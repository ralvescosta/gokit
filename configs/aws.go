// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// AWSConfigs defines authentication and credential settings for AWS services.
// It contains the necessary parameters to authenticate with AWS APIs.
type AWSConfigs struct {
	// AccessKeyID is the AWS access key part of the credential pair
	AccessKeyID string
	// SecretAccessKey is the AWS secret key part of the credential pair
	SecretAccessKey string
	// SessionToken provides temporary credentials when using AWS STS (Security Token Service)
	SessionToken string
}

// DynamoDBConfigs provides configuration settings specific to Amazon DynamoDB.
// It contains connection and targeting parameters for DynamoDB operations.
type DynamoDBConfigs struct {
	// Endpoint specifies the DynamoDB service endpoint URL (useful for local development)
	Endpoint string
	// Region defines the AWS region where the DynamoDB table is located
	Region string
	// Table specifies the default DynamoDB table name to use
	Table string
}

// AWSSecretManagerConfigs contains settings for AWS Secret Manager service.
// It provides configuration needed to retrieve secrets from AWS Secret Manager.
type AWSSecretManagerConfigs struct {
	// Region defines the AWS region where the Secret Manager service is located
	Region string
}
