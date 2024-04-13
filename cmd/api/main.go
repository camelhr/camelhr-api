package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/camelhr/camelhr-api/internal/config"
	"github.com/camelhr/camelhr-api/internal/web"
	"github.com/camelhr/log"
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

	handler := web.SetupRoutes()
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

func trapTerminateSignal(onTerminate func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	onTerminate()
}
