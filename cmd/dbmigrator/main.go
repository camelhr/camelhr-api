package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	_ "github.com/camelhr/camelhr-api/migrations/datafix"
	_ "github.com/camelhr/camelhr-api/migrations/schema"
	"github.com/camelhr/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	driver     = "postgres"
	schemaDir  = "migrations/schema"
	datafixDir = "migrations/datafix"

	usagePrefix = `Usage: dbmigrator [OPTIONS] COMMAND

You can also provide the database connection string as an environment variable named DB_CONN.
If both the -db_conn flag and the DB_CONN environment variable are provided, the -db_conn flag will take precedence.

Examples:
  dbmigrator -db_conn=postgres://user:password@localhost:5432/dbname?sslmode=disable up
  DB_CONN=postgres://user:password@localhost:5432/dbname?sslmode=disable dbmigrator up

Options:
`

	usageCommands = `
Commands:
  up                   Migrate the DB to the most recent version available
  up-by-one            Migrate the DB up by 1
  up-to VERSION        Migrate the DB to a specific VERSION
  down                 Roll back the version by 1
  down-to VERSION      Roll back to a specific VERSION
  redo                 Re-run the latest migration (down, then up)
  status               Dump the migration status for the current DB
  version              Print the current version of the database
  create NAME [sql|go] Creates new migration file with the current timestamp
`
)

var (
	// ErrNoDBConn is returned when the db connection string is missing.
	ErrNoDBConn = errors.New("db connection string is missing")

	// ErrUnknownCommand is returned when the command is unknown.
	ErrUnknownCommand = errors.New("unknown command")

	// ErrCommandArgMissing is returned when the command argument is missing.
	ErrCommandArgMissing = errors.New("command argument is missing")

	// ErrDBConnStringMissing is returned when the db connection string is missing.
	ErrDBConnStringMissing = errors.New("db connection string is missing")

	// ErrCreateCommandArgsMissing is returned when the create command arguments are missing.
	ErrCreateCommandArgsMissing = errors.New("create command arguments are missing")
)

func main() {
	flagSet := prepareFlags()

	// if db connection string is provided as a flag then use it. otherwise, use it from the environment variable
	dbConn := flagSet.String("db_conn", os.Getenv("DB_CONN"), "database connection string")
	operationType := flagSet.String("type", "all", "migration type (schema|datafix|all)")

	allowedCommands := []string{"up", "up-by-one", "up-to", "down", "down-to", "redo", "status", "version", "create"}
	dbIndependentCommands := []string{"create"}
	dirIndependentCommands := []string{"version"}

	flagSet.Parse(os.Args[1:]) //nolint:errcheck // flag.ExitOnError will exit on error
	log.InitGlobalLogger("dbmigrator", "info")
	log.Info("migration started")

	args := flagSet.Args()

	command, err := extractCommand(args, allowedCommands)
	if err != nil {
		flagSet.Usage()
		log.Fatal("failed to extract command: %v", err)
	}

	if err := validateDBConnectionString(*dbConn, command, dbIndependentCommands); err != nil {
		flagSet.Usage()
		log.Fatal("failed to validate db connection string: %v", err)
	}

	db, err := goose.OpenDBWithDriver(driver, *dbConn)
	if err != nil {
		log.Fatal("failed to connect db: %v", err)
	}
	defer db.Close()

	arguments := args[1:]

	if contains(dirIndependentCommands, command) {
		if err := runMigration(context.Background(), command, db, "", arguments); err != nil {
			log.Error("failed to run migration: %v", err)
		}

		log.Info("migration completed")

		return
	}

	if shouldMigrateSchema(*operationType) {
		if err := runMigration(context.Background(), command, db, schemaDir, arguments); err != nil {
			log.Error("failed to run schema migration: %v", err)
			return
		}
	}

	if shouldMigrateDatafix(*operationType) {
		if err := runMigration(context.Background(), command, db, datafixDir, arguments); err != nil {
			log.Error("failed to run datafix migration: %v", err)
			return
		}
	}

	log.Info("migration completed")
}

func shouldMigrateSchema(migrationType string) bool {
	return migrationType == "schema" || migrationType == "all"
}

func shouldMigrateDatafix(migrationType string) bool {
	return migrationType == "datafix" || migrationType == "all"
}

// runMigration runs the migration command with the provided arguments.
func runMigration(ctx context.Context, command string, db *sql.DB, dir string, args []string) error {
	if command == "create" {
		const minArgs = 2
		if len(args) < minArgs {
			return ErrCreateCommandArgsMissing
		}

		migrationType := args[1] // go or sql
		migrationDir := filepath.Base(dir)

		if migrationType == "go" {
			return goose.CreateWithTemplate(db, dir, goMigrationTemplate(migrationDir), args[0], migrationType)
		}

		return goose.Create(db, dir, args[0], migrationType)
	}

	// run sql migrations from the provided directory
	// please note that the path to the migrations directory (schemaDir) is used only for SQL migrations
	// the Go migrations are found via the Go import path, not the file system path
	if err := goose.RunWithOptionsContext(ctx, command, db, dir, args, goose.WithAllowMissing()); err != nil {
		if !errors.Is(err, goose.ErrNoNextVersion) && !errors.Is(err, goose.ErrNoCurrentVersion) &&
			!errors.Is(err, goose.ErrNoMigrationFiles) {
			return fmt.Errorf("failed to execute migration with goose %v: %w", command, err)
		}
	}

	return nil
}

// prepareFlags prepares the flag set with the required flags and usage.
func prepareFlags() *flag.FlagSet {
	flagSet := flag.NewFlagSet("dbmigrator", flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Print(usagePrefix) //nolint:forbidigo // structured logger is not needed here
		flagSet.PrintDefaults()
		fmt.Print(usageCommands) //nolint:forbidigo // structured logger is not needed here
	}

	return flagSet
}

// extractCommand extracts the command from the arguments,
// validates it against the allowed commands and returns the command.
// It returns an error if the command is invalid.
func extractCommand(args []string, allowedCommands []string) (string, error) {
	if len(args) < 1 {
		return "", ErrCommandArgMissing
	}

	command := args[0]
	if !contains(allowedCommands, strings.ToLower(command)) {
		return "", fmt.Errorf("command: %s not found: %w", command, ErrUnknownCommand)
	}

	return command, nil
}

// validateDBConnectionString ensures that the db connection string is not empty for the commands that require it.
// It returns error if db connection string is missing for the required commands.
func validateDBConnectionString(dbConn string, command string, dbIndependentCommands []string) error {
	if dbConn == "" {
		// if the command does not require a db connection string then return
		if contains(dbIndependentCommands, command) {
			return nil
		}

		return fmt.Errorf("db_conn flag or DB_CONN env variable is required for '%s' command: %w",
			command, ErrDBConnStringMissing)
	}

	return nil
}

func goMigrationTemplate(pkg string) *template.Template {
	tpl := template.Must(template.New("goose.go-migration").Parse(fmt.Sprintf(`package %s

	import (
		"database/sql"
		"github.com/pressly/goose/v3"
	)
	
	func init() {
		goose.AddMigration(up{{.CamelName}}, down{{.CamelName}})
	}
	
	func up{{.CamelName}}(tx *sql.Tx) error {
		// This code is executed when the migration is applied.
		return nil
	}
	
	func down{{.CamelName}}(tx *sql.Tx) error {
		// This code is executed when the migration is rolled back.
		return nil
	}
	`, pkg)))

	return tpl
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
