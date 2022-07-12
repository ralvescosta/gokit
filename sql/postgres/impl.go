package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ralvescosta/toolkit/env"
	"github.com/ralvescosta/toolkit/logging"
	pkgSql "github.com/ralvescosta/toolkit/sql"

	_ "github.com/lib/pq"
)

func New(logger logging.ILogger, cfg *env.Configs, shotdown chan bool) pkgSql.ISqlConnection {
	connString := pkgSql.GetConnectionString(cfg)

	return &PostgresSqlConnection{
		logger:           logger,
		connectionString: connString,
		cfg:              cfg,
		shotdown:         shotdown,
	}
}

func (pg *PostgresSqlConnection) Connect() pkgSql.ISqlConnection {
	db, err := open("postgres", pg.connectionString)
	if err != nil {
		pg.logger.Error(FailureConnErrorMessage, logging.ErrorField(err))
		pg.Err = fmt.Errorf(FailureConnErrorMessage, err.Error())
		return pg
	}

	err = db.Ping()
	if err != nil {
		pg.logger.Error(FailureConnErrorMessage, logging.ErrorField(err))
		pg.Err = fmt.Errorf(FailureConnErrorMessage, err.Error())
		return pg
	}

	pg.conn = db

	return pg
}

func (pg *PostgresSqlConnection) ShotdownSignal() pkgSql.ISqlConnection {
	if pg.Err != nil {
		return pg
	}

	if pg.shotdown == nil || pg.cfg.SQL_DB_SECONDS_TO_PING == 0 {
		pg.Err = errors.New("[PostgreSQL::Connect] shotdown channel and SQL_DB_SECONDS_TO_PING is required")
		return pg
	}

	go pkgSql.ShotdownSignal(pg.cfg.SQL_DB_SECONDS_TO_PING, pg.conn, pg.logger, pg.shotdown, "[PostgreSQL::Connect] - connection failure : %s")

	return pg
}

func (pg *PostgresSqlConnection) Build() (*sql.DB, error) {
	if pg.Err != nil {
		return nil, pg.Err
	}

	return pg.conn, nil
}
