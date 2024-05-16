package database

//go:generate mockery --name=Database --structname=DatabaseMock --inpackage --filename=database_mock.go

import (
	"context"
	"database/sql"
	"errors"
)

// Database is an interface that defines the methods that a database should implement.
type Database interface {
	// Exec executes a query. Should be used for write operations.
	Exec(ctx context.Context, dest any, query string, args ...any) error
	// Get executes a query that is expected to return at most one row. Should be used for read operations.
	Get(ctx context.Context, dest any, query string, args ...any) error
	// List executes a query that is expected to return multiple rows. Should be used for read operations.
	List(ctx context.Context, dest any, query string, args ...any) error

	Transactor
}

// SuppressNoRowsError suppresses the sql.ErrNoRows error.
func SuppressNoRowsError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return err
}
