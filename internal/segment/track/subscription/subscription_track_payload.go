package subscription

import (
	"fmt"
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/segment"
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
	Platform         *string
	Status           *string
}

// ToTrack maps this payload into a Segment analytics.Track call.
// It mirrors the PHP SubscriptionModel::toArray() structure.
func (p *SubscriptionTrackPayload) ToTrack() analytics.Track {
	traits := segment.NewFieldBuilder()
	props := segment.NewFieldBuilder()

	// ─── Traits ─────────────────────────────────────
	traits.SetIfNotEmpty(fmt.Sprintf("acc_%s_guid", p.BrandCode), p.AccountGuid)
	traits.SetIfNotEmpty(fmt.Sprintf("app_%s_sub_id", p.BrandCode), p.SubscriptionID)
	traits.SetIfNotEmpty(fmt.Sprintf("%s_airship_channel_id", p.BrandCode), p.AirshipChannelID)
	traits.SetIfNotEmpty(fmt.Sprintf("acc_%s_airship_id", p.BrandCode), p.AirshipID)

	// ─── Properties ────────────────────────────────
	props.SetIfNotEmpty("airship_channel_id", p.AirshipChannelID)
	props.SetIfNotEmpty("airship_id", p.AirshipID)
	props.SetIfNotEmpty("sub_id", p.SubscriptionID)
	props.SetIfNotEmpty("brand_code", string(p.BrandCode))
	props.SetIfNotEmpty("sub_auto_renew_status", p.AutoRenewEnabled)
	props.SetIfNotEmpty("sub_currency", p.Currency)
	props.SetIfNotEmpty("sub_frequency", p.Frequency)
	props.SetIfNotEmpty("sub_in_trial", p.InTrial)
	props.SetIfNotEmpty("sub_product_name", p.ProductName)
	props.SetIfNotEmpty("sub_renewal_date", p.RenewalDate)
	props.SetIfNotEmpty("sub_start_date", p.StartDate)
	props.SetIfNotEmpty("sub_type", p.SubType)
	props.SetIfNotEmpty("sub_with_offer", p.WithOffer)
	props.SetIfNotEmpty("platform", p.Platform)
	props.SetIfNotEmpty("sub_status", p.Status)

	c := &analytics.Context{
		Extra: map[string]interface{}{
			"brand_code": string(p.BrandCode),
			"traits":     traits.ToMap(),
		},
	}

	return analytics.Track{
		Event:      p.Event,
		UserId:     p.UserID,
		Properties: props.Properties(),
		Context:    c,
	}
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
