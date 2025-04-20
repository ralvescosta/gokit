// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

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
	PostgresSqlConnection struct {
		Err              error
		logger           logging.Logger
		connectionString string
		conn             *sql.DB
		cfg              *configs.Configs
	}
)

var sqlOpen = sql.Open

var otelOpen = otelsql.Open

const (
	FailureConnErrorMessage = "[PostgreSQL::Connect] failure to connect to the database"
)

func New(cfgs *configs.Configs) *PostgresSqlConnection {
	connString := pkgSql.GetConnectionString(cfgs.SQLConfigs)

	return &PostgresSqlConnection{
		logger:           cfgs.Logger,
		connectionString: connString,
		cfg:              cfgs,
	}
}

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
