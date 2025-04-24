package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/testutil"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
)

func TestDecodeSubscriptionWebhook_ValidPayload(t *testing.T) {
	// Build a valid payload using the shared test utility function.
	payload := testutil.BuildSubNotesPayload(
		"SUBSCRIBED",
		"INITIAL_BUY",
		"100000123456789",
		"uk.co.bbc.goodfood2",
		false,
		1,
	)

	// Decode the payload.
	sub, err := decodeSubscriptionWebhook(payload)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, brand.GF, sub.Brand)

	// Properties.
	assert.Equal(t, "active", sub.Properties.MemberStatus)
	assert.Equal(t, "APP - iOS", sub.Properties.Entitlement)
	assert.Equal(t, "Annual", sub.Properties.TimePeriod)

	// JwsTransaction.
	assert.NotNil(t, sub.JwsTransaction)
	assert.Equal(t, "uk.co.bbc.goodfood2", sub.JwsTransaction.BundleID)
	assert.Equal(t, "INITIAL_BUY", sub.JwsTransaction.Type)

	// JwsRenewalInfo.
	assert.NotNil(t, sub.JwsRenewalInfo)
	assert.Equal(t, "Sandbox", sub.JwsRenewalInfo.Environment)

	// ServerData.
	assert.NotNil(t, sub.ServerData)
	assert.Equal(t, "SUBSCRIBED", sub.ServerData.NotificationType)
	assert.Equal(t, "INITIAL_BUY", *sub.ServerData.SubType)
}

func TestDecodeSubscriptionWebhook_InvalidPayload(t *testing.T) {
	invalidPayload := `{"invalid": "data"}`
	_, err := decodeSubscriptionWebhook(invalidPayload)
	if err == nil {
		t.Fatal("expected error for invalid payload, got nil")
	}
}

func TestCreateAppleEvent_Success(t *testing.T) {
	// Set up default subscription data.
	defaultNotif := "subscribed"
	initialBuy := "initial_buy"
	sub := &subnotes.Subscription{
		ServerData: &subnotes.ServerData{
			NotificationType: defaultNotif,
			SubType:          testutil.PtrStr(initialBuy),
		},
		Properties: subnotes.SubscriptionProperties{
			PromotionalOfferApplied: true,
		},
	}

	// Call createAppleEvent, which now internally retrieves the lookup map.
	subEvent, err := createAppleEvent(sub)
	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, subEvent)
	assert.Equal(t, "subscription_started", subEvent.Name)
	assert.Equal(t, "Trialist", *subEvent.SubStatus)
	assert.Equal(t, "CATEGORY_START", subEvent.Category)
	assert.Equal(t, "SUBSCRIBED", subEvent.NotificationType)
	assert.Equal(t, "INITIAL_BUY", *subEvent.SubType)
	assert.NotNil(t, subEvent.Subscription)
}

func TestCreateAppleEvent_NoMatch(t *testing.T) {
	// Set up a subscription with a notification type that doesn't exist in the lookup data.
	sub := &subnotes.Subscription{
		ServerData: &subnotes.ServerData{
			NotificationType: "nonexistent",
			SubType:          testutil.PtrStr("initial_buy"),
		},
		Properties: subnotes.SubscriptionProperties{
			PromotionalOfferApplied: true,
		},
	}

	// Call createAppleEvent, expecting an error because no matching event is found.
	_, err := createAppleEvent(sub)
	if err == nil {
		t.Fatal("expected error when no matching event is found, got nil")
	}
}

