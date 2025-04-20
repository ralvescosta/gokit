// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package sql

import (
	"context"
	"database/sql/driver"

	"github.com/stretchr/testify/mock"
)

type (
	MockPingDriver struct {
		mock.Mock
	}

	MockPingDriverConn struct {
		MockSQLDbConn
		driver *MockPingDriver
		mock.Mock
	}

	MockRows struct {
		mock.Mock
	}

	MockResult struct {
		mock.Mock
	}

	MockStmt struct {
		mock.Mock
	}

	MockSQLDbConn struct {
		mock.Mock
	}

	MockConnector struct {
		mock.Mock
	}
)

func (m MockPingDriverConn) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPingDriver) Open(name string) (driver.Conn, error) {
	args := m.Called(name)
	c := args.Get(0).(driver.Conn)
	return c, args.Error(1)
}

func (m *MockRows) Columns() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockRows) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) Next(dest []driver.Value) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockStmt) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStmt) NumInput() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockStmt) Exec(args []driver.Value) (driver.Result, error) {
	mArgs := m.Called(args)
	d := mArgs.Get(0).(driver.Result)
	return d, mArgs.Error(1)
}

func (m *MockStmt) Query(args []driver.Value) (driver.Rows, error) {
	mArgs := m.Called(args)
	d := mArgs.Get(0).(driver.Rows)
	return d, mArgs.Error(1)
}

func (m MockSQLDbConn) Prepare(query string) (driver.Stmt, error) {
	args := m.Called(query)
	stmt := args.Get(0).(driver.Stmt)
	return stmt, args.Error(1)
}

func (m MockSQLDbConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m MockSQLDbConn) Begin() (driver.Tx, error) {
	args := m.Called()
	tx := args.Get(0).(driver.Tx)
	return tx, args.Error(1)
}

func (m MockSQLDbConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	mArgs := m.Called(query, args)
	r := mArgs.Get(0).(driver.Result)
	return r, mArgs.Error(1)
}

func (m *MockConnector) Connect(ctx context.Context) (driver.Conn, error) {
	args := m.Called(ctx)
	c := args.Get(0).(driver.Conn)
	return c, args.Error(1)
}

func (m *MockConnector) Driver() driver.Driver {
	args := m.Called()
	d := args.Get(0).(driver.Driver)
	return d
}

func NewMockMockPingDriver() *MockPingDriver {
	return new(MockPingDriver)
}

func NewMockPingDriverConn() *MockPingDriverConn {
	return new(MockPingDriverConn)
}

func NewMockRows() *MockRows {
	return new(MockRows)
}

func NewMockResult() *MockResult {
	return new(MockResult)
}

func NewMockStmt() *MockStmt {
	return new(MockStmt)
}

func NewMockSQLDbConn() *MockSQLDbConn {
	return new(MockSQLDbConn)
}

func NewMockConnector() *MockConnector {
	return new(MockConnector)
}
