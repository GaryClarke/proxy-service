package forwarder

import (
	"github.com/garyclarke/proxy-service/internal/event"
)

// CategoryStart is the constant representing the start category.
// In a real application, this might be defined in a lookup or constants file.
const CategoryStart = "CATEGORY_START"

// AppleSubscriptionStartForwarder is a minimal implementation of the EventForwarder interface
// for subscription start events.
type AppleSubscriptionStartForwarder struct {
	// Future dependencies (e.g., ModelValidator, SegmentClient) can be added here.
}

// Supports returns true if the event category matches the subscription start category.
func (f *AppleSubscriptionStartForwarder) Supports(e *event.SubscriptionEvent) bool {
	return e.Category == CategoryStart
}

// Forward is a stub implementation that simply returns nil.
// Real forwarding logic (e.g., mapping to a model, validating, and sending to an external service)
// will be implemented in a future branch.
func (f *AppleSubscriptionStartForwarder) Forward(e *event.SubscriptionEvent) error {
	return nil
}
