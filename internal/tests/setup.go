package tests

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/camelhr/camelhr-api/migrations/schema"
	"github.com/camelhr/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
)

const resourceExpirationTime = 900

type PostgreSQLContainer struct {
	resource *dockertest.Resource
	pool     *dockertest.Pool
}

const (
	pgUser     = "postgres"
	pgPassword = "postgres"
	pgDBName   = "camelhr_test_db"
)

// NewPostgresContainer creates a new postgres database docker container for testing.
func NewPostgresContainer() (*PostgreSQLContainer, error) {
	container := &PostgreSQLContainer{}

	// create a new pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not prepare docker image: %w", err)
	}

	if err = pool.Client.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping docker server: %w", err)
	}

	container.pool = pool

	// start a new postgres container
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.2-alpine",
		Cmd: []string{
			"postgres",
			"-c", "log_statement=all",
			"-c", "log_error_verbosity=verbose",
		},
		Env: []string{
			"POSTGRES_USER=" + pgUser,
			"POSTGRES_PASSWORD=" + pgPassword,
			"POSTGRES_DB=" + pgDBName,
		},
		ExposedPorts: []string{"5432/tcp"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %w", err)
	}

	// tell docker to hard kill the container after specified time
	// this is useful in case a test panics and suite teardown is not called to remove the container
	if err := resource.Expire(resourceExpirationTime); err != nil {
		return nil, fmt.Errorf("could not set resource expiration time: %w", err)
	}

	container.resource = resource

	return container, nil
}

// Connect returns a new database connection to the PostgreSQLcontainer.
func (c *PostgreSQLContainer) Connect() (*sqlx.DB, error) {
	var db *sqlx.DB

	if c.pool == nil || c.resource == nil {
		return nil, errors.New("container is not initialized")
	}

	databaseURI := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		pgUser,
		pgPassword,
		c.resource.GetHostPort("5432/tcp"),
		pgDBName,
	)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	// even if the PostgreSQLcontainer is up and running
	if err := c.pool.Retry(func() error {
		postgresDB, err := sqlx.Open("pgx", databaseURI)
		if err != nil {
			return err
		}

		if err := postgresDB.Ping(); err != nil {
			return fmt.Errorf("could not ping postgres database: %w", err)
		}

		db = postgresDB

		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	return db, nil
}

// Purge purges the PostgreSQLcontainer.
func (c *PostgreSQLContainer) Purge() error {
	if c.pool == nil || c.resource == nil {
		return nil
	}

	err := c.pool.Purge(c.resource)

	return err
}

// RunMigrations runs the database schema migration.
func RunMigrations(db *sql.DB) error {
	schemaMigrationsDir := "../../../migrations/schema"

	// get the current version
	oldVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}

	// run schema migrations
	err = goose.Up(db, schemaMigrationsDir)
	if err != nil {
		return err
	}

	// get the new migration version
	newVersion, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}

	log.Info("migrations successful. migrated db from version: %d to version: %d", oldVersion, newVersion)

	return nil
}
