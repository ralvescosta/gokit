package pg

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logger"
	pkgSql "github.com/ralvescostati/pkgs/sql"

	_ "github.com/lib/pq"
)

var open = sql.Open

const (
	FailureConnErrorMessage = "[PostgreSQL::Connect] failure to connect to the database: %s"
)

type PostgresSqlConnection struct {
	Err              error
	logger           logger.ILogger
	connectionString string
	conn             *sql.DB
	cfg              *env.Configs
	shotdown         chan bool
}

func New(log logger.ILogger, cfg *env.Configs, shotdown chan bool) pkgSql.ISqlConnection {
	connString := pkgSql.GetConnectionString(cfg)

	return &PostgresSqlConnection{
		logger:           log,
		connectionString: connString,
		cfg:              cfg,
		shotdown:         shotdown,
	}
}

func (pg *PostgresSqlConnection) Connect() pkgSql.ISqlConnection {
	db, err := open("postgres", pg.connectionString)
	if err != nil {
		pg.logger.Error(fmt.Sprintf(FailureConnErrorMessage, err.Error()))
		pg.Err = errors.New(fmt.Sprintf(FailureConnErrorMessage, err.Error()))
		return pg
	}

	err = db.Ping()
	if err != nil {
		pg.logger.Error(fmt.Sprintf(FailureConnErrorMessage, err.Error()))
		pg.Err = errors.New(fmt.Sprintf(FailureConnErrorMessage, err.Error()))
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
