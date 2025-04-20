// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

type Auth0Configs struct {
	ClientID               string
	ClientSecret           string
	GrantType              string
	MillisecondsBetweenJWK int64
}
