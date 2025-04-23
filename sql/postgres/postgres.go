// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package pg provides PostgreSQL database integration for the GoKit framework.
// It implements connection and interaction with PostgreSQL databases using both
// standard and OpenTelemetry-instrumented connections.
package pg

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.8.0"
	"go.uber.org/zap"

	pkgSql "github.com/ralvescosta/gokit/sql"
)

type (
	// PostgresSqlConnection handles connection to PostgreSQL databases.
	// It supports both standard and OpenTelemetry-instrumented connections.
	PostgresSqlConnection struct {
		// Err holds any error that occurred during connection operations
		Err error
		// logger handles structured logging
		logger logging.Logger
		// connectionString holds the formatted PostgreSQL connection string
		connectionString string
		// conn holds the active database connection
		conn *sql.DB
		// cfg holds the application configurations
		cfg *configs.Configs
	}
)

// Variables for dependency injection during testing
var sqlOpen = sql.Open
var otelOpen = otelsql.Open

const (
	// FailureConnErrorMessage is the standard error message used when connection fails
	FailureConnErrorMessage = "[PostgreSQL::Connect] failure to connect to the database"
)

// New creates a new PostgreSQL connection instance with the provided configurations.
// It prepares the connection string but does not establish the connection.
//
// Parameters:
//   - cfgs: Application configurations including SQL and tracing settings
//
// Returns:
//   - A new PostgresSqlConnection instance ready to connect
func New(cfgs *configs.Configs) *PostgresSqlConnection {
	connString := pkgSql.GetConnectionString(cfgs.SQLConfigs)

	return &PostgresSqlConnection{
		logger:           cfgs.Logger,
		connectionString: connString,
		cfg:              cfgs,
	}
}

// open establishes a database connection using either standard or
// OpenTelemetry-instrumented connection methods based on configuration.
//
// Returns:
//   - A database connection and any error that occurred
func (pg *PostgresSqlConnection) open() (*sql.DB, error) {
	if pg.cfg.TracingConfigs.Enabled {
		return otelOpen(
			"postgres",
			pg.connectionString,
			otelsql.WithAttributes(semconv.DBSystemSqlite),
			otelsql.WithDBName(pg.cfg.SQLConfigs.DbName),
		)
	}

	return sqlOpen("postgres", pg.connectionString)
}

// Connect establishes a connection to the PostgreSQL database and
// verifies connectivity with a ping.
//
// Returns:
//   - A connected database instance and any error that occurred
func (pg *PostgresSqlConnection) Connect() (*sql.DB, error) {
	db, err := pg.open()
	if err != nil {
		pg.logger.Error(FailureConnErrorMessage, zap.Error(err))
		return db, err
	}

	err = db.Ping()
	if err != nil {
		pg.logger.Error(FailureConnErrorMessage, zap.Error(err))
		return db, err
	}

	return db, nil
}
