// Package pg contains PostgreSQL database integration tests for the GoKit framework.
// These tests verify the proper functionality of PostgreSQL connection and operations.
package pg

import (
	"database/sql"
	"testing"

	"github.com/ralvescosta/gokit/configs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	mSQL "github.com/ralvescosta/gokit/sql"
)

// PostgresSqlTestSuite defines the test suite for PostgreSQL database operations.
// It provides setup and teardown for testing PostgreSQL connections with mocks.
type PostgresSqlTestSuite struct {
	suite.Suite

	connector  *mSQL.MockConnector
	driverConn *mSQL.MockPingDriverConn
	driver     *mSQL.MockPingDriver
}

// TestPostgresSqlTestSuite runs the PostgreSQL database test suite.
func TestPostgresSqlTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresSqlTestSuite))
}

// SetupTest initializes the mock objects before each test.
func (s *PostgresSqlTestSuite) SetupTest() {
	s.connector = &mSQL.MockConnector{}
	s.driverConn = &mSQL.MockPingDriverConn{}
	s.driver = &mSQL.MockPingDriver{}
}

// TestNew verifies the New function creates a PostgreSQL connection instance correctly.
func (s *PostgresSqlTestSuite) TestNew() {
	conn := New(&configs.Configs{SQLConfigs: &configs.SQLConfigs{}})

	s.IsType(&PostgresSqlConnection{}, conn)
}

// TestOpen tests the OpenTelemetry-enabled connection opening process.
func (s *PostgresSqlTestSuite) TestOpen() {
	s.driverConn.On("Ping", mock.AnythingOfType("context.backgroundCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("context.backgroundCtx")).Return(s.driverConn, nil)

	otelOpen = func(driverName, dsn string, opts ...otelsql.Option) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	conn := New(&configs.Configs{TracingConfigs: &configs.TracingConfigs{Enabled: true}, SQLConfigs: &configs.SQLConfigs{}})

	db, err := conn.Connect()

	s.NoError(err)
	s.IsType(&sql.DB{}, db)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

// TestConnectionPing tests the database ping functionality.
// Currently disabled/commented out.
func (s *PostgresSqlTestSuite) TestConnectionPing() {
	// s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	// s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	// conn := New(&logging.MockLogger{}, &env.Configs{SqlConfigs: &env.SqlConfigs{}})

	// sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
	// 	return sql.OpenDB(s.connector), nil
	// }

	// db, err := conn.Connect()

	// s.NoError(err)
	// s.IsType(&sql.DB{}, db)
	// s.driverConn.AssertExpectations(s.T())
	// s.connector.AssertExpectations(s.T())
}

// TestConnectionOpenErr tests error handling when connection opening fails.
// Currently disabled/commented out.
func (s *PostgresSqlTestSuite) TestConnectionOpenErr() {
	// conn := New(&logging.MockLogger{}, &env.Configs{SqlConfigs: &env.SqlConfigs{}})

	// sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
	// 	return nil, errors.New("")
	// }

	// _, err := conn.Connect()

	// s.Error(err)
}

// TestConnectionPingErr tests error handling when database ping fails.
// Currently disabled/commented out.
func (s *PostgresSqlTestSuite) TestConnectionPingErr() {
	// s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(errors.New("ping err"))
	// s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	// conn := New(&logging.MockLogger{}, &env.Configs{SqlConfigs: &env.SqlConfigs{}})

	// sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
	// 	return sql.OpenDB(s.connector), nil
	// }

	// _, err := conn.Connect()

	// s.Error(err)
	// s.driverConn.AssertExpectations(s.T())
	// s.connector.AssertExpectations(s.T())
}

// The following tests are commented out as they appear to be for functionality
// that may have been removed or is not yet implemented

// func (s *PostgresSqlTestSuite) TestShotdownSignalSignal() { ... }
// func (s *PostgresSqlTestSuite) TestShotdownSignalSignalIfSomeErrOccurBefore() { ... }
// func (s *PostgresSqlTestSuite) TestShotdownSignalSignalWithoutRequirements() { ... }
