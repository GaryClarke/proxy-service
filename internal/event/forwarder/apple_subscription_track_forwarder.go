package forwarder

import (
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/segment"
	track "github.com/garyclarke/proxy-service/internal/segment/track/subscription"
)

// AppleSubscriptionTrackForwarder emits a Track call for subscription events.
type AppleSubscriptionTrackForwarder struct {
	client segment.Client
}

// NewAppleSubscriptionTrackForwarder creates a new track forwarder.
func NewAppleSubscriptionTrackForwarder(client segment.Client) *AppleSubscriptionTrackForwarder {
	return &AppleSubscriptionTrackForwarder{client: client}
}

func (f *AppleSubscriptionTrackForwarder) Supports(e *event.SubscriptionEvent) bool {
	switch e.Category {
	case event.CategoryStart, event.CategoryStop, event.CategoryRenew, event.CategoryChange:
		return true
	default:
		return false
	}
}

// Forward is currently a stub that dumps the event.
// Later we'll map e → analytics.Track and call f.client.Track(...)
func (f *AppleSubscriptionTrackForwarder) Forward(e *event.SubscriptionEvent) error {
	// 1) Map: e → analytics.Track payload
	// 2) Validate: ensure required fields are set
	// 3) Send: f.client.Track(...)
	return nil
}

func mapToSubscriptionTrackPayload(e *event.SubscriptionEvent) track.SubscriptionTrackPayload {
	return track.SubscriptionTrackPayload{}
}
