package main

import (
	"bytes"
	"github.com/garyclarke/proxy-service/internal/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	// Create a new instance of our application struct. For now, this just
	// contains a structured logger (which discards anything written to it).
	app := &application{
		config: config{env: "testing"},
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	// Use the httptest.NewServer() function to create a new test server
	ts := httptest.NewServer(app.routes())
	defer ts.Close()

	// The network address that the test server is listening on is contained in
	// the ts.URL field. We can use this along with the ts.Client().Get() method
	// to make a GET /healthcheck request against the test server. This returns a
	// http.Response struct containing the response.
	rs, err := ts.Client().Get(ts.URL + "/healthcheck")
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	assert.Equal(t, rs.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, string(body), `{"environment":"testing","status":"available","version":"1.0.0"}`)
}
