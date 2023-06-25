package errors

import "errors"

var (
	ErrSignatureNotAllowed = errors.New("signature not allowed")
	ErrJWKRetrieving       = errors.New("error to retrieve jwk from identity server")
	ErrKeysFunc            = errors.New("error getting the keys from the key func")
	ErrClaimsRetrieving    = errors.New("could not get token claims")
)
