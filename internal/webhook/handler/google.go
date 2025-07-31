package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/webhook"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
	"sort"
	"time"
)

const GoogleNotification = "GoogleIAPNotification"

// GoogleHandler is a placeholder for Google RTDN support.
type GoogleHandler struct{}

// NewGoogleHandler returns a bare‐bones GoogleHandler.
func NewGoogleHandler() *GoogleHandler {
	return &GoogleHandler{}
}

func (h *GoogleHandler) supports(wh webhook.Webhook) bool {
	return wh.Meta.Notification == GoogleNotification
}

func (h *GoogleHandler) handle(ctx context.Context, wh webhook.Webhook) error {
	_, err := decodeGoogleWebhook(wh.Payload)
	if err != nil {
		return err
	}
	// todo validate brand
	// todo create event(s)
	// todo forward event(s)

	return nil
}

// decodeGoogleWebhook unmarshals and normalizes the Google RTDN payload.
// Steps to implement:
//  1. Unmarshal into subnotes.SubscriptionPayload
//  2. Map sub.Brand via brand.FromPlatformBrandID
//  3. Sort purchase lineItems by expiry desc
//  4. Pick top item’s autoRenewEnabled → sub.Properties.FreeTrial
//  5. Return &sub, nil
func decodeGoogleWebhook(payload string) (*subnotes.Subscription, error) {
	var innerPayload subnotes.SubscriptionPayload
	if err := json.Unmarshal([]byte(payload), &innerPayload); err != nil {
		return nil, fmt.Errorf("the subscription payload could not be decoded. Reason: %w", err)
	}

	brandValue, err := brand.FromPlatformBrandID(innerPayload.Payload.Subscription.DeveloperNotification.PackageName)
	if err != nil {
		return nil, fmt.Errorf("invalid brand in Google payload: %w", err)
	}
	innerPayload.Payload.Subscription.Brand = brandValue

	items := innerPayload.Payload.Subscription.SubscriptionPurchase.LineItems
	// sort line items descending by expiry time (i.e. latest first)
	sort.Slice(items, func(i, j int) bool {
		ti, err1 := time.Parse(time.RFC3339, items[i].ExpiryTime)
		tj, err2 := time.Parse(time.RFC3339, items[j].ExpiryTime)
		if err1 != nil || err2 != nil {
			// fallback to lexicographical
			return items[i].ExpiryTime > items[j].ExpiryTime
		}
		return ti.After(tj)
	})

	// Now populate the new field:
	if len(items) > 0 {
		innerPayload.Payload.Subscription.SubscriptionPurchase.AutoRenewing =
			items[0].AutoRenewingPlan.AutoRenewEnabled
	} else {
		innerPayload.Payload.Subscription.SubscriptionPurchase.AutoRenewing = false
	}

	return &innerPayload.Payload.Subscription, nil
}
