package mock

import (
	"context"
	"database/sql/driver"
	"errors"
)

type BadConn struct{}

func (bc BadConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("BadConn Prepare")
}

func (bc BadConn) Close() error {
	return nil
}

func (bc BadConn) Begin() (driver.Tx, error) {
	return nil, errors.New("BadConn Begin")
}

func (bc BadConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	panic("BadConn.Exec")
}

type BadDriver struct{}

func (bd BadDriver) Open(name string) (driver.Conn, error) {
	return BadConn{}, nil
}

type PingDriver struct {
	Fails bool
}

type pingConn struct {
	BadConn
	driver *PingDriver
}

var pingError = errors.New("Ping failed")

func (pc pingConn) Ping(ctx context.Context) error {
	if pc.driver.Fails {
		return pingError
	}
	return nil
}

var _ driver.Pinger = pingConn{}

func (pd *PingDriver) Open(name string) (driver.Conn, error) {
	return pingConn{driver: pd}, nil
}