func TestResolveAppleSubscriptionEvent(t *testing.T) {
	// Define default values.
	defaultNotif := "subscribed" // will be converted to uppercase by our candidate function.
	initialBuy := "initial_buy"

	// Default subscription with NotificationType and SubType set.
	defaultSub := &subnotes.Subscription{
		ServerData: &subnotes.ServerData{
			NotificationType: defaultNotif,
			SubType:          testutil.PtrStr(initialBuy),
		},
		Properties: subnotes.SubscriptionProperties{
			PromotionalOfferApplied: true,
		},
	}

	// Fake lookup map for Apple events.
	lookupMap := map[string]event.SubscriptionEvent{
		// Exact match candidate: "SUBSCRIBED|INITIAL_BUY|true"
		"SUBSCRIBED|INITIAL_BUY|true": {
			Name:             "subscription_started",
			SubStatus:        testutil.PtrStr("First time subscriber"),
			Category:         "CATEGORY_START",
			NotificationType: "SUBSCRIBED",
			SubType:          testutil.PtrStr("INITIAL_BUY"),
			InTrial:          testutil.PtrBool(true),
		},
		// Fallback candidate: sub_type generic, key "SUBSCRIBED|null|true"
		"SUBSCRIBED|null|true": {
			Name:             "subscription_started_generic_subtype",
			SubStatus:        testutil.PtrStr("Subscriber"),
			Category:         "CATEGORY_START",
			NotificationType: "SUBSCRIBED",
			SubType:          nil,
			InTrial:          testutil.PtrBool(true),
		},
		// Fallback candidate: in_trial generic, key "SUBSCRIBED|INITIAL_BUY|null"
		"SUBSCRIBED|INITIAL_BUY|null": {
			Name:             "subscription_started_generic_in_trial",
			SubStatus:        testutil.PtrStr("First time subscriber"),
			Category:         "CATEGORY_START",
			NotificationType: "SUBSCRIBED",
			SubType:          testutil.PtrStr("INITIAL_BUY"),
			InTrial:          nil,
		},
		// Most generic fallback: key "SUBSCRIBED|null|null"
		"SUBSCRIBED|null|null": {
			Name:             "subscription_started_generic_both",
			SubStatus:        testutil.PtrStr("Subscriber"),
			Category:         "CATEGORY_START",
			NotificationType: "SUBSCRIBED",
			SubType:          nil,
			InTrial:          nil,
		},
	}

	// Define table-driven test cases.
	tests := []struct {
		name         string
		modify       func(sub *subnotes.Subscription)
		expectedName string
		expectError  bool
	}{
		{
			name: "Exact match (SUBSCRIBED|INITIAL_BUY|true)",
			modify: func(s *subnotes.Subscription) {
				// Use default values (sub_type provided, in_trial true).
			},
			expectedName: "subscription_started",
			expectError:  false,
		},
		{
			name: "Fallback for sub_type (SUBSCRIBED|null|true)",
			modify: func(s *subnotes.Subscription) {
				s.ServerData.SubType = nil // sub_type nil should trigger fallback.
			},
			expectedName: "subscription_started_generic_subtype",
			expectError:  false,
		},
		{
			name: "Fallback for in_trial (SUBSCRIBED|INITIAL_BUY|null)",
			modify: func(s *subnotes.Subscription) {
				s.Properties.PromotionalOfferApplied = false // in_trial false should not match existing, triggering fallback.
			},
			expectedName: "subscription_started_generic_in_trial",
			expectError:  false,
		},
		{
			name: "Generic fallback (SUBSCRIBED|null|null)",
			modify: func(s *subnotes.Subscription) {
				s.ServerData.SubType = nil
				s.Properties.PromotionalOfferApplied = false // both fallback.
			},
			expectedName: "subscription_started_generic_both",
			expectError:  false,
		},
		{
			name: "No match found",
			modify: func(s *subnotes.Subscription) {
				// Change notification type so that no composite key exists.
				s.ServerData.NotificationType = "nonexistent"
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the default subscription.
			sub := &subnotes.Subscription{
				ServerData: &subnotes.ServerData{
					NotificationType: defaultSub.ServerData.NotificationType,
					SubType:          defaultSub.ServerData.SubType,
				},
				Properties: subnotes.SubscriptionProperties{
					PromotionalOfferApplied: defaultSub.Properties.PromotionalOfferApplied,
				},
			}
			// Apply modifications for this test case.
			tt.modify(sub)
			// Resolve the event.
			subEvent, err := resolveAppleSubscriptionEvent(sub, lookupMap)
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
