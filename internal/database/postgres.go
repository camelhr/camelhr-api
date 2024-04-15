package database

import (
	"context"
	"database/sql"
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

func (p *postgresDatabase) Select(ctx context.Context, dest any, query string, args ...any) error {
	return p.db.SelectContext(ctx, dest, query, args...)
}

func (p *postgresDatabase) Transact( //nolint:nonamedreturns // named return is used to simplify the error handling
	ctx context.Context,
	fn func(*sql.Tx) error,
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
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func isSliceOrArray(dest any) bool {
	kind := reflect.Indirect(reflect.ValueOf(dest)).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}
