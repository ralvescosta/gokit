package sql

import (
	"fmt"

	"github.com/ralvescosta/gokit/env"
)

func GetConnectionString(cfg *env.Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.SQL_DB_HOST,
		cfg.SQL_DB_PORT,
		cfg.SQL_DB_USER,
		cfg.SQL_DB_PASSWORD,
		cfg.SQL_DB_NAME,
	)
}
