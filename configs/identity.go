// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

// IdentityConfigs contains configuration settings for identity and authentication services.
// It provides the necessary parameters for OAuth2/OIDC flows and JWT token verification.
type IdentityConfigs struct {
	// ClientID represents the application's client identifier for authentication with identity providers
	ClientID string
	// ClientSecret contains the secret key used to authenticate the application with the identity provider
	ClientSecret string
	// GrantType specifies the OAuth2 grant type to be used (e.g., "client_credentials", "authorization_code")
	GrantType string
	// MillisecondsBetweenJWK defines the interval between JWK (JSON Web Key) refresh operations
	MillisecondsBetweenJWK int64
	// Domain specifies the identity provider's domain (often used with Auth0 or similar services)
	Domain string
	// Audience identifies the intended recipient of the authentication tokens
	Audience string
	// Issuer identifies the entity that issued the JWT token
	Issuer string
	// Signature contains cryptographic material used to verify token signatures
	Signature string
}
