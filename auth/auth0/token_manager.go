package auth0

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/configs"
)

type (
	auth0nManager struct {
		cfg *configs.Auth0Configs
	}
)

func NewAuth0TokenManger(cfg *configs.Auth0Configs) auth.IdentityManager {

	return &auth0nManager{
		cfg: cfg,
	}
}

func (m *auth0nManager) Validate(ctx context.Context, token string) (*auth.Session, error) {
	sig, err := jose.ParseSigned(token)
	if err != nil {
		return nil, err
	}

	headers := make([]jose.Header, len(sig.Signatures))
	for i, signature := range sig.Signatures {
		headers[i] = signature.Header
	}

	// validate algorithm signature
	if headers[0].Algorithm != string(auth.HS256) {
		return nil, errors.New("")
	}

	claims, err := m.deserializeClaims(ctx, nil)
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

func (m *auth0nManager) deserializeClaims(ctx context.Context, token *jwt.JSONWebToken) (jwt.Claims, error) {
	// key, err := v.keyFunc(ctx)
	// if err != nil {
	// 	return jwt.Claims{}, nil, fmt.Errorf("error getting the keys from the key func: %w", err)
	// }

	claims := []interface{}{&jwt.Claims{}}
	// if v.customClaimsExist() {
	// 	claims = append(claims, v.customClaims())
	// }

	if err := token.Claims("", claims...); err != nil {
		return jwt.Claims{}, fmt.Errorf("could not get token claims: %w", err)
	}

	registeredClaims := *claims[0].(*jwt.Claims)

	// var customClaims CustomClaims
	// if len(claims) > 1 {
	// 	customClaims = claims[1].(CustomClaims)
	// }

	return registeredClaims, nil
}

func (v *auth0nManager) validateClaimsWithLeeway(actualClaims jwt.Claims, expected jwt.Expected, leeway time.Duration) error {
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
