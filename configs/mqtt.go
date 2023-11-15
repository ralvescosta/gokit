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
