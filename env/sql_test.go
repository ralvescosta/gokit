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
	os.Setenv(SQL_DB_SECONDS_TO_PING_ENV_KEY, "1")

	cfg, err := New().SqlDatabase().Build()

	s.NoError(err)
	s.NotNil(cfg.SqlConfigs)
	s.Equal(cfg.SqlConfigs.Host, "host")
	s.Equal(cfg.SqlConfigs.Port, "port")
	s.Equal(cfg.SqlConfigs.User, "user")
	s.Equal(cfg.SqlConfigs.Password, "password")
	s.Equal(cfg.SqlConfigs.DbName, "name")
	s.Equal(cfg.SqlConfigs.SecondsToPing, 1)

	cfg, err = New().Build()

	s.Nil(err)
	s.Nil(cfg.SqlConfigs)
}

func (s *DatabaseTestSuite) TestDatabaseErr() {
	os.Setenv(GO_ENV_KEY, "")

	_, err := New().SqlDatabase().Build()
	s.Error(err)

	os.Setenv(GO_ENV_KEY, "dev")
	os.Setenv(SQL_DB_HOST_ENV_KEY, "")

	_, err = New().SqlDatabase().Build()
	s.Error(err)

	os.Setenv(SQL_DB_HOST_ENV_KEY, "host")
	os.Setenv(SQL_DB_PORT_ENV_KEY, "")

	_, err = New().SqlDatabase().Build()
	s.Error(err)

	os.Setenv(SQL_DB_PORT_ENV_KEY, "port")
	os.Setenv(SQL_DB_USER_ENV_KEY, "")

	_, err = New().SqlDatabase().Build()
	s.Error(err)

	os.Setenv(SQL_DB_USER_ENV_KEY, "user")
	os.Setenv(SQL_DB_PASSWORD_ENV_KEY, "")

	_, err = New().SqlDatabase().Build()
	s.Error(err)

	os.Setenv(SQL_DB_PASSWORD_ENV_KEY, "password")
	os.Setenv(SQL_DB_NAME_ENV_KEY, "")

	_, err = New().SqlDatabase().Build()
	s.Error(err)

	os.Setenv(SQL_DB_NAME_ENV_KEY, "name")
	os.Setenv(SQL_DB_SECONDS_TO_PING_ENV_KEY, "")

	_, err = New().SqlDatabase().Build()
	s.Error(err)
}
