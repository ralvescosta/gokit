// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// Auth0Configs provides Auth0-specific configuration settings.
// It contains the necessary parameters to authenticate and interact with Auth0 services.
type Auth0Configs struct {
	// ClientID represents the Auth0 application's client identifier
	ClientID string
	// ClientSecret contains the secret key used to authenticate the application with Auth0
	ClientSecret string
	// GrantType specifies the OAuth2 grant type to be used (e.g., "client_credentials")
	GrantType string
	// MillisecondsBetweenJWK defines the interval between JWK (JSON Web Key) refresh operations
	MillisecondsBetweenJWK int64
}
