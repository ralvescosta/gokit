package pg

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	pkgSql "github.com/ralvescosta/gokit/sql"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.8.0"
	"go.uber.org/zap"
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

func New(logger logging.Logger, cfg *configs.Configs) *PostgresSqlConnection {
	connString := pkgSql.GetConnectionString(cfg.SQLConfigs)

	return &PostgresSqlConnection{
		logger:           logger,
		connectionString: connString,
		cfg:              cfg,
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
