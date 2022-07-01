package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logging"
)

type ISqlConnection interface {
	Connect() ISqlConnection
	ShotdownSignal() ISqlConnection
	Build() (*sql.DB, error)
}

func GetConnectionString(cfg *env.Configs) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.SQL_DB_HOST,
		cfg.SQL_DB_PORT,
		cfg.SQL_DB_USER,
		cfg.SQL_DB_PASSWORD,
		cfg.SQL_DB_NAME,
	)
}

func ShotdownSignal(timeToPing int, conn *sql.DB, log logging.ILogger, shotdown chan bool, connFailureLogMsg string) {
	for {
		time.Sleep(time.Duration(timeToPing) * time.Millisecond)
		err := conn.Ping()
		if err != nil {
			log.Error(connFailureLogMsg, logging.ErrorField(err))
			shotdown <- true
			break
		}
	}
}
