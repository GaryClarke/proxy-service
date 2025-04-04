package handler

import (
	"github.com/garyclarke/proxy-service/internal/assert"
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/events"
	"github.com/garyclarke/proxy-service/internal/testutil"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
	"testing"
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
	assert.NilFatalError(t, err)
	assert.NotNil(t, sub)
	assert.Equal(t, sub.Brand, brand.GF)

	// Assert some values in nested structs.
	assert.Equal(t, sub.Properties.MemberStatus, "active")
	assert.Equal(t, sub.Properties.Entitlement, "APP - iOS")
	assert.Equal(t, sub.Properties.TimePeriod, "Annual")

	// Assert a value in jwsTransaction.
	assert.NotNil(t, sub.JwsTransaction)
	assert.Equal(t, sub.JwsTransaction.BundleID, "uk.co.bbc.goodfood2")
	assert.Equal(t, sub.JwsTransaction.Type, "INITIAL_BUY")

	// Assert a value in jwsRenewalInfo.
	assert.NotNil(t, sub.JwsRenewalInfo)
	assert.Equal(t, sub.JwsRenewalInfo.Environment, "Sandbox")

	// Assert a value in serverData.
	assert.NotNil(t, sub.ServerData)
	assert.Equal(t, sub.ServerData.NotificationType, "SUBSCRIBED")
	assert.Equal(t, *sub.ServerData.SubType, "INITIAL_BUY")
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
	event, err := createAppleEvent(sub)
	assert.NilFatalError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, event.Name, "subscription_started")
	assert.Equal(t, *event.SubStatus, "Trialist")
	assert.Equal(t, event.Category, "CATEGORY_START")
	assert.Equal(t, event.NotificationType, "SUBSCRIBED")
	assert.Equal(t, *event.SubType, "INITIAL_BUY")
	assert.NotNil(t, event.Subscription)
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
	lookupMap := map[string]events.SubscriptionEvent{
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
			event, err := resolveAppleSubscriptionEvent(sub, lookupMap)
			if tt.expectError {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
			} else {
				assert.NilFatalError(t, err)
				assert.NotNil(t, event)
				assert.Equal(t, event.Name, tt.expectedName)
			}
		})
	}
}
