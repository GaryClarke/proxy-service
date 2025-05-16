package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/segment"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	app := newTestApplication(t, false, &segment.SpyClient{})

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
	assert.NoError(t, err)

	// Status and headers
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	assert.Equal(t, "test-value", rr.Header().Get("X-Custom-Header"))

	// Body (JSON + newline)
	expectedBody, err := json.Marshal(data)
	assert.NoError(t, err)
	expectedBody = append(expectedBody, '\n')
	assert.Equal(
		t,
		strings.TrimSpace(string(expectedBody)),
		strings.TrimSpace(rr.Body.String()),
	)
}

type testReadStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestReadJSON(t *testing.T) {
	app := newTestApplication(t, false, &segment.SpyClient{})

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
		},
		{
			name:          "Unknown Field",
			input:         `{"name": "Alice", "age": 30, "unknown": "field"}`,
			expectedError: "body contains unknown key",
		},
		{
			name:          "Empty Body",
			input:         ``,
			expectedError: "body must not be empty",
		},
		{
			name:          "Multiple top‐level JSON Objects",
			input:         `{"name": "Alice", "age": 30} {"name": "Bob", "age": 25}`,
			expectedError: "body must only contain a single JSON value",
		},
		{
			name:          "Exceed max bytes",
			input:         fmt.Sprintf(`{"name": "%s", "age": 30}`, strings.Repeat("a", 1_048_577)),
			expectedError: "body must not be larger than 1048576 bytes",
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Build request
			req, err := http.NewRequest("POST", "/", strings.NewReader(tc.input))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			var dst testReadStruct

			err = app.readJSON(rr, req, &dst)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, *tc.expectedData, dst)
			}
		})
	}
}
