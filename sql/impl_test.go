package sql

import (
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/ralvescostati/pkgs/env"
	loggerMock "github.com/ralvescostati/pkgs/logger/mock"

	// sqlMock "github.com/ralvescostati/pkgs/sql/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SqlTestSuite struct {
	suite.Suite

	connector  *MockConnector
	driverConn *MockPingDriverConn
	driver     *MockPingDriver
}

func TestSqlTestSuite(t *testing.T) {
	suite.Run(t, new(SqlTestSuite))
}

func (s *SqlTestSuite) SetupTest() {
	s.connector = &MockConnector{}
	s.driverConn = &MockPingDriverConn{}
	s.driver = &MockPingDriver{}
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
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	db := sql.OpenDB(s.connector)

	channel := make(chan bool)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go ShotdownSignal(1, db, &loggerMock.MockLogger{}, channel, "%s")
	time.Sleep(1 * time.Second)
	wg.Done()

	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *SqlTestSuite) TestShotdownSignalErr() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(errors.New("ping err"))
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	db := sql.OpenDB(s.connector)

	channel := make(chan bool)

	go ShotdownSignal(1, db, &loggerMock.MockLogger{}, channel, "%s")

	res := <-channel

	s.True(res)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}
