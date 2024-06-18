package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/camelhr/camelhr-api/internal/config"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/web"
	"github.com/camelhr/log"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	serverShutdownTimeout   = 30 * time.Second
	serverReadTimeout       = 15 * time.Second
	serverWriteTimeout      = 15 * time.Second
	serverReadHeaderTimeout = 15 * time.Second
	serverMaxHeaderBytes    = 1 << 20 // 1 MB
)

func main() {
	configs := config.LoadConfig()
	log.InitGlobalLogger("api", configs.LogLevel)
	log.Debug("debug logging enabled") // printed only if log level is set to debug
	log.Info("config loaded successfully")

	// connect to the database, set configurations and check if the connection is successful
	db, err := connectToDatabase(configs)
	if err != nil {
		log.Fatal("failed to connect to database: %v", err)
	}
	defer db.Close()

	// create a new instance of the postgresDatabase
	pgDB := database.NewPostgresDatabase(db)

	// connect to redis
	redisClient, err := connectToRedis(configs)
	if err != nil {
		log.Error("failed to connect to redis: %v", err)
		return
	}
	defer redisClient.Close()

	// setup routes and start the server
	handler := web.SetupRoutes(pgDB, redisClient, configs)
	server := &http.Server{
		Addr:              configs.HTTPAddress,
		Handler:           handler,
		ReadTimeout:       serverReadTimeout,
		WriteTimeout:      serverWriteTimeout,
		ReadHeaderTimeout: serverReadHeaderTimeout,
		MaxHeaderBytes:    serverMaxHeaderBytes,
	}

	// start the server in separate go routine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("failed to start server: %v", err)
		}
	}()

	// trap terminate signals and gracefully shutdown the server
	trapTerminateSignal(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
		defer stopCancel()

		if err := server.Shutdown(stopCtx); err != nil {
			log.Error("failed to gracefully shutdown server: %v", err)
		}
	})
}

func connectToDatabase(c config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", c.DBConn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(c.DBMaxOpen)
	db.SetMaxIdleConns(c.DBMaxIdle)
	db.SetConnMaxIdleTime(time.Duration(c.DBMaxIdleConnTime) * time.Minute)
	db.SetConnMaxLifetime(0)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func connectToRedis(c config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(c.RedisConn)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func trapTerminateSignal(onTerminate func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	onTerminate()
}
