package database

import "context"

//go:generate mockery --name=Transactor --structname=TransactorMock --inpackage --filename=transactor_mock.go

// Transactor defines the methods related to database transaction.
type Transactor interface {
	// WithTx executes the given function inside a transaction.
	WithTx(ctx context.Context, txFn func(ctx context.Context) error) error
}
