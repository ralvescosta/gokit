// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package sql provides utilities and interfaces for working with SQL databases.
// It includes functionality for connection string generation and supports various
// SQL database drivers through the standard database/sql package.
package sql

import (
	"fmt"

	"github.com/ralvescosta/gokit/configs"
)

// GetConnectionString creates a formatted connection string using the provided SQL configurations.
// This function returns a connection string compatible with PostgreSQL drivers.
//
// Parameters:
//   - cfg: SQL configuration containing host, port, user, password, and database name.
//
// Returns:
//   - A formatted connection string ready to be used with SQL drivers.
func GetConnectionString(cfg *configs.SQLConfigs) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
	)
}
