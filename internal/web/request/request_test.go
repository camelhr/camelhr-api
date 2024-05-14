package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/camelhr/camelhr-api/internal/web/request"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	// create a sample JSON payload
	payload := []byte(`{"name": "John Doe", "age": 30}`)

	// create a request with the payload
	req, err := http.NewRequest(http.MethodPost, "/path", bytes.NewBuffer(payload))
	require.NoError(t, err)

	// call the DecodeJSON function
	var data map[string]any
	err = request.DecodeAndValidateJSON(req.Body, &data)
	require.NoError(t, err)

	// assert that the decoded JSON matches the expected values
	assert.Equal(t, "John Doe", data["name"])
	assert.InEpsilon(t, float64(30), data["age"], 0.1)
}

func TestURLParam(t *testing.T) {
	t.Parallel()

	t.Run("should return the url param", func(t *testing.T) {
		t.Parallel()

		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/profile/{username}", func(_ http.ResponseWriter, r *http.Request) {
			// get the value of the "username" URL parameter using the URLParam function
			username := request.URLParam(r, "username")
			assert.Equal(t, "john_doe", username)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest(http.MethodGet, "/profile/john_doe", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return unescaped url param", func(t *testing.T) {
		t.Parallel()

		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/profile/{username}", func(_ http.ResponseWriter, r *http.Request) {
			// get the value of the "username" URL parameter using the URLParam function
			username := request.URLParam(r, "username")
			assert.Equal(t, "john_doe", username)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest(http.MethodGet, "/profile/john%5fdoe", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestURLParamID(t *testing.T) {
	t.Parallel()

	t.Run("should return the url param as an int64", func(t *testing.T) {
		t.Parallel()

		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/users/{userID}", func(_ http.ResponseWriter, r *http.Request) {
			// get the value of the "userID" URL parameter using the URLParam function
			id := request.URLParam(r, "userID")
			assert.Equal(t, "123", id)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest(http.MethodGet, "/users/123", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return an error if the url param is not a valid int64", func(t *testing.T) {
		t.Parallel()

		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/users/{userID}", func(_ http.ResponseWriter, r *http.Request) {
			// get the value of the "userID" URL parameter using the URLParam function
			id, err := request.URLParamID(r, "userID")
			require.Error(t, err)
			assert.Equal(t, int64(0), id)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest(http.MethodGet, "/users/abc", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return an error if the url param is zero", func(t *testing.T) {
		t.Parallel()

		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/users/{userID}", func(_ http.ResponseWriter, r *http.Request) {
			// get the value of the "userID" URL parameter using the URLParam function
			id, err := request.URLParamID(r, "userID")
			require.Error(t, err)
			assert.Equal(t, int64(0), id)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest(http.MethodGet, "/users/0", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return an error if the url param is negative", func(t *testing.T) {
		t.Parallel()

		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/users/{userID}", func(_ http.ResponseWriter, r *http.Request) {
			// get the value of the "userID" URL parameter using the URLParam function
			id, err := request.URLParamID(r, "userID")
			require.Error(t, err)
			assert.Equal(t, int64(0), id)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest(http.MethodGet, "/users/-1", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
