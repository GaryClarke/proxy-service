package subscription

import (
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/segment/track"
	"github.com/segmentio/analytics-go"
)

// Compile‐time assertion that this implements TrackConverter.
var _ track.TrackConverter = (*SubscriptionTrackPayload)(nil)

type SubscriptionTrackPayload struct {
	Event            string
	UserID           string
	BrandCode        brand.Brand
	AccountGuid      string
	SubscriptionID   string
	AirshipChannelID *string
	AirshipID        *string
	AutoRenewEnabled *bool
	Currency         *string
	Frequency        *string
	InTrial          *bool
	ProductName      *string
	WithOffer        *bool
	RenewalDate      *string
	StartDate        *string
	NotificationType string
	SubType          *string
	Category         string
}

func (p *SubscriptionTrackPayload) ToTrack() analytics.Track {
	panic("implement me")
}

func (p *SubscriptionTrackPayload) Validate() error {
	panic("implement me")
}
