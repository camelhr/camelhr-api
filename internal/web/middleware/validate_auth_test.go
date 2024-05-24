package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/domains/user"
	"github.com/camelhr/camelhr-api/internal/web/middleware"
	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//nolint:maintidx // all these tests belong to the single method
func TestAuthMiddleware_ValidateAuth(t *testing.T) {
	t.Parallel()

	t.Run("should validate the request with a jwt bearer token", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(appSecret, nil)
		require.NotNil(t, m)

		// generate a new jwt token
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		subdomain := gofakeit.Word()
		token, err := auth.GenerateJWT(appSecret, userID, orgID, subdomain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userIDCtx := ctx.Value(request.CtxUserIDKey)
			require.NotNil(t, userIDCtx)
			uid, ok := userIDCtx.(int64)
			assert.True(t, ok)
			assert.Equal(t, userID, uid)

			orgIDCtx := ctx.Value(request.CtxOrgIDKey)
			require.NotNil(t, orgIDCtx)
			oid, ok := orgIDCtx.(int64)
			assert.True(t, ok)
			assert.Equal(t, orgID, oid)

			orgSubdomainCtx := ctx.Value(request.CtxOrgSubdomainKey)
			require.NotNil(t, orgSubdomainCtx)
			s, ok := orgSubdomainCtx.(string)
			assert.True(t, ok)
			assert.Equal(t, subdomain, s)
		})).ServeHTTP(rr, req)

		// assert that the response status code is 200
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should validate the request with a jwt cookie", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(appSecret, nil)
		require.NotNil(t, m)

		// generate a new jwt token
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		subdomain := gofakeit.Word()
		token, err := auth.GenerateJWT(appSecret, userID, orgID, subdomain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.AddCookie(&http.Cookie{
			Name:  auth.JWTCookieName,
			Value: token,
		})

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userIDCtx := ctx.Value(request.CtxUserIDKey)
			require.NotNil(t, userIDCtx)
			uid, ok := userIDCtx.(int64)
			assert.True(t, ok)
			assert.Equal(t, userID, uid)

			orgIDCtx := ctx.Value(request.CtxOrgIDKey)
			require.NotNil(t, orgIDCtx)
			oid, ok := orgIDCtx.(int64)
			assert.True(t, ok)
			assert.Equal(t, orgID, oid)

			orgSubdomainCtx := ctx.Value(request.CtxOrgSubdomainKey)
			require.NotNil(t, orgSubdomainCtx)
			s, ok := orgSubdomainCtx.(string)
			assert.True(t, ok)
			assert.Equal(t, subdomain, s)
		})).ServeHTTP(rr, req)

		// assert that the response status code is 200
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should validate the request with a api token", func(t *testing.T) {
		t.Parallel()

		userService := user.NewMockService(t)
		subdomain := gofakeit.Word()
		apiToken := gofakeit.UUID()
		u := user.User{
			ID:             gofakeit.Int64(),
			OrganizationID: gofakeit.Int64(),
		}

		userService.On("GetUserByOrgSubdomainAPIToken", mock.Anything, subdomain, apiToken).
			Return(u, nil).Once()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware("", userService)
		require.NotNil(t, m)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.SetBasicAuth(
			apiToken,
			auth.APITokenBasicAuthPassword,
		)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userIDCtx := ctx.Value(request.CtxUserIDKey)
			require.NotNil(t, userIDCtx)
			uid, ok := userIDCtx.(int64)
			assert.True(t, ok)
			assert.Equal(t, u.ID, uid)

			orgIDCtx := ctx.Value(request.CtxOrgIDKey)
			require.NotNil(t, orgIDCtx)
			oid, ok := orgIDCtx.(int64)
			assert.True(t, ok)
			assert.Equal(t, u.OrganizationID, oid)

			orgSubdomainCtx := ctx.Value(request.CtxOrgSubdomainKey)
			require.NotNil(t, orgSubdomainCtx)
			s, ok := orgSubdomainCtx.(string)
			assert.True(t, ok)
			assert.Equal(t, subdomain, s)
		})).ServeHTTP(rr, req)

		// assert that the response status code is 200
		require.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return unauthorized response if subdomain path param is missing", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(appSecret, nil)
		require.NotNil(t, m)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Fail(t, "should not be called")
		})).ServeHTTP(rr, req)

		// assert that the response status code is 401
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.JSONEq(t, `{"error":"subdomain is required for user authentication"}`, rr.Body.String())
	})

	t.Run("should return unauthorized response for an invalid jwt token", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(appSecret, nil)
		require.NotNil(t, m)

		// create random user id, org id and org subdomain
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		orgSubdomain := gofakeit.Username()

		// generate jwt token
		token, err := auth.GenerateJWT(appSecret, userID, orgID, orgSubdomain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// modify the token
		token += "invalid"

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", "test")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Fail(t, "should not be called")
		})).ServeHTTP(rr, req)

		// assert that the response
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.JSONEq(t, `{"error":"invalid token"}`, rr.Body.String())
	})

	t.Run("should return unauthorized response if user not found for api token & subdomain", func(t *testing.T) {
		t.Parallel()

		userService := user.NewMockService(t)
		subdomain := gofakeit.Word()
		apiToken := gofakeit.UUID()

		userService.On("GetUserByOrgSubdomainAPIToken", mock.Anything, subdomain, apiToken).
			Return(user.User{}, base.NewNotFoundError("not found")).Once()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware("", userService)
		require.NotNil(t, m)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.SetBasicAuth(
			apiToken,
			auth.APITokenBasicAuthPassword,
		)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Fail(t, "should not be called")
		})).ServeHTTP(rr, req)

		// assert that the response
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.JSONEq(t, `{"error":"invalid api token"}`, rr.Body.String())
	})

	t.Run("should return unauthorized response for a api token of disabled user", func(t *testing.T) {
		t.Parallel()

		userService := user.NewMockService(t)
		subdomain := gofakeit.Word()
		apiToken := gofakeit.UUID()
		now := time.Now()
		u := user.User{
			ID:         gofakeit.Int64(),
			DisabledAt: &now,
		}

		userService.On("GetUserByOrgSubdomainAPIToken", mock.Anything, subdomain, apiToken).
			Return(u, nil).Once()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware("", userService)
		require.NotNil(t, m)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.SetBasicAuth(
			apiToken,
			auth.APITokenBasicAuthPassword,
		)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Fail(t, "should not be called")
		})).ServeHTTP(rr, req)

		// assert that the response
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.JSONEq(t, `{"error":"user is disabled"}`, rr.Body.String())
	})

	t.Run("should return unauthorized response for a request without auth details", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		appSecret := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(appSecret, nil)
		require.NotNil(t, m)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", "test")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Fail(t, "should not be called")
		})).ServeHTTP(rr, req)

		// assert that the response status code is 401
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		require.Empty(t, rr.Body.String())
	})
}
