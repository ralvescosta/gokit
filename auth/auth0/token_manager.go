package auth0

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/auth/errors"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	auth0nManager struct {
		logger          logging.Logger
		jwtConfigs      *configs.IdentityConfigs
		auth0Configs    *configs.Auth0Configs
		jwks            *jose.JSONWebKeySet
		lastJWKRetrieve int64
	}

	Auth0Claims struct {
		UserData *map[string]interface{} `json:"user_data,omitempty"`
		*auth.Claims
	}
)

func NewAuth0TokenManger(logger logging.Logger, cfg *configs.Configs) auth.IdentityManager {
	return &auth0nManager{
		logger:       logger,
		jwtConfigs:   cfg.IdentityConfigs,
		auth0Configs: cfg.Auth0Configs,
	}
}

func (m *auth0nManager) Validate(ctx context.Context, token string) (*auth.Claims, error) {
	if err := m.manageJWK(); err != nil {
		return nil, err
	}

	jsonToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	if _, ok := auth.AllowedSigningAlgorithms[auth.SignatureAlgorithm(jsonToken.Headers[0].Algorithm)]; !ok {
		m.logger.Error("[auth/auth0] - signature not allowed", zap.String("alg", jsonToken.Headers[0].Algorithm))
		return nil, errors.ErrSignatureNotAllowed
	}

	claims, err := m.deserializeClaims(ctx, jsonToken)
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

func (m *auth0nManager) manageJWK() error {
	if m.jwks == nil {
		return m.getJWK()
	}

	now := time.Now().UnixMilli()
	if now-m.lastJWKRetrieve >= m.auth0Configs.MillisecondsBetweenJWK {
		return m.getJWK()
	}

	return nil
}

func (m *auth0nManager) wellKnownURI() string {
	return fmt.Sprintf("https://%s/.well-known/jwks.json", m.jwtConfigs.Domain)
}

func (m *auth0nManager) getJWK() error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, m.wellKnownURI(), nil)
	if err != nil {
		m.logger.Error("[auth/auth0] - error creating jwk request", zap.Error(err))
		return errors.ErrJWKRetrieving
	}

	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		m.logger.Error("[auth/auth0] - error getting jwk", zap.Error(err))
		return errors.ErrJWKRetrieving
	}
	defer resp.Body.Close()

	keySet := &jose.JSONWebKeySet{}
	if err := json.NewDecoder(resp.Body).Decode(keySet); err != nil {
		m.logger.Error("[auth/auth0] - error unmarshaling jwk", zap.Error(err))
		return errors.ErrJWKRetrieving
	}

	m.jwks = keySet
	return nil
}

func (m *auth0nManager) deserializeClaims(ctx context.Context, token *jwt.JSONWebToken) (*Auth0Claims, error) {
	if m.jwks.Keys == nil || m.jwks.Keys[0].Public().Key == nil {
		return nil, errors.ErrClaimsRetrieving
	}

	claims := []interface{}{&Auth0Claims{}}
	if err := token.Claims(m.jwks.Keys[0].Public().Key, claims...); err != nil {
		m.logger.Error("[auth/auth0] - could not get token claims", zap.Error(err))
		return nil, errors.ErrClaimsRetrieving
	}

	registeredClaims := claims[0].(*Auth0Claims)

	return registeredClaims, nil
}

func (v *auth0nManager) validateClaimsWithLeeway(actualClaims *auth.Claims, expected jwt.Expected, leeway time.Duration) error {
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
