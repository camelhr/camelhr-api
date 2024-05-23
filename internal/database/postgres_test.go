package database //nolint:testpackage // since ctxTxKey is not exported

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) { //nolint:maintidx // test function
	t.Parallel()

	t.Run("should call underlying db.ExecContext when dest is nil", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = pgDB.Exec(context.Background(), nil, "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should call underlying tx.ExecContext when dest is nil", func(t *testing.T) {
		t.Parallel()

		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		txMockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer txMockDB.Close()

		sqlxTxDB := sqlx.NewDb(txMockDB, "sqlmock")
		defer sqlxTxDB.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnResult(sqlmock.NewResult(1, 1))

		tx, err := sqlxTxDB.BeginTxx(context.Background(), nil)
		require.NoError(t, err)
		// inject transaction to context
		ctx := context.WithValue(context.Background(), ctxTxKey, tx)

		err = pgDB.Exec(ctx, nil, "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should call underlying db.SelectContext when dest is slice", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("John Doe", 30)
		mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnRows(rows)

		var users []struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		err = pgDB.Exec(context.Background(), &users,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING name, age", "John Doe", 30)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
		assert.Equal(t, 30, users[0].Age)
	})

	t.Run("should call underlying tx.SelectContext when dest is slice", func(t *testing.T) {
		t.Parallel()

		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		txMockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer txMockDB.Close()

		sqlxTxDB := sqlx.NewDb(txMockDB, "sqlmock")
		defer sqlxTxDB.Close()

		mock.ExpectBegin()

		rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("John Doe", 30)
		mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnRows(rows)

		var users []struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		tx, err := sqlxTxDB.BeginTxx(context.Background(), nil)
		require.NoError(t, err)
		// inject transaction to context
		ctx := context.WithValue(context.Background(), ctxTxKey, tx)

		err = pgDB.Exec(ctx, &users,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING name, age", "John Doe", 30)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
		assert.Equal(t, 30, users[0].Age)
	})

	t.Run("should call underlying db.Get when dest is non-slice non-array pointer", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnRows(rows)

		var id *int64
		err = pgDB.Exec(context.Background(), &id,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", "John Doe", 30)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		require.NotNil(t, id)
		assert.Equal(t, int64(1), *id)
	})

	t.Run("should call underlying tx.Get when dest is non-slice non-array pointer", func(t *testing.T) {
		t.Parallel()

		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		txMockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer txMockDB.Close()

		sqlxTxDB := sqlx.NewDb(txMockDB, "sqlmock")
		defer sqlxTxDB.Close()

		mock.ExpectBegin()

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnRows(rows)

		var id *int64

		tx, err := sqlxTxDB.BeginTxx(context.Background(), nil)
		require.NoError(t, err)
		// inject transaction to context
		ctx := context.WithValue(context.Background(), ctxTxKey, tx)

		err = pgDB.Exec(ctx, &id,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", "John Doe", 30)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		require.NotNil(t, id)
		assert.Equal(t, int64(1), *id)
	})

	t.Run("should return an error when the underlying db.ExecContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)

		err = pgDB.Exec(context.Background(), nil, "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return an error when the underlying tx.ExecContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		txMockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer txMockDB.Close()

		sqlxTxDB := sqlx.NewDb(txMockDB, "sqlmock")
		defer sqlxTxDB.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)

		tx, err := sqlxTxDB.BeginTxx(context.Background(), nil)
		require.NoError(t, err)
		// inject transaction to context
		ctx := context.WithValue(context.Background(), ctxTxKey, tx)

		err = pgDB.Exec(ctx, nil, "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return an error when the underlying db.SelectContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)

		var users []struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		err = pgDB.Exec(context.Background(), &users,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING name, age", "John Doe", 30)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return an error when the underlying tx.SelectContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		txMockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer txMockDB.Close()

		sqlxTxDB := sqlx.NewDb(txMockDB, "sqlmock")
		defer sqlxTxDB.Close()

		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)

		var users []struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		tx, err := sqlxTxDB.BeginTxx(context.Background(), nil)
		require.NoError(t, err)
		// inject transaction to context
		ctx := context.WithValue(context.Background(), ctxTxKey, tx)

		err = pgDB.Exec(ctx, &users,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING name, age", "John Doe", 30)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return an error when the underlying db.GetContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)

		var id *int64
		err = pgDB.Exec(context.Background(), &id,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", "John Doe", 30)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return an error when the underlying tx.GetContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		txMockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer txMockDB.Close()

		sqlxTxDB := sqlx.NewDb(txMockDB, "sqlmock")
		defer sqlxTxDB.Close()

		mock.ExpectBegin()

		mock.ExpectQuery("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)

		var id *int64

		tx, err := sqlxTxDB.BeginTxx(context.Background(), nil)
		require.NoError(t, err)
		// inject transaction to context
		ctx := context.WithValue(context.Background(), ctxTxKey, tx)

		err = pgDB.Exec(ctx, &id,
			"INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id", "John Doe", 30)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("should populate provided type with result", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("John Doe", 30)
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = (.+)").
			WithArgs(1).
			WillReturnRows(rows)

		var user struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		err = pgDB.Get(context.Background(), &user, "SELECT name, age FROM users WHERE id = $1", 1)
		require.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, 30, user.Age)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should populate provided type with result using context injected tx", func(t *testing.T) {
		t.Parallel()

		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		txMockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer txMockDB.Close()

		sqlxTxDB := sqlx.NewDb(txMockDB, "sqlmock")
		defer sqlxTxDB.Close()

		mock.ExpectBegin()

		rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("John Doe", 30)
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id = (.+)").
			WithArgs(1).
			WillReturnRows(rows)

		tx, err := sqlxTxDB.BeginTxx(context.Background(), nil)
		require.NoError(t, err)
		// inject transaction to context
		ctx := context.WithValue(context.Background(), ctxTxKey, tx)

		var user struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		err = pgDB.Get(ctx, &user, "SELECT name, age FROM users WHERE id = $1", 1)
		require.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, 30, user.Age)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("should populate provided slice with result", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("John Doe", 30)
		mock.ExpectQuery("SELECT (.+) FROM users WHERE age > (.+)").
			WithArgs(25).
			WillReturnRows(rows)

		var users []struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		err = pgDB.List(context.Background(), &users, "SELECT name, age FROM users WHERE age > $1", 25)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
		assert.Equal(t, 30, users[0].Age)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should populate provided slice with result using context injected tx", func(t *testing.T) {
		t.Parallel()

		mockDB, _, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		txMockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer txMockDB.Close()

		sqlxTxDB := sqlx.NewDb(txMockDB, "sqlmock")
		defer sqlxTxDB.Close()

		mock.ExpectBegin()

		rows := sqlmock.NewRows([]string{"name", "age"}).AddRow("John Doe", 30)
		mock.ExpectQuery("SELECT (.+) FROM users WHERE age > (.+)").
			WithArgs(25).
			WillReturnRows(rows)

		tx, err := sqlxTxDB.BeginTxx(context.Background(), nil)
		require.NoError(t, err)
		// inject transaction to context
		ctx := context.WithValue(context.Background(), ctxTxKey, tx)

		var users []struct {
			Name string `db:"name"`
			Age  int    `db:"age"`
		}

		err = pgDB.List(ctx, &users, "SELECT name, age FROM users WHERE age > $1", 25)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
		assert.Equal(t, 30, users[0].Age)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestWithTx(t *testing.T) {
	t.Parallel()

	t.Run("should execute transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = pgDB.WithTx(ctx, func(ctx context.Context) error {
			require.NotNil(t, ctx.Value(ctxTxKey))
			return pgDB.Exec(ctx, nil, "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		})
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should rollback the transaction if an error occurs", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err = pgDB.WithTx(context.Background(), func(ctx context.Context) error {
			require.NotNil(t, ctx.Value(ctxTxKey))
			return pgDB.Exec(ctx, nil, "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		})
		require.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should rollback the transaction if a panic occurs", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := NewPostgresDatabase(sqlxDB)

		mock.ExpectBegin()
		mock.ExpectRollback()

		require.Panics(t, func() {
			_ = pgDB.WithTx(context.Background(), func(ctx context.Context) error {
				panic("panic")
			})
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
