package sql

import (
	"fmt"

	"github.com/ralvescosta/gokit/configs"
)

func GetConnectionString(cfg *configs.SQLConfigs) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DbName,
	)
}
