package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	app := &application{
		config: config{env: "testing"},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	req := httptest.NewRequest("GET", "/healthcheck", nil)
	w := httptest.NewRecorder()

	app.healthcheckHandler(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Define expected JSON as a map
	expected := map[string]string{
		"status":      "available",
		"environment": "testing",
		"version":     "1.0.0",
	}

	// Decode actual response into a map
	var actual map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Compare maps instead of raw JSON strings
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected JSON %+v, got %+v", expected, actual)
	}
}
