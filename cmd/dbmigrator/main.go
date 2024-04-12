package main

import (
	"context"
	"flag"
	"os"

	_ "github.com/camelhr/camelhr-api/migrations"
	"github.com/camelhr/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var (
	gooseFlagSet = flag.NewFlagSet("goose", flag.ExitOnError)
	dir          = gooseFlagSet.String("dir", "migrations", "directory with migration files")
	dbConn       = gooseFlagSet.String("db_conn", "", "database connection string")
)

func main() {
	log.InitGlobalLogger("goose", "info")

	gooseFlagSet.Parse(os.Args[1:])

	// if db connection string is provided as a flag then use it
	// otherwise, check if it is provided as an environment variable
	// if it is not provided in either of the ways then exit the program
	if dbConn == nil || *dbConn == "" {
		envDBConn := os.Getenv("DB_CONN")
		if envDBConn == "" {
			log.Fatal("db connection string is required")
		}
		dbConn = &envDBConn
	}
	args := gooseFlagSet.Args()

	if len(args) < 1 {
		gooseFlagSet.Usage()
		return
	}
	command := args[0]

	db, err := goose.OpenDBWithDriver("postgres", *dbConn)
	if err != nil {
		log.Fatal("failed to connect db: %v\n", err)
	}
	log.Info("connected to db")

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("failed to close db: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.RunContext(context.Background(), command, db, *dir, arguments...); err != nil {
		log.Fatal("failed to execute migration with goose %v: %v", command, err)
	}

	log.Info("migration completed successfully")
}
