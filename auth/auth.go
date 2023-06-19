package auth

type (
	Session struct {
		Iss     string
		Sub     string
		Aud     []string
		Iat     uint32
		Exp     uint32
		Scope   string
		Session map[string]interface{}
	}

	TokenManager interface {
		Validate(token string) (*Session, error)
	}
)
