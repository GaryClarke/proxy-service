package subscription_test

import (
	"testing"

	"github.com/segmentio/analytics-go"
	"github.com/stretchr/testify/assert"

	identify "github.com/garyclarke/proxy-service/internal/segment/identify/subscription"
)

func ptrStr(s string) *string { return &s }
func ptrBool(b bool) *bool    { return &b }

func TestSubscriptionStartPayload_ToIdentify(t *testing.T) {
	tests := []struct {
		name          string
		payload       identify.SubscriptionStartPayload
		wantTraits    analytics.Traits
		wantBrandCode string
	}{
		{
			name: "fully populated",
			payload: identify.SubscriptionStartPayload{
				UserID:           "user-123",
				BrandCode:        "gf",
				AccountGuid:      "acct-abc",
				Subscribed:       true,
				SubscriptionID:   "sub-xyz",
				AirshipChannelID: ptrStr("chan-1"),
				AirshipID:        ptrStr("aid-1"),
				AutoRenewEnabled: ptrBool(false),
			},
			wantBrandCode: "gf",
			wantTraits: analytics.Traits{
				"acc_gf_guid":              "acct-abc",
				"app_gf_sub":               true,
				"app_gf_sub_id":            "sub-xyz",
				"app_gf_auto_renew_status": false,
				"gf_airship_channel_id":    "chan-1",
				"acc_gf_airship_id":        "aid-1",
			},
		},
		{
			name: "minimal (nil pointers)",
			payload: identify.SubscriptionStartPayload{
				UserID:         "user-123",
				BrandCode:      "gf",
				AccountGuid:    "acct-abc",
				Subscribed:     true,
				SubscriptionID: "sub-xyz",
				// pointer fields left nil
			},
			wantBrandCode: "gf",
			wantTraits: analytics.Traits{
				"acc_gf_guid":   "acct-abc",
				"app_gf_sub":    true,
				"app_gf_sub_id": "sub-xyz",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.payload.ToIdentify()

			// user ID
			assert.Equal(t, tt.payload.UserID, got.UserId)
			// context
			assert.Equal(t, tt.wantBrandCode, got.Context.Extra["brand_code"])
			// traits
			assert.Equal(t, tt.wantTraits, analytics.Traits(got.Traits))
		})
	}
}
