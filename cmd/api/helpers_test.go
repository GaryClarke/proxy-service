package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	app := newTestApplication(t, false)

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

type testReadStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestReadJSON(t *testing.T) {
	app := newTestApplication(t, false)

	tests := []struct {
		name          string
		input         string
		expectedError string
		expectedData  *testReadStruct
	}{
		{
			name:          "Valid JSON",
			input:         `{"name": "Alice", "age": 30}`,
			expectedError: "",
			expectedData:  &testReadStruct{Name: "Alice", Age: 30},
		},
		{
			name:          "Malformed JSON",
			input:         `{"name": "Alice", "age":30,}`,
			expectedError: "body contains badly-formed JSON",
			expectedData:  nil,
		},
		{
			name:          "Unknown Field",
			input:         `{"name": "Alice", "age": 30, "unknown": "field"}`,
			expectedError: "body contains unknown key",
			expectedData:  nil,
		},
		{
			name:          "Empty Body",
			input:         ``,
			expectedError: "body must not be empty",
			expectedData:  nil,
		},
		{
			name:          "Multiple top-level JSON Objects",
			input:         `{"name": "Alice", "age": 30} {"name": "Bob", "age": 25}`,
			expectedError: "body must only contain a single JSON value",
			expectedData:  nil,
		},
		{
			name:          "Exceed max bytes",
			input:         fmt.Sprintf(`{"name": "%s", "age": 30}`, strings.Repeat("a", 1_048_577)),
			expectedError: "body must not be larger than 1048576 bytes",
			expectedData:  nil,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request with the test input as the body.
			req, err := http.NewRequest("POST", "/", strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			// Use ResponseRecorder as a dummy writer.
			rr := httptest.NewRecorder()

			var dst testReadStruct
			err = app.readJSON(rr, req, &dst) // WHEN

			if tt.expectedError == "" {
				assert.NilFatalError(t, err)
				assert.Equal(t, dst, *tt.expectedData)
			} else {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.expectedError)
				}
				assert.StringContains(t, err.Error(), tt.expectedError)
			}
		})
	}
}
