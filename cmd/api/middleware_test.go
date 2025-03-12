package main

import (
	"github.com/garyclarke/proxy-service/internal/assert"
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
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Serve the request.
	handler.ServeHTTP(rr, req)

	// Check that the response has a 500 status code.
	assert.Equal(t, rr.Code, http.StatusInternalServerError)

	// Check the response body for a generic error message.
	// In tests using httptest.ResponseRecorder, you don’t need to explicitly close the Body.
	// The ResponseRecorder simulates an HTTP response and doesn’t hold any
	// open resources that require closing.
	assert.StringContains(t, rr.Body.String(), `{"error":"the server encountered a problem and could not process your request"}`)

	// Check that the "Connection" header is set to "close".
	assert.Equal(t, rr.Header().Get("Connection"), "close")
}
