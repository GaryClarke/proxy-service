package handler

import (
	"context"
	"github.com/garyclarke/proxy-service/internal/webhook"
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
	// TODO: decode payload, lookup event, forward to Segment
	return nil
}
