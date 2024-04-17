package organization_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationHandler_GetOrganizationByID(t *testing.T) {
	t.Parallel()

	t.Run("should return the organization by ID", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodGet, "/organizations/{orgID}", nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("orgID", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		expectedBody := `{"id": 1, "name": "Test Organization"}`
		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewOrganizationHandler(mockService)

		// mock the GetOrganizationByID function
		mockService.On("GetOrganizationByID", req.Context(), int64(1)).
			Return(organization.Organization{ID: 1, Name: "Test Organization"}, nil)

		// call the GetOrganizationByID function
		handler.GetOrganizationByID(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("should return an error when the organization ID is invalid", func(t *testing.T) {
		t.Parallel()
		// create a new request with an invalid URL parameter
		req, err := http.NewRequest(http.MethodGet, "/organizations/invalid", nil)
		require.NoError(t, err)

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewOrganizationHandler(mockService)

		// call the GetOrganizationByID function
		handler.GetOrganizationByID(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})
}

func TestOrganizationHandler_CreateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should create the organization", func(t *testing.T) {
		t.Parallel()
		// create a new request with a JSON payload
		payload := `{"name": "Test Organization"}`
		req, err := http.NewRequest(http.MethodPost, "/organizations", strings.NewReader(payload))
		require.NoError(t, err)

		expectedBody := `{"id": 1, "name": "Test Organization"}`
		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewOrganizationHandler(mockService)

		// mock the CreateOrganization function
		mockService.On("CreateOrganization", req.Context(), organization.Organization{Name: "Test Organization"}).
			Return(int64(1), nil)

		// call the CreateOrganization function
		handler.CreateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("should return an error when the request payload is invalid", func(t *testing.T) {
		t.Parallel()
		// create a new request with an invalid JSON payload
		payload := `{"invalid": "payload"}`
		req, err := http.NewRequest(http.MethodPost, "/organizations", strings.NewReader(payload))
		require.NoError(t, err)

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewOrganizationHandler(mockService)

		// call the CreateOrganization function
		handler.CreateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})
}

func TestOrganizationHandler_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should update the organization", func(t *testing.T) {
		t.Parallel()
		// create a new request with a JSON payload
		payload := `{"name": "Test Organization"}`
		req, err := http.NewRequest(http.MethodPut, "/organizations", strings.NewReader(payload))
		require.NoError(t, err)

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewOrganizationHandler(mockService)

		// mock the UpdateOrganization function
		mockService.On("UpdateOrganization", req.Context(), organization.Organization{Name: "Test Organization"}).
			Return(nil)

		// call the UpdateOrganization function
		handler.UpdateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("should return an error when the request payload is invalid", func(t *testing.T) {
		t.Parallel()
		// create a new request with an invalid JSON payload
		payload := `{"invalid": "payload"}`
		req, err := http.NewRequest(http.MethodPut, "/organizations", strings.NewReader(payload))
		require.NoError(t, err)

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewOrganizationHandler(mockService)

		// call the UpdateOrganization function
		handler.UpdateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})
}

func TestOrganizationHandler_DeleteOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should delete the organization", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodDelete, "/organizations/{orgID}", nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("orgID", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewOrganizationHandler(mockService)

		// mock the DeleteOrganization function
		mockService.On("DeleteOrganization", req.Context(), int64(1)).
			Return(nil)

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("should return an error when the organization ID is invalid", func(t *testing.T) {
		t.Parallel()
		// create a new request with an invalid URL parameter
		req, err := http.NewRequest(http.MethodDelete, "/organizations/invalid", nil)
		require.NoError(t, err)

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewOrganizationHandler(mockService)

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})
}
