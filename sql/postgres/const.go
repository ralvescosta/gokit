package pg

import "database/sql"

var open = sql.Open

const (
	FailureConnErrorMessage = "[PostgreSQL::Connect] failure to connect to the database: %s"
)
