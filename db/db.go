package db

import (
	"context"
	"database/sql"
)

// SQLDatabase ...
type SQLDatabase interface {
	PingContext(ctx context.Context) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
