package web

import (
	"net/http"
	"strings"

	"github.com/camelhr/camelhr-api/internal/config"
	"github.com/camelhr/camelhr-api/internal/database"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/web/middleware"
	"github.com/camelhr/camelhr-api/internal/web/response"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
)

// SetupRoutes initializes the routes for the web server.
//
//nolint:funlen // ignore function length since this is a setup function
func SetupRoutes(db database.Database, redisClient *redis.Client, conf config.Config) http.Handler {
	// initialize dependencies
	sessionManager := session.NewRedisSessionManager(redisClient)
	orgRepo := organization.NewRepository(db)
	orgService := organization.NewService(orgRepo, sessionManager)
	orgHandler := organization.NewHandler(orgService)
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, sessionManager)
	authService := auth.NewService(conf.AppSecret, db, orgService, userService, sessionManager)
	authHandler := auth.NewHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(conf.AppSecret, userService, sessionManager)

	// create a default router
	r := chi.NewRouter()

	// add middlewares
	r.Use(cors.Handler(corsOptions()))
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

		r.Post("/auth/register", authHandler.Register)
	})

	// create a sub-router for v1 subdomain endpoints
	v1Subdomain := chi.NewRouter()
	v1.Mount("/subdomains/{subdomain}", v1Subdomain)

	v1Subdomain.Route("/auth", func(r chi.Router) {
		// open routes. no auth required
		r.Post("/login", authHandler.Login)

		// protected routes. auth required
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.ValidateAuth)

			r.Post("/logout", authHandler.Logout)
		})
	})

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

func corsOptions() cors.Options {
	const corsMaxAgeInSeconds = 300 // 5 minutes

	return cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			// allow exact match
			if origin == "https://camelhr.com" {
				return true
			}
			// allow any subdomain of camelhr.com
			if strings.HasSuffix(origin, ".camelhr.com") {
				return true
			}

			return false
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           corsMaxAgeInSeconds,
	}
}
