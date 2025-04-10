package forwarder

import (
	"testing"

	"github.com/garyclarke/proxy-service/internal/event"
)

// TestAppleSubscriptionStartForwarder_Supports verifies that the Supports method
// returns true when the event category matches CategoryStart and false otherwise.
func TestAppleSubscriptionStartForwarder_Supports(t *testing.T) {
	// Create a new forwarder instance.
	fwd := &AppleSubscriptionStartForwarder{}

	// Define test cases.
	tests := []struct {
		name     string
		category string
		expected bool
	}{
		{
			name:     "Happy path - category matches",
			category: "CATEGORY_START",
			expected: true,
		},
		{
			name:     "Failure - category does not match",
			category: "SOMETHING_ELSE",
			expected: false,
		},
	}

	// Run tests in a table-driven manner.
	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a sample event with the specified category.
			evt := &event.SubscriptionEvent{
				Category: tt.category,
			}
			// Call the Supports method.
			got := fwd.Supports(evt)
			if got != tt.expected {
				t.Errorf("Supports() = %v; want %v", got, tt.expected)
			}
		})
	}
}
