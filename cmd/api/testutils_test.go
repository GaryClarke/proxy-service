package main

import (
	"bytes"
	"github.com/garyclarke/proxy-service/internal/config"
	"github.com/garyclarke/proxy-service/internal/event/forwarder"
	"github.com/garyclarke/proxy-service/internal/webhook/handler"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// testServer type embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// newTestServer helper initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

// newTestApplication returns an instance of the
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T, debug bool) *application {
	var logger *slog.Logger
	if debug {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	// Initialize the WebhookHandlers and initialize the delegator
	appleHandler := handler.NewAppleHandler([]forwarder.EventForwarder{&forwarder.AppleSubscriptionStartForwarder{}})

	handlerDelegator := handler.NewHandlerDelegator(appleHandler)

	return &application{
		config: config.Config{
			Env:             "testing",
			DebugMode:       debug,
			SegmentKey:      "segment-key",
			SegmentEndpoint: "https://segment-endpoint",
			Version:         "1.0.0",
		},
		logger:           logger,
		handlerDelegator: handlerDelegator,
	}
}

func (ts *testServer) postJSON(t *testing.T, urlPath, payload string) (int, http.Header, string) {
	req, err := http.NewRequest("POST", ts.URL+urlPath, strings.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to execute POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	body = bytes.TrimSpace(body)
	return resp.StatusCode, resp.Header, string(body)
}
