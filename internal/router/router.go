package router

import (
	"net/http"

	customMiddleware "github.com/camelhr/camelhr-api/internal/router/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(customMiddleware.ChiRequestLoggerMiddleware()) // <--<< Logger should come before Recoverer
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	return r
}
