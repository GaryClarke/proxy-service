package subscription

import (
	"fmt"
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/segment/track"
	"github.com/garyclarke/proxy-service/internal/validator"
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
	v := validator.Validator{}

	// Required, non‐blank strings
	v.CheckField(validator.NotBlank(p.Event), "event", "must not be blank")
	v.CheckField(validator.NotBlank(p.UserID), "userId", "must not be blank")
	v.CheckField(validator.NotBlank(p.AccountGuid), "accountGuid", "must not be blank")
	v.CheckField(validator.NotBlank(p.SubscriptionID), "subscriptionId", "must not be blank")
	v.CheckField(validator.NotBlank(p.NotificationType), "notificationType", "must not be blank")
	v.CheckField(validator.NotBlank(p.Category), "category", "must not be blank")

	// Brand must be one of the allowed brands
	v.CheckField(validator.NotBlank(p.BrandCode), "brandCode", "must not be blank")
	v.CheckField(validator.PermittedValue(p.BrandCode, brand.AllBrands...),
		"brandCode", fmt.Sprintf("must be one of %v, got %q", brand.AllBrands, p.BrandCode),
	)

	// ProductName is required in the PHP model
	if p.ProductName == nil || !validator.NotBlank(*p.ProductName) {
		v.AddFieldError("productName", "must not be blank")
	}

	if !v.Valid() {
		return fmt.Errorf("validation failed: %v", v.FieldErrors)
	}
	return nil

}
