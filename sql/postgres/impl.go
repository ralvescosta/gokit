package pg

import (
	"database/sql"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	pkgSql "github.com/ralvescosta/gokit/sql"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.8.0"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

func New(logger logging.Logger, cfg *env.Configs) *PostgresSqlConnection {
	connString := pkgSql.GetConnectionString(cfg.SqlConfigs)

	return &PostgresSqlConnection{
		logger:           logger,
		connectionString: connString,
		cfg:              cfg,
	}
}

func (pg *PostgresSqlConnection) open() (*sql.DB, error) {
	if pg.cfg.OtelConfigs.TracingEnabled {
		return otelOpen(
			"postgres",
			pg.connectionString,
			otelsql.WithAttributes(semconv.DBSystemSqlite),
			otelsql.WithDBName(pg.cfg.SqlConfigs.DbName),
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
