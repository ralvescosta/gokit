package sql

import (
	"fmt"

	"github.com/ralvescosta/gokit/env"
)

func GetConnectionString(cfg *env.SqlConfigs) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
	)
}
