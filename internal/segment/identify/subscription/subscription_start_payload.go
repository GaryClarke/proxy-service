package subscription

import (
	"fmt"
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/segment"
	"github.com/garyclarke/proxy-service/internal/segment/identify"
	"github.com/garyclarke/proxy-service/internal/validator"
	"github.com/segmentio/analytics-go"
)

// Compile-time assertion: ensure *SubscriptionStartPayload implements IdentifyConverter.
var _ identify.IdentifyConverter = (*SubscriptionStartPayload)(nil)

// SubscriptionStartPayload represents the data for a subscription start identify event.
type SubscriptionStartPayload struct {
	UserID           string
	BrandCode        brand.Brand
	AccountGuid      string
	Subscribed       bool
	SubscriptionID   string
	AirshipChannelID *string
	AirshipID        *string
	AutoRenewEnabled *bool
}

func (p *SubscriptionStartPayload) ToIdentify() analytics.Identify {
	fb := segment.NewFieldBuilder()
	fb.SetIfNotEmpty(fmt.Sprintf("acc_%s_guid", p.BrandCode), p.AccountGuid)
	fb.SetIfNotEmpty(fmt.Sprintf("app_%s_sub", p.BrandCode), p.Subscribed)
	fb.SetIfNotEmpty(fmt.Sprintf("app_%s_sub_id", p.BrandCode), p.SubscriptionID)
	fb.SetIfNotEmpty(fmt.Sprintf("app_%s_auto_renew_status", p.BrandCode), p.AutoRenewEnabled)
	fb.SetIfNotEmpty(fmt.Sprintf("%s_airship_channel_id", p.BrandCode), p.AirshipChannelID)
	fb.SetIfNotEmpty(fmt.Sprintf("acc_%s_airship_id", p.BrandCode), p.AirshipID)

	c := &analytics.Context{
		Extra: map[string]interface{}{
			"brand_code": string(p.BrandCode),
		},
	}

	return analytics.Identify{
		UserId:  p.UserID,
		Traits:  fb.Traits(),
		Context: c,
	}
}

// Validate applies custom rules using the shared validator.Validator.
//
// - userId, brandCode, accountGuid & subscriptionId must be non‐blank.
// - brandCode must be one of the defined allowed values
func (p *SubscriptionStartPayload) Validate() error {
	v := validator.Validator{}

	v.CheckField(validator.NotBlank(p.UserID), "userId", "must not be blank")
	v.CheckField(validator.NotBlank(p.BrandCode), "brandCode", "must not be blank")
	v.CheckField(validator.PermittedValue(p.BrandCode, brand.AllBrands...),
		"brandCode", fmt.Sprintf("must be one of %v, got %q", brand.AllBrands, p.BrandCode))
	v.CheckField(validator.NotBlank(p.AccountGuid), "accountGuid", "must not be blank")
	v.CheckField(validator.NotBlank(p.SubscriptionID), "subscriptionId", "must not be blank")

	if !v.Valid() {
		return fmt.Errorf("validation failed: %v", v.FieldErrors)
	}
	return nil
}
