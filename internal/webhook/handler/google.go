package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/brand"
	"github.com/garyclarke/proxy-service/internal/debug"
	"github.com/garyclarke/proxy-service/internal/event"
	"github.com/garyclarke/proxy-service/internal/webhook"
	"github.com/garyclarke/proxy-service/internal/webhook/dto/subnotes"
	"sort"
	"time"
)

const GoogleNotification = "GoogleIAPNotification"

// notification codes from Google → human‐readable names matching our lookup JSON keys
var googleCodeToName = map[int]string{
	1:  "RECOVERED",
	2:  "RENEWED",
	3:  "CANCELLED",
	4:  "PURCHASED",
	5:  "ON_HOLD",
	7:  "RESTARTED",
	10: "PAUSED",
	12: "REVOKED",
	13: "EXPIRED",
}

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
	sub, err := decodeGoogleWebhook(wh.Payload)
	if err != nil {
		return err
	}
	if err = brand.ValidateBrand(sub.Brand); err != nil {
		return err
	}
	// todo create subscription event
	subEvent, err := createGoogleEvent(sub)
	if err != nil {
		return err
	}
	debug.DD(subEvent)
	// todo create change event

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

	// grab a pointer to make later code cleaner
	sub := &innerPayload.Payload.Subscription

	// --- nil‐guard the required blocks ---
	if sub.DeveloperNotification == nil {
		return nil, fmt.Errorf("invalid Google payload: missing developerNotification")
	}
	if sub.DeveloperNotification.SubscriptionNotification == nil {
		return nil, fmt.Errorf("invalid Google payload: missing subscriptionNotification")
	}
	if sub.SubscriptionPurchase == nil {
		return nil, fmt.Errorf("invalid Google payload: missing subscriptionPurchase")
	}

	// now safe to map brand
	brandValue, err := brand.FromPlatformBrandID(sub.DeveloperNotification.PackageName)
	if err != nil {
		return nil, fmt.Errorf("invalid brand in Google payload: %w", err)
	}
	sub.Brand = brandValue

	// extract and sort line items...
	items := sub.SubscriptionPurchase.LineItems
	sort.Slice(items, func(i, j int) bool {
		ti, err1 := time.Parse(time.RFC3339, items[i].ExpiryTime)
		tj, err2 := time.Parse(time.RFC3339, items[j].ExpiryTime)
		if err1 != nil || err2 != nil {
			return items[i].ExpiryTime > items[j].ExpiryTime
		}
		return ti.After(tj)
	})

	// pick top item or default
	if len(items) > 0 {
		sub.SubscriptionPurchase.AutoRenewing = items[0].AutoRenewingPlan.AutoRenewEnabled
	} else {
		sub.SubscriptionPurchase.AutoRenewing = false
	}

	return sub, nil
}

func createGoogleEvent(sub *subnotes.Subscription) (*event.SubscriptionEvent, error) {
	// 1. Retrieve the lookup map for Google events.
	lookupMap, err := event.GetLookupData("google")
	if err != nil {
		return nil, err
	}
	// 2. Resolve via composite key
	subEvent, err := resolveGoogleSubscriptionEvent(sub, lookupMap)
	if err != nil {
		return nil, err
	}

	// Attach the decoded subscription to the subEvent definition.
	subEvent.Subscription = sub

	return subEvent, err
}

// resolveGoogleSubscriptionEvent returns the first matching event from the lookup map.
func resolveGoogleSubscriptionEvent(
	sub *subnotes.Subscription,
	lookupMap map[string]event.SubscriptionEvent,
) (*event.SubscriptionEvent, error) {
	candidates := googleCompositeKeyCandidates(sub)

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

			// Copy so we can take a stable address
			eventCopy := subEvent
			return &eventCopy, nil
		}
	}

	return nil, fmt.Errorf(
		"no valid Google subscription event found for notificationType=%d inTrial=%t",
		sub.DeveloperNotification.SubscriptionNotification.NotificationType,
		sub.Properties.PromotionalOfferApplied,
	)
}

// googleCompositeKeyCandidates builds the lookup keys for a Google subscription.
// It uses the NotificationType code → name mapping, and the in_trial flag.
// Fallback order:
//  1. "<NAME>|null|<true|false>"
//  2. "<NAME>|null|null"
func googleCompositeKeyCandidates(sub *subnotes.Subscription) []string {
	notifCode := sub.DeveloperNotification.
		SubscriptionNotification.
		NotificationType

	// map code to the string used in google.json; fall back to the raw code if missing
	notifName, ok := googleCodeToName[notifCode]
	if !ok {
		notifName = fmt.Sprintf("%d", notifCode)
	}

	inTrial := sub.Properties.PromotionalOfferApplied

	return []string{
		fmt.Sprintf("%s|null|%t", notifName, inTrial),
		fmt.Sprintf("%s|null|null", notifName),
	}
}
