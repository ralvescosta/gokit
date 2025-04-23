// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package secretsmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

// awsSecretClient is an implementation of the SecretClient interface that uses
// AWS Secrets Manager to store and retrieve secrets.
type awsSecretClient struct {
	logger      logging.Logger
	client      *secretsmanager.Client
	appSecretId string
	secrets     map[string]string
}

// NewAwsSecretClient creates a new instance of AWS Secrets Manager client.
// It initializes the AWS configuration and prepares the secret identifier based on
// the application environment and secret key.
// Returns a SecretClient interface and any error encountered during initialization.
func NewAwsSecretClient(cfgs *configs.Configs) (SecretClient, error) {
	logger := cfgs.Logger

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("error get aws configs from env", zap.Error(err))
		return nil, err
	}

	appSecretId := fmt.Sprintf("%s/%s", cfgs.AppConfigs.GoEnv.ToString(), cfgs.AppConfigs.SecretKey)

	return &awsSecretClient{client: secretsmanager.NewFromConfig(awsCfg), appSecretId: appSecretId}, nil
}

// LoadSecrets retrieves all secrets from AWS Secrets Manager for the configured secret ID.
// It fetches the secret value as a JSON blob and unmarshals it into a map of string keys to string values.
// The secrets are stored in memory for fast access by the GetSecret method.
// Returns an error if the secret cannot be fetched or parsed.
func (c *awsSecretClient) LoadSecrets(ctx context.Context) error {
	res, err := c.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &c.appSecretId,
	})

	if err != nil {
		c.logger.Error("error to get secret", zap.Error(err))
		return err
	}

	c.secrets = map[string]string{}

	err = json.Unmarshal(res.SecretBinary, &c.secrets)
	if err != nil {
		c.logger.Error("error get secret from aws", zap.Error(err))
		return err
	}

	return nil
}

// GetSecret retrieves a specific secret value by its key from the in-memory cache.
// This method is more efficient than fetching from AWS Secrets Manager directly
// for each request, as it uses the cached values loaded by LoadSecrets.
// Returns the secret value as a string if found, otherwise returns an error.
// The context parameter is not used in this implementation but is included to satisfy the interface.
func (c *awsSecretClient) GetSecret(_ context.Context, key string) (string, error) {
	value, ok := c.secrets[key]
	if !ok {
		return "", errors.New("secret was not found")
	}

	return value, nil
}
