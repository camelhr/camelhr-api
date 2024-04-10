package request_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/camelhr/camelhr-api/internal/web/request"
)

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	// create a sample JSON payload
	payload := []byte(`{"name": "John Doe", "age": 30}`)

	// create a request with the payload
	req, err := http.NewRequest("POST", "/path", bytes.NewBuffer(payload))
	assert.NoError(t, err)

	// call the DecodeJSON function
	var data map[string]interface{}
	err = request.DecodeJSON(req.Body, &data)
	assert.NoError(t, err)

	// assert that the decoded JSON matches the expected values
	assert.Equal(t, "John Doe", data["name"])
	assert.Equal(t, float64(30), data["age"])
}

func TestURLParam(t *testing.T) {
	t.Parallel()

	t.Run("should return the url param", func(t *testing.T) {
		t.Parallel()
		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/profile/{username}", func(w http.ResponseWriter, r *http.Request) {
			// get the value of the "username" URL parameter using the URLParam function
			username := request.URLParam(r, "username")
			assert.Equal(t, "john_doe", username)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest("GET", "/profile/john_doe", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return unescaped url param", func(t *testing.T) {
		t.Parallel()
		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/profile/{username}", func(w http.ResponseWriter, r *http.Request) {
			// get the value of the "username" URL parameter using the URLParam function
			username := request.URLParam(r, "username")
			assert.Equal(t, "john_doe", username)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest("GET", "/profile/john%5fdoe", nil)
		assert.NoError(t, err)

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
		r.Get("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
			// get the value of the "userID" URL parameter using the URLParam function
			id := request.URLParam(r, "userID")
			assert.Equal(t, "123", id)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest("GET", "/users/123", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return an error if the url param is not a valid int64", func(t *testing.T) {
		t.Parallel()
		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
			// get the value of the "userID" URL parameter using the URLParam function
			id, err := request.URLParamID(r, "userID")
			assert.Error(t, err)
			assert.Equal(t, int64(0), id)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest("GET", "/users/abc", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return an error if the url param is zero", func(t *testing.T) {
		t.Parallel()
		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
			// get the value of the "userID" URL parameter using the URLParam function
			id, err := request.URLParamID(r, "userID")
			assert.Error(t, err)
			assert.Equal(t, int64(0), id)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest("GET", "/users/0", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return an error if the url param is negative", func(t *testing.T) {
		t.Parallel()
		r := chi.NewRouter()

		// define a route that expects a URL parameter
		r.Get("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
			// get the value of the "userID" URL parameter using the URLParam function
			id, err := request.URLParamID(r, "userID")
			assert.Error(t, err)
			assert.Equal(t, int64(0), id)
		})

		// create a request with the URL parameter
		req, err := http.NewRequest("GET", "/users/-1", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}
