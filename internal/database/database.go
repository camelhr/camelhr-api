package database

//go:generate mockery --name=Database --structname=DatabaseMock --inpackage --filename=database_mock.go

import (
	"context"
	"database/sql"
)

// Database is an interface that defines the methods that a database should implement.
type Database interface {
	// Exec executes a query. Should be used for write operations.
	Exec(ctx context.Context, dest any, query string, args ...any) error
	// Get executes a query that is expected to return at most one row. Should be used for read operations.
	Get(ctx context.Context, dest any, query string, args ...any) error
	// Select executes a query that is expected to return multiple rows. Should be used for read operations.
	Select(ctx context.Context, dest any, query string, args ...any) error
	// Transact executes the given function inside a transaction.
	Transact(ctx context.Context, fn func(*sql.Tx) error) error
}
