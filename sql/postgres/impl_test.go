package pg

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logging"
	mSQL "github.com/ralvescostati/pkgs/sql"

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
	var sh chan bool
	conn := New(&logging.MockLogger{}, &env.Configs{}, sh)

	s.IsType(&PostgresSqlConnection{}, conn)
}

func (s *PostgresSqlTestSuite) TestConnectionPing() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	sh := make(chan bool)
	conn := New(&logging.MockLogger{}, &env.Configs{}, sh)

	open = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	db, err := conn.Connect().Build()

	s.NoError(err)
	s.IsType(&sql.DB{}, db)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestConnectionOpenErr() {
	var sh chan bool
	conn := New(&logging.MockLogger{}, &env.Configs{}, sh)

	open = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("")
	}

	_, err := conn.Connect().Build()

	s.Error(err)
}

func (s *PostgresSqlTestSuite) TestConnectionPingErr() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(errors.New("ping err"))
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	sh := make(chan bool)
	conn := New(&logging.MockLogger{}, &env.Configs{}, sh)

	open = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	_, err := conn.Connect().Build()

	s.Error(err)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestShotdownSignalSignal() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	sh := make(chan bool)
	conn := New(&logging.MockLogger{}, &env.Configs{
		SQL_DB_SECONDS_TO_PING: 10,
	}, sh)

	open = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	db, err := conn.Connect().ShotdownSignal().Build()

	s.NoError(err)
	s.IsType(&sql.DB{}, db)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestShotdownSignalSignalIfSomeErrOccurBefore() {
	sh := make(chan bool)
	conn := New(&logging.MockLogger{}, &env.Configs{
		SQL_DB_SECONDS_TO_PING: 10,
	}, sh)

	open = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("some err")
	}

	_, err := conn.Connect().ShotdownSignal().Build()

	s.Error(err)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}

func (s *PostgresSqlTestSuite) TestShotdownSignalSignalWithoutRequirements() {
	s.driverConn.On("Ping", mock.AnythingOfType("*context.emptyCtx")).Return(nil)
	s.connector.On("Connect", mock.AnythingOfType("*context.emptyCtx")).Return(s.driverConn, nil)

	conn := New(&logging.MockLogger{}, &env.Configs{}, nil)

	open = func(driverName, dataSourceName string) (*sql.DB, error) {
		return sql.OpenDB(s.connector), nil
	}

	_, err := conn.Connect().ShotdownSignal().Build()

	s.Error(err)
	s.driverConn.AssertExpectations(s.T())
	s.connector.AssertExpectations(s.T())
}
