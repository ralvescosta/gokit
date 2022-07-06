package pg

import (
	"database/sql"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logging"
)

type PostgresSqlConnection struct {
	Err              error
	logger           logging.ILogger
	connectionString string
	conn             *sql.DB
	cfg              *env.Configs
	shotdown         chan bool
}
