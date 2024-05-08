package sql

import (
	"testing"

	// sqlMock "github.com/ralvescosta/gokit/sql/mock"

	"github.com/ralvescosta/gokit/configs"
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
	cfg := &configs.SQLConfigs{
		Host:     "host",
		Port:     "port",
		User:     "user",
		Password: "password",
		DbName:   "name",
	}

	connStr := GetConnectionString(cfg)

	s.Equal("host=host port=port user=user password=password dbname=name sslmode=disable", connStr)
}
