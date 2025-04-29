package forwarder

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/garyclarke/proxy-service/internal/event"
	identify "github.com/garyclarke/proxy-service/internal/segment/identify/subscription"
	"github.com/garyclarke/proxy-service/internal/testutil"
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
	// 1) Map
	payload := mapToSubscriptionStartPayload(e)

	// dump as a Go-style struct
	spew.Dump(payload)

	// 2) Validate

	// 3) Send to Segment

	return nil
}

// mapToSubscriptionStartPayload translates our internal SubscriptionEvent
// into the identify.SubscriptionStartPayload that Segment expects.
func mapToSubscriptionStartPayload(e *event.SubscriptionEvent) identify.SubscriptionStartPayload {
	return identify.SubscriptionStartPayload{
		UserID:           *e.Subscription.Properties.IdentityID,
		BrandCode:        e.Subscription.Brand,
		AccountGuid:      *e.Subscription.Properties.IdentityID,
		Subscribed:       true,
		SubscriptionID:   e.Subscription.JwsTransaction.OriginalTransactionID,
		AirshipChannelID: e.Subscription.AirshipChannelID,
		AirshipID:        e.Subscription.AirshipClaim,
		AutoRenewEnabled: testutil.PtrBoolFromInt(e.Subscription.JwsRenewalInfo.AutoRenewStatus),
	}
}
