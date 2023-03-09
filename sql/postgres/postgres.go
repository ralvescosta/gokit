package pg

import (
	"database/sql"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

type (
	PostgresSqlConnection struct {
		Err              error
		logger           logging.Logger
		connectionString string
		conn             *sql.DB
		cfg              *env.Configs
	}
)

var sqlOpen = sql.Open
var otelOpen = otelsql.Open

const (
	FailureConnErrorMessage = "[PostgreSQL::Connect] failure to connect to the database"
)
