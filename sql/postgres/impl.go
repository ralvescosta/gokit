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

func New(logger logging.Logger, cfg *env.Config) *PostgresSqlConnection {
	connString := pkgSql.GetConnectionString(cfg)

	return &PostgresSqlConnection{
		logger:           logger,
		connectionString: connString,
		cfg:              cfg,
	}
}

func (pg *PostgresSqlConnection) open() (*sql.DB, error) {
	if pg.cfg.TRACING_ENABLED {
		return otelOpen(
			"postgres",
			pg.connectionString,
			otelsql.WithAttributes(semconv.DBSystemSqlite),
			otelsql.WithDBName(pg.cfg.SQL_DB_NAME),
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
