package configs

type IdentityConfigs struct {
	ClientID               string
	ClientSecret           string
	GrantType              string
	MillisecondsBetweenJWK int64
	Domain                 string
	Audience               string
	Issuer                 string
	Signature              string
}
