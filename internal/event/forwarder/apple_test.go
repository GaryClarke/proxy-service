package forwarder

import (
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/testutil"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/garyclarke/proxy-service/internal/event"
)

// TestAppleSubscriptionStartForwarder_Supports verifies that the Supports method
// returns true when the event category matches CategoryStart and false otherwise.
func TestAppleSubscriptionStartForwarder_Supports(t *testing.T) {
	// Create a new forwarder instance.
	fwd := &AppleSubscriptionStartForwarder{}

	// Define test cases.
	tests := []struct {
		name     string
		category string
		expected bool
	}{
		{
			name:     "Happy path - category matches",
			category: "CATEGORY_START",
			expected: true,
		},
		{
			name:     "Failure - category does not match",
			category: "SOMETHING_ELSE",
			expected: false,
		},
	}

	// Run tests in a table-driven manner.
	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a sample event with the specified category.
			evt := &event.SubscriptionEvent{
				Category: tt.category,
			}
			// Call the Supports method.
			got := fwd.Supports(evt)
			if got != tt.expected {
				t.Errorf("Supports() = %v; want %v", got, tt.expected)
			}
		})
	}
}

func Test_mapToSubscriptionStartPayload(t *testing.T) {
	const (
		identityID    = "user-123"
		origTransID   = "orig-trans"
		airshipChanID = "chan-abc"
		airshipID     = "aid-xyz"
	)

	tests := []struct {
		name                 string
		autoRenewRaw         *int  // input JwsRenewalInfo.AutoRenewStatus
		wantAutoRenewEnabled *bool // expected pointer bool in payload
	}{
		{"autoRenew=1", testutil.PtrInt(1), testutil.PtrBool(true)},
		{"autoRenew=0", testutil.PtrInt(0), testutil.PtrBool(false)},
		{"autoRenew=nil", nil, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// build a minimal Subscription for the test
			sub := &subnotes.Subscription{
				Properties: subnotes.SubscriptionProperties{
					IdentityID: testutil.PtrStr(identityID),
				},
				JwsTransaction: &subnotes.JwsTransaction{
					OriginalTransactionID: origTransID,
				},
				AirshipChannelID: testutil.PtrStr(airshipChanID),
				AirshipClaim:     testutil.PtrStr(airshipID),
				JwsRenewalInfo: &subnotes.JwsRenewalInfo{
					AutoRenewStatus: tc.autoRenewRaw,
				},
				Brand: brand.GF,
			}

			evt := event.SubscriptionEvent{
				Subscription: sub,
			}

			got := mapToSubscriptionStartPayload(&evt)

			// verify all fields
			assert.Equal(t, identityID, got.UserID)
			assert.Equal(t, brand.GF, got.BrandCode)
			assert.Equal(t, identityID, got.AccountGuid)
			assert.True(t, got.Subscribed)
			assert.Equal(t, origTransID, got.SubscriptionID)
			assert.Equal(t, airshipChanID, *got.AirshipChannelID)
			assert.Equal(t, airshipID, *got.AirshipID)

			if tc.wantAutoRenewEnabled == nil {
				assert.Nil(t, got.AutoRenewEnabled)
			} else {
				assert.NotNil(t, got.AutoRenewEnabled)
				assert.Equal(t, *tc.wantAutoRenewEnabled, *got.AutoRenewEnabled)
			}
		})
	}
}
