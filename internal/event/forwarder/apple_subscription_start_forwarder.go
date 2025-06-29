package forwarder

import (
	"fmt"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/ptr"
	"github.com/garyclarke/proxy-service/internal/segment"
	identify "github.com/garyclarke/proxy-service/internal/segment/identify/subscription"
)

// AppleSubscriptionStartForwarder is a minimal implementation of the EventForwarder interface
// for subscription start events.
type AppleSubscriptionStartForwarder struct {
	client segment.Client
}

// NewAppleSubscriptionStartForwarder constructs a forwarder that will send
// Identify calls to the provided Segment client.
func NewAppleSubscriptionStartForwarder(client segment.Client) *AppleSubscriptionStartForwarder {
	return &AppleSubscriptionStartForwarder{client: client}
}

// Supports returns true if the event category matches the subscription start category.
func (f *AppleSubscriptionStartForwarder) Supports(e *event.SubscriptionEvent) bool {
	return e.Category == event.CategoryStart
}

// Forward is a stub implementation that simply returns nil.
// Real forwarding logic (e.g., mapping to a model, validating, and sending to an external service)
// will be implemented in a future branch.
func (f *AppleSubscriptionStartForwarder) Forward(e *event.SubscriptionEvent) error {
	// 1) Map
	payload := mapToSubscriptionStartPayload(e)

	// 2) Validate the mapped payload
	if err := payload.Validate(); err != nil {
		// Log and surface a 4xx-style error, or wrap and return
		return fmt.Errorf("payload validation failed: %w", err)
	}

	// 3) Send to Segment
	if err := f.client.Identify(payload.ToIdentify()); err != nil {
		return fmt.Errorf("segment identify failed: %w", err)
	}

	return nil
}

// mapToSubscriptionStartPayload translates our internal SubscriptionEvent
// into the identify.SubscriptionStartPayload that Segment expects.
func mapToSubscriptionStartPayload(e *event.SubscriptionEvent) identify.SubscriptionStartPayload {
	sub := e.Subscription
	return identify.SubscriptionStartPayload{
		UserID:           *sub.Properties.IdentityID,
		BrandCode:        sub.Brand,
		AccountGuid:      *sub.Properties.IdentityID,
		Subscribed:       true,
		SubscriptionID:   sub.JwsTransaction.OriginalTransactionID,
		AirshipChannelID: sub.AirshipChannelID,
		AirshipID:        sub.AirshipClaim,
		AutoRenewEnabled: ptr.BoolFromIntPtr(sub.JwsRenewalInfo.AutoRenewStatus),
	}
}
