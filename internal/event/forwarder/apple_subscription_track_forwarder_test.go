package forwarder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/ptr"
	"github.com/garyclarke/proxy-service/internal/timeutil"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
)

func Test_mapToSubscriptionTrackPayload(t *testing.T) {
	// Build a representative SubscriptionEvent
	sub := &subnotes.Subscription{
		Properties: subnotes.SubscriptionProperties{
			IdentityID:              ptr.Str("user-123"),
			Currency:                ptr.Str("USD"),
			TimePeriod:              "P1M",
			PromotionalOfferApplied: true,
			ProductName:             "goodfood+",
			EndDate:                 "2106170000",
			StartDate:               "2103170000",
		},
		JwsTransaction: &subnotes.JwsTransaction{
			OriginalTransactionID: "txn-789",
		},
		JwsRenewalInfo: &subnotes.JwsRenewalInfo{
			AutoRenewStatus: ptr.Int(1),
		},
		AirshipChannelID: ptr.Str("chan-X"),
		AirshipClaim:     ptr.Str("aid-Y"),
		Brand:            brand.GF,
	}

	evt := &event.SubscriptionEvent{
		Name:             "subscription_started",
		NotificationType: "SUBSCRIBED",
		SubType:          ptr.Str("INITIAL_BUY"),
		Category:         event.CategoryStart,
		Subscription:     sub,
	}

	got := mapToSubscriptionTrackPayload(evt)

	// Scalars
	assert.Equal(t, "subscription_started", got.Event)
	assert.Equal(t, "user-123", got.UserID)
	assert.Equal(t, brand.GF, got.BrandCode)
	assert.Equal(t, "txn-789", got.SubscriptionID)
	assert.Equal(t, ptr.Str("USD"), got.Currency)
	assert.Equal(t, event.CategoryStart, got.Category)

	// Pointers
	assert.Equal(t, "user-123", got.AccountGuid)
	assert.Equal(t, ptr.Str("chan-X"), got.AirshipChannelID)
	assert.Equal(t, ptr.Str("aid-Y"), got.AirshipID)
	assert.Equal(t, ptr.Bool(true), got.InTrial)
	assert.Equal(t, ptr.BoolFromIntPtr(ptr.Int(1)), got.AutoRenewEnabled)
	assert.Equal(t, ptr.Str("P1M"), got.Frequency)
	assert.Equal(t, ptr.Str("goodfood+"), got.ProductName)
	assert.Equal(t, ptr.Bool(true), got.WithOffer)
	assert.Equal(t, ptr.Str("INITIAL_BUY"), got.SubType)

	// Dates via timeutil
	wantStart := timeutil.ParseSubscriptionDate("2103170000")
	wantRenew := timeutil.ParseSubscriptionDate("2106170000")
	assert.Equal(t, wantStart, got.StartDate)
	assert.Equal(t, wantRenew, got.RenewalDate)
}
