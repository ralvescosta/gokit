// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package auth

import (
	"context"
	"time"
)

type (
	Claims struct {
		Issuer      string    `json:"iss,omitempty"`
		Subject     string    `json:"sub,omitempty"`
		Audience    []string  `json:"aud,omitempty"`
		Expiry      *int64    `json:"exp,omitempty"`
		NotBefore   *int64    `json:"nbf,omitempty"`
		IssuedAt    *int64    `json:"iat,omitempty"`
		ID          string    `json:"jti,omitempty"`
		Scope       *string   `json:"scope,omitempty"`
		Permissions []*string `json:"permissions,omitempty"`
	}

	IdentityManager interface {
		Validate(ctx context.Context, token string) (*Claims, error)
	}

	SignatureAlgorithm string
)

var (
	EdDSA = SignatureAlgorithm("EdDSA")
	HS256 = SignatureAlgorithm("HS256") // HMAC using SHA-256
	HS384 = SignatureAlgorithm("HS384") // HMAC using SHA-384
	HS512 = SignatureAlgorithm("HS512") // HMAC using SHA-512
	RS256 = SignatureAlgorithm("RS256") // RSASSA-PKCS-v1.5 using SHA-256
	RS384 = SignatureAlgorithm("RS384") // RSASSA-PKCS-v1.5 using SHA-384
	RS512 = SignatureAlgorithm("RS512") // RSASSA-PKCS-v1.5 using SHA-512
	ES256 = SignatureAlgorithm("ES256") // ECDSA using P-256 and SHA-256
	ES384 = SignatureAlgorithm("ES384") // ECDSA using P-384 and SHA-384
	ES512 = SignatureAlgorithm("ES512") // ECDSA using P-521 and SHA-512
	PS256 = SignatureAlgorithm("PS256") // RSASSA-PSS using SHA256 and MGF1-SHA256
	PS384 = SignatureAlgorithm("PS384") // RSASSA-PSS using SHA384 and MGF1-SHA384
	PS512 = SignatureAlgorithm("PS512") // RSASSA-PSS using SHA512 and MGF1-SHA512

	AllowedSigningAlgorithms = map[SignatureAlgorithm]uint8{
		EdDSA: 1,
		HS256: 1,
		HS384: 1,
		HS512: 1,
		RS256: 1,
		RS384: 1,
		RS512: 1,
		ES256: 1,
		ES384: 1,
		ES512: 1,
		PS256: 1,
		PS384: 1,
		PS512: 1,
	}
)

func (c *Claims) ExpiryTime() *time.Time {
	if c.Expiry == nil {
		return nil
	}

	v := time.Unix(*c.Expiry, 0)

	return &v
}

func (c *Claims) NotBeforeTime() *time.Time {
	if c.NotBefore == nil {
		return nil
	}

	v := time.Unix(*c.NotBefore, 0)

	return &v
}

func (c *Claims) IssuedAtTime() *time.Time {
	if c.IssuedAt == nil {
		return nil
	}

	v := time.Unix(*c.IssuedAt, 0)

	return &v
}
