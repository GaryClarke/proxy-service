package handler

import (
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/testutil"
	"testing"

	"github.com/garyclarke/proxy-service/internal/webhook"
	"github.com/stretchr/testify/assert"
)

func TestGoogleHandler_Supports(t *testing.T) {
	h := NewGoogleHandler()

	cases := []struct {
		name  string
		notif string
		want  bool
	}{
		{"exact match", GoogleNotification, true},
		{"wrong notif", "FooBar", false},
		{"apple notif", AppleNotification, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wh := webhook.Webhook{
				Meta: webhook.Meta{Notification: tc.notif},
			}
			got := h.supports(wh)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestDecodeGoogleWebhook_ValidPayload verifies that decodeGoogleWebhook
// correctly unmarshals and normalizes a well-formed Google payload.
func TestDecodeGoogleWebhook_ValidPayload(t *testing.T) {
	// 1) build the raw inner JSON for subnotes.SubscriptionPayload
	raw := testutil.BuildGoogleInnerPayload(
		4,                     // notificationType = PURCHASED
		false,                 // inTrial = false
		true,                  // autoRenew = true
		"100000123456789",     // subscriptionID
		"uk.co.bbc.goodfood2", // brand
	)

	// 2) decode it
	sub, err := decodeGoogleWebhook(raw)

	// 3) basic unmarshalling assertions
	assert.NoError(t, err, "expected no error on valid payload")
	assert.NotNil(t, sub, "expected a non-nil Subscription result")

	// 4) brand was mapped correctly
	assert.Equal(t, brand.GF, sub.Brand, "expected brand.GF from packageName")

	// 5) properties
	assert.Equal(t, "active", sub.Properties.MemberStatus)
	assert.Equal(t, "APP - Android", sub.Properties.Entitlement)
	assert.Equal(t, "Annual", sub.Properties.TimePeriod)

	// 6) developerNotification
	dn := sub.DeveloperNotification
	if assert.NotNil(t, dn, "expected DeveloperNotification to be non-nil") {
		assert.Equal(t, "1.0", dn.Version)
		assert.Equal(t, "uk.co.bbc.goodfood2", dn.PackageName)
	}

	// 7) subscriptionPurchase + autoRenewing
	sp := sub.SubscriptionPurchase
	if assert.NotNil(t, sp, "expected SubscriptionPurchase to be non-nil") {
		assert.Equal(t, "100000123456789", sp.LatestOrderID)
		assert.True(t, sp.AutoRenewing, "expected AutoRenewing to come from top line-item")
	}
}

// TestDecodeGoogleWebhook_InvalidPayload verifies that we get an error on bad JSON.
func TestDecodeGoogleWebhook_InvalidPayload(t *testing.T) {
	_, err := decodeGoogleWebhook(`{"not": "right"}`)
	if err == nil {
		t.Fatal("expected error for invalid payload, got nil")
	}
}
