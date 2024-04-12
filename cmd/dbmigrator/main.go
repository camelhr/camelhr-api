package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/camelhr/camelhr-api/migrations"
	"github.com/camelhr/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	driver = "postgres"

	usagePrefix = `Usage: dbmigrator [OPTIONS] COMMAND

You can also provide the database connection string as an environment variable named DB_CONN.
If both the -db_conn flag and the DB_CONN environment variable are provided, the -db_conn flag will take precedence.

Examples:
  dbmigrator -dir=./migrations -db_conn=postgres://user:password@localhost:5432/dbname?sslmode=disable up
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
	flagSet                     = flag.NewFlagSet("dbmigrator", flag.ExitOnError)
	dir                         = flagSet.String("dir", "migrations", "directory with migration files")
	dbConn                      = flagSet.String("db_conn", "", "database connection string")
	allowedCommands             = []string{"up", "up-by-one", "up-to", "down", "down-to", "redo", "status", "version", "create"}
	dbConnectionNotNeedCommands = []string{"create"}
)

func main() {
	ctx := context.Background()
	log.InitGlobalLogger("dbmigrator", "info")
	flagSet.Usage = usage      // override the default usage function
	flagSet.Parse(os.Args[1:]) // nolint:errcheck

	log.Info("migration started")

	// if db connection string is provided as a flag then use it
	// otherwise, get it from the environment variable
	if dbConn == nil || *dbConn == "" {
		envDBConn := os.Getenv("DB_CONN")
		dbConn = &envDBConn
	}

	args := flagSet.Args()
	command := extractCommand(args)
	validateDBConnectionString(*dbConn, command)

	db, err := goose.OpenDBWithDriver(driver, *dbConn)
	if err != nil {
		log.Fatal("failed to connect db: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("failed to close db: %v", err)
		}
	}()

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.RunContext(ctx, command, db, *dir, arguments...); err != nil {
		log.Fatal("failed to execute migration with goose %v: %v", command, err)
	}

	log.Info("migration completed successfully")
}

func usage() {
	fmt.Print(usagePrefix)
	flagSet.PrintDefaults()
	fmt.Print(usageCommands)
}

// extractCommand extracts the command from the arguments
// validates it against the allowed commands and returns the command
// it exits the program with an error message if the command is invalid
func extractCommand(args []string) string {
	if len(args) < 1 {
		flagSet.Usage()
		log.Fatal("command argument is required")
	}

	command := args[0]
	if !contains(allowedCommands, strings.ToLower(command)) {
		flagSet.Usage()
		log.Fatal("invalid command: %s", command)
	}

	return command
}

// validateDBConnectionString ensures that the db connection string is not empty for the commands that require it
// it exits the program with an error message if db connection string is missing for the required commands
func validateDBConnectionString(dbConn string, command string) {
	if dbConn == "" {
		// if the command does not require a db connection string then return
		if contains(dbConnectionNotNeedCommands, command) {
			return
		}
		flagSet.Usage()
		log.Fatal("db_conn flag or DB_CONN environment variable is required for '%s' command", command)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
