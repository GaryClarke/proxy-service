package subscription_test

import (
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/ptr"
	"github.com/stretchr/testify/assert"
	"testing"

	track "github.com/garyclarke/proxy-service/internal/segment/track/subscription"
)

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
