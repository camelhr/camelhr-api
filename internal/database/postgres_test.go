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

	t.Run("should call the Exec method", func(t *testing.T) {
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

		_, err = pgDB.Exec(context.TODO(), "INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
		assert.NoError(t, err)
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("should call the Get method", func(t *testing.T) {
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

		err = pgDB.Get(context.TODO(), &user, "SELECT name, age FROM users WHERE id = $1", 1)
		require.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, 30, user.Age)
	})
}

func TestSelect(t *testing.T) {
	t.Parallel()

	t.Run("should call the Select method", func(t *testing.T) {
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

		err = pgDB.Select(context.TODO(), &users, "SELECT name, age FROM users WHERE age > $1", 25)
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
		assert.Equal(t, 30, users[0].Age)
	})
}

func TestTransact(t *testing.T) {
	t.Parallel()

	t.Run("should call the Transact method", func(t *testing.T) {
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

		err = pgDB.Transact(context.TODO(), func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
			return err
		})
		assert.NoError(t, err)
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

		err = pgDB.Transact(context.TODO(), func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
			return err
		})
		assert.Error(t, err)
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

		assert.Panics(t, func() {
			_ = pgDB.Transact(context.TODO(), func(tx *sql.Tx) error {
				panic("panic")
			})
		})
	})

	t.Run("should rollback the transaction if an error occurs during commit", func(t *testing.T) {
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
		mock.ExpectCommit().WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err = pgDB.Transact(context.TODO(), func(tx *sql.Tx) error {
			_, err := tx.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "John Doe", 30)
			return err
		})
		assert.Error(t, err)
	})
}
