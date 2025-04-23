// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

/**
 * Package auth0 provides an implementation of the IdentityManager interface for Auth0.
 * It includes functionality for managing JSON Web Keys (JWKs) and validating JWT tokens.
 */
package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"

	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/auth/errors"
)

/**
 * auth0nManager is the implementation of the IdentityManager interface for Auth0.
 * It manages JSON Web Keys (JWKs) and validates JWT tokens.
 *
 * Fields:
 * - logger: Logger instance for logging errors and information.
 * - jwtConfigs: Configuration for JWT validation, including issuer and audience.
 * - auth0Configs: Configuration for Auth0, including JWK retrieval intervals.
 * - jwks: JSON Web Key Set used for token validation.
 * - lastJWKRetrieve: Timestamp of the last JWK retrieval.
 */
type auth0nManager struct {
	logger          logging.Logger
	jwtConfigs      *configs.IdentityConfigs
	auth0Configs    *configs.Auth0Configs
	jwks            *jose.JSONWebKeySet
	lastJWKRetrieve int64
}

/**
 * Claims represents the JWT claims structure.
 *
 * Fields:
 * - UserData: Optional user data included in the token.
 * - Claims: Embedded standard claims from the auth package.
 */
type Claims struct {
	UserData *map[string]interface{} `json:"user_data,omitempty"`
	*auth.Claims
}

/**
 * NewAuth0TokenManager creates a new instance of auth0nManager.
 *
 * Parameters:
 * - cfg: Configuration object containing logger, identity, and Auth0 configurations.
 *
 * Returns:
 * - An instance of auth.IdentityManager.
 */
func NewAuth0TokenManager(cfg *configs.Configs) auth.IdentityManager {
	return &auth0nManager{
		logger:       cfg.Logger,
		jwtConfigs:   cfg.IdentityConfigs,
		auth0Configs: cfg.Auth0Configs,
	}
}

/**
 * Validate validates a JWT token and returns the claims if valid.
 *
 * Parameters:
 * - ctx: Context for managing request lifecycle.
 * - token: JWT token string to validate.
 *
 * Returns:
 * - Pointer to auth.Claims if the token is valid.
 * - Error if the token is invalid or validation fails.
 */
func (m *auth0nManager) Validate(ctx context.Context, token string) (*auth.Claims, error) {
	if err := m.manageJWK(ctx); err != nil {
		return nil, err
	}

	jsonToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	if _, ok := auth.AllowedSigningAlgorithms[auth.SignatureAlgorithm(jsonToken.Headers[0].Algorithm)]; !ok {
		m.logger.Error("[gokit:auth/auth0] - signature not allowed", zap.String("alg", jsonToken.Headers[0].Algorithm))
		return nil, errors.ErrSignatureNotAllowed
	}

	claims, err := m.deserializeClaims(jsonToken)
	if err != nil {
		return nil, err
	}

	expectedClaims := jwt.Expected{
		Issuer:   m.jwtConfigs.Issuer,
		Audience: []string{m.jwtConfigs.Audience},
	}

	if err := m.validateClaimsWithLeeway(claims.Claims, expectedClaims, time.Hour); err != nil {
		return nil, err
	}

	return claims.Claims, nil
}

/**
 * manageJWK ensures the JWK set is up-to-date.
 *
 * Parameters:
 * - ctx: Context for managing request lifecycle.
 *
 * Returns:
 * - Error if JWK retrieval fails.
 */
func (m *auth0nManager) manageJWK(ctx context.Context) error {
	if m.jwks == nil {
		return m.getJWK(ctx)
	}

	now := time.Now().UnixMilli()
	if now-m.lastJWKRetrieve >= m.auth0Configs.MillisecondsBetweenJWK {
		return m.getJWK(ctx)
	}

	return nil
}

/**
 * wellKnownURI constructs the URI for retrieving the JWK set.
 *
 * Returns:
 * - A string representing the well-known JWK URI.
 */
func (m *auth0nManager) wellKnownURI() string {
	return fmt.Sprintf("https://%s/.well-known/jwks.json", m.jwtConfigs.Domain)
}

/**
 * getJWK retrieves the JWK set from the Auth0 well-known URI.
 *
 * Parameters:
 * - ctx: Context for managing request lifecycle.
 *
 * Returns:
 * - Error if JWK retrieval or unmarshaling fails.
 */
func (m *auth0nManager) getJWK(ctx context.Context) error {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.wellKnownURI(), nil)
	if err != nil {
		m.logger.Error("[gokit:auth/auth0] - error creating jwk request", zap.Error(err))
		return errors.ErrJWKRetrieving
	}

	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		m.logger.Error("[gokit:auth/auth0] - error getting jwk", zap.Error(err))
		return errors.ErrJWKRetrieving
	}
	defer resp.Body.Close()

	keySet := &jose.JSONWebKeySet{}
	if err := json.NewDecoder(resp.Body).Decode(keySet); err != nil {
		m.logger.Error("[gokit:auth/auth0] - error unmarshaling jwk", zap.Error(err))
		return errors.ErrJWKRetrieving
	}

	m.jwks = keySet
	return nil
}

/**
 * deserializeClaims extracts claims from a JWT token.
 *
 * Parameters:
 * - token: Parsed JWT token.
 *
 * Returns:
 * - Pointer to Claims if deserialization is successful.
 * - Error if deserialization fails.
 */
func (m *auth0nManager) deserializeClaims(token *jwt.JSONWebToken) (*Claims, error) {
	if m.jwks.Keys == nil || m.jwks.Keys[0].Public().Key == nil {
		return nil, errors.ErrClaimsRetrieving
	}

	claims := []interface{}{&Claims{}}
	if err := token.Claims(m.jwks.Keys[0].Public().Key, claims...); err != nil {
		m.logger.Error("[gokit:auth/auth0] - could not get token claims", zap.Error(err))
		return nil, errors.ErrClaimsRetrieving
	}

	registeredClaims := claims[0].(*Claims)

	return registeredClaims, nil
}

/**
 * validateClaimsWithLeeway validates the claims of a JWT token with a leeway.
 *
 * Parameters:
 * - actualClaims: Claims extracted from the token.
 * - expected: Expected claims for validation.
 * - leeway: Allowed time difference for validation.
 *
 * Returns:
 * - Error if validation fails.
 */
func (m *auth0nManager) validateClaimsWithLeeway(actualClaims *auth.Claims, expected jwt.Expected, leeway time.Duration) error {
	expectedClaims := expected
	expectedClaims.Time = time.Now()

	if actualClaims.Issuer != expectedClaims.Issuer {
		return errors.ErrInvalidIssuer
	}

	foundAudience := false
	for _, value := range actualClaims.Audience {
		if expectedClaims.Audience.Contains(value) {
			foundAudience = true
			break
		}
	}
	if !foundAudience {
		return errors.ErrInvalidAudience
	}

	if actualClaims.NotBefore != nil && expectedClaims.Time.Add(leeway).Before(*actualClaims.NotBeforeTime()) {
		return errors.ErrNotValidYet
	}

	if actualClaims.Expiry != nil && expectedClaims.Time.Add(-leeway).After(*actualClaims.ExpiryTime()) {
		return errors.ErrExpired
	}

	if actualClaims.IssuedAt != nil && expectedClaims.Time.Add(leeway).Before(*actualClaims.IssuedAtTime()) {
		return errors.ErrIssuedInTheFuture
	}

	return nil
}
