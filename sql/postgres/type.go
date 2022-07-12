package pg

import (
	"database/sql"

	"github.com/ralvescosta/toolkit/env"
	"github.com/ralvescosta/toolkit/logging"
)

type PostgresSqlConnection struct {
	Err              error
	logger           logging.ILogger
	connectionString string
	conn             *sql.DB
	cfg              *env.Configs
	shotdown         chan bool
}
