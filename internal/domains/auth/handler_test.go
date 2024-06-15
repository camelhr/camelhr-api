package auth_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/auth"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	registerPath = "/api/v1/auth/register"
	loginPath    = "/api/v1/subdomains/{subdomain}/auth/login"
	logoutPath   = "/api/v1/subdomains/{subdomain}/auth/logout"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		email := "invalid email"
		password := validPassword
		subdomain := gofakeit.LetterN(30)
		orgName := gofakeit.Company()
		payload := fmt.Sprintf(`{"email": "%s","password":"%s","organization_subdomain":"%s","organization_name":"%s"}`,
			email, password, subdomain, orgName,
		)
		req, err := http.NewRequest(http.MethodPost, registerPath, strings.NewReader(payload))
		require.NoError(t, err)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Register(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"email must be a valid email address"}`, rr.Body.String())
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := "@2nR"
		subdomain := gofakeit.LetterN(30)
		orgName := gofakeit.Company()
		payload := fmt.Sprintf(`{"email": "%s","password":"%s","organization_subdomain":"%s","organization_name":"%s"}`,
			email, password, subdomain, orgName,
		)
		req, err := http.NewRequest(http.MethodPost, registerPath, strings.NewReader(payload))
		require.NoError(t, err)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Register(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"password must be at least 8 characters in length"}`, rr.Body.String())
	})

	t.Run("should return error when subdomain is invalid", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := validPassword
		subdomain := "invalid subdomain"
		orgName := gofakeit.Company()
		payload := fmt.Sprintf(`{"email": "%s","password":"%s","organization_subdomain":"%s","organization_name":"%s"}`,
			email, password, subdomain, orgName,
		)
		req, err := http.NewRequest(http.MethodPost, registerPath, strings.NewReader(payload))
		require.NoError(t, err)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Register(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"organization_subdomain can only contain alphanumeric characters"}`, rr.Body.String())
	})

	t.Run("should return error when organization name is invalid", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.LetterN(30)
		orgName := gofakeit.LetterN(61)
		payload := fmt.Sprintf(`{"email": "%s","password":"%s","organization_subdomain":"%s","organization_name":"%s"}`,
			email, password, subdomain, orgName,
		)
		req, err := http.NewRequest(http.MethodPost, registerPath, strings.NewReader(payload))
		require.NoError(t, err)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Register(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"organization_name must be a maximum of 60 characters in length"}`, rr.Body.String())
	})

	t.Run("should return error when service call fails", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.LetterN(30)
		orgName := gofakeit.Company()
		payload := fmt.Sprintf(`{"email": "%s","password":"%s","organization_subdomain":"%s","organization_name":"%s"}`,
			email, password, subdomain, orgName,
		)
		req, err := http.NewRequest(http.MethodPost, registerPath, strings.NewReader(payload))
		require.NoError(t, err)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Register", fake.MockContext, email, password, subdomain, orgName).Return(assert.AnError)

		// call the handler
		handler.Register(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":""}`, rr.Body.String())
	})

	t.Run("should register", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.LetterN(30)
		orgName := gofakeit.Company()
		payload := fmt.Sprintf(`{"email": "%s","password":"%s","organization_subdomain":"%s","organization_name":"%s"}`,
			email, password, subdomain, orgName,
		)
		req, err := http.NewRequest(http.MethodPost, registerPath, strings.NewReader(payload))
		require.NoError(t, err)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Register", fake.MockContext, email, password, subdomain, orgName).Return(nil)

		// call the handler
		handler.Register(rr, req)

		// check the result
		require.Equal(t, http.StatusCreated, rr.Code)
		assert.Empty(t, rr.Body.String())
	})
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	t.Run("should return error when subdomain is invalid", func(t *testing.T) {
		t.Parallel()

		subdomain := "invalid subdomain"
		req, err := http.NewRequest(http.MethodPost, loginPath, nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Login(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"subdomain can only contain alphanumeric characters"}`, rr.Body.String())
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		email := "invalid email"
		password := validPassword
		subdomain := gofakeit.LetterN(30)

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Login(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"email must be a valid email address"}`, rr.Body.String())
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := "@2nR"
		subdomain := gofakeit.LetterN(30)

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Login(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"password must be at least 8 characters in length"}`, rr.Body.String())
	})

	t.Run("should return error when service call fails", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.LetterN(30)

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Login", fake.MockContext, subdomain, email, password).Return("", assert.AnError)

		// call the handler
		handler.Login(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":""}`, rr.Body.String())
	})

	t.Run("should return unauthorized when invalid credentials", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.LetterN(30)

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Login", fake.MockContext, subdomain, email, password).Return("", auth.ErrInvalidCredentials)

		// call the handler
		handler.Login(rr, req)

		// check the result
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.JSONEq(t, `{"error":"email or password is invalid"}`, rr.Body.String())
	})

	t.Run("should return unauthorized when user is disabled", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.LetterN(30)

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Login", fake.MockContext, subdomain, email, password).Return("", auth.ErrUserDisabled)

		// call the handler
		handler.Login(rr, req)

		// check the result
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.JSONEq(t, `{"error":"user is disabled"}`, rr.Body.String())
	})

	t.Run("should return unauthorized when organization is disabled", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.LetterN(30)

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Login", fake.MockContext, subdomain, email, password).Return("", organization.ErrOrganizationDisabled)

		// call the handler
		handler.Login(rr, req)

		// check the result
		require.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.JSONEq(t, `{"error":"organization is disabled"}`, rr.Body.String())
	})

	t.Run("should login", func(t *testing.T) {
		t.Parallel()

		jwt := gofakeit.UUID()
		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.LetterN(30)

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Login", fake.MockContext, subdomain, email, password).Return(jwt, nil)

		// call the handler
		handler.Login(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
		assert.Equal(
			t,
			fmt.Sprintf("jwt_session_id=%s; Max-Age=%d; HttpOnly; Secure; SameSite=Strict",
				jwt, int(auth.SessionTTLDuration.Seconds())),
			rr.Header().Get("Set-Cookie"),
		)
	})
}

