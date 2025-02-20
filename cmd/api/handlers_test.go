package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	expectedJSON := `{"status":"available","environment":"testing","version":"1.0.0"}`
	if w.Body.String() != expectedJSON {
		t.Errorf("expected JSON %q, got %q", expectedJSON, w.Body.String())
	}
}
