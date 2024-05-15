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
	region string
}
