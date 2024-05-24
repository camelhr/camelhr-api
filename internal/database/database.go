package database

import (
	"context"
	"database/sql"
	"errors"
)

// Database is an interface that defines the methods that a database should implement.
type Database interface {
	// Exec executes a query. Should be used for write operations.
	// When used inside WithTx, the query is executed within the transaction.
	// If dest is nil, the query is executed without expecting any result.
	// If dest is a pointer to a slice or an array, the query is executed and the result is stored in the slice or array.
	// If dest is a pointer to a struct, the query is executed and the result is stored in the struct.
	Exec(ctx context.Context, dest any, query string, args ...any) error

	// Get executes a query that is expected to return at most one row. Should be used for read operations.
	// When used inside WithTx, the query is executed within the transaction.
	// The result is stored in the dest.
	Get(ctx context.Context, dest any, query string, args ...any) error

	// List executes a query that is expected to return multiple rows. Should be used for read operations.
	// When used inside WithTx, the query is executed within the transaction.
	// The result is stored in the dest.
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
