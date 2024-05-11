package organization_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/camelhr/camelhr-api/internal/domains/organization"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetOrganizationByID(t *testing.T) {
	t.Parallel()

	t.Run("should return the organization by id", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodGet, "/organizations/{orgID}", nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("orgID", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		expectedBody := `{"id": 1, "subdomain": "sub1", "name": "org1",
		"suspended_at": null, "blacklisted_at": null,
		"created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}`
		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the GetOrganizationByID function
		mockService.On("GetOrganizationByID", req.Context(), int64(1)).
			Return(organization.Organization{ID: 1, Subdomain: "sub1", Name: "org1"}, nil)

		// call the GetOrganizationByID function
		handler.GetOrganizationByID(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("should return an error when the service call fails", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodGet, "/organizations/{orgID}", nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("orgID", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the GetOrganizationByID function
		mockService.On("GetOrganizationByID", req.Context(), int64(1)).
			Return(organization.Organization{}, assert.AnError)

		// call the GetOrganizationByID function
		handler.GetOrganizationByID(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})

	t.Run("should return an error when the organization ID is invalid", func(t *testing.T) {
		t.Parallel()
		// create a new request with an invalid URL parameter
		req, err := http.NewRequest(http.MethodGet, "/organizations/invalid", nil)
		require.NoError(t, err)

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// call the GetOrganizationByID function
		handler.GetOrganizationByID(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})
}

func TestHandler_GetOrganizationBySubdomain(t *testing.T) {
	t.Parallel()

	t.Run("should return the organization by subdomain", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodGet, "/organizations/subdomain/{subdomain}", nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", "sub1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		expectedBody := `{"id": 1, "subdomain": "sub1", "name": "org1",
		"suspended_at": null, "blacklisted_at": null,
		"created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}`
		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the GetOrganizationBySubdomain function
		mockService.On("GetOrganizationBySubdomain", req.Context(), "sub1").
			Return(organization.Organization{ID: 1, Subdomain: "sub1", Name: "org1"}, nil)

		// call the GetOrganizationBySubdomain function
		handler.GetOrganizationBySubdomain(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("should return an error when the service call fails", func(t *testing.T) {
		t.Parallel()
		// create a new request with a URL parameter
		req, err := http.NewRequest(http.MethodGet, "/organizations/subdomain/{subdomain}", nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("subdomain", "sub1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		mockService := organization.NewServiceMock(t)
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

func TestHandler_CreateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should create the organization", func(t *testing.T) {
		t.Parallel()
		// create a new request with a JSON payload
		payload := `{"subdomain": "sub1", "name": "org_create_test"}`
		req, err := http.NewRequest(http.MethodPost, "/organizations", strings.NewReader(payload))
		require.NoError(t, err)

		expectedBody := `{"id": 1, "subdomain": "sub1", "name": "org_create_test", 
		"suspended_at": null, "blacklisted_at": null,
		"created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}`
		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the CreateOrganization function
		mockService.On("CreateOrganization", req.Context(), organization.Organization{
			Subdomain: "sub1",
			Name:      "org_create_test",
		}).Return(int64(1), nil)

		// call the CreateOrganization function
		handler.CreateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusCreated, rr.Code)
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})

	t.Run("should return an error when the service call fails", func(t *testing.T) {
		t.Parallel()
		// create a new request with a JSON payload
		payload := `{"subdomain": "sub1", "name": "org_create_test"}`
		req, err := http.NewRequest(http.MethodPost, "/organizations", strings.NewReader(payload))
		require.NoError(t, err)

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the CreateOrganization function
		mockService.On("CreateOrganization", req.Context(), organization.Organization{
			Subdomain: "sub1",
			Name:      "org_create_test",
		}).Return(int64(0), assert.AnError)

		// call the CreateOrganization function
		handler.CreateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})

	//nolint:lll // ignore long line length
	t.Run("should return an error when the request payload is invalid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			payload string
			err     string
		}{
			{
				name:    "subdomain is missing",
				payload: `{"name": "test org pvt ltd."}`,
				err:     `{"error": "Key: 'Request.subdomain' Error:Field validation for 'subdomain' failed on the 'required' tag"}`,
			},
			{
				name:    "name is missing",
				payload: `{"subdomain": "sub1"}`,
				err:     `{"error": "Key: 'Request.name' Error:Field validation for 'name' failed on the 'required' tag"}`,
			},
			{
				name:    "subdomain is too long",
				payload: fmt.Sprintf(`{"subdomain": "%s", "name": "test org pvt ltd."}`, gofakeit.LetterN(31)),
				err:     `{"error": "Key: 'Request.subdomain' Error:Field validation for 'subdomain' failed on the 'max' tag"}`,
			},
			{
				name:    "name is too long",
				payload: fmt.Sprintf(`{"subdomain": "sub", "name": "%s"}`, gofakeit.LetterN(61)),
				err:     `{"error": "Key: 'Request.name' Error:Field validation for 'name' failed on the 'max' tag"}`,
			},
			{
				name:    "subdomain is not alphanumeric",
				payload: `{"subdomain": "sub1!", "name": "test org pvt ltd."}`,
				err:     `{"error": "Key: 'Request.subdomain' Error:Field validation for 'subdomain' failed on the 'alphanum' tag"}`,
			},
			{
				name:    "name contains non ascii characters",
				payload: `{"subdomain": "sub1", "name": "€€"}`,
				err:     `{"error": "Key: 'Request.name' Error:Field validation for 'name' failed on the 'ascii' tag"}`,
			},
		}

		for _, tt := range tests {
			// avoid loop closure issue by defining the variables here
			payload := tt.payload
			errResponse := tt.err

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				req, err := http.NewRequest(http.MethodPost, "/organizations", strings.NewReader(payload))
				require.NoError(t, err)

				mockService := organization.NewServiceMock(t)
				rr := httptest.NewRecorder()
				handler := organization.NewHandler(mockService)

				handler.CreateOrganization(rr, req)

				require.Equal(t, http.StatusBadRequest, rr.Code)
				assert.JSONEq(t, errResponse, rr.Body.String())
			})
		}
	})
}

func TestHandler_UpdateOrganization(t *testing.T) {
	t.Parallel()

	t.Run("should update the organization", func(t *testing.T) {
		t.Parallel()
		// create a new request with a JSON payload
		payload := `{"subdomain": "sub1", "name": "org_update_test"}`
		req, err := http.NewRequest(http.MethodPut, "/organizations/{orgID}", strings.NewReader(payload))
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("orgID", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the UpdateOrganization function
		mockService.On("UpdateOrganization", req.Context(), organization.Organization{
			ID:        1,
			Subdomain: "sub1",
			Name:      "org_update_test",
		}).Return(nil)

		// call the UpdateOrganization function
		handler.UpdateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("should return an error when the service call fails", func(t *testing.T) {
		t.Parallel()
		// create a new request with a JSON payload
		payload := `{"subdomain": "sub1", "name": "org_update_test"}`
		req, err := http.NewRequest(http.MethodPut, "/organizations/{orgID}", strings.NewReader(payload))
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("orgID", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the UpdateOrganization function
		mockService.On("UpdateOrganization", req.Context(), organization.Organization{
			ID:        1,
			Subdomain: "sub1",
			Name:      "org_update_test",
		}).Return(assert.AnError)

		// call the UpdateOrganization function
		handler.UpdateOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})

	//nolint:lll // ignore long line length
	t.Run("should return an error when the request payload is invalid", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			payload string
			err     string
		}{
			{
				name:    "subdomain is missing",
				payload: `{"name": "test org pvt ltd."}`,
				err:     `{"error": "Key: 'Request.subdomain' Error:Field validation for 'subdomain' failed on the 'required' tag"}`,
			},
			{
				name:    "name is missing",
				payload: `{"subdomain": "sub1"}`,
				err:     `{"error": "Key: 'Request.name' Error:Field validation for 'name' failed on the 'required' tag"}`,
			},
			{
				name:    "subdomain is too long",
				payload: fmt.Sprintf(`{"subdomain": "%s", "name": "test org pvt ltd."}`, gofakeit.LetterN(31)),
				err:     `{"error": "Key: 'Request.subdomain' Error:Field validation for 'subdomain' failed on the 'max' tag"}`,
			},
			{
				name:    "name is too long",
				payload: fmt.Sprintf(`{"subdomain": "sub", "name": "%s"}`, gofakeit.LetterN(61)),
				err:     `{"error": "Key: 'Request.name' Error:Field validation for 'name' failed on the 'max' tag"}`,
			},
			{
				name:    "subdomain is not alphanumeric",
				payload: `{"subdomain": "sub1!", "name": "test org pvt ltd."}`,
				err:     `{"error": "Key: 'Request.subdomain' Error:Field validation for 'subdomain' failed on the 'alphanum' tag"}`,
			},
			{
				name:    "name contains non ascii characters",
				payload: `{"subdomain": "sub1", "name": "€€"}`,
				err:     `{"error": "Key: 'Request.name' Error:Field validation for 'name' failed on the 'ascii' tag"}`,
			},
		}

		for _, tt := range tests {
			// avoid loop closure issue by defining the variables here
			payload := tt.payload
			errResponse := tt.err

			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				req, err := http.NewRequest(http.MethodPut, "/organizations/{orgID}", strings.NewReader(payload))
				require.NoError(t, err)

				// simulate chi's URL parameters
				reqContext := chi.NewRouteContext()
				reqContext.URLParams.Add("orgID", "1")
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

				mockService := organization.NewServiceMock(t)
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
		req, err := http.NewRequest(http.MethodDelete, "/organizations/{orgID}", nil)
		require.NoError(t, err)

		// simulate chi's URL parameters
		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("orgID", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// mock the DeleteOrganization function
		mockService.On("DeleteOrganization", req.Context(), int64(1)).
			Return(nil)

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})

	t.Run("should return an error when the service call fails", func(t *testing.T) {
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
		handler := organization.NewHandler(mockService)

		// mock the DeleteOrganization function
		mockService.On("DeleteOrganization", req.Context(), int64(1)).
			Return(assert.AnError)

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})

	t.Run("should return an error when the organization ID is invalid", func(t *testing.T) {
		t.Parallel()
		// create a new request with an invalid URL parameter
		req, err := http.NewRequest(http.MethodDelete, "/organizations/invalid", nil)
		require.NoError(t, err)

		mockService := organization.NewServiceMock(t)
		rr := httptest.NewRecorder()
		handler := organization.NewHandler(mockService)

		// call the DeleteOrganization function
		handler.DeleteOrganization(rr, req)

		// check the result
		require.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, `{"error": ""}`, rr.Body.String())
	})
}
