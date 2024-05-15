package database_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	t.Parallel()

	t.Run("should call underlying ExecContext when dest is nil", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := database.NewPostgresDatabase(sqlxDB)

		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = pgDB.Exec(context.Background(), nil, "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should call underlying SelectContext when dest is slice", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := database.NewPostgresDatabase(sqlxDB)

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

	t.Run("should call underlying Get when dest is non-slice non-array pointer", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := database.NewPostgresDatabase(sqlxDB)

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

	t.Run("should return an error when the underlying ExecContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := database.NewPostgresDatabase(sqlxDB)

		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)

		err = pgDB.Exec(context.Background(), nil, "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		require.Error(t, err)
		require.ErrorIs(t, err, assert.AnError)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return an error when the underlying SelectContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := database.NewPostgresDatabase(sqlxDB)

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

	t.Run("should return an error when the underlying GetContext fails", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := database.NewPostgresDatabase(sqlxDB)

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
		pgDB := database.NewPostgresDatabase(sqlxDB)

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
		pgDB := database.NewPostgresDatabase(sqlxDB)

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
}

func TestTransact(t *testing.T) {
	t.Parallel()

	t.Run("should execute transaction", func(t *testing.T) {
		t.Parallel()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
		defer sqlxDB.Close()
		pgDB := database.NewPostgresDatabase(sqlxDB)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = pgDB.Transact(context.Background(), func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
			return err
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
		pgDB := database.NewPostgresDatabase(sqlxDB)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO users (.+) VALUES (.+)").
			WithArgs("John Doe", 30).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err = pgDB.Transact(context.Background(), func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
			return err
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
		pgDB := database.NewPostgresDatabase(sqlxDB)

		mock.ExpectBegin()
		mock.ExpectRollback()

		require.Panics(t, func() {
			_ = pgDB.Transact(context.Background(), func(tx *sql.Tx) error {
				panic("panic")
			})
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
