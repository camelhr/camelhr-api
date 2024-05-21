package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/camelhr/camelhr-api/internal/web/response"
)

type authMiddleware struct {
	secretKey string
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(secretKey string) *authMiddleware {
	return &authMiddleware{secretKey: secretKey}
}

// ValidateAuth is a middleware that authenticates the request.
func (m *authMiddleware) ValidateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// try to get jwt from authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			processJWT(next, w, r, strings.TrimPrefix(authHeader, "Bearer "), m.secretKey)
			return
		}

		// try to get jwt from cookie
		if cookie, err := r.Cookie("jwt_session_id"); err == nil {
			decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
			if err == nil && len(decoded) > 0 {
				processJWT(next, w, r, string(decoded), m.secretKey)
				return
			}
		}

		response.Empty(w, http.StatusUnauthorized)
	})
}

// processJWT parses and validates the jwt token.
// If the token is valid, it sets the user-id, org-id and org-subdomain in the request context.
func processJWT(next http.Handler, w http.ResponseWriter, r *http.Request, jwtString, secretKey string) {
	token, claims, err := auth.ParseAndValidateJWT(jwtString, secretKey)
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

	// validate the subdomain(if available) from the request path
	if request.URLParam(r, "subdomain") != "" && request.URLParam(r, "subdomain") != claims.OrgSubdomain {
		response.ErrorResponse(w, http.StatusUnauthorized, base.NewAPIError("user doesn't belong to the organization"))
		return
	}

	next.ServeHTTP(w, r.WithContext(ctx))
}
