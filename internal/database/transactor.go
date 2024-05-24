package database

import "context"

// Transactor defines the methods related to database transaction.
type Transactor interface {
	// WithTx executes the given function inside a transaction.
	WithTx(ctx context.Context, txFn func(ctx context.Context) error) error
}
