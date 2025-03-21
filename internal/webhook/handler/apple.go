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
	fmt.Println("Apple handler processing webhook:", wh)
	// TODO: decode into subscription struct
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
	subscription.Brand = brand.FromPlatformBrandID(subscription.JwsTransaction.BundleID)

	return &subscription, nil
}
