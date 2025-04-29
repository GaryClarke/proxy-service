package subscription_test

import (
	"github.com/garyclarke/proxy-service/internal/brand"
	"testing"

	"github.com/segmentio/analytics-go"
	"github.com/stretchr/testify/assert"

	identify "github.com/garyclarke/proxy-service/internal/segment/identify/subscription"
	"github.com/garyclarke/proxy-service/internal/testutil"
)

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
				AirshipChannelID: testutil.PtrStr("chan-1"),
				AirshipID:        testutil.PtrStr("aid-1"),
				AutoRenewEnabled: testutil.PtrBool(false),
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

func TestSubscriptionStartPayload_Validate(t *testing.T) {
	tests := []struct {
		name       string
		payload    identify.SubscriptionStartPayload
		wantErr    bool
		wantFields []string
	}{
		{
			name: "happy path",
			payload: identify.SubscriptionStartPayload{
				UserID:         "user-123",
				BrandCode:      "gf",
				AccountGuid:    "acct-abc",
				Subscribed:     true,
				SubscriptionID: "sub-xyz",
				// optional pointers omitted
			},
			wantErr: false,
		},
		{
			name: "blank UserID",
			payload: identify.SubscriptionStartPayload{
				UserID:         "",
				BrandCode:      brand.GF,
				AccountGuid:    "acct-abc",
				Subscribed:     true,
				SubscriptionID: "sub-xyz",
			},
			wantErr:    true,
			wantFields: []string{"userId"},
		},
		{
			name: "invalid BrandCode",
			payload: identify.SubscriptionStartPayload{
				UserID:         "user-123",
				BrandCode:      "invalid-brand",
				AccountGuid:    "acct-abc",
				Subscribed:     true,
				SubscriptionID: "sub-xyz",
			},
			wantErr:    true,
			wantFields: []string{"brandCode"},
		},
		{
			name: "blank AccountGuid and SubscriptionID",
			payload: identify.SubscriptionStartPayload{
				UserID:         "user-123",
				BrandCode:      brand.GF,
				AccountGuid:    "",
				Subscribed:     true,
				SubscriptionID: "",
			},
			wantErr:    true,
			wantFields: []string{"accountGuid", "subscriptionId"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate()
			if tt.wantErr {
				assert.Error(t, err, "expected validation to fail")
				for _, field := range tt.wantFields {
					assert.Contains(t, err.Error(), field, "error should mention field %q", field)
				}
			} else {
				assert.NoError(t, err, "expected validation to pass")
			}
		})
	}
}
