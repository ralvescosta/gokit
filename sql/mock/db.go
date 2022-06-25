package mock

import (
	"context"
	"database/sql"
	"database/sql/driver"

	"github.com/stretchr/testify/mock"
)

type MockPingDriver struct {
	mock.Mock
}

type MockPingDriverConn struct {
	MockSqlDbConn
	driver *MockPingDriver
	mock.Mock
}

func (m MockPingDriverConn) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockPingDriver) Open(name string) (driver.Conn, error) {
	args := m.Called(name)
	c := args.Get(0).(driver.Conn)
	return c, args.Error(1)
}

type MockRows struct {
	mock.Mock
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

type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	args := m.Called()
	return int64(args.Int(0)), args.Error(1)
}

type MockStmt struct {
	mock.Mock
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

type MockSqlDbConn struct {
	mock.Mock
}

func (m MockSqlDbConn) Prepare(query string) (driver.Stmt, error) {
	args := m.Called(query)
	stmt := args.Get(0).(driver.Stmt)
	return stmt, args.Error(1)
}

func (m MockSqlDbConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m MockSqlDbConn) Begin() (driver.Tx, error) {
	args := m.Called()
	tx := args.Get(0).(driver.Tx)
	return tx, args.Error(1)
}

func (m MockSqlDbConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	mArgs := m.Called(query, args)
	r := mArgs.Get(0).(driver.Result)
	return r, mArgs.Error(1)
}

type MockConnector struct {
	mock.Mock
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

func Oi() {
	driver := &MockConnector{}
	sql.OpenDB(driver)
}
