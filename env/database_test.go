package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) SetupTest() {
	dotEnvConfig = func(path string) error {
		return nil
	}
}

func (s *DatabaseTestSuite) TestDatabase() {
	os.Setenv(GO_ENV_KEY, "dev")
	os.Setenv(SQL_DB_HOST_ENV_KEY, "host")
	os.Setenv(SQL_DB_PORT_ENV_KEY, "port")
	os.Setenv(SQL_DB_USER_ENV_KEY, "user")
	os.Setenv(SQL_DB_PASSWORD_ENV_KEY, "password")
	os.Setenv(SQL_DB_NAME_ENV_KEY, "name")

	cfg, err := New().Database().Build()

	s.NoError(err)
	s.Equal(cfg.SQL_DB_HOST, "host")
	s.Equal(cfg.SQL_DB_PORT, "port")
	s.Equal(cfg.SQL_DB_USER, "user")
	s.Equal(cfg.SQL_DB_PASSWORD, "password")
	s.Equal(cfg.SQL_DB_NAME, "name")
}

func (s *DatabaseTestSuite) TestDatabaseErr() {
	os.Setenv(GO_ENV_KEY, "dev")
	os.Setenv(SQL_DB_HOST_ENV_KEY, "")

	_, err := New().Database().Build()
	s.Error(err)

	os.Setenv(SQL_DB_HOST_ENV_KEY, "host")
	os.Setenv(SQL_DB_PORT_ENV_KEY, "")

	_, err = New().Database().Build()
	s.Error(err)

	os.Setenv(SQL_DB_PORT_ENV_KEY, "post")
	os.Setenv(SQL_DB_USER_ENV_KEY, "")

	_, err = New().Database().Build()
	s.Error(err)

	os.Setenv(SQL_DB_USER_ENV_KEY, "user")
	os.Setenv(SQL_DB_PASSWORD_ENV_KEY, "")

	_, err = New().Database().Build()
	s.Error(err)

	os.Setenv(SQL_DB_PASSWORD_ENV_KEY, "password")
	os.Setenv(SQL_DB_NAME_ENV_KEY, "")

	_, err = New().Database().Build()
	s.Error(err)
}
