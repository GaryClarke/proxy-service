package main

import (
	"bytes"
	"github.com/garyclarke/proxy-service/internal/assert"
	"github.com/garyclarke/proxy-service/internal/testutil"
	"io"
	"net/http"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	// Create a new instance of our application struct.
	app := newTestApplication(t, false)

	// create a new test server
	ts := newTestServer(t, app.routes())
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

func TestWebhookHandler_AppleScenarios(t *testing.T) {
	// Define a set of test cases for various Apple subscription scenarios.
	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		//expectedIdentify interface{} // Replace with your concrete type
		//expectedTrack    interface{} // Replace with your concrete type
		// Optionally, expected error message etc.
	}{
		{
			name: "Subscription Start|Initial Buy|Not in trial",
			payload: testutil.BuildSubNotesPayload(
				"SUBSCRIBED",
				"INITIAL_BUY",
				"false",
				"100000123456789",
				"uk.co.bbc.goodfood2",
				1,
			),
			expectedStatus: http.StatusNoContent,
			// expectedIdentify: expectedIdentifySubStartModel(true), // Your helper function
			// expectedTrack:    expectedTrackModel("SUB_EVENT_STARTED", "Subscriber", false),
		},
	}

	// Create your test server using the application's routes.
	app := newTestApplication(t, false)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			code, _, _ := ts.postJSON(t, "/webhook", tt.payload)

			// Check the HTTP status code.
			if code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, code)
			}

			// Next, you can assert that your mock segment client (or any other dependency)
			// was called with the expected identify and track models.
			//
			// For instance, if your segment client is a dependency in your application and
			// you can inspect its state after the request, you could do:
			//
			// gotIdentify := app.segmentClient.GetLatestIdentifyModel()
			// assert.Equal(t, gotIdentify, tt.expectedIdentify)
			//
			// gotTrack := app.segmentClient.GetLatestTrackModel()
			// assert.Equal(t, gotTrack, tt.expectedTrack)
			//
			// These functions would be part of your fake/mock segment client that you inject
			// into your application for testing.
		})
	}
}
