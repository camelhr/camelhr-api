package organization_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	getOrganizationBySubdomainPath = "/api/v1/subdomains/{subdomain}/organizations"
	updateOrganizationPath         = "/api/v1/subdomains/{subdomain}/organizations"
	deleteOrganizationPath         = "/api/v1/subdomains/{subdomain}/organizations"
)

func TestHandler_GetOrganizationBySubdomain(t *testing.T) {
	t.Parallel()

	t.Run("should return the organization by subdomain", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodGet, getOrganizationBySubdomainPath, nil)
		require.NoError(t, err)

		org := organization.Organization{
			ID:        1,
			Subdomain: "sub1",
			Name:      "org1",
			Timestamps: base.Timestamps{
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
		}
		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", org.Subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		expectedBody := fmt.Sprintf(`{"id": %d, "subdomain": "%s", "name": "%s",
		"suspended_at": null, "disabled_at": null, "created_at": "%s", "updated_at": "%s"}`,
			org.ID, org.Subdomain, org.Name, org.CreatedAt.Format(time.RFC3339Nano), org.UpdatedAt.Format(time.RFC3339Nano))
		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the GetOrganizationBySubdomain function
		mockService.On("GetOrganizationBySubdomain", req.Context(), org.Subdomain).
			Return(org, nil)

		// call the GetOrganizationBySubdomain function
		handler.GetOrganizationBySubdomain(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("should return an error when the organization is not found", func(t *testing.T) {
		t.Parallel()

		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodGet, getOrganizationBySubdomainPath, nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", "sub1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the GetOrganizationBySubdomain function
		mockService.On("GetOrganizationBySubdomain", req.Context(), "sub1").
			Return(organization.Organization{}, base.NewNotFoundError("organization not found for the given subdomain"))

		// call the GetOrganizationBySubdomain function
		handler.GetOrganizationBySubdomain(rr, req)

		// check the result
		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error": "organization not found for the given subdomain"}`, rr.Body.String())
	})

	t.Run("should return an error when the service call fails", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodGet, getOrganizationBySubdomainPath, nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", "sub1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the GetOrganizationBySubdomain function
		mockService.On("GetOrganizationBySubdomain", req.Context(), "sub1").
			Return(organization.Organization{}, assert.AnError)

		// call the GetOrganizationBySubdomain function
		handler.GetOrganizationBySubdomain(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})
}

