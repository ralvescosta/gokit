# SQL Package

The SQL package provides utilities and interfaces for working with SQL databases in the GoKit framework. It includes functionality for connection string generation, mock implementations for testing, and database-specific implementations for PostgreSQL.

## Overview

This package provides tools and utilities for:

- Connection string formatting
- Database connection management
- Mocks for SQL database testing
- PostgreSQL-specific implementation

## Usage

### Getting a Connection String

```go
import (
    "github.com/ralvescosta/gokit/configs"
    "github.com/ralvescosta/gokit/sql"
)

func main() {
    // Create SQL configurations
    cfg := &configs.SQLConfigs{
        Host:     "localhost",
        Port:     "5432",
        User:     "postgres",
        Password: "password",
        DbName:   "mydb",
    }

    // Generate connection string
    connStr := sql.GetConnectionString(cfg)
    // connStr will be: "host=localhost port=5432 user=postgres password=password dbname=mydb sslmode=disable"
}
```

### PostgreSQL Connection

The package provides a PostgreSQL-specific implementation in the `pg` subpackage:

```go
import (
    "github.com/ralvescosta/gokit/configs"
    "github.com/ralvescosta/gokit/sql/postgres"
)

func main() {
    // Create application configurations
    cfgs := &configs.Configs{
        SQLConfigs: &configs.SQLConfigs{
            Host:     "localhost",
            Port:     "5432",
            User:     "postgres",
            Password: "password",
            DbName:   "mydb",
        },
        TracingConfigs: &configs.TracingConfigs{
            Enabled: true, // Set to true to enable OpenTelemetry tracing
        },
        Logger: logger, // Your logger implementation
    }

    // Create PostgreSQL connection
    pgConn := pg.New(cfgs)

    // Connect to database
    db, err := pgConn.Connect()
    if err != nil {
        // Handle error
    }

    // Use the database connection
    // ...
}
```

## Testing

The package provides mock implementations for testing SQL database code:

```go
import (
    "testing"

    "github.com/ralvescosta/gokit/sql"
    "github.com/stretchr/testify/suite"
)

type MyTestSuite struct {
    suite.Suite

    connector  *sql.MockConnector
    driverConn *sql.MockPingDriverConn
    driver     *sql.MockPingDriver
}

func TestMyTestSuite(t *testing.T) {
    suite.Run(t, new(MyTestSuite))
}

func (s *MyTestSuite) SetupTest() {
    s.connector = sql.NewMockConnector()
    s.driverConn = sql.NewMockPingDriverConn()
    s.driver = sql.NewMockMockPingDriver()

    // Set up expectations
    // ...
}

func (s *MyTestSuite) TestMyFunction() {
    // Test using mocks
    // ...
}
```

## Features

### Connection Management

- Connection string generation compatible with PostgreSQL
- Connection handling with proper error reporting
- OpenTelemetry instrumented connections for tracing

### Mock Objects for Testing

The package provides mock implementations for testing:

- `MockPingDriver`: Mocks driver.Driver for testing
- `MockPingDriverConn`: Mocks a pingable connection
- `MockRows`: Mocks driver.Rows for testing row operations
- `MockResult`: Mocks driver.Result for testing query results
- `MockStmt`: Mocks driver.Stmt for testing prepared statements
- `MockSQLDbConn`: Mocks driver.Conn for testing connections
- `MockConnector`: Mocks driver.Connector for testing connection creation

### Database-Specific Implementations

Currently, the package provides:

- PostgreSQL implementation with OpenTelemetry support
- Empty MySQL placeholder for future implementation

## OpenTelemetry Integration

The package provides built-in support for OpenTelemetry tracing:

- When tracing is enabled in the configuration, database connections use the OpenTelemetry-instrumented version
- Adds database attributes to spans for better visibility
- Includes database name and system in the traces

## License

MIT License
