package handler

import (
	"context"
	"fmt"
	"github.com/garyclarke/proxy-service/internal/webhook"
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
