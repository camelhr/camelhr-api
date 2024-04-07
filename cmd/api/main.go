package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/camelhr/camelhr-api/internal/config"
	"github.com/camelhr/camelhr-api/internal/router"
	"github.com/camelhr/log"
)

func main() {
	config := config.LoadConfig()
	log.InitGlobalLogger("api", config.LogLevel)
	log.Info("config loaded successfully")

	handler := router.SetupRoutes()
	server := &http.Server{
		Addr:              config.HTTPAddress,
		Handler:           handler,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	// start the server in separate go routine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start server: %v", err)
		}
	}()

	// trap terminate signals and gracefully shutdown the server
	trapTerminateSignal(func() {
		stopCtx, stopCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer stopCancel()

		if err := server.Shutdown(stopCtx); err != nil {
			log.Error("failed to gracefully shutdown server: %v", err)
		}
	})
}

func trapTerminateSignal(onTerminate func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	onTerminate()
}
