package router

import (
	"net/http"

	customMiddleware "github.com/camelhr/camelhr-api/internal/router/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() http.Handler {
	// create a default router
	r := chi.NewRouter()

	// add middlewares
	r.Use(middleware.RequestID)
	r.Use(customMiddleware.ChiRequestLoggerMiddleware()) // <--<< logger should come before recoverer
	r.Use(middleware.Recoverer)

	// create a sub-router for v1 api endpoints
	v1 := chi.NewRouter()
	r.Mount("/api/v1", v1)

	// open routes. no auth required
	v1.Group(func(r chi.Router) {
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})
	})

	return r
}
