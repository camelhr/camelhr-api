package main

import (
	"context"
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

	// setup routes and start the server
	handler := web.SetupRoutes(db)
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
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

func connectToDatabase(configs config.Config) (database.Database, error) {
	db, err := sqlx.Open("pgx", configs.DBConn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(configs.DBMaxOpen)
	db.SetMaxIdleConns(configs.DBMaxIdle)
	db.SetConnMaxIdleTime(time.Duration(configs.DBMaxIdleConnTime) * time.Minute)
	db.SetConnMaxLifetime(0)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return database.NewPostgresDatabase(db), nil
}

func trapTerminateSignal(onTerminate func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	onTerminate()
}
