package sql

import (
	"context"
	"database/sql"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type SqlConnBuilder interface {
	WthShotdownSig() SqlConnBuilder
	Build() (*sql.DB, error)
}

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

func ShotdownSignal(timeToPing int, conn *sql.DB, log logging.ILogger) {
	for {
		time.Sleep(time.Duration(timeToPing) * time.Millisecond)
		err := conn.Ping()
		if err != nil {
			log.Error("[gokit::sql]", logging.ErrorField(err))
			signal.NotifyContext(context.Background(), syscall.SIGQUIT)
			break
		}
	}
}
