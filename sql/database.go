package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logger"

	_ "github.com/lib/pq"
)

var open = sql.Open

func Connect(log logger.ILogger, cfg *env.Configs, shotdown chan bool) (*sql.DB, error) {
	connString, err := getConnectionString(cfg)
	if err != nil {
		log.Error(fmt.Sprintf("[Sql::Connect] - wrong database credentials %s", err.Error()))
		return nil, err
	}

	db, err := open("postgres", connString)
	if err != nil {
		log.Error(fmt.Sprintf("[Sql::Connect] - error while connect to database: %s", err.Error()))
		return nil, errors.New(fmt.Sprintf("failure to connect to the database: %s", err.Error()))
	}

	err = db.Ping()
	if err != nil {
		log.Error(fmt.Sprintf("[Sql::Connect] - error while check database connection: %s", err.Error()))
		return nil, errors.New(fmt.Sprintf("failure to connect to the database: %s", err.Error()))
	}

	secondsToSleep, err := strconv.Atoi(os.Getenv("DB_SECONDS_TO_PING"))
	if err != nil {
		log.Error(fmt.Sprintf("[Sql::Connect] - DB_SECONDS_TO_PING is required: %s", err.Error()))
		return nil, err
	}

	go signalShotdown(db, log, secondsToSleep, shotdown)

	return db, nil
}

func getConnectionString(cfg *env.Configs) (string, error) {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.SQL_DB_HOST,
		cfg.SQL_DB_PORT,
		cfg.SQL_DB_USER,
		cfg.SQL_DB_PASSWORD,
		cfg.SQL_DB_NAME,
	), nil
}

func signalShotdown(db *sql.DB, log logger.ILogger, secondsToSleep int, shotdown chan bool) {
	time.Sleep(time.Duration(secondsToSleep) * time.Second)
	err := db.Ping()
	if err != nil {
		log.Error(fmt.Sprintf("[Database::Connection] - Connection failure : %s", err.Error()))
		shotdown <- true
	}
}
