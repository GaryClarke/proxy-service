package subscription

import (
	"fmt"
	"github.com/garyclarke/proxy-service/internal/segment"
	"github.com/garyclarke/proxy-service/internal/segment/identify"
	"github.com/segmentio/analytics-go"
)

// Compile-time assertion: ensure *SubscriptionStartPayload implements IdentifyConverter.
var _ identify.IdentifyConverter = (*SubscriptionStartPayload)(nil)

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
	t.SetIfNotEmpty(fmt.Sprintf("acc_%s_guid", p.BrandCode), p.AccountGuid)
	t.SetIfNotEmpty(fmt.Sprintf("app_%s_sub", p.BrandCode), p.Subscribed)
	t.SetIfNotEmpty(fmt.Sprintf("app_%s_sub_id", p.BrandCode), p.SubscriptionID)
	t.SetIfNotEmpty(fmt.Sprintf("app_%s_auto_renew_status", p.BrandCode), p.AutoRenewEnabled)
	t.SetIfNotEmpty(fmt.Sprintf("%s_airship_channel_id", p.BrandCode), p.AirshipChannelID)
	t.SetIfNotEmpty(fmt.Sprintf("acc_%s_airship_id", p.BrandCode), p.AirshipID)

	c := &analytics.Context{
		Extra: map[string]interface{}{
			"brand_code": p.BrandCode,
		},
	}

	return analytics.Identify{
		UserId:  p.UserID,
		Traits:  t.Traits,
		Context: c,
	}
}
