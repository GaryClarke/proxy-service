package identify

import (
	"github.com/garyclarke/proxy-service/internal/segment"
	"github.com/segmentio/analytics-go"
)

// SubscriptionStartPayload represents the data for a subscription start identify event.
type SubscriptionStartPayload struct {
	UserID           string
	BrandCode        string
	AccountGuid      string
	Subscribed       bool
	SubscriptionID   string
	AirshipChannelID *string
	AirshipID        *string
	AutoRenewEnabled *bool
}

func (p *SubscriptionStartPayload) ToIdentify() analytics.Identify {
	t := &segment.Traits{
		Traits: make(analytics.Traits),
	}
}
