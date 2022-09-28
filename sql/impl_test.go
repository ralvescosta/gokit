package sql

import (
	"testing"

	"github.com/ralvescosta/gokit/env"

	// sqlMock "github.com/ralvescosta/gokit/sql/mock"

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
	cfg := &env.Config{
		SQL_DB_HOST:     "host",
		SQL_DB_PORT:     "port",
		SQL_DB_USER:     "user",
		SQL_DB_PASSWORD: "password",
		SQL_DB_NAME:     "name",
	}

	connStr := GetConnectionString(cfg)

	s.Equal("host=host port=port user=user password=password dbname=name sslmode=disable", connStr)
}
