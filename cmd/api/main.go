package main

import (
	"net/http"

	"github.com/camelhr/camelhr-api/internal/router/middlewares"
	"github.com/camelhr/log"
	"github.com/go-chi/chi/v5"
)

func main() {
	log.InitGlobalLogger("api", "info")
	r := chi.NewRouter()
	r.Use(middlewares.ChiRequestLoggerMiddleware())
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":8080", r)
}
