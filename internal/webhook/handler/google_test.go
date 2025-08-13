package handler

import (
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/ptr"
	"github.com/garyclarke/proxy-service/internal/testutil"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
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

func TestCreateGoogleEvent_Success(t *testing.T) {
	// Set up default subscription data.
	sub := &subnotes.Subscription{
		DeveloperNotification: &subnotes.DeveloperNotification{
			PackageName: "com.example.app",
			SubscriptionNotification: &subnotes.SubscriptionNotification{
				NotificationType: 4, // PURCHASED
			},
		},
		SubscriptionPurchase: &subnotes.SubscriptionPurchase{},
		Properties: subnotes.SubscriptionProperties{
			PromotionalOfferApplied: true, // in_trial = true
		},
	}

	// Call createGoogleEvent, which retrieves the google lookup and resolves by composite key.
	subEvent, err := createGoogleEvent(sub)

	// Assertions (from lookup/google.json)
	assert.NoError(t, err)
	assert.NotNil(t, subEvent)
	assert.Equal(t, "subscription_started", subEvent.Name)
	if assert.NotNil(t, subEvent.SubStatus) {
		assert.Equal(t, "First time subscriber", *subEvent.SubStatus) // |true variant
	}
	assert.Equal(t, "CATEGORY_START", subEvent.Category)
	assert.EqualValues(t, 4, subEvent.NotificationType) // interface{} from JSON → float64; EqualValues handles it
	assert.Nil(t, subEvent.SubType)                     // google keys use "null" subtype
	assert.NotNil(t, subEvent.Subscription)
}

func TestCreateGoogleEvent_NoMatch(t *testing.T) {
	// Set up a subscription with a notification code that doesn't exist in the lookup data.
	sub := &subnotes.Subscription{
		DeveloperNotification: &subnotes.DeveloperNotification{
			PackageName: "com.example.app",
			SubscriptionNotification: &subnotes.SubscriptionNotification{
				NotificationType: 9999, // unknown code
			},
		},
		Properties: subnotes.SubscriptionProperties{
			PromotionalOfferApplied: true,
		},
	}

	// Call createGoogleEvent, expecting an error because no matching event is found.
	_, err := createGoogleEvent(sub)
	if err == nil {
		t.Fatal("expected error when no matching event is found, got nil")
	}
}

func TestResolveGoogleSubscriptionEvent(t *testing.T) {
	// Default subscription: PURCHASED (code 4), in_trial = true.
	defaultSub := &subnotes.Subscription{
		DeveloperNotification: &subnotes.DeveloperNotification{
			PackageName: "com.example.app",
			SubscriptionNotification: &subnotes.SubscriptionNotification{
				NotificationType: 4, // PURCHASED
			},
		},
		Properties: subnotes.SubscriptionProperties{
			PromotionalOfferApplied: true,
		},
	}

	// Fake lookup map for Google events.
	lookupMap := map[string]event.SubscriptionEvent{
		// Exact matches for PURCHASED
		"PURCHASED|null|true": {
			Name:             "subscription_started_trial",
			SubStatus:        ptr.Str("First time subscriber"),
			Category:         "CATEGORY_START",
			NotificationType: 4,
			SubType:          nil,
			InTrial:          ptr.Bool(true),
		},
		"PURCHASED|null|false": {
			Name:             "subscription_started",
			SubStatus:        ptr.Str("Subscriber"),
			Category:         "CATEGORY_START",
			NotificationType: 4,
			SubType:          nil,
			InTrial:          ptr.Bool(false),
		},
		// Only generic entry for RENEWED to force fallback.
		"RENEWED|null|null": {
			Name:             "subscription_renewed_generic",
			SubStatus:        ptr.Str("Subscriber"),
			Category:         "CATEGORY_RENEW",
			NotificationType: 2,
			SubType:          nil,
			InTrial:          nil,
		},
	}

	tests := []struct {
		name         string
		modify       func(s *subnotes.Subscription)
		expectedName string
		expectError  bool
	}{
		{
			name: "Exact match (PURCHASED|null|true)",
			modify: func(s *subnotes.Subscription) {
				// keep defaults
			},
			expectedName: "subscription_started_trial",
			expectError:  false,
		},
		{
			name: "Exact match (PURCHASED|null|false)",
			modify: func(s *subnotes.Subscription) {
				s.Properties.PromotionalOfferApplied = false
			},
			expectedName: "subscription_started",
			expectError:  false,
		},
		{
			name: "Fallback to generic (RENEWED|null|null)",
			modify: func(s *subnotes.Subscription) {
				// Switch to RENEWED (code 2); only generic key exists in the map.
				s.DeveloperNotification.SubscriptionNotification.NotificationType = 2
				s.Properties.PromotionalOfferApplied = true // should fall back regardless
			},
			expectedName: "subscription_renewed_generic",
			expectError:  false,
		},
		{
			name: "No match found",
			modify: func(s *subnotes.Subscription) {
				s.DeveloperNotification.SubscriptionNotification.NotificationType = 9999
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Copy defaults
			sub := &subnotes.Subscription{
				DeveloperNotification: &subnotes.DeveloperNotification{
					PackageName: defaultSub.DeveloperNotification.PackageName,
					SubscriptionNotification: &subnotes.SubscriptionNotification{
						NotificationType: defaultSub.DeveloperNotification.SubscriptionNotification.NotificationType,
					},
				},
				Properties: subnotes.SubscriptionProperties{
					PromotionalOfferApplied: defaultSub.Properties.PromotionalOfferApplied,
				},
			}

			tt.modify(sub)

			subEvent, err := resolveGoogleSubscriptionEvent(sub, lookupMap)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, subEvent)
				assert.Equal(t, tt.expectedName, subEvent.Name)
			}
		})
	}
}
