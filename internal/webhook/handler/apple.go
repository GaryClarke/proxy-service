// Package handler contains webhook handlers for different platforms (Apple, Google, etc.).
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/event/forwarder"
	"github.com/garyclarke/proxy-service/internal/webhook"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
)

const AppleNotification = "AppleIAPNotification"

// AppleHandler is a basic Apple webhook handler.
type AppleHandler struct {
	forwarders []forwarder.EventForwarder
}

// NewAppleHandler returns a new instance of the AppleHandler with the provided forwarders.
func NewAppleHandler(forwarders []forwarder.EventForwarder) *AppleHandler {
	return &AppleHandler{forwarders: forwarders}
}

// supports checks whether this handler should process the webhook.
func (h *AppleHandler) supports(wh webhook.Webhook) bool {
	return wh.Meta.Notification == AppleNotification
}

func (h *AppleHandler) handle(ctx context.Context, wh webhook.Webhook) error {
	sub, err := decodeSubscriptionWebhook(wh.Payload)
	if err != nil {
		return err
	}

	if err = brand.ValidateBrand(sub.Brand); err != nil {
		return err
	}

	subEvent, err := createAppleEvent(sub)
	if err != nil {
		return err
	}

	// Forward the event using the appropriate forwarder.
	for _, f := range h.forwarders {
		if f.Supports(subEvent) {
			if err := f.Forward(subEvent); err != nil {
				return err
			}
		}
	}

	return nil
}

// decodeSubscriptionWebhook decodes the inner JSON payload into a Subscription.
// It expects the payload string to be in the format:
// {"payload": {"subscription": { ... }}}
func decodeSubscriptionWebhook(payload string) (*subnotes.Subscription, error) {
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

func createAppleEvent(sub *subnotes.Subscription) (*event.SubscriptionEvent, error) {
	// Retrieve the lookup map for Apple events.
	lookupMap, err := event.GetLookupData("apple")
	if err != nil {
		return nil, err
	}

	// Resolve the event definition using composite key fallback.
	subEvent, err := resolveAppleSubscriptionEvent(sub, lookupMap)
	if err != nil {
		return nil, err
	}

	// Attach the decoded subscription to the subEvent definition.
	subEvent.Subscription = sub

	return subEvent, nil
}

// appleCompositeKeyCandidates generates a slice of candidate composite keys
// for resolving an Apple subscription event from the given subscription data.
// The composite key is constructed using the notification type (converted to uppercase),
// the subscription's sub_type (if provided; otherwise "null"), and the in_trial flag
// (formatted as "true" or "false").
// The candidates are returned in order from most specific to most generic,
// allowing a lookup function to try each key in turn until a match is found.
func appleCompositeKeyCandidates(sub *subnotes.Subscription) []string {
	notifType := strings.ToUpper(sub.ServerData.NotificationType)

	// Handle sub_type, which is a pointer. If nil, use "null".
	var specificSubType string
	if sub.ServerData.SubType != nil {
		specificSubType = strings.ToUpper(*sub.ServerData.SubType)
	} else {
		specificSubType = "null"
	}

	inTrial := sub.Properties.PromotionalOfferApplied

	// Generate candidate keys in order from most specific to most generic.
	return []string{
		fmt.Sprintf("%s|%s|%t", notifType, specificSubType, inTrial),
		fmt.Sprintf("%s|null|%t", notifType, inTrial),
		fmt.Sprintf("%s|%s|null", notifType, specificSubType),
		fmt.Sprintf("%s|null|null", notifType),
	}
}

// resolveAppleSubscriptionEvent returns the first matching SubscriptionEvent from the lookup map.
func resolveAppleSubscriptionEvent(
	sub *subnotes.Subscription,
	lookupMap map[string]event.SubscriptionEvent,
) (*event.SubscriptionEvent, error) {
	candidates := appleCompositeKeyCandidates(sub)
	for _, key := range candidates {
		if subEvent, ok := lookupMap[key]; ok {
			// Create a copy of subEvent so we can return its address.
			//
			// In Go, when you retrieve a value from a map (e.g. "subEvent := lookupMap[key]"),
			// the returned value is not "addressable." This means you cannot take its address
			// directly (i.e. you cannot write &lookupMap[key] or &subEvent).
			// The value is a copy and does not have a stable memory address you can refer to.
			//
			// To work around this, we assign the value to a local variable (eventCopy).
			// This local variable is addressable, so we can take its address (i.e. &eventCopy)
			// and return a pointer to it.
			eventCopy := subEvent
			return &eventCopy, nil
		}
	}
	return nil, fmt.Errorf("no valid composite event found for subscription")
}
