package pg

import (
	"database/sql"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type PostgresSqlConnection struct {
	Err              error
	logger           logging.ILogger
	connectionString string
	conn             *sql.DB
	cfg              *env.Configs
	withShotdownSig  bool
}
