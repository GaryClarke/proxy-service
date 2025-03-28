package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/webhook"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
)

const AppleNotification = "AppleIAPNotification"

// AppleHandler is a basic Apple webhook handler.
type AppleHandler struct{}

// NewAppleHandler returns a new instance of the Apple Handler.
func NewAppleHandler() *AppleHandler {
	return &AppleHandler{}
}

// supports checks whether this handler should process the webhook.
func (h *AppleHandler) supports(wh webhook.Webhook) bool {
	return wh.Meta.Notification == AppleNotification
}

func (h *AppleHandler) handle(ctx context.Context, wh webhook.Webhook) error {
	sub, err := DecodeSubscriptionWebhook(wh.Payload)
	fmt.Println("Apple decoded payload:", sub)
	if err != nil {
		return err
	}

	// TODO: validate brand
	// TODO: create event struct
	// TODO: forward the event
	return nil
}

// DecodeSubscriptionWebhook decodes the inner JSON payload into a Subscription.
// It expects the payload string to be in the format:
// {"payload": {"subscription": { ... }}}
func DecodeSubscriptionWebhook(payload string) (*subnotes.Subscription, error) {
	var innerPayload subnotes.SubscriptionPayload
	if err := json.Unmarshal([]byte(payload), &innerPayload); err != nil {
		return nil, fmt.Errorf("the subscription payload could not be decoded. Reason: %w", err)
	}

	subscription := innerPayload.Payload.Subscription
	if subscription.JwsTransaction == nil {
		return nil, fmt.Errorf("missing JwsTransaction in subscription")
	}

	brandValue, err := brand.FromPlatformBrandID(subscription.JwsTransaction.BundleID)
	if err != nil {
		return nil, err
	}
	subscription.Brand = brandValue

	return &subscription, nil
}
