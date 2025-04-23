// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package sql

import (
	"context"
	"database/sql/driver"

	"github.com/stretchr/testify/mock"
)

// Mock types for SQL database testing

type (
	// MockPingDriver implements a mock driver.Driver for testing purposes.
	// It allows mocking the Open method.
	MockPingDriver struct {
		mock.Mock
	}

	// MockPingDriverConn extends MockSQLDbConn and implements a pingable connection
	// for testing database ping functionality.
	MockPingDriverConn struct {
		MockSQLDbConn
		driver *MockPingDriver
		mock.Mock
	}

	// MockRows implements driver.Rows interface for testing purposes.
	// It allows mocking database row operations during testing.
	MockRows struct {
		mock.Mock
	}

	// MockResult implements driver.Result interface for testing purposes.
	// It allows mocking database execution results during testing.
	MockResult struct {
		mock.Mock
	}

	// MockStmt implements driver.Stmt interface for testing purposes.
	// It allows mocking prepared statement operations during testing.
	MockStmt struct {
		mock.Mock
	}

	// MockSQLDbConn implements driver.Conn interface for testing purposes.
	// It allows mocking database connection operations during testing.
	MockSQLDbConn struct {
		mock.Mock
	}

	// MockConnector implements driver.Connector interface for testing purposes.
	// It allows mocking the database connector functionality during testing.
	MockConnector struct {
		mock.Mock
	}
)

// Ping mocks the Ping method for testing ping functionality.
func (m MockPingDriverConn) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Open mocks the Open method for testing driver connection opening.
func (m *MockPingDriver) Open(name string) (driver.Conn, error) {
	args := m.Called(name)
	c := args.Get(0).(driver.Conn)
	return c, args.Error(1)
}

// Columns mocks the Columns method for testing row column names retrieval.
func (m *MockRows) Columns() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// Close mocks the Close method for testing row closing.
func (m *MockRows) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Next mocks the Next method for testing row iteration.
func (m *MockRows) Next(dest []driver.Value) error {
	args := m.Called(dest)
	return args.Error(0)
}

// LastInsertId mocks the LastInsertId method for testing last inserted ID retrieval.
func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

// RowsAffected mocks the RowsAffected method for testing row count retrieval.
func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

// Close mocks the Close method for testing statement closing.
func (m *MockStmt) Close() error {
	args := m.Called()
	return args.Error(0)
}

// NumInput mocks the NumInput method for testing parameter count retrieval.
func (m *MockStmt) NumInput() int {
	args := m.Called()
	return args.Int(0)
}

// Exec mocks the Exec method for testing statement execution.
func (m *MockStmt) Exec(args []driver.Value) (driver.Result, error) {
	mArgs := m.Called(args)
	d := mArgs.Get(0).(driver.Result)
	return d, mArgs.Error(1)
}

// Query mocks the Query method for testing statement querying.
func (m *MockStmt) Query(args []driver.Value) (driver.Rows, error) {
	mArgs := m.Called(args)
	d := mArgs.Get(0).(driver.Rows)
	return d, mArgs.Error(1)
}

// Prepare mocks the Prepare method for testing statement preparation.
func (m MockSQLDbConn) Prepare(query string) (driver.Stmt, error) {
	args := m.Called(query)
	stmt := args.Get(0).(driver.Stmt)
	return stmt, args.Error(1)
}

// Close mocks the Close method for testing connection closing.
func (m MockSQLDbConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Begin mocks the Begin method for testing transaction initiation.
func (m MockSQLDbConn) Begin() (driver.Tx, error) {
	args := m.Called()
	tx := args.Get(0).(driver.Tx)
	return tx, args.Error(1)
}

// Exec mocks the Exec method for testing direct query execution.
func (m MockSQLDbConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	mArgs := m.Called(query, args)
	r := mArgs.Get(0).(driver.Result)
	return r, mArgs.Error(1)
}

// Connect mocks the Connect method for testing connector connection.
func (m *MockConnector) Connect(ctx context.Context) (driver.Conn, error) {
	args := m.Called(ctx)
	c := args.Get(0).(driver.Conn)
	return c, args.Error(1)
}

// Driver mocks the Driver method for testing driver retrieval.
func (m *MockConnector) Driver() driver.Driver {
	args := m.Called()
	d := args.Get(0).(driver.Driver)
	return d
}

// NewMockMockPingDriver creates a new instance of MockPingDriver for testing.
func NewMockMockPingDriver() *MockPingDriver {
	return new(MockPingDriver)
}

// NewMockPingDriverConn creates a new instance of MockPingDriverConn for testing.
func NewMockPingDriverConn() *MockPingDriverConn {
	return new(MockPingDriverConn)
}

// NewMockRows creates a new instance of MockRows for testing.
func NewMockRows() *MockRows {
	return new(MockRows)
}

// NewMockResult creates a new instance of MockResult for testing.
func NewMockResult() *MockResult {
	return new(MockResult)
}

// NewMockStmt creates a new instance of MockStmt for testing.
func NewMockStmt() *MockStmt {
	return new(MockStmt)
}

// NewMockSQLDbConn creates a new instance of MockSQLDbConn for testing.
func NewMockSQLDbConn() *MockSQLDbConn {
	return new(MockSQLDbConn)
}

// NewMockConnector creates a new instance of MockConnector for testing.
func NewMockConnector() *MockConnector {
	return new(MockConnector)
}
