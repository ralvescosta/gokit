package auth

import "context"

type (
	Session struct {
		Iss     string
		Sub     string
		Aud     []string
		Iat     int64
		Exp     int64
		Scope   string
		Session map[string]interface{}
	}

	IdentityManager interface {
		Validate(ctx context.Context, token string) (*Session, error)
	}
)
