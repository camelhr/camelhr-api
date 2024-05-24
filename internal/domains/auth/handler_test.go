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
	"github.com/camelhr/camelhr-api/internal/tests/fake"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	registerPath = "/api/v1/auth/register"
	loginPath    = "/api/v1/subdomains/{subdomain}/auth/login"
	logoutPath   = "/api/v1/auth/logout"
)

func TestHandler_Register(t *testing.T) {
	t.Parallel()

	t.Run("should return error when email is invalid", func(t *testing.T) {
		t.Parallel()

		email := "invalid email"
		password := validPassword
		subdomain := gofakeit.Word()
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
		subdomain := gofakeit.Word()
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
		subdomain := gofakeit.Word()
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
		subdomain := gofakeit.Word()
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
		subdomain := gofakeit.Word()
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
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

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
		subdomain := gofakeit.Word()

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

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
		subdomain := gofakeit.Word()

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

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
		subdomain := gofakeit.Word()

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

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
		subdomain := gofakeit.Word()

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

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

	t.Run("should login", func(t *testing.T) {
		t.Parallel()

		jwt := gofakeit.UUID()
		email := gofakeit.Email()
		password := validPassword
		subdomain := gofakeit.Word()

		// create url-encoded form data
		form := url.Values{}
		form.Add("email", email)
		form.Add("password", password)
		req, err := http.NewRequest(http.MethodPost, loginPath, strings.NewReader(form.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

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
			fmt.Sprintf("jwt_session_id=%s; Max-Age=7200; HttpOnly; Secure; SameSite=Strict", jwt),
			rr.Header().Get("Set-Cookie"),
		)
	})
}

func TestHandler_Logout(t *testing.T) {
	t.Parallel()

	t.Run("should logout", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequest(http.MethodPost, loginPath, nil)
		require.NoError(t, err)

		mockService := auth.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := auth.NewHandler(mockService)

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
}
