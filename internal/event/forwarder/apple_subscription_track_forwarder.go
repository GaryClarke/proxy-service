package forwarder

import (
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/ptr"
	"github.com/garyclarke/proxy-service/internal/segment"
	track "github.com/garyclarke/proxy-service/internal/segment/track/subscription"
	"github.com/garyclarke/proxy-service/internal/timeutil"
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
	sub := e.Subscription
	return track.SubscriptionTrackPayload{
		Event:            e.Name,
		UserID:           *sub.Properties.IdentityID,
		BrandCode:        sub.Brand,
		AccountGuid:      *sub.Properties.IdentityID,
		SubscriptionID:   sub.JwsTransaction.OriginalTransactionID,
		AirshipChannelID: sub.AirshipChannelID,
		AirshipID:        sub.AirshipClaim,
		AutoRenewEnabled: ptr.BoolFromIntPtr(sub.JwsRenewalInfo.AutoRenewStatus),
		Currency:         sub.Properties.Currency,
		Frequency:        ptr.Str(sub.Properties.TimePeriod),
		InTrial:          ptr.Bool(sub.Properties.PromotionalOfferApplied),
		ProductName:      ptr.Str(sub.Properties.ProductName),
		WithOffer:        ptr.Bool(sub.Properties.PromotionalOfferApplied),
		RenewalDate:      timeutil.ParseSubscriptionDate(sub.Properties.EndDate),
		StartDate:        timeutil.ParseSubscriptionDate(sub.Properties.StartDate),
		NotificationType: e.NotificationType,
		SubType:          e.SubType,
		Category:         e.Category,
	}
}
