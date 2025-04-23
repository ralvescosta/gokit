# Auth Package

The `auth` package provides tools and abstractions for managing authentication and authorization in Go applications. It includes support for token validation, claims management, and integration with identity providers like Auth0.

## Features

- **Claims Management**:
  - Define and manage JWT claims, including standard fields like `Issuer`, `Subject`, `Audience`, and custom fields like `Scope` and `Permissions`.

- **Token Validation**:
  - Validate JWT tokens using various signature algorithms.
  - Support for integration h identity providers like Auth0.

- **Error Handling**:
  - Predefined error types for common authentication and token validation issues.

## Package Structure

### `auth.go`
Defines the core types and interfaces for the `auth` package:
- **`Claims`**: Represents JWT claims, including standard and custom fields.
- **`IdentityManager`**: Interface for validating tokens and retrieving claims.
- **`SignatureAlgorithm`**: Enum-like type for supported JWT signature algorithms.

### `auth0/token_manager.go`
Provides an implementation of the `IdentityManager` interface for Auth0:
- **`auth0nManager`**: Manages token validation and integration with Auth0's JSON Web Key Set (JWKS).
- **`NewAuth0TokenManager`**: Factory function to create a new Auth0 token manager.

### `errors/errors.go`
Defines common errors used in the `auth` package:
- **`ErrSignatureNotAllowed`**: Indicates an unsupported signature algorithm.
- **`ErrJWKRetrieving`**: Error retrieving JWKS from the identity server.
- **`ErrInvalidIssuer`**: Validation failed due to an invalid `iss` claim.
- **`ErrExpired`**: Token is expired, among others.

## Usage

### Validating a Token
To validate a JWT token and retrieve its claims:
```go
package main

import (
	"context"
	"fmt"
	"github.com/ralvescosta/gokit/auth"
	"github.com/ralvescosta/gokit/auth/auth0"
	configsBuilder "github.com/ralvescosta/gokit/configs_builder"
)

func main() {
	// Load configurations
	cfgs, err := configsBuilder.
		NewConfigsBuilder().
		Identity().
		Build()

  if err != nil {
		panic(err)
	}

	// Create an Auth0 token manager
	manager := auth0.NewAuth0TokenManager(cfg)

	// Validate a token
	token := "your-jwt-token"
	claims, err := manager.Validate(context.Background(), token)
	if err != nil {
		fmt.Println("Token validation failed:", err)
		return
	}

	fmt.Println("Token is valid. Claims:", claims)
}
```

### Handling Errors
The `auth/errors` package provides predefined errors for common scenarios:
```go
if err == autherrors.ErrInvalidIssuer {
	fmt.Println("Invalid issuer in token")
}
```

## Supported Signature Algorithms
The `auth` package supports the following JWT signature algorithms:
- `HS256`, `HS384`, `HS512` (HMAC)
- `RS256`, `RS384`, `RS512` (RSA)
- `ES256`, `ES384`, `ES512` (ECDSA)
- `PS256`, `PS384`, `PS512` (RSA-PSS)
- `EdDSA`

## License
This package is licensed under the MIT License. See the LICENSE file for details.
