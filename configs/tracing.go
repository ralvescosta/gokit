package configs

type TracingConfigs struct {
	Enabled bool

	OtlpEndpoint string
	OtlpApiKey   string
}
