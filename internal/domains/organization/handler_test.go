package organization

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOrganization(t *testing.T) {
	t.Parallel()
	t.Run("should create the organization", func(t *testing.T) {
		t.Parallel()
		// create a new request with a JSON payload
		payload := `{"name": "Test Organization"}`
		req, err := http.NewRequest("POST", "/organizations", strings.NewReader(payload))
		assert.NoError(t, err)

		// create a response recorder to capture the response
		rr := httptest.NewRecorder()

		// create a new organization handler
		handler := NewOrganizationHandler()

		// call the CreateOrganization function
		handler.CreateOrganization(rr, req)

		// check the response status code
		require.Equal(t, http.StatusCreated, rr.Code)

		// check the response body
		expectedBody := `{"id": 1, "name": "Test Organization"}`
		assert.JSONEq(t, expectedBody, rr.Body.String())
	})
}
