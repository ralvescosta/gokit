package pg

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
	mSQL "github.com/ralvescosta/gokit/sql"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PostgresSqlTestSuite struct {
	suite.Suite

	connector  *mSQL.MockConnector
	driverConn *mSQL.MockPingDriverConn
	driver     *mSQL.MockPingDriver
}

func TestPostgresSqlTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresSqlTestSuite))
}

func (s *PostgresSqlTestSuite) SetupTest() {
	s.connector = &mSQL.MockConnector{}
	s.driverConn = &mSQL.MockPingDriverConn{}
	s.driver = &mSQL.MockPingDriver{}
}

func (s *PostgresSqlTestSuite) TestNew() {
	conn := New(&logging.MockLogger{}, &env.Configs{})

	s.IsType(&PostgresSqlConnection{}, conn)
}

func (s *PostgresSqlTestSuite) TestOpen() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	otelOpen = func(driverName, dsn string, opts ...otelsql.Option) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	conn := New(&logging.MockLogger{}, &env.Configs{IS_TRACING_ENABLED: true})

	db, err := conn.Build()

	s.NoError(err)
	s.IsType(&sql.DB{}, db)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestConnectionPing() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	conn := New(&logging.MockLogger{}, &env.Configs{})

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	db, err := conn.Build()

	s.NoError(err)
	s.IsType(&sql.DB{}, db)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestConnectionOpenErr() {
	conn := New(&logging.MockLogger{}, &env.Configs{})

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("")
	}

	_, err := conn.Build()

	s.Error(err)
}

func (s *PostgresSqlTestSuite) TestConnectionPingErr() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(errors.New("ping err"))
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	conn := New(&logging.MockLogger{}, &env.Configs{})

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	_, err := conn.Build()

	s.Error(err)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestShotdownSignalSignal() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	conn := New(&logging.MockLogger{}, &env.Configs{
		SQL_DB_SECONDS_TO_PING: 10,
	})

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	db, err := conn.WthShotdownSig().Build()

	s.NoError(err)
	s.IsType(&sql.DB{}, db)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestShotdownSignalSignalIfSomeErrOccurBefore() {
	conn := New(&logging.MockLogger{}, &env.Configs{
		SQL_DB_SECONDS_TO_PING: 10,
	})

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("some err")
	}

	_, err := conn.WthShotdownSig().Build()

	s.Error(err)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestShotdownSignalSignalWithoutRequirements() {
	s.T().Skip()
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	conn := New(&logging.MockLogger{}, &env.Configs{})

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	_, err := conn.WthShotdownSig().Build()

	s.Error(err)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}
