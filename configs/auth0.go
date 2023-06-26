package configs

type Auth0Configs struct {
	ClientID               string
	ClientSecret           string
	GrantType              string
	MillisecondsBetweenJWK int64
}
