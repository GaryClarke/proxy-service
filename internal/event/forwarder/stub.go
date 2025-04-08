package forwarder

import (
	"fmt"
	"github.com/garyclarke/proxy-service/internal/event"
)

// StubForwarder is a dummy forwarder used for testing and development.
// It simply returns true in the Supports method and logs the event in Forward.
type StubForwarder struct{}

// Supports always returns true, allowing the stub to process any event.
func (s *StubForwarder) Supports(e *event.SubscriptionEvent) bool {
	return true
}

// Forward logs the event and returns nil to indicate success.
func (s *StubForwarder) Forward(e *event.SubscriptionEvent) error {
	fmt.Println("StubForwarder: forwarding event:", e)
	return nil
}
