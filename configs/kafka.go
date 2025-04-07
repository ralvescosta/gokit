package configs

type KafkaConfigs struct {
	Host             string
	Port             int
	SecurityProtocol string
	SASLMechanisms   string
	User             string
	Password         string
}
