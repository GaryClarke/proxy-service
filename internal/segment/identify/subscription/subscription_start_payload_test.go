package subscription_test

import (
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/ptr"
	track "github.com/garyclarke/proxy-service/internal/segment/track/subscription"
	"testing"

	"github.com/segmentio/analytics-go"
	"github.com/stretchr/testify/assert"

	identify "github.com/garyclarke/proxy-service/internal/segment/identify/subscription"
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
				AirshipChannelID: ptr.Str("chan-1"),
				AirshipID:        ptr.Str("aid-1"),
				AutoRenewEnabled: ptr.Bool(false),
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

func TestSubscriptionTrackPayload_Validate(t *testing.T) {
	tests := []struct {
		name       string
		payload    track.SubscriptionTrackPayload
		wantErr    bool
		wantFields []string
	}{
		{
			name: "happy path",
			payload: track.SubscriptionTrackPayload{
				Event:            "ev",
				UserID:           "u1",
				BrandCode:        brand.GF,
				AccountGuid:      "u1",
				SubscriptionID:   "s1",
				NotificationType: "nt",
				Category:         event.CategoryStart,
				ProductName:      ptr.Str("goodfood+"),
			},
			wantErr: false,
		},
		{
			name:    "missing required",
			payload: track.SubscriptionTrackPayload{},
			wantErr: true,
			wantFields: []string{
				"event", "userId", "accountGuid", "subscriptionId",
				"notificationType", "category", "brandCode", "productName",
			},
		},
		{
			name: "invalid brandCode",
			payload: track.SubscriptionTrackPayload{
				Event:            "ev",
				UserID:           "u1",
				BrandCode:        "not-a-brand", // invalid
				AccountGuid:      "u1",
				SubscriptionID:   "s1",
				NotificationType: "nt",
				Category:         event.CategoryStart,
				ProductName:      ptr.Str("goodfood+"),
			},
			wantErr:    true,
			wantFields: []string{"brandCode"},
		},
		{
			name: "missing productName only",
			payload: track.SubscriptionTrackPayload{
				Event:            "ev",
				UserID:           "u1",
				BrandCode:        brand.GF,
				AccountGuid:      "u1",
				SubscriptionID:   "s1",
				NotificationType: "nt",
				Category:         event.CategoryStart,
				// ProductName omitted
			},
			wantErr:    true,
			wantFields: []string{"productName"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				for _, f := range tt.wantFields {
					assert.Contains(t, err.Error(), f)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
