package pg

import (
	"database/sql"
	"testing"

	"github.com/ralvescostati/pkgs/env"
	"github.com/ralvescostati/pkgs/logger/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type PostgresSqlTestSuite struct {
	suite.Suite
}

func TestPostgresSqlTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresSqlTestSuite))
}

func (s *PostgresSqlTestSuite) TestNew() {
	var sh chan bool
	conn := New(mock.NewMockLogger(), &env.Configs{}, sh)

	s.IsType(&PostgresSqlConnection{}, conn)
}

func (s *PostgresSqlTestSuite) TestConnection() {
	var sh chan bool
	conn := New(mock.NewMockLogger(), &env.Configs{}, sh)
	db, mock, _ := sqlmock.New()

	open = func(driverName, dataSourceName string) (*sql.DB, error) {
		return db, nil
	}
	mock.ExpectPing()

	conn.Connect()
}

func (s *PostgresSqlTestSuite) TestConnectionErr() {
	// var sh chan bool
	// conn := New(mock.NewMockLogger(), &env.Configs{}, sh)

	// open = func(driverName, dataSourceName string) (*sql.DB, error) {
	// 	return nil, errors.New("")
	// }

	// _, err := conn.Connect().Build()

	// s.Error(err)
}
