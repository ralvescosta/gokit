// Package sql contains tests for the SQL utilities in the GoKit framework.
// These tests verify the proper functionality of SQL connection string formatting.
package sql

import (
	"testing"

	// sqlMock "github.com/ralvescosta/gokit/sql/mock"

	"github.com/ralvescosta/gokit/configs"
	"github.com/stretchr/testify/suite"
)

// SqlTestSuite defines the test suite for SQL database operations.
// It provides setup for testing SQL utilities with mocks.
type SqlTestSuite struct {
	suite.Suite

	connector  *MockConnector
	driverConn *MockPingDriverConn
	driver     *MockPingDriver
}

// TestSqlTestSuite runs the SQL test suite.
func TestSqlTestSuite(t *testing.T) {
	suite.Run(t, new(SqlTestSuite))
}

// SetupTest initializes the mock objects before each test.
func (s *SqlTestSuite) SetupTest() {
	s.connector = &MockConnector{}
	s.driverConn = &MockPingDriverConn{}
	s.driver = &MockPingDriver{}
}

// TestGetConnection tests the connection string formatting for SQL databases.
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
