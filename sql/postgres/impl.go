package pg

import (
	"database/sql"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	pkgSql "github.com/ralvescosta/gokit/sql"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.8.0"

	_ "github.com/lib/pq"
)

func New(logger logging.ILogger, cfg *env.Config) pkgSql.SqlConnBuilder {
	connString := pkgSql.GetConnectionString(cfg)

	return &PostgresSqlConnection{
		logger:           logger,
		connectionString: connString,
		cfg:              cfg,
	}
}

func (pg *PostgresSqlConnection) WthShotdownSig() pkgSql.SqlConnBuilder {
	pg.withShotdownSig = true

	return pg
}

func (pg *PostgresSqlConnection) open() (*sql.DB, error) {
	var db *sql.DB
	var err error

	if pg.cfg.IS_TRACING_ENABLED {
		db, err = otelOpen(
			"postgres",
			pg.connectionString,
			otelsql.WithAttributes(semconv.DBSystemSqlite),
			otelsql.WithDBName(pg.cfg.SQL_DB_NAME),
		)

		return db, err
	}

	db, err = sqlOpen("postgres", pg.connectionString)
	return db, err
}

func (pg *PostgresSqlConnection) connect() (*sql.DB, error) {
	db, err := pg.open()
	if err != nil {
		pg.logger.Error(FailureConnErrorMessage, logging.ErrorField(err))
		return db, err
	}

	err = db.Ping()
	if err != nil {
		pg.logger.Error(FailureConnErrorMessage, logging.ErrorField(err))
		return db, err
	}

	return db, nil
}

func (pg *PostgresSqlConnection) Build() (*sql.DB, error) {
	conn, err := pg.connect()
	if err != nil {
		return nil, err
	}

	if pg.withShotdownSig {
		go pkgSql.ShotdownSignal(pg.cfg.SQL_DB_SECONDS_TO_PING, pg.conn, pg.logger)
	}

	return conn, nil
}
