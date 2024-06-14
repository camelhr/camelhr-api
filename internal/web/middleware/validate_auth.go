package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/domains/session"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/camelhr/camelhr-api/internal/web/response"
)

type authMiddleware struct {
	appSecret      string
	userService    user.Service
	sessionManager session.SessionManager
}

// NewAuthMiddleware creates a new auth middleware.
func NewAuthMiddleware(
	appSecret string,
	userService user.Service,
	sessionManager session.SessionManager,
) *authMiddleware {
	return &authMiddleware{appSecret, userService, sessionManager}
}

// ValidateAuth is a middleware that authenticates the request.
// Before using this middleware, make sure that associated endpoint contains the subdomain path parameter.
func (m *authMiddleware) ValidateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// subdomain is required for user authentication
		if request.URLParam(r, "subdomain") == "" {
			response.ErrorResponse(w, base.NewAPIError("subdomain is required for user authentication",
				base.ErrorHTTPStatus(http.StatusUnauthorized)))

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
// It then ensures that the token is present in the session.
// If the token is valid, it sets the user-id, org-id and org-subdomain in the request context.
func (m *authMiddleware) processJWT(next http.Handler, w http.ResponseWriter, r *http.Request, jwtString string) {
	token, claims, err := auth.ParseAndValidateJWT(jwtString, m.appSecret)
	if err != nil {
		response.ErrorResponse(w, base.NewAPIError("invalid token", base.ErrorCause(err),
			base.ErrorHTTPStatus(http.StatusUnauthorized)))

		return
	}

	if token == nil || !token.Valid || claims == nil {
		response.ErrorResponse(w, base.NewAPIError("invalid token", base.ErrorHTTPStatus(http.StatusUnauthorized)))
		return
	}

	// only one most recent jwt is stored in the session upon login
	// validate the session to ensure that the same jwt is present in the session
	if err := m.sessionManager.ValidateJWTSession(r.Context(), claims.UserID, claims.OrgID, jwtString); err != nil {
		response.ErrorResponse(w, base.NewAPIError("invalid token", base.ErrorCause(err),
			base.ErrorHTTPStatus(http.StatusUnauthorized)))

		return
	}

	// set claims values to request context
	ctx := context.WithValue(r.Context(), request.CtxUserIDKey, claims.UserID)
	ctx = context.WithValue(ctx, request.CtxOrgIDKey, claims.OrgID)
	ctx = context.WithValue(ctx, request.CtxOrgSubdomainKey, claims.OrgSubdomain)

	// validate the subdomain from the request path against the subdomain in the jwt claims
	if request.URLParam(r, "subdomain") != claims.OrgSubdomain {
		response.ErrorResponse(w, base.NewAPIError("user doesn't belong to the organization",
			base.ErrorHTTPStatus(http.StatusUnauthorized)))

		return
	}

	next.ServeHTTP(w, r.WithContext(ctx))
}

// processAPIToken validates the api token from the basic auth header.
// It first checks the api-token in the session.
// If the token is not present in the session, it queries the database to get the user.
// It sets the user-id, org-id and org-subdomain in the request context.
func (m *authMiddleware) processAPIToken(
	next http.Handler,
	w http.ResponseWriter,
	r *http.Request,
	basicAuthCred string,
) {
	// decode the basic auth header
	decoded, err := base64.StdEncoding.DecodeString(basicAuthCred)
	if err != nil {
		response.ErrorResponse(w, base.NewAPIError("invalid basic auth header", base.ErrorCause(err),
			base.ErrorHTTPStatus(http.StatusUnauthorized)))

		return
	}

	// extract the api token
	apiToken, found := strings.CutSuffix(string(decoded), ":"+auth.APITokenBasicAuthPassword)
	if !found {
		response.ErrorResponse(w, base.NewAPIError(fmt.Sprintf("invalid basic auth header. expected %s as password",
			auth.APITokenBasicAuthPassword), base.ErrorHTTPStatus(http.StatusUnauthorized)))

		return
	}

	subdomain := request.URLParam(r, "subdomain")

	userID, orgID, err := m.getUserForAPIToken(r.Context(), apiToken, subdomain)
	if err != nil {
		response.ErrorResponse(w, err)
		return
	}

	// set user-id, org-id and org-subdomain in the request context
	ctx := context.WithValue(r.Context(), request.CtxUserIDKey, userID)
	ctx = context.WithValue(ctx, request.CtxOrgIDKey, orgID)
	ctx = context.WithValue(ctx, request.CtxOrgSubdomainKey, subdomain)

	next.ServeHTTP(w, r.WithContext(ctx))
}

func (m *authMiddleware) getUserForAPIToken(ctx context.Context, apiToken, subdomain string) (int64, int64, error) {
	// check if the api token is present in the session and get associated ids
	if userID, orgID, err := m.sessionManager.ValidateAPITokenSession(ctx, apiToken); err == nil {
		return userID, orgID, nil
	}

	// if not, get user for the given subdomain and api token from the database
	u, err := m.userService.GetUserByOrgSubdomainAPIToken(ctx, subdomain, apiToken)
	if err != nil {
		if base.IsNotFoundError(err) {
			return 0, 0, base.NewAPIError("invalid api token", base.ErrorCause(err),
				base.ErrorHTTPStatus(http.StatusUnauthorized))
		}

		return 0, 0, err
	}

	// ensure that the user is not disabled
	if u.DisabledAt != nil {
		return 0, 0, base.WrapError(auth.ErrUserDisabled, base.ErrorHTTPStatus(http.StatusUnauthorized))
	}

	// set the api token in the session
	if err := m.sessionManager.CreateSession(ctx, u.ID, u.OrganizationID, "",
		apiToken, auth.SessionTTLDuration); err != nil {
		return 0, 0, err
	}

	return u.ID, u.OrganizationID, nil
}
