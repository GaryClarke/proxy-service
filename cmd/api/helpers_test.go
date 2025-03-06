package main

import (
	"encoding/json"
	"github.com/garyclarke/proxy-service/internal/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	app := newTestApplication(t)

	// ResponseRecorder captures the response.
	rr := httptest.NewRecorder()

	// Define some test headers.
	headers := http.Header{}
	headers.Set("X-Custom-Header", "test-value")

	// Define the data we want to write.
	data := map[string]string{
		"foo": "bar",
	}

	err := app.writeJSON(rr, http.StatusOK, data, headers)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, rr.Header().Get("X-Custom-Header"), "test-value")

	// Check the response body.
	expectedBody, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal expected JSON: %v", err)
	}
	// Append a newline since writeJSON adds one.
	expectedBody = append(expectedBody, '\n')

	assert.Equal(t, strings.TrimSpace(rr.Body.String()), strings.TrimSpace(string(expectedBody)))
}
