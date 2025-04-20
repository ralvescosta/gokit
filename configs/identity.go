// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

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
