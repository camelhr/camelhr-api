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
	if dbConn == nil || *dbConn == "" {
		gooseFlagSet.Usage()
		log.Fatal("db_conn flag is mandatory")
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
