package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/camelhr/camelhr-api/internal/web/response"
)

type authMiddleware struct {
	appSecret   string
	userService user.Service
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(appSecret string, userService user.Service) *authMiddleware {
	return &authMiddleware{appSecret: appSecret, userService: userService}
}

// ValidateAuth is a middleware that authenticates the request.
// Before using this middleware, make sure that associated endpoint contains the subdomain path parameter.
func (m *authMiddleware) ValidateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// subdomain is required for user authentication
		if request.URLParam(r, "subdomain") == "" {
			response.ErrorResponse(w, http.StatusUnauthorized,
				base.NewAPIError("subdomain is required for user authentication"))

			return
		}

		// try to get jwt from bearer authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			m.processJWT(next, w, r, strings.TrimPrefix(authHeader, "Bearer "))
			return
		}

		// try to get api token from basic auth header
		if strings.HasPrefix(authHeader, "Basic ") {
			m.processAPIToken(next, w, r, strings.TrimPrefix(authHeader, "Basic "))
			return
		}

		// try to get jwt from cookie
		if cookie, err := r.Cookie(auth.JWTCookieName); err == nil {
			if cookie.Value != "" {
				m.processJWT(next, w, r, cookie.Value)
				return
			}
		}

		response.Empty(w, http.StatusUnauthorized)
	})
}

// processJWT parses and validates the jwt token.
// If the token is valid, it sets the user-id, org-id and org-subdomain in the request context.
func (m *authMiddleware) processJWT(next http.Handler, w http.ResponseWriter, r *http.Request, jwtString string) {
	token, claims, err := auth.ParseAndValidateJWT(jwtString, m.appSecret)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, base.NewAPIError("invalid token", base.ErrorCause(err)))
		return
	}

	if token == nil || !token.Valid || claims == nil {
		response.ErrorResponse(w, http.StatusUnauthorized, base.NewAPIError("invalid token"))
		return
	}

	// set claims values to request context
	ctx := context.WithValue(r.Context(), request.CtxUserIDKey, claims.UserID)
	ctx = context.WithValue(ctx, request.CtxOrgIDKey, claims.OrgID)
	ctx = context.WithValue(ctx, request.CtxOrgSubdomainKey, claims.OrgSubdomain)

	// validate the subdomain from the request path against the subdomain in the jwt claims
	if request.URLParam(r, "subdomain") != claims.OrgSubdomain {
		response.ErrorResponse(w, http.StatusUnauthorized, base.NewAPIError("user doesn't belong to the organization"))
		return
	}

	next.ServeHTTP(w, r.WithContext(ctx))
}

// processAPIToken validates the api token from the basic auth header.
// If the token is valid, it sets the user-id, org-id and org-subdomain in the request context.
func (m *authMiddleware) processAPIToken(
	next http.Handler,
	w http.ResponseWriter,
	r *http.Request,
	basicAuthCred string,
) {
	// decode the basic auth header
	decoded, err := base64.StdEncoding.DecodeString(basicAuthCred)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized,
			base.NewAPIError("invalid basic auth header", base.ErrorCause(err)))

		return
	}

	// extract the api token
	apiToken, found := strings.CutSuffix(string(decoded), ":"+auth.APITokenPassword)
	if !found {
		response.ErrorResponse(w, http.StatusUnauthorized,
			base.NewAPIError(fmt.Sprintf("invalid basic auth header. expected %s as password", auth.APITokenPassword)))

		return
	}

	subdomain := request.URLParam(r, "subdomain")

	// get user for the given subdomain and api token
	u, err := m.userService.GetUserByOrgSubdomainAPIToken(r.Context(), subdomain, apiToken)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized,
			base.NewAPIError("invalid api token", base.ErrorCause(err)))

		return
	}

	if u.DisabledAt != nil {
		response.ErrorResponse(w, http.StatusUnauthorized,
			base.NewAPIError("user is disabled"))

		return
	}

	// set user-id, org-id and org-subdomain in the request context
	ctx := context.WithValue(r.Context(), request.CtxUserIDKey, u.ID)
	ctx = context.WithValue(ctx, request.CtxOrgIDKey, u.OrganizationID)
	ctx = context.WithValue(ctx, request.CtxOrgSubdomainKey, subdomain)

	next.ServeHTTP(w, r.WithContext(ctx))
}
