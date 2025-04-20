// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package errors

import "errors"

var (
	ErrSignatureNotAllowed = errors.New("signature not allowed")
	ErrJWKRetrieving       = errors.New("error to retrieve jwk from identity server")
	ErrKeysFunc            = errors.New("error getting the keys from the key func")
	ErrClaimsRetrieving    = errors.New("could not get token claims")
	ErrInvalidIssuer       = errors.New("validation failed, invalid issuer claim (iss)")
	ErrInvalidAudience     = errors.New("validation failed, invalid audience claim (aud)")
	ErrNotValidYet         = errors.New("validation failed, token not valid yet (nbf)")
	ErrExpired             = errors.New("validation failed, token is expired (exp)")
	ErrIssuedInTheFuture   = errors.New("validation field, token issued in the future (iat)")
)
