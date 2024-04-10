package web

import (
	"net/http"

	"github.com/camelhr/camelhr-api/internal/domains/organization"
	customMiddleware "github.com/camelhr/camelhr-api/internal/web/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() http.Handler {
	// initialize dependencies
	orgHandler := organization.NewOrganizationHandler()

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

	// protected routes. auth required
	v1.Group(func(r chi.Router) {
		// TODO: add auth middleware

		r.Route("/organizations", func(r chi.Router) {
			r.Post("/", orgHandler.CreateOrganization)
			r.Route("/{orgID:[0-9]+}", func(r chi.Router) {
				r.Get("/", orgHandler.GetOrganization)
				r.Put("/", orgHandler.UpdateOrganization)
				r.Delete("/", orgHandler.DeleteOrganization)
			})
		})
	})

	return r
}
