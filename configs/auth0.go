package configs

type Auth0Configs struct {
	JWTConfigs
	ClientID               string
	ClientSecret           string
	GrantType              string
	MillisecondsBetweenJWK int64
}
