package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// dummyHandler is a handler that intentionally panics.
func dummyHandler(w http.ResponseWriter, r *http.Request) {
	panic("something went wrong")
}

func TestRecoverPanicMiddleware(t *testing.T) {
	// Create a dummy application instance.
	app := newTestApplication(t, false)

	// Wrap the dummy handler with the recoverPanic middleware.
	// The HandlerFunc type is an adapter to allow the use of
	// ordinary functions as HTTP handlers...it is casting dummyHandler
	// to the type http.Handler 🧐
	handler := app.recoverPanic(http.HandlerFunc(dummyHandler))

	// Create a ResponseRecorder to capture the output.
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)

	// Building the *http.Request is a precondition for the rest of the test.
	// If this fails, there’s no point continuing—all subsequent calls would panic
	// or be meaningless. So we use require.NoError instead of assert.NoError
	// to abort immediately on failure.
	require.NoError(t, err, "should be able to build request")

	// Invoke the middleware (which should recover the panic).
	handler.ServeHTTP(rr, req)

	// It should return a 500 Internal Server Error.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check the response body for a generic error message.
	// In tests using httptest.ResponseRecorder, you don’t need to explicitly close the Body.
	// The ResponseRecorder simulates an HTTP response and doesn’t hold any
	// open resources that require closing.
	assert.Contains(t,
		rr.Body.String(),
		`{"error":"the server encountered a problem and could not process your request"}`,
	)

	// And the Connection header should be set to "close".
	assert.Equal(t, "close", rr.Header().Get("Connection"))
}
