package subscription_test

import (
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/ptr"
	"github.com/segmentio/analytics-go"
	"github.com/stretchr/testify/assert"
	"testing"

	track "github.com/garyclarke/proxy-service/internal/segment/track/subscription"
)

func TestSubscriptionTrackPayload_ToTrack(t *testing.T) {
	tests := []struct {
		name          string
		payload       track.SubscriptionTrackPayload
		wantProps     analytics.Properties
		wantBrandCode string
		wantTraits    map[string]interface{}
	}{
		{
			name: "fully populated",
			payload: track.SubscriptionTrackPayload{
				Event:            "Subscription Renewed",
				UserID:           "user-123",
				BrandCode:        brand.GF,
				AccountGuid:      "acct-abc",
				SubscriptionID:   "sub-xyz",
				AirshipChannelID: ptr.Str("chan-1"),
				AirshipID:        ptr.Str("aid-1"),
				AutoRenewEnabled: ptr.Bool(false),
				Currency:         ptr.Str("GBP"),
				Frequency:        ptr.Str("monthly"),
				InTrial:          ptr.Bool(true),
				ProductName:      ptr.Str("goodfood+"),
				WithOffer:        ptr.Bool(true),
				RenewalDate:      ptr.Str("2025-06-01"),
				StartDate:        ptr.Str("2025-05-01"),
				NotificationType: "renewal",
				SubType:          ptr.Str("app"),
				Category:         "subscription",
			},
			wantBrandCode: "gf",
			wantTraits: map[string]interface{}{
				"acc_gf_guid":           "acct-abc",
				"app_gf_sub_id":         "sub-xyz",
				"gf_airship_channel_id": "chan-1",
				"acc_gf_airship_id":     "aid-1",
			},
			wantProps: analytics.Properties{
				"airship_channel_id":    "chan-1",
				"airship_id":            "aid-1",
				"sub_id":                "sub-xyz",
				"brand_code":            "gf",
				"sub_auto_renew_status": false,
				"sub_currency":          "GBP",
				"sub_frequency":         "monthly",
				"sub_in_trial":          true,
				"sub_product_name":      "goodfood+",
				"sub_renewal_date":      "2025-06-01",
				"sub_start_date":        "2025-05-01",
				"sub_type":              "app",
				"sub_with_offer":        true,
			},
		},
		{
			name: "minimal (nil pointers)",
			payload: track.SubscriptionTrackPayload{
				Event:            "Subscription Cancelled",
				UserID:           "user-123",
				BrandCode:        brand.GF,
				AccountGuid:      "acct-abc",
				SubscriptionID:   "sub-xyz",
				NotificationType: "cancel",
				Category:         "subscription",
				ProductName:      ptr.Str("goodfood+"),
			},
			wantBrandCode: "gf",
			wantTraits: map[string]interface{}{
				"acc_gf_guid":   "acct-abc",
				"app_gf_sub_id": "sub-xyz",
			},
			wantProps: analytics.Properties{
				"sub_id":           "sub-xyz",
				"brand_code":       "gf",
				"sub_product_name": "goodfood+",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.payload.ToTrack()

			assert.Equal(t, tt.payload.UserID, got.UserId)
			assert.Equal(t, tt.payload.Event, got.Event)
			assert.Equal(t, tt.wantBrandCode, got.Context.Extra["brand_code"])
			assert.Equal(t, tt.wantTraits, got.Context.Extra["traits"])
			assert.Equal(t, tt.wantProps, analytics.Properties(got.Properties))
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
			name: "missing required",
			payload: track.SubscriptionTrackPayload{
				// leave everything blank
			},
			wantErr:    true,
			wantFields: []string{"event", "userId", "accountGuid", "subscriptionId", "notificationType", "category", "brandCode", "productName"},
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
