package auth0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"go.uber.org/zap"
)

type (
	auth0nManager struct {
		logger          logging.Logger
		cfg             *configs.Auth0Configs
		jwks            *jose.JSONWebKeySet
		lastJWKRetrieve int64
	}
)

func NewAuth0TokenManger(logger logging.Logger, cfg *configs.Auth0Configs) auth.IdentityManager {
	return &auth0nManager{
		logger: logger,
		cfg:    cfg,
	}
}

func (m *auth0nManager) Validate(ctx context.Context, token string) (*auth.Session, error) {
	if err := m.manageJWK(); err != nil {
		return nil, err
	}

	jsonToken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	if _, ok := auth.AllowedSigningAlgorithms[auth.SignatureAlgorithm(jsonToken.Headers[0].Algorithm)]; !ok {
		return nil, errors.New("")
	}

	claims, err := m.deserializeClaims(ctx, jsonToken)
	if err != nil {
		return nil, err
	}

	expectedClaims := jwt.Expected{
		Issuer:   m.cfg.Issuer,
		Audience: []string{m.cfg.Audience},
	}

	if err := m.validateClaimsWithLeeway(claims, expectedClaims, time.Hour); err != nil {
		return nil, err
	}

	return &auth.Session{
		Issuer:   claims.Issuer,
		Subject:  claims.Subject,
		Audience: claims.Audience,
		IssuedAt: m.numericDateToUnixTime(claims.IssuedAt),
		Expiry:   m.numericDateToUnixTime(claims.Expiry),
	}, nil
}

func (m *auth0nManager) manageJWK() error {
	if m.jwks == nil {
		return m.getJWK()
	}

	now := time.Now().UnixMilli()
	if now-m.lastJWKRetrieve >= m.cfg.MillisecondsBetweenJWK {
		return m.getJWK()
	}

	return nil
}

func (m *auth0nManager) wellKnownURI() string {
	return fmt.Sprintf("https://%s/.well-known/jwks.json", m.cfg.Domain)
}

func (m *auth0nManager) getJWK() error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, m.wellKnownURI(), nil)
	if err != nil {
		m.logger.Error("error creating jwk request", zap.Error(err))
		return err
	}

	req.Header.Set("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		m.logger.Error("error getting jwk", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	keySet := &jose.JSONWebKeySet{}
	if err := json.NewDecoder(resp.Body).Decode(keySet); err != nil {
		m.logger.Error("error unmarshaling jwk", zap.Error(err))
		return err
	}

	m.jwks = keySet
	return nil
}

func (m *auth0nManager) deserializeClaims(ctx context.Context, token *jwt.JSONWebToken) (*jwt.Claims, error) {
	key, err := func(ctx context.Context) (interface{}, error) {
		return m.jwks.Keys[0].Public().Key, nil
	}(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting the keys from the key func: %w", err)
	}

	claims := []interface{}{&jwt.Claims{}}
	if err := token.Claims(key, claims); err != nil {
		return nil, fmt.Errorf("could not get token claims: %w", err)
	}

	registeredClaims := *claims[0].(*jwt.Claims)

	// var customClaims CustomClaims
	// if len(claims) > 1 {
	// 	customClaims = claims[1].(CustomClaims)
	// }

	return &registeredClaims, nil
}

func (v *auth0nManager) validateClaimsWithLeeway(actualClaims *jwt.Claims, expected jwt.Expected, leeway time.Duration) error {
	expectedClaims := expected
	expectedClaims.Time = time.Now()

	if actualClaims.Issuer != expectedClaims.Issuer {
		return jwt.ErrInvalidIssuer
	}

	foundAudience := false
	for _, value := range expectedClaims.Audience {
		if actualClaims.Audience.Contains(value) {
			foundAudience = true
			break
		}
	}
	if !foundAudience {
		return jwt.ErrInvalidAudience
	}

	if actualClaims.NotBefore != nil && expectedClaims.Time.Add(leeway).Before(actualClaims.NotBefore.Time()) {
		return jwt.ErrNotValidYet
	}

	if actualClaims.Expiry != nil && expectedClaims.Time.Add(-leeway).After(actualClaims.Expiry.Time()) {
		return jwt.ErrExpired
	}

	if actualClaims.IssuedAt != nil && expectedClaims.Time.Add(leeway).Before(actualClaims.IssuedAt.Time()) {
		return jwt.ErrIssuedInTheFuture
	}

	return nil
}

func (v *auth0nManager) numericDateToUnixTime(date *jwt.NumericDate) int64 {
	if date != nil {
		return date.Time().Unix()
	}
	return 0
}
