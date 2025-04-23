// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package errors defines custom error types for the auth package.
//
// These errors are used to handle specific failure cases during token validation and JWK retrieval.
//
// Error Variables:
//   - ErrSignatureNotAllowed: Returned when a JWT uses a disallowed signature algorithm.
//   - ErrJWKRetrieving: Returned when there is a failure retrieving JWKs from the identity server.
//   - ErrKeysFunc: Returned when there is an error obtaining keys from the key function.
//   - ErrClaimsRetrieving: Returned when token claims cannot be extracted or parsed.
//   - ErrInvalidIssuer: Returned when the issuer (iss) claim is invalid during validation.
//   - ErrInvalidAudience: Returned when the audience (aud) claim is invalid during validation.
//   - ErrNotValidYet: Returned when the token is not valid yet (nbf claim).
//   - ErrExpired: Returned when the token is expired (exp claim).
//   - ErrIssuedInTheFuture: Returned when the token's issued-at (iat) claim is in the future.
//
// These errors are intended for use by the auth package and its subpackages.
//
// Example usage:
//
//	if err == errors.ErrInvalidIssuer {
//	    // handle invalid issuer error
//	}
//
// Author: The GoKit Authors
// License: MIT
package errors

import "errors"

// ErrSignatureNotAllowed is returned when a JWT uses a disallowed signature algorithm.
var ErrSignatureNotAllowed = errors.New("signature not allowed")

// ErrJWKRetrieving is returned when there is a failure retrieving JWKs from the identity server.
var ErrJWKRetrieving = errors.New("error to retrieve jwk from identity server")

// ErrKeysFunc is returned when there is an error obtaining keys from the key function.
var ErrKeysFunc = errors.New("error getting the keys from the key func")

// ErrClaimsRetrieving is returned when token claims cannot be extracted or parsed.
var ErrClaimsRetrieving = errors.New("could not get token claims")

// ErrInvalidIssuer is returned when the issuer (iss) claim is invalid during validation.
var ErrInvalidIssuer = errors.New("validation failed, invalid issuer claim (iss)")

// ErrInvalidAudience is returned when the audience (aud) claim is invalid during validation.
var ErrInvalidAudience = errors.New("validation failed, invalid audience claim (aud)")

// ErrNotValidYet is returned when the token is not valid yet (nbf claim).
var ErrNotValidYet = errors.New("validation failed, token not valid yet (nbf)")

// ErrExpired is returned when the token is expired (exp claim).
var ErrExpired = errors.New("validation failed, token is expired (exp)")

// ErrIssuedInTheFuture is returned when the token's issued-at (iat) claim is in the future.
var ErrIssuedInTheFuture = errors.New("validation field, token issued in the future (iat)")
