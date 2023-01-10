package pg

import (
	"database/sql"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

var sqlOpen = sql.Open
var otelOpen = otelsql.Open

const (
	FailureConnErrorMessage = "[PostgreSQL::Connect] failure to connect to the database"
)
