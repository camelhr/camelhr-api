package database

import (
	"context"
	"reflect"

	"github.com/jmoiron/sqlx"
)

type postgresDatabase struct {
	db *sqlx.DB
}

// NewPostgresDatabase creates a new instance of the postgresDatabase.
// Currently we are executing the queries directly.
// The idea to use different connections for read and write operations.
// The Exec method will be used for write operations. Get & Select methods will be used for readonly operations.
func NewPostgresDatabase(db *sqlx.DB) Database {
	return &postgresDatabase{db: db}
}

func (p *postgresDatabase) Exec(ctx context.Context, dest any, query string, args ...any) error {
	if dest == nil {
		_, err := p.db.ExecContext(ctx, query, args...)
		return err
	}

	if isSliceOrArray(dest) {
		return p.db.SelectContext(ctx, dest, query, args...)
	}

	return p.db.GetContext(ctx, dest, query, args...)
}

func (p *postgresDatabase) Get(ctx context.Context, dest any, query string, args ...any) error {
	return p.db.GetContext(ctx, dest, query, args...)
}

func (p *postgresDatabase) List(ctx context.Context, dest any, query string, args ...any) error {
	return p.db.SelectContext(ctx, dest, query, args...)
}

func (p *postgresDatabase) WithTx(
	ctx context.Context,
	txFn func(ctx context.Context) error,
) (err error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback() //nolint:errcheck // ignore error inside recover
			panic(p)      // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() //nolint:errcheck // err is non-nil, don't update the err
		} else {
			// err is nil, commit the transaction
			// if commit fails, return the error but do not call rollback
			// since the transaction won't be valid after commit is called
			// irrespective of the commit result
			err = tx.Commit()
		}
	}()

	err = txFn(ctx)

	return err
}

func isSliceOrArray(dest any) bool {
	kind := reflect.Indirect(reflect.ValueOf(dest)).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}
