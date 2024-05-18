package middleware_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/web/middleware"
	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware_ValidateAuth(t *testing.T) {
	t.Parallel()

	t.Run("should validate the request with a jwt bearer token", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		secretKey := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(secretKey)
		require.NotNil(t, m)

		// generate a new jwt token
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		subdomain := gofakeit.Word()
		token, err := auth.GenerateJWT(secretKey, userID, orgID, subdomain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.Header.Set("Authorization", "Bearer "+token)

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
		secretKey := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(secretKey)
		require.NotNil(t, m)

		// generate a new jwt token
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		subdomain := gofakeit.Word()
		token, err := auth.GenerateJWT(secretKey, userID, orgID, subdomain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.AddCookie(&http.Cookie{
			Name:  "jwt_session_id",
			Value: base64.StdEncoding.EncodeToString([]byte(token)),
		})

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

	t.Run("should return unauthorized response for an invalid jwt token", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		secretKey := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(secretKey)
		require.NotNil(t, m)

		// create random user id, org id and org subdomain
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		orgSubdomain := gofakeit.Username()

		// generate jwt token
		token, err := auth.GenerateJWT(secretKey, userID, orgID, orgSubdomain)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		// modify the token
		token += "invalid"

		// create a new request with jwt bearer token
		req := httptest.NewRequest(http.MethodGet, "/api/some-endpoint", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		// create a new response recorder
		rr := httptest.NewRecorder()

		m.ValidateAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Fail(t, "should not be called")
		})).ServeHTTP(rr, req)

		// assert that the response status code is 401
		require.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return unauthorized response for a request without auth details", func(t *testing.T) {
		t.Parallel()

		// create a random secret key
		secretKey := gofakeit.UUID()

		// create a new auth middleware
		m := middleware.NewAuthMiddleware(secretKey)
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
	})
}
