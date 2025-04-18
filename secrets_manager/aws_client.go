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

type awsSecretClient struct {
	logger      logging.Logger
	client      *secretsmanager.Client
	appSecretId string
	secrets     map[string]string
}

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

func (c *awsSecretClient) GetSecret(_ context.Context, key string) (string, error) {
	value, ok := c.secrets[key]
	if !ok {
		return "", errors.New("secret was not found")
	}

	return value, nil
}