func TestHandler_Logout(t *testing.T) {
	t.Parallel()

	t.Run("should logout successfully", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodPost, logoutPath, nil)
		require.NoError(t, err)

		// set required values in request context
		subdomain := gofakeit.LetterN(30)
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		ctx := req.Context()
		ctx = context.WithValue(ctx, request.CtxOrgSubdomainKey, subdomain)
		ctx = context.WithValue(ctx, request.CtxUserIDKey, userID)
		ctx = context.WithValue(ctx, request.CtxOrgIDKey, orgID)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
		req = req.WithContext(ctx)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Logout", fake.MockContext, userID, orgID).Return(nil)

		// call the handler
		handler.Logout(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(
			t,
			"jwt_session_id=; Max-Age=0; HttpOnly; Secure; SameSite=Strict",
			rr.Header().Get("Set-Cookie"),
		)
	})

	t.Run("should return error when service call fails", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodPost, logoutPath, nil)
		require.NoError(t, err)

		// set required values in request context
		subdomain := gofakeit.LetterN(30)
		userID := gofakeit.Int64()
		orgID := gofakeit.Int64()
		ctx := req.Context()
		ctx = context.WithValue(ctx, request.CtxOrgSubdomainKey, subdomain)
		ctx = context.WithValue(ctx, request.CtxUserIDKey, userID)
		ctx = context.WithValue(ctx, request.CtxOrgIDKey, orgID)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		ctx = context.WithValue(ctx, chi.RouteCtxKey, routeContext)
		req = req.WithContext(ctx)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// mock the service calls
		mockService.On("Logout", fake.MockContext, userID, orgID).Return(assert.AnError)

		// call the handler
		handler.Logout(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error":""}`, rr.Body.String())
	})

	t.Run("should return bad request when user-id is not found in the request context", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodPost, logoutPath, nil)
		require.NoError(t, err)

		// set required values in request context
		subdomain := gofakeit.LetterN(30)
		orgID := gofakeit.Int64()
		ctx := req.Context()
		ctx = context.WithValue(ctx, request.CtxOrgSubdomainKey, subdomain)
		ctx = context.WithValue(ctx, request.CtxOrgIDKey, orgID)
		req = req.WithContext(ctx)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Logout(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"user id not found in the request context: invalid context"}`, rr.Body.String())
	})

	t.Run("should return bad request when org-id is not found in the request context", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodPost, logoutPath, nil)
		require.NoError(t, err)

		// set required values in request context
		subdomain := gofakeit.LetterN(30)
		userID := gofakeit.Int64()
		ctx := req.Context()
		ctx = context.WithValue(ctx, request.CtxOrgSubdomainKey, subdomain)
		ctx = context.WithValue(ctx, request.CtxUserIDKey, userID)
		req = req.WithContext(ctx)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

		// call the handler
		handler.Logout(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error":"org id not found in the request context: invalid context"}`, rr.Body.String())
	})
}
