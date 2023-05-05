package configs

type AWSConfigs struct {
	AccessKeyId     string
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
