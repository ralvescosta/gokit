package configs

type MQTTConfigs struct {
	Host       string
	Port       uint32
	User       string
	Password   string
	DeviceName string

	RootCaPath     string
	CertPath       string
	PrivateKeyPath string
}
