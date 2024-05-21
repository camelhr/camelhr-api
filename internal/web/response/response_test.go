package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/camelhr/camelhr-api/internal/base"
	"github.com/camelhr/camelhr-api/internal/web/response"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	t.Parallel()

	t.Run("should write an error response with status code and empty error message for nil error",
		func(t *testing.T) {
			t.Parallel()

			// create a new recorder
			rr := httptest.NewRecorder()

			// call the ErrorResponse function
			response.ErrorResponse(rr, http.StatusBadRequest, nil)

			// assert that the response is correct
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.JSONEq(t, `{"error":""}`, rr.Body.String())
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})

	t.Run("should write empty error message to response for general error and log the error message",
		func(t *testing.T) {
			t.Parallel()

			// create a new recorder
			rr := httptest.NewRecorder()

			// call the ErrorResponse function with an error
			response.ErrorResponse(rr, http.StatusInternalServerError, assert.AnError)

			// assert that the response is correct
			assert.Equal(t, http.StatusInternalServerError, rr.Code)
			assert.JSONEq(t, `{"error":""}`, rr.Body.String())
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})

	t.Run("should write error message to response but not log the error for api error without cause",
		func(t *testing.T) {
			t.Parallel()

			// create a new recorder
			rr := httptest.NewRecorder()

			// call the ErrorResponse function with an APIError
			response.ErrorResponse(rr, http.StatusBadRequest, base.NewAPIError("test error"))

			// assert that the response is correct
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.JSONEq(t, `{"error":"test error"}`, rr.Body.String())
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})

	t.Run("should write error message to response but not log the error for not-found error",
		func(t *testing.T) {
			t.Parallel()

			// create a new recorder
			rr := httptest.NewRecorder()

			// call the ErrorResponse function with a NotFoundError
			response.ErrorResponse(rr, http.StatusNotFound, base.NewNotFoundError("test not found error"))

			// assert that the response is correct
			assert.Equal(t, http.StatusNotFound, rr.Code)
			assert.JSONEq(t, `{"error":"test not found error"}`, rr.Body.String())
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})

	t.Run("should write error message to response and log the cause when error has a cause",
		func(t *testing.T) {
			t.Parallel()

			// create a new recorder
			rr := httptest.NewRecorder()

			// call the ErrorResponse function
			response.ErrorResponse(rr, http.StatusBadRequest, base.NewAPIError("test error", base.ErrorCause(assert.AnError)))

			// assert that the response is correct
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			assert.JSONEq(t, `{"error":"test error"}`, rr.Body.String())
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})

	t.Run("should write error message to response and log the error when error is wrapped",
		func(t *testing.T) {
			t.Parallel()

			// create a new recorder
			rr := httptest.NewRecorder()

			// call the ErrorResponse function with an error wrapped with WrapError
			response.ErrorResponse(rr, http.StatusInternalServerError, base.WrapError(assert.AnError))

			// assert that the response is correct
			assert.Equal(t, http.StatusInternalServerError, rr.Code)
			assert.JSONEq(t, `{"error":"assert.AnError general error for testing"}`, rr.Body.String())
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		})
}

func TestJSON(t *testing.T) {
	t.Parallel()

	t.Run("should write a JSON response with the given status code and value", func(t *testing.T) {
		t.Parallel()

		// create a new recorder
		rr := httptest.NewRecorder()

		// call the JSON function
		response.JSON(rr, http.StatusOK, map[string]string{"message": "Hello, World!"})

		// assert that the response is correct
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"message":"Hello, World!"}`, rr.Body.String())
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})

	t.Run("should return internal server error if the response cannot be encoded", func(t *testing.T) {
		t.Parallel()

		// create a new recorder
		rr := httptest.NewRecorder()

		// call the JSON function with an invalid value
		response.JSON(rr, http.StatusOK, make(chan int))

		// assert that the response is correct
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "\n", rr.Body.String()) // http.Error writes a newline char at the end of response
	})
}

func TestOK(t *testing.T) {
	t.Parallel()

	t.Run("should write an empty response with status code 200", func(t *testing.T) {
		t.Parallel()

		// create a new recorder
		rr := httptest.NewRecorder()

		// call the OK function
		response.OK(rr)

		// assert that the response is correct
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Empty(t, rr.Body.String())
	})
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	t.Run("should write an empty response with the given status code", func(t *testing.T) {
		t.Parallel()

		// create a new recorder
		rr := httptest.NewRecorder()

		// call the Empty function
		response.Empty(rr, http.StatusAccepted)

		// assert that the response is correct
		assert.Equal(t, http.StatusAccepted, rr.Code)
		assert.Empty(t, rr.Body.String())
	})
}

func TestText(t *testing.T) {
	t.Parallel()

	t.Run("should write a text response with the given status code and value", func(t *testing.T) {
		t.Parallel()

		// create a new recorder
		rr := httptest.NewRecorder()

		// call the Text function
		response.Text(rr, http.StatusOK, "Hello, World!")

		// assert that the response is correct
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "Hello, World!", rr.Body.String())
		assert.Equal(t, "text/plain; charset=utf-8", rr.Header().Get("Content-Type"))
	})
}
