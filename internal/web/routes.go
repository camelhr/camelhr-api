package web

import (
	"net/http"

	"github.com/camelhr/camelhr-api/internal/config"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/web/middleware"
	"github.com/camelhr/camelhr-api/internal/web/response"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// SetupRoutes initializes the routes for the web server.
func SetupRoutes(db database.Database, conf config.Config) http.Handler {
	// initialize dependencies
	orgRepo := organization.NewRepository(db)
	orgService := organization.NewService(orgRepo)
	orgHandler := organization.NewHandler(orgService)
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	authMiddleware := middleware.NewAuthMiddleware(conf.AppSecret, userService)

	// create a default router
	r := chi.NewRouter()

	// add middlewares
	r.Use(chimiddleware.RequestID)
	r.Use(middleware.ChiRequestLoggerMiddleware()) // <--<< logger should come before recoverer
	r.Use(chimiddleware.Recoverer)

	// create a sub-router for v1 api endpoints
	v1 := chi.NewRouter()
	r.Mount("/api/v1", v1)

	// v1 open routes. no auth required
	v1.Group(func(r chi.Router) {
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			response.Text(w, http.StatusOK, "OK")
		})
	})

	// create a sub-router for v1 subdomain endpoints
	v1Subdomain := chi.NewRouter()
	v1.Mount("/subdomains/{subdomain}", v1Subdomain)

	v1Subdomain.Route("/organizations", func(r chi.Router) {
		// open routes. no auth required
		r.Get("/", orgHandler.GetOrganizationBySubdomain)

		// protected routes. auth required
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.ValidateAuth)

			r.Put("/", orgHandler.UpdateOrganization)
			r.Delete("/", orgHandler.DeleteOrganization)
		})
	})

	return r
}