func TestHandler_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should update the organization", func(t *testing.T) {
		t.Parallel()

		subdomain := randomOrganizationSubdomain()
		currentOrg := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: subdomain,
		}
		newOrgName := randomOrganizationName()
		// create a new request with a JSON payload
		payload := fmt.Sprintf(`{"name": "%s"}`, newOrgName)
		req, err := http.NewRequest(http.MethodPut, updateOrganizationPath, strings.NewReader(payload))
		require.NoError(t, err)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the service calls
		mockService.On("GetOrganizationBySubdomain", req.Context(), subdomain).Return(currentOrg, nil)
		mockService.On("UpdateOrganization", req.Context(), currentOrg.ID, newOrgName).Return(nil)

		// call the UpdateOrganization function
		handler.UpdateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("should return an error when get service call fails", func(t *testing.T) {
		t.Parallel()

		subdomain := randomOrganizationSubdomain()
		// create a new request with a JSON payload
		payload := fmt.Sprintf(`{"name": "%s"}`, randomOrganizationName())
		req, err := http.NewRequest(http.MethodPut, updateOrganizationPath, strings.NewReader(payload))
		require.NoError(t, err)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the service calls
		mockService.On("GetOrganizationBySubdomain", req.Context(), subdomain).
			Return(organization.Organization{}, assert.AnError)

		// call the UpdateOrganization function
		handler.UpdateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})

	t.Run("should return an error when the organization is not found", func(t *testing.T) {
		t.Parallel()

		subdomain := randomOrganizationSubdomain()
		// create a new request with a JSON payload
		payload := fmt.Sprintf(`{"name": "%s"}`, randomOrganizationName())
		req, err := http.NewRequest(http.MethodPut, updateOrganizationPath, strings.NewReader(payload))
		require.NoError(t, err)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the service calls
		mockService.On("GetOrganizationBySubdomain", req.Context(), subdomain).
			Return(organization.Organization{}, base.NewNotFoundError("organization not found for the given subdomain"))

		// call the UpdateOrganization function
		handler.UpdateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error": "organization not found for the given subdomain"}`, rr.Body.String())
	})

	t.Run("should return an error when update service call fails", func(t *testing.T) {
		t.Parallel()

		subdomain := randomOrganizationSubdomain()
		currentOrg := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: subdomain,
		}
		newOrgName := randomOrganizationName()
		// create a new request with a JSON payload
		payload := fmt.Sprintf(`{"name": "%s"}`, newOrgName)
		req, err := http.NewRequest(http.MethodPut, updateOrganizationPath, strings.NewReader(payload))
		require.NoError(t, err)

		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the service calls
		mockService.On("GetOrganizationBySubdomain", req.Context(), subdomain).Return(currentOrg, nil)
		mockService.On("UpdateOrganization", req.Context(), currentOrg.ID, newOrgName).Return(assert.AnError)

		// call the UpdateOrganization function
		handler.UpdateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})

	t.Run("should return an error when the request is invalid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			testName  string
			payload   string
			subdomain string
			err       string
		}{
			{
				testName:  "subdomain is missing",
				payload:   `{"name": "test org pvt ltd."}`,
				subdomain: "",
				err:       `{"error": "subdomain is required"}`,
			},
			{
				testName:  "name is missing",
				payload:   `{}`,
				subdomain: "sub1",
				err:       `{"error": "name is a required field"}`,
			},
			{
				testName:  "subdomain is too long",
				payload:   `{"name": "test org pvt ltd."}`,
				subdomain: gofakeit.LetterN(31),
				err:       `{"error": "subdomain must be a maximum of 30 characters in length"}`,
			},
			{
				testName:  "name is too long",
				payload:   fmt.Sprintf(`{"name": "%s"}`, gofakeit.LetterN(61)),
				subdomain: "sub1",
				err:       `{"error": "name must be a maximum of 60 characters in length"}`,
			},
			{
				testName:  "subdomain is not alphanumeric",
				payload:   `{"name": "test org pvt ltd."}`,
				subdomain: "sub1!",
				err:       `{"error": "subdomain can only contain alphanumeric characters"}`,
			},
			{
				testName:  "name contains non ascii characters",
				payload:   `{"name": "€€"}`,
				subdomain: "sub1",
				err:       `{"error": "name must contain only ascii characters"}`,
			},
		}

		for _, tt := range tests {
			// avoid loop closure issue by defining the variables here
			payload := tt.payload
			errResponse := tt.err
			subdomain := tt.subdomain

			t.Run(tt.testName, func(t *testing.T) {
				t.Parallel()

				req, err := http.NewRequest(http.MethodPut, updateOrganizationPath, strings.NewReader(payload))
				require.NoError(t, err)

				// simulate chi's URL parameters
				routeContext := chi.NewRouteContext()
				routeContext.URLParams.Add("subdomain", subdomain)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

				mockService := organization.NewMockService(t)
				rr := httptest.NewRecorder()
				handler := organization.NewHandler(mockService)

				handler.UpdateOrganization(rr, req)

				require.Equal(t, http.StatusBadRequest, rr.Code)
				assert.JSONEq(t, errResponse, rr.Body.String())
			})
		}
	})
}

func TestHandler_DeleteOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should delete the organization", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodDelete, deleteOrganizationPath, nil)
		require.NoError(t, err)

		org := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: randomOrganizationSubdomain(),
		}
		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", org.Subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the DeleteOrganization function
		mockService.On("GetOrganizationBySubdomain", req.Context(), org.Subdomain).Return(org, nil)
		mockService.On("DeleteOrganization", req.Context(), org.ID).Return(nil)

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("should return an error when get service call fails", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodDelete, deleteOrganizationPath, nil)
		require.NoError(t, err)

		org := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: randomOrganizationSubdomain(),
		}
		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", org.Subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the GetOrganizationBySubdomain function
		mockService.On("GetOrganizationBySubdomain", req.Context(), org.Subdomain).
			Return(organization.Organization{}, assert.AnError)

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})

	t.Run("should return an error when the organization is not found", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodDelete, deleteOrganizationPath, nil)
		require.NoError(t, err)

		org := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: randomOrganizationSubdomain(),
		}
		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", org.Subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the GetOrganizationBySubdomain function
		mockService.On("GetOrganizationBySubdomain", req.Context(), org.Subdomain).
			Return(organization.Organization{}, base.NewNotFoundError("organization not found for the given subdomain"))

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusNotFound, rr.Code)
		assert.JSONEq(t, `{"error": "organization not found for the given subdomain"}`, rr.Body.String())
	})

	t.Run("should return an error when delete service call fails", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodDelete, deleteOrganizationPath, nil)
		require.NoError(t, err)

		org := organization.Organization{
			ID:        gofakeit.Int64(),
			Subdomain: randomOrganizationSubdomain(),
		}
		// simulate chi's URL parameters
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("subdomain", org.Subdomain)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

		mockService := organization.NewMockService(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the DeleteOrganization function
		mockService.On("GetOrganizationBySubdomain", req.Context(), org.Subdomain).Return(org, nil)
		mockService.On("DeleteOrganization", req.Context(), org.ID).
			Return(assert.AnError)

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})
}
