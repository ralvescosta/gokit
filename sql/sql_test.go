package sql

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/ralvescostati/pkgs/env"
	loggerMock "github.com/ralvescostati/pkgs/logger/mock"
	sqlMock "github.com/ralvescostati/pkgs/sql/mock"
	"github.com/stretchr/testify/suite"

	"github.com/DATA-DOG/go-sqlmock"
)

type SqlTestSuite struct {
	suite.Suite
}

func TestSqlTestSuite(t *testing.T) {
	suite.Run(t, new(SqlTestSuite))
}

func (s *SqlTestSuite) TestGetConnection() {
	cfg := &env.Configs{
		SQL_DB_HOST:     "host",
		SQL_DB_PORT:     "port",
		SQL_DB_USER:     "user",
		SQL_DB_PASSWORD: "password",
		SQL_DB_NAME:     "name",
	}

	connStr := GetConnectionString(cfg)

	s.Equal("host=host port=port user=user password=password dbname=name sslmode=disable", connStr)
}

func (s *SqlTestSuite) TestShotdownSignal() {
	db, _, _ := sqlmock.New()
	var channel chan bool

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		ShotdownSignal(0, db, loggerMock.NewMockLogger(), channel, "%s")
	}()
	wg.Done()
}

func (s *SqlTestSuite) TestShotdownSignalErr() {
	driver := &sqlMock.PingDriver{}
	sql.Register("ping", driver)

	db, _ := sql.Open("ping", "ignored")
	driver.Fails = true
	channel := make(chan bool)

	go ShotdownSignal(0, db, loggerMock.NewMockLogger(), channel, "%s")

	res := <-channel

	s.True(res)
}
