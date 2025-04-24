// Package handler tests the Delegator’s behavior in dispatching to WebhookHandlers.
package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/garyclarke/proxy-service/internal/webhook"
)

// fakeHandler implements the WebhookHandler interface.
type fakeHandler struct {
	supportsCalled bool
	handleCalled   bool
	shouldSupport  bool
	handleErr      error
}

func (h *fakeHandler) supports(wh webhook.Webhook) bool {
	h.supportsCalled = true
	return h.shouldSupport
}

func (h *fakeHandler) handle(ctx context.Context, wh webhook.Webhook) error {
	h.handleCalled = true
	return h.handleErr
}

func TestDelegator_Delegate(t *testing.T) {
	// Create a fake webhook to test with.
	testWebhook := webhook.Webhook{
		Payload: "fake-payload",
		Meta: webhook.Meta{
			Client:       "Sub Notifications API",
			Notification: "AppleIAPNotification",
		},
	}

	// Define table-driven test cases
	tests := []struct {
		name                 string
		handlers             []*fakeHandler
		expectedHandleCalled []bool
		expectedError        error
	}{
		{
			name: "Handler supports and returns no error",
			handlers: []*fakeHandler{
				{shouldSupport: true, handleErr: nil},
			},
			expectedHandleCalled: []bool{true},
			expectedError:        nil,
		},
		{
			name: "Handler supports and returns an error",
			handlers: []*fakeHandler{
				{shouldSupport: true, handleErr: errors.New("oops")},
			},
			expectedHandleCalled: []bool{true},
			expectedError:        errors.New("oops"),
		},
		{
			name: "Handler does not support",
			handlers: []*fakeHandler{
				{shouldSupport: false},
			},
			expectedHandleCalled: []bool{false},
			expectedError:        nil,
		},
		{
			name: "Multiple handlers|all support",
			handlers: []*fakeHandler{
				{shouldSupport: true, handleErr: nil},
				{shouldSupport: true, handleErr: nil},
				{shouldSupport: false, handleErr: nil},
			},
			expectedHandleCalled: []bool{true, true, false},
			expectedError:        nil,
		},
	}

	for _, tt := range tests {
		tt := tt // capture the range var
		t.Run(tt.name, func(t *testing.T) {
			// Convert the slice of *fakeHandler to a slice of handler.WebhookHandler.
			// Note: Although each element of type *fakeHandler implements the WebhookHandler interface,
			// a slice of *fakeHandler (i.e. []*fakeHandler) is not automatically treated as a slice of
			// WebhookHandler ([]WebhookHandler). In Go, slices are not covariant, meaning you cannot
			// directly convert one slice type to another even if their element types are compatible.
			// Instead, we need to iterate over the slice and append each element to a new slice of
			// the interface type...which does kinda suck
			var handlers []WebhookHandler
			for _, fh := range tt.handlers {
				handlers = append(handlers, fh)
			}

			delegator := NewHandlerDelegator(handlers...)

			err := delegator.Delegate(context.Background(), testWebhook)

			if tt.expectedError != nil {
				// Since an unexpected error means there’s no point in continuing the rest of the sub‑test
				// we use require instead of assert
				require.EqualError(t, err, tt.expectedError.Error())

			} else {
				assert.NoError(t, err)
			}

			// Check that each fake's handle method was called as expected.
			for i, fh := range tt.handlers {
				assert.Equal(
					t,
					tt.expectedHandleCalled[i],
					fh.handleCalled,
					"handler #%d handleCalled mismatch", i,
				)
			}
		})
	}
}
