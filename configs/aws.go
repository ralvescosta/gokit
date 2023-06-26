package configs

type AWSConfigs struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

type DynamoDBConfigs struct {
	Endpoint string
	region   string
	table    string
}

type AWSSecretManagerConfigs struct {
	region string
}
